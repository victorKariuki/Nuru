package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/NuruProgramming/Nuru/ast"
	"github.com/NuruProgramming/Nuru/evaluator"
	"github.com/NuruProgramming/Nuru/lexer"
	"github.com/NuruProgramming/Nuru/object"
	"github.com/NuruProgramming/Nuru/parser"
)

const contentLengthPrefix = "Content-Length: "

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "nuru-dap: %v\n", err)
		os.Exit(1)
	}
}

type dapRequest struct {
	Seq       int             `json:"seq"`
	Type      string          `json:"type"`
	Command   string          `json:"command"`
	Arguments json.RawMessage `json:"arguments"`
}

type dapResponse struct {
	Seq         int             `json:"seq"`
	Type        string          `json:"type"`
	RequestSeq  int             `json:"request_seq"`
	Success     bool            `json:"success"`
	Command     string          `json:"command"`
	Body        json.RawMessage `json:"body,omitempty"`
	Message    string           `json:"message,omitempty"`
}

type dapEvent struct {
	Seq    int             `json:"seq"`
	Type   string          `json:"type"`
	Event  string          `json:"event"`
	Body   json.RawMessage `json:"body,omitempty"`
}

func run(r io.Reader, w io.Writer) error {
	bufr := bufio.NewReader(r)
	for {
		body, err := readMessage(bufr)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		var req dapRequest
		if err := json.Unmarshal(body, &req); err != nil {
			continue
		}
		if req.Type != "request" {
			continue
		}
		var respBody json.RawMessage
		var success = true
		var errMsg string
		switch req.Command {
		case "initialize":
			respBody = []byte(`{"supportsConditionalBreakpoints":false,"supportsConfigurationDoneRequest":true}`)
			sendEvent(w, "initialized", json.RawMessage(`{}`))
		case "launch":
			var params struct {
				Program string `json:"program"`
			}
			json.Unmarshal(req.Arguments, &params)
			if err := launchDebug(params.Program, w); err != nil {
				success = false
				errMsg = err.Error()
			} else {
				respBody = []byte(`{}`)
			}
		case "setBreakpoints":
			var params struct {
				Source struct {
					Path string `json:"path"`
				} `json:"source"`
				Breakpoints []struct {
					Line int `json:"line"`
				} `json:"breakpoints"`
			}
			json.Unmarshal(req.Arguments, &params)
			lines := make([]int, len(params.Breakpoints))
			for i, b := range params.Breakpoints {
				lines[i] = b.Line
			}
			setBreakpoints(params.Source.Path, lines)
			respBody = []byte(`{"breakpoints":[]}`)
		case "configurationDone":
			respBody = []byte(`{}`)
		case "continue":
			respBody = []byte(`{"allThreadsContinued":true}`)
			unblockHook()
		case "next", "stepIn", "stepOut":
			respBody = []byte(`{}`)
			unblockHook()
		case "disconnect":
			respBody = []byte(`{}`)
			unblockHook()
			sendResponse(w, req.Seq, req.Command, respBody, true, "")
			return nil
		case "threads":
			respBody = []byte(`[{"id":1,"name":"main"}]`)
		case "stackTrace":
			frames := getStackTrace()
			b, _ := json.Marshal(map[string]interface{}{"stackFrames": frames})
			respBody = b
		case "scopes":
			var params struct {
				FrameID int `json:"frameId"`
			}
			json.Unmarshal(req.Arguments, &params)
			scopes := getScopes(params.FrameID)
			b, _ := json.Marshal(map[string]interface{}{"scopes": scopes})
			respBody = b
		case "variables":
			var params struct {
				VariablesReference int `json:"variablesReference"`
			}
			json.Unmarshal(req.Arguments, &params)
			vars := getVariables(params.VariablesReference)
			b, _ := json.Marshal(map[string]interface{}{"variables": vars})
			respBody = b
		default:
			respBody = []byte(`{}`)
		}
		sendResponse(w, req.Seq, req.Command, respBody, success, errMsg)
	}
}

func readMessage(r *bufio.Reader) ([]byte, error) {
	var contentLen int
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSuffix(line, "\r\n")
		line = strings.TrimSuffix(line, "\n")
		if line == "" {
			break
		}
		if strings.HasPrefix(line, contentLengthPrefix) {
			fmt.Sscanf(line[len(contentLengthPrefix):], "%d", &contentLen)
		}
	}
	if contentLen <= 0 {
		return nil, fmt.Errorf("invalid Content-Length")
	}
	body := make([]byte, contentLen)
	_, err := io.ReadFull(r, body)
	return body, err
}

func sendResponse(w io.Writer, seq int, command string, body json.RawMessage, success bool, errMsg string) {
	resp := dapResponse{
		Seq:        seq,
		Type:       "response",
		RequestSeq: seq,
		Success:    success,
		Command:    command,
		Body:       body,
		Message:    errMsg,
	}
	b, _ := json.Marshal(resp)
	send(w, b)
}

