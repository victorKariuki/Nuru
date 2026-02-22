package lsp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// jsonrpcMessage is the raw JSON-RPC 2 message.
type jsonrpcMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"` // number or string
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *jsonrpcError   `json:"error,omitempty"`
}

type jsonrpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const contentLengthPrefix = "Content-Length: "

// readMessage reads one LSP message (Content-Length: N + \r\n\r\n + body).
func readMessage(r *bufio.Reader) ([]byte, error) {
	var contentLength int
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
			n, err := strconv.Atoi(strings.TrimSpace(line[len(contentLengthPrefix):]))
			if err != nil {
				return nil, fmt.Errorf("invalid Content-Length: %w", err)
			}
			contentLength = n
		}
	}
	if contentLength <= 0 {
		return nil, fmt.Errorf("missing or invalid Content-Length")
	}
	body := make([]byte, contentLength)
	if _, err := io.ReadFull(r, body); err != nil {
		return nil, err
	}
	return body, nil
}

// writeMessage writes one LSP message.
func writeMessage(w io.Writer, body []byte) error {
	header := fmt.Sprintf("%s%d\r\n\r\n", contentLengthPrefix, len(body))
	if _, err := io.WriteString(w, header); err != nil {
		return err
	}
	_, err := w.Write(body)
	return err
}

// runReadLoop reads JSON-RPC messages from r and calls handle for each.
// handle(msg) returns the response body to send (nil for notifications).
func runReadLoop(r *bufio.Reader, w io.Writer, handle func([]byte) ([]byte, error)) error {
	for {
		body, err := readMessage(r)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		respBody, err := handle(body)
		if err != nil {
			// Send error response if we have an id
			var msg jsonrpcMessage
			if jerr := json.Unmarshal(body, &msg); jerr == nil && msg.ID != nil {
				errResp := jsonrpcMessage{
					JSONRPC: "2.0",
					ID:      msg.ID,
					Error:   &jsonrpcError{Code: -32603, Message: err.Error()},
				}
				if b, jerr := json.Marshal(errResp); jerr == nil {
					_ = writeMessage(w, b)
				}
			}
			continue
		}
		if respBody != nil {
			if err := writeMessage(w, respBody); err != nil {
				return err
			}
		}
	}
}

// respondResult builds a JSON-RPC success response.
func respondResult(id interface{}, result interface{}) ([]byte, error) {
	return json.Marshal(jsonrpcMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	})
}

// notify sends a notification (no id) to the client.
func notify(w io.Writer, method string, params interface{}) error {
	msg := struct {
		JSONRPC string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params,omitempty"`
	}{JSONRPC: "2.0", Method: method, Params: params}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return writeMessage(w, body)
}

// decodeParams decodes params into the given value.
func decodeParams(raw json.RawMessage, v interface{}) error {
	if len(raw) == 0 {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	return dec.Decode(v)
}
