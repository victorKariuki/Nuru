package lsp

import (
	"github.com/NuruProgramming/Nuru/ast"
	"github.com/NuruProgramming/Nuru/token"
	"github.com/NuruProgramming/Nuru/lsp/symbols"
)

// tokenContains returns true if the 0-based (line, character) is inside the token's range.
// Token uses 1-based Line and Column.
func tokenContains(tok token.Token, line, character int) bool {
	tokLine := tok.Line - 1
	tokCol := tok.Column - 1
	if tokCol < 0 {
		tokCol = 0
	}
	endCol := tokCol + len(tok.Literal)
	return tokLine == line && character >= tokCol && character < endCol
}

// identifierAtPosition finds the identifier node at the given 0-based position in the program.
// Returns the identifier's name and true if found.
func identifierAtPosition(program *ast.Program, line, character int) (name string, ok bool) {
	if program == nil {
		return "", false
	}
	var found string
	var foundOk bool
	walkAST(program, line, character, &found, &foundOk)
	return found, foundOk
}

func walkAST(node ast.Node, line, character int, out *string, ok *bool) {
	// node can be *ast.Program, ast.Statement, or ast.Expression (concrete types)
	if node == nil || *ok {
		return
	}
	switch n := node.(type) {
	case *ast.Identifier:
		if tokenContains(n.Token, line, character) {
			*out = n.Value
			*ok = true
		}
		return
	case *ast.Program:
		for _, stmt := range n.Statements {
			walkAST(stmt, line, character, out, ok)
			if *ok {
				return
			}
		}
	case *ast.LetStatement:
		walkAST(n.Name, line, character, out, ok)
		if !*ok && n.Value != nil {
			walkAST(n.Value, line, character, out, ok)
		}
	case *ast.ReturnStatement:
		if n.ReturnValue != nil {
			walkAST(n.ReturnValue, line, character, out, ok)
		}
	case *ast.ExpressionStatement:
		if n.Expression != nil {
			walkAST(n.Expression, line, character, out, ok)
		}
	case *ast.BlockStatement:
		for _, stmt := range n.Statements {
			walkAST(stmt, line, character, out, ok)
			if *ok {
				return
			}
		}
	case *ast.FunctionLiteral:
		for _, p := range n.Parameters {
			walkAST(p, line, character, out, ok)
			if *ok {
				return
			}
		}
		if n.Body != nil {
			walkAST(n.Body, line, character, out, ok)
		}
	case *ast.CallExpression:
		walkAST(n.Function, line, character, out, ok)
		if !*ok {
			for _, arg := range n.Arguments {
				walkAST(arg, line, character, out, ok)
				if *ok {
					return
				}
			}
		}
	case *ast.InfixExpression:
		walkAST(n.Left, line, character, out, ok)
		if !*ok {
			walkAST(n.Right, line, character, out, ok)
		}
	case *ast.PrefixExpression:
		walkAST(n.Right, line, character, out, ok)
	case *ast.IfExpression:
		walkAST(n.Condition, line, character, out, ok)
		if !*ok && n.Consequence != nil {
			walkAST(n.Consequence, line, character, out, ok)
		}
		if !*ok && n.Alternative != nil {
			walkAST(n.Alternative, line, character, out, ok)
		}
	case *ast.IndexExpression:
		walkAST(n.Left, line, character, out, ok)
		if !*ok {
			walkAST(n.Index, line, character, out, ok)
		}
	case *ast.AssignmentExpression:
		walkAST(n.Left, line, character, out, ok)
		if !*ok {
			walkAST(n.Value, line, character, out, ok)
		}
	case *ast.MethodExpression:
		walkAST(n.Object, line, character, out, ok)
		if !*ok {
			walkAST(n.Method, line, character, out, ok)
		}
		if !*ok {
			for _, arg := range n.Arguments {
				walkAST(arg, line, character, out, ok)
				if *ok {
					return
				}
			}
		}
	case *ast.PropertyExpression:
		walkAST(n.Object, line, character, out, ok)
		if !*ok {
			walkAST(n.Property, line, character, out, ok)
		}
	case *ast.WhileExpression:
		walkAST(n.Condition, line, character, out, ok)
		if !*ok && n.Consequence != nil {
			walkAST(n.Consequence, line, character, out, ok)
		}
	default:
		// Other node types: try to walk children if they implement Node
		walkChildren(node, line, character, out, ok)
	}
}

func walkChildren(node ast.Node, line, character int, out *string, ok *bool) {
	// Minimal child walking for other expression types
	switch n := node.(type) {
	case *ast.For:
		if n.Block != nil {
			walkAST(n.Block, line, character, out, ok)
		}
	case *ast.ForIn:
		walkAST(n.Iterable, line, character, out, ok)
		if !*ok && n.Block != nil {
			walkAST(n.Block, line, character, out, ok)
		}
	}
}

// tokenToRange converts a token (1-based Line, Column) to LSP Range (0-based).
func tokenToRange(tok token.Token) Range {
	line := tok.Line - 1
	if line < 0 {
		line = 0
	}
	col := tok.Column - 1
	if col < 0 {
		col = 0
	}
	endCol := col + len(tok.Literal)
	return Range{
		Start: Position{Line: line, Character: col},
		End:   Position{Line: line, Character: endCol},
	}
}

// defToRange converts a symbol Def (1-based) to LSP Range (0-based).
func defToRange(d symbols.Def) Range {
	line := d.Line - 1
	if line < 0 {
		line = 0
	}
	col := d.Column - 1
	if col < 0 {
		col = 0
	}
	endCol := d.EndCol - 1
	if endCol < col {
		endCol = col + 1
	}
	return Range{
		Start: Position{Line: line, Character: col},
		End:   Position{Line: line, Character: endCol},
	}
}