func sendEvent(w io.Writer, event string, body json.RawMessage) {
	ev := dapEvent{Seq: 0, Type: "event", Event: event, Body: body}
	b, _ := json.Marshal(ev)
	send(w, b)
}

func send(w io.Writer, body []byte) {
	writeMu.Lock()
	defer writeMu.Unlock()
	fmt.Fprintf(w, "Content-Length: %d\r\n\r\n", len(body))
	w.Write(body)
}

var (
	writeMu      sync.Mutex
	hookUnblock  chan struct{}
	breakpointsMu sync.Mutex
	breakpoints  map[string]map[int]bool // path -> line (1-based) -> true
	stackMu      sync.Mutex
	stack        []frame
	programPath  string
)

type frame struct {
	ID   int
	Name string
	Line int
	Path string
	Env  *object.Environment
}

func setBreakpoints(path string, lines []int) {
	breakpointsMu.Lock()
	defer breakpointsMu.Unlock()
	if breakpoints == nil {
		breakpoints = make(map[string]map[int]bool)
	}
	m := make(map[int]bool)
	for _, line := range lines {
		m[line+1] = true // DAP uses 0-based, Nuru uses 1-based
	}
	breakpoints[path] = m
}

func launchDebug(program string, w io.Writer) error {
	programPath, _ = filepath.Abs(program)
	content, err := os.ReadFile(programPath)
	if err != nil {
		return err
	}
	l := lexer.New(string(content))
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		return fmt.Errorf("parse: %s", strings.Join(p.Errors(), "; "))
	}
	env := object.NewEnvironment()
	env.Set("__FILE__", &object.String{Value: programPath})
	dir := filepath.Dir(programPath)
	env.Set("__DIR__", &object.String{Value: dir})

	hookUnblock = make(chan struct{})
	stackMu.Lock()
	stack = []frame{{ID: 0, Name: "global", Line: 1, Path: programPath, Env: env}}
	stackMu.Unlock()

	evaluator.DebugStackPush = func(name string, line int, env *object.Environment) {
		stackMu.Lock()
		stack = append(stack, frame{ID: len(stack), Name: name, Line: line, Path: programPath, Env: env})
		stackMu.Unlock()
	}
	evaluator.DebugStackPop = func() {
		stackMu.Lock()
		if len(stack) > 1 {
			stack = stack[:len(stack)-1]
		}
		stackMu.Unlock()
	}
	evaluator.DebugHook = func(node ast.Node, env *object.Environment) {
		line, ok := evaluator.NodeLine(node)
		if !ok {
			return
		}
		breakpointsMu.Lock()
		lines, _ := breakpoints[programPath]
		_, hit := lines[line]
		breakpointsMu.Unlock()
		if !hit {
			return
		}
		body, _ := json.Marshal(map[string]interface{}{"reason": "breakpoint", "threadId": 1, "line": line - 1})
		sendEvent(w, "stopped", body)
		<-hookUnblock
	}

	go func() {
		evaluator.Eval(prog, env)
		sendEvent(w, "terminated", json.RawMessage(`{}`))
	}()
	return nil
}

func unblockHook() {
	breakpointsMu.Lock()
	ch := hookUnblock
	breakpointsMu.Unlock()
	if ch != nil {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func getStackTrace() []map[string]interface{} {
	stackMu.Lock()
	defer stackMu.Unlock()
	var out []map[string]interface{}
	for i, f := range stack {
		out = append(out, map[string]interface{}{
			"id":     i,
			"name":   f.Name,
			"line":   f.Line - 1,
			"source": map[string]string{"path": f.Path},
		})
	}
	return out
}

func getScopes(frameID int) []map[string]interface{} {
	stackMu.Lock()
	var env *object.Environment
	if frameID >= 0 && frameID < len(stack) {
		env = stack[frameID].Env
	}
	stackMu.Unlock()
	if env == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"name":               "Locals",
			"variablesReference": frameID + 1000,
		},
	}
}

func getVariables(ref int) []map[string]interface{} {
	if ref < 1000 {
		return nil
	}
	frameID := ref - 1000
	stackMu.Lock()
	var env *object.Environment
	if frameID >= 0 && frameID < len(stack) {
		env = stack[frameID].Env
		}
	stackMu.Unlock()
	if env == nil {
		return nil
	}
	var out []map[string]interface{}
	for _, name := range env.Names() {
		val, _ := env.Get(name)
		value := "nil"
		if val != nil {
			value = val.Inspect()
		}
		out = append(out, map[string]interface{}{
			"name":  name,
			"value": value,
		})
	}
	return out
}
