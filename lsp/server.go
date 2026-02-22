package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/NuruProgramming/Nuru/analysis"
	"github.com/NuruProgramming/Nuru/ast"
	"github.com/NuruProgramming/Nuru/format"
	"github.com/NuruProgramming/Nuru/lexer"
	"github.com/NuruProgramming/Nuru/lsp/symbols"
	"github.com/NuruProgramming/Nuru/parser"
)

// Server is the Nuru LSP server.
type Server struct {
	mu   sync.Mutex
	docs map[string]*document
	w    io.Writer // for notifications (set in Run)
}

type document struct {
	content string
	program *ast.Program
	symbols *symbols.Table
}

// NewServer creates a new LSP server.
func NewServer() *Server {
	return &Server{
		docs: make(map[string]*document),
	}
}

// Run runs the server loop reading from r and writing to w (stdio).
func (s *Server) Run(r io.Reader, w io.Writer) error {
	s.mu.Lock()
	s.w = w
	s.mu.Unlock()
	bufr := bufio.NewReader(r)
	return runReadLoop(bufr, w, s.handleMessage)
}

func (s *Server) handleMessage(body []byte) ([]byte, error) {
	var msg jsonrpcMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, err
	}
	if msg.JSONRPC != "2.0" {
		return nil, nil
	}

	switch msg.Method {
	case "initialize":
		return s.handleInitialize(msg.ID, msg.Params)
	case "initialized":
		return nil, nil
	case "textDocument/didOpen":
		return s.handleDidOpen(msg.Params)
	case "textDocument/didChange":
		return s.handleDidChange(msg.Params)
	case "textDocument/didClose":
		return s.handleDidClose(msg.Params)
	case "textDocument/definition":
		return s.handleDefinition(msg.ID, msg.Params)
	case "textDocument/hover":
		return s.handleHover(msg.ID, msg.Params)
	case "textDocument/completion":
		return s.handleCompletion(msg.ID, msg.Params)
	case "textDocument/documentSymbol":
		return s.handleDocumentSymbol(msg.ID, msg.Params)
	case "textDocument/rename":
		return s.handleRename(msg.ID, msg.Params)
	case "textDocument/codeAction":
		return s.handleCodeAction(msg.ID, msg.Params)
	case "textDocument/formatting":
		return s.handleFormatting(msg.ID, msg.Params)
	case "shutdown":
		return respondResult(msg.ID, nil)
	case "exit":
		return nil, nil
	default:
		// Unknown method: respond with method not found if request has id
		if msg.ID != nil {
			errResp := jsonrpcMessage{
				JSONRPC: "2.0",
				ID:      msg.ID,
				Error:   &jsonrpcError{Code: -32601, Message: "method not found: " + msg.Method},
			}
			return json.Marshal(errResp)
		}
		return nil, nil
	}
}

func (s *Server) handleInitialize(id interface{}, params json.RawMessage) ([]byte, error) {
	var initParams InitializeParams
	if err := decodeParams(params, &initParams); err != nil {
		return nil, err
	}
	result := InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync:       intPtr(1), // full sync
			DefinitionProvider:     boolPtr(true),
			HoverProvider:          boolPtr(true),
			CompletionProvider:     &CompletionOptions{},
			DocumentSymbolProvider: boolPtr(true),
			RenameProvider:         boolPtr(true),
			CodeActionProvider:         boolPtr(true),
			DocumentFormattingProvider: boolPtr(true),
		},
		ServerInfo: &struct {
			Name    string `json:"name"`
			Version string `json:"version,omitempty"`
		}{Name: "nuru-lsp", Version: "0.1.0"},
	}
	return respondResult(id, result)
}

func (s *Server) handleDidOpen(params json.RawMessage) ([]byte, error) {
	var p DidOpenTextDocumentParams
	if err := decodeParams(params, &p); err != nil {
		return nil, err
	}
	s.mu.Lock()
	content := p.TextDocument.Text
	program, parseErrs := parseContentAndErrors(content)
	doc := &document{content: content, program: program}
	if program != nil && len(parseErrs) == 0 {
		doc.symbols = symbols.Build(program)
	}
	s.docs[p.TextDocument.URI] = doc
	s.mu.Unlock()
	s.publishDiagnostics(p.TextDocument.URI, content, program, parseErrs)
	return nil, nil
}

