package lsp

import (
	"testing"
)

func TestParserErrorsToDiagnostics(t *testing.T) {
	content := "fanya x =\nandika x"
	errs := []string{
		"Mstari 1: Tulitegemea kupata ;, badala yake tumepata NENO",
		"Mstari 2: Tumeshindwa kupata ;",
	}
	diags := parserErrorsToDiagnostics(content, errs)
	if len(diags) != 2 {
		t.Fatalf("expected 2 diagnostics, got %d", len(diags))
	}
	if diags[0].Range.Start.Line != 0 {
		t.Errorf("first diagnostic line: expected 0, got %d", diags[0].Range.Start.Line)
	}
	if diags[0].Message != errs[0] {
		t.Errorf("message mismatch: %q", diags[0].Message)
	}
	if diags[1].Range.Start.Line != 1 {
		t.Errorf("second diagnostic line: expected 1, got %d", diags[1].Range.Start.Line)
	}
}

func TestParserErrorsToDiagnostics_NoLine(t *testing.T) {
	content := "x"
	errs := []string{"Some error without Mstari"}
	diags := parserErrorsToDiagnostics(content, errs)
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
	if diags[0].Range.Start.Line != 0 {
		t.Errorf("expected line 0 when no Mstari, got %d", diags[0].Range.Start.Line)
	}
}
