package lsp

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/NuruProgramming/Nuru/analysis"
)

// parserErrorLineRe extracts "Mstari N:" from parser error messages.
var parserErrorLineRe = regexp.MustCompile(`Mstari (\d+):`)

// Diagnostic codes for code actions.
const (
	DiagnosticCodeMissingSemicolon = "missing_semicolon"
)

// parserErrorsToDiagnostics converts parser error strings to LSP diagnostics.
// Parser messages are like "Mstari 3: Tulitegemea...". Line is 1-based; we convert to 0-based.
func parserErrorsToDiagnostics(content string, errs []string) []Diagnostic {
	lines := strings.Split(content, "\n")
	out := make([]Diagnostic, 0, len(errs))
	for _, msg := range errs {
		line := 0
		if m := parserErrorLineRe.FindStringSubmatch(msg); len(m) >= 2 {
			if n, err := strconv.Atoi(m[1]); err == nil {
				line = n - 1 // 1-based to 0-based
				if line < 0 {
					line = 0
				}
			}
		}
		endChar := 0
		if line < len(lines) {
			endChar = len(lines[line])
		}
		sev := DiagnosticSeverityError
		d := Diagnostic{
			Range: Range{
				Start: Position{Line: line, Character: 0},
				End:   Position{Line: line, Character: endChar},
			},
			Message:  msg,
			Severity: &sev,
		}
		// Fixable: parser expected semicolon
		if strings.Contains(msg, "Tulitegemea") && strings.Contains(msg, ";") {
			d.Code = strPtr(DiagnosticCodeMissingSemicolon)
		}
		out = append(out, d)
	}
	return out
}

func strPtr(s string) *string { return &s }

// analysisWarningsToDiagnostics converts analysis.Warning to LSP diagnostics.
// Warning uses 1-based Line and Column; LSP uses 0-based.
func analysisWarningsToDiagnostics(warnings []analysis.Warning) []Diagnostic {
	out := make([]Diagnostic, 0, len(warnings))
	sev := DiagnosticSeverityWarning
	for _, w := range warnings {
		line := w.Line - 1
		if line < 0 {
			line = 0
		}
		col := w.Column - 1
		if col < 0 {
			col = 0
		}
		out = append(out, Diagnostic{
			Range: Range{
				Start: Position{Line: line, Character: col},
				End:   Position{Line: line, Character: col + 1},
			},
			Message:  w.Message,
			Severity: &sev,
		})
	}
	return out
}