func (s *Server) handleDidChange(params json.RawMessage) ([]byte, error) {
	var p DidChangeTextDocumentParams
	if err := decodeParams(params, &p); err != nil {
		return nil, err
	}
	if len(p.ContentChanges) == 0 {
		return nil, nil
	}
	content := p.ContentChanges[len(p.ContentChanges)-1].Text
	s.mu.Lock()
	program, parseErrs := parseContentAndErrors(content)
	doc := &document{content: content, program: program}
	if program != nil && len(parseErrs) == 0 {
		doc.symbols = symbols.Build(program)
	}
	s.docs[p.TextDocument.URI] = doc
	s.mu.Unlock()
	s.publishDiagnostics(p.TextDocument.URI, content, program, parseErrs)
	return nil, nil
}

func (s *Server) handleDidClose(params json.RawMessage) ([]byte, error) {
	var p DidCloseTextDocumentParams
	if err := decodeParams(params, &p); err != nil {
		return nil, err
	}
	s.mu.Lock()
	delete(s.docs, p.TextDocument.URI)
	s.mu.Unlock()
	s.publishDiagnostics(p.TextDocument.URI, "", nil, nil)
	return nil, nil
}

func (s *Server) handleDefinition(id interface{}, params json.RawMessage) ([]byte, error) {
	var p TextDocumentPositionParams
	if err := decodeParams(params, &p); err != nil {
		return nil, err
	}
	s.mu.Lock()
	doc := s.docs[p.TextDocument.URI]
	s.mu.Unlock()
	if doc == nil || doc.symbols == nil {
		return respondResult(id, []Location{})
	}
	name, ok := identifierAtPosition(doc.program, p.Position.Line, p.Position.Character)
	if !ok {
		return respondResult(id, []Location{})
	}
	def, found := doc.symbols.Lookup(name)
	if !found {
		return respondResult(id, []Location{})
	}
	loc := Location{URI: p.TextDocument.URI, Range: defToRange(def)}
	return respondResult(id, []Location{loc})
}

func (s *Server) handleHover(id interface{}, params json.RawMessage) ([]byte, error) {
	var p TextDocumentPositionParams
	if err := decodeParams(params, &p); err != nil {
		return nil, err
	}
	s.mu.Lock()
	doc := s.docs[p.TextDocument.URI]
	s.mu.Unlock()
	if doc == nil || doc.symbols == nil {
		return respondResult(id, nil)
	}
	name, ok := identifierAtPosition(doc.program, p.Position.Line, p.Position.Character)
	if !ok {
		return respondResult(id, nil)
	}
	def, found := doc.symbols.Lookup(name)
	if !found {
		return respondResult(id, nil)
	}
	contents := fmt.Sprintf("**%s** `%s`", def.Kind, name)
	return respondResult(id, Hover{Contents: contents})
}

func (s *Server) handleCompletion(id interface{}, params json.RawMessage) ([]byte, error) {
	var p CompletionParams
	if err := decodeParams(params, &p); err != nil {
		return nil, err
	}
	s.mu.Lock()
	doc := s.docs[p.TextDocument.URI]
	s.mu.Unlock()
	if doc == nil {
		return respondResult(id, []CompletionItem{})
	}
	items := completionItems(doc)
	return respondResult(id, items)
}

func (s *Server) handleDocumentSymbol(id interface{}, params json.RawMessage) ([]byte, error) {
	var p DocumentSymbolParams
	if err := decodeParams(params, &p); err != nil {
		return nil, err
	}
	s.mu.Lock()
	doc := s.docs[p.TextDocument.URI]
	s.mu.Unlock()
	if doc == nil || doc.symbols == nil {
		return respondResult(id, []DocumentSymbol{})
	}
	var out []DocumentSymbol
	for _, sym := range doc.symbols.ListFileScopeDefs() {
		kind := SymbolKindVariable
		if sym.Def.Kind == symbols.DefFunction {
			kind = SymbolKindFunction
		}
		out = append(out, DocumentSymbol{
			Name:   sym.Name,
			Kind:   kind,
			Range:  defToRange(sym.Def),
			Detail: string(sym.Def.Kind),
		})
	}
	return respondResult(id, out)
}

