package lsp

import (
	"github.com/NuruProgramming/Nuru/lsp/symbols"
)

// Nuru keywords for completion (from token package).
var completionKeywords = []string{
	"unda", "fanya", "kweli", "sikweli", "kama", "au", "sivyo",
	"wakati", "rudisha", "vunja", "endelea", "tupu", "ktk", "kwa",
	"badili", "ikiwa", "kawaida", "tumia", "pakeji", "jaribu", "shika", "bila",
}

// completionItems builds completion items for the given document (symbol table + keywords).
func completionItems(doc *document) []CompletionItem {
	var items []CompletionItem
	kindKeyword := CompletionItemKindKeyword
	kindVar := CompletionItemKindVariable
	kindFunc := CompletionItemKindFunction
	if doc.symbols != nil {
		for _, s := range doc.symbols.ListFileScope() {
			k := &kindVar
			if s.Kind == symbols.DefFunction {
				k = &kindFunc
			}
			items = append(items, CompletionItem{Label: s.Name, Kind: k})
		}
	}
	for _, kw := range completionKeywords {
		items = append(items, CompletionItem{Label: kw, Kind: &kindKeyword})
	}
	return items
}
