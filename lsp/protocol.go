package lsp

// Minimal LSP protocol types for Nuru LSP server.
// LSP uses 0-based line and character (UTF-16 code units); we use 0-based line and character.

// Position in a document (0-based).
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Range in a document.
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Location (uri + range).
type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

// Diagnostic severity.
const (
	DiagnosticSeverityError   = 1
	DiagnosticSeverityWarning = 2
)

// Diagnostic represents a diagnostic (error or warning).
type Diagnostic struct {
	Range   Range  `json:"range"`
	Message string `json:"message"`
	// Severity optional; 1=Error, 2=Warning
	Severity *int `json:"severity,omitempty"`
	// Code identifies the diagnostic for code actions (e.g. "missing_semicolon").
	Code *string `json:"code,omitempty"`
}

// InitializeParams for initialize request.
type InitializeParams struct {
	ProcessID *int `json:"processId,omitempty"`
	RootURI   *string `json:"rootUri,omitempty"`
	Capabilities struct{} `json:"capabilities,omitempty"`
}

// InitializeResult returned by initialize.
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *struct {
		Name    string `json:"name"`
		Version string `json:"version,omitempty"`
	} `json:"serverInfo,omitempty"`
}

// ServerCapabilities advertised by the server.
type ServerCapabilities struct {
	TextDocumentSync        *int              `json:"textDocumentSync,omitempty"` // 1 = full
	DefinitionProvider      *bool             `json:"definitionProvider,omitempty"`
	HoverProvider           *bool             `json:"hoverProvider,omitempty"`
	CompletionProvider      *CompletionOptions `json:"completionProvider,omitempty"`
	DocumentSymbolProvider  *bool             `json:"documentSymbolProvider,omitempty"`
	RenameProvider          *bool             `json:"renameProvider,omitempty"`
	CodeActionProvider         *bool             `json:"codeActionProvider,omitempty"`
	DocumentFormattingProvider *bool             `json:"documentFormattingProvider,omitempty"`
}

// CompletionOptions for completion provider.
type CompletionOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

// TextDocumentItem for didOpen.
type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

// DidOpenTextDocumentParams for textDocument/didOpen.
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

// TextDocumentContentChangeEvent for didChange (full sync).
type TextDocumentContentChangeEvent struct {
	Text string `json:"text"`
}

// DidChangeTextDocumentParams for textDocument/didChange.
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// VersionedTextDocumentIdentifier identifies a document with version.
type VersionedTextDocumentIdentifier struct {
	URI     string `json:"uri"`
	Version int    `json:"version"`
}

// DidCloseTextDocumentParams for textDocument/didClose.
type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// TextDocumentIdentifier identifies a document.
type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

// TextDocumentPositionParams for definition/hover (position in document).
type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

// PublishDiagnosticsParams sent to client (notification).
type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// Hover result.
type Hover struct {
	Contents string `json:"contents"` // string or MarkupContent; we use string
	Range    *Range `json:"range,omitempty"`
}

// MarkupKind for hover contents.
const MarkupKindPlainText = "plaintext"
const MarkupKindMarkdown = "markdown"

// CompletionItem for textDocument/completion.
type CompletionItem struct {
	Label  string `json:"label"`
	Kind   *int   `json:"kind,omitempty"` // CompletionItemKind (e.g. 6 = Variable, 3 = Function)
	Detail string `json:"detail,omitempty"`
}

// CompletionItemKind (LSP enum).
const (
	CompletionItemKindVariable = 6
	CompletionItemKindFunction = 3
	CompletionItemKindKeyword  = 14
)

// CompletionParams for textDocument/completion.
type CompletionParams struct {
	TextDocumentPositionParams
}

// DocumentSymbolParams for textDocument/documentSymbol.
type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// DocumentSymbol for textDocument/documentSymbol (optional).
type DocumentSymbol struct {
	Name     string           `json:"name"`
	Kind     int              `json:"kind"` // SymbolKind
	Range    Range            `json:"range"`
	Detail   string           `json:"detail,omitempty"`
	Children []DocumentSymbol `json:"children,omitempty"`
}

// SymbolKind (LSP enum) for document symbols.
const (
	SymbolKindFunction = 12
	SymbolKindVariable = 13
)

// TextEdit is a single edit (range + newText).
type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

// WorkspaceEdit holds edits per document URI.
type WorkspaceEdit struct {
	Changes map[string][]TextEdit `json:"changes,omitempty"`
}

// RenameParams for textDocument/rename.
type RenameParams struct {
	TextDocumentPositionParams
	NewName string `json:"newName"`
}

// CodeActionKind (LSP) for code action kind.
const CodeActionKindQuickFix = "quickfix"

// CodeActionParams for textDocument/codeAction.
type CodeActionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                 `json:"range"`
	Context      CodeActionContext     `json:"context"`
}

// CodeActionContext contains diagnostics for the code action request.
type CodeActionContext struct {
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// CodeAction for textDocument/codeAction response.
type CodeAction struct {
	Title       string         `json:"title"`
	Kind        string         `json:"kind,omitempty"`
	Diagnostics []Diagnostic   `json:"diagnostics,omitempty"`
	Edit        *WorkspaceEdit `json:"edit,omitempty"`
}

// DocumentFormattingParams for textDocument/formatting.
type DocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

// FormattingOptions for formatting (tabSize, insertSpaces).
type FormattingOptions struct {
	TabSize      int  `json:"tabSize"`
	InsertSpaces bool `json:"insertSpaces"`
}