func (s *Server) handleRename(id interface{}, params json.RawMessage) ([]byte, error) {
	var p RenameParams
	if err := decodeParams(params, &p); err != nil {
		return nil, err
	}
	s.mu.Lock()
	doc := s.docs[p.TextDocument.URI]
	s.mu.Unlock()
	if doc == nil || doc.symbols == nil {
		return respondResult(id, WorkspaceEdit{Changes: map[string][]TextEdit{}})
	}
	name, ok := identifierAtPosition(doc.program, p.Position.Line, p.Position.Character)
	if !ok {
		return respondResult(id, WorkspaceEdit{Changes: map[string][]TextEdit{}})
	}
	if _, found := doc.symbols.Lookup(name); !found {
		return respondResult(id, WorkspaceEdit{Changes: map[string][]TextEdit{}})
	}
	ranges := referencesInProgram(doc.program, name)
	edits := make([]TextEdit, len(ranges))
	for i, r := range ranges {
		edits[i] = TextEdit{Range: r, NewText: p.NewName}
	}
	return respondResult(id, WorkspaceEdit{
		Changes: map[string][]TextEdit{p.TextDocument.URI: edits},
	})
}

func (s *Server) handleCodeAction(id interface{}, params json.RawMessage) ([]byte, error) {
	var p CodeActionParams
	if err := decodeParams(params, &p); err != nil {
		return nil, err
	}
	s.mu.Lock()
	doc := s.docs[p.TextDocument.URI]
	s.mu.Unlock()
	var actions []CodeAction
	for _, d := range p.Context.Diagnostics {
		if d.Code == nil {
			continue
		}
		switch *d.Code {
		case DiagnosticCodeMissingSemicolon:
			// Insert ";" at end of the line (diagnostic range is the whole line).
			edit := TextEdit{
				Range:   Range{Start: d.Range.End, End: d.Range.End},
				NewText: ";",
			}
			actions = append(actions, CodeAction{
				Title:       "Insert semicolon",
				Kind:        CodeActionKindQuickFix,
				Diagnostics: []Diagnostic{d},
				Edit: &WorkspaceEdit{
					Changes: map[string][]TextEdit{p.TextDocument.URI: {edit}},
				},
			})
		}
	}
	if doc != nil && len(actions) == 0 {
		// No fixable diagnostics; return empty array (not null).
		return respondResult(id, []CodeAction{})
	}
	return respondResult(id, actions)
}

func (s *Server) handleFormatting(id interface{}, params json.RawMessage) ([]byte, error) {
	var p DocumentFormattingParams
	if err := decodeParams(params, &p); err != nil {
		return nil, err
	}
	s.mu.Lock()
	doc := s.docs[p.TextDocument.URI]
	s.mu.Unlock()
	if doc == nil {
		return respondResult(id, []TextEdit{})
	}
	program, parseErrs := parseContentAndErrors(doc.content)
	if program == nil || len(parseErrs) > 0 {
		return respondResult(id, []TextEdit{})
	}
	opts := format.DefaultOptions()
	opts.TabSize = p.Options.TabSize
	opts.InsertSpaces = p.Options.InsertSpaces
	formatted := format.Program(program, opts)
	lines := strings.Split(doc.content, "\n")
	endLine := len(lines) - 1
	if endLine < 0 {
		endLine = 0
	}
	endChar := 0
	if endLine < len(lines) {
		endChar = len(lines[endLine])
	}
	fullRange := Range{
		Start: Position{Line: 0, Character: 0},
		End:   Position{Line: endLine, Character: endChar},
	}
	edit := TextEdit{Range: fullRange, NewText: formatted}
	return respondResult(id, []TextEdit{edit})
}

// parseContentAndErrors parses content and returns the program and any parser errors.
func parseContentAndErrors(content string) (*ast.Program, []string) {
	l := lexer.New(content)
	p := parser.New(l)
	program := p.ParseProgram()
	return program, p.Errors()
}

// publishDiagnostics sends diagnostics to the client.
func (s *Server) publishDiagnostics(uri, content string, program *ast.Program, parseErrs []string) {
	s.mu.Lock()
	w := s.w
	s.mu.Unlock()
	if w == nil {
		return
	}
	diags := parserErrorsToDiagnostics(content, parseErrs)
	if len(parseErrs) == 0 && program != nil {
		opts := analysis.DefaultAnalysisOptions()
		opts.AnalyzeModules = false
		a := analysis.NewAnalyzer(opts)
		if warnings, err := a.AnalyzeFile(uri, content); err == nil {
			diags = append(diags, analysisWarningsToDiagnostics(warnings)...)
		}
	}
	params := PublishDiagnosticsParams{URI: uri, Diagnostics: diags}
	_ = notify(w, "textDocument/publishDiagnostics", params)
}

func intPtr(i int) *int    { return &i }
func boolPtr(b bool) *bool { return &b }
