package symbols

import (
	"github.com/NuruProgramming/Nuru/ast"
	"github.com/NuruProgramming/Nuru/token"
)

// DefKind is the kind of a symbol definition.
type DefKind string

const (
	DefVariable  DefKind = "variable"
	DefFunction  DefKind = "function"
	DefParameter DefKind = "parameter"
)

// Def holds the location and kind of a symbol definition.
type Def struct {
	Line   int     // 1-based
	Column int     // 1-based
	EndCol int     // 1-based end column (exclusive)
	Kind   DefKind
}

// Table is a symbol table with a scope stack.
type Table struct {
	scopes []map[string]Def
}

// NewTable returns a new symbol table.
func NewTable() *Table {
	return &Table{
		scopes: []map[string]Def{make(map[string]Def)},
	}
}

// PushScope pushes a new scope.
func (t *Table) PushScope() {
	t.scopes = append(t.scopes, make(map[string]Def))
}

// PopScope pops the innermost scope.
func (t *Table) PopScope() {
	if len(t.scopes) > 1 {
		t.scopes = t.scopes[:len(t.scopes)-1]
	}
}

// Bind binds a name to a definition in the current scope.
func (t *Table) Bind(name string, tok token.Token, kind DefKind) {
	if len(t.scopes) == 0 {
		return
	}
	scope := t.scopes[len(t.scopes)-1]
	endCol := tok.Column
	if tok.Column > 0 {
		endCol = tok.Column + len(tok.Literal)
	}
	scope[name] = Def{Line: tok.Line, Column: tok.Column, EndCol: endCol, Kind: kind}
}

// Lookup looks up a name from innermost to outermost scope. Returns (def, true) if found.
func (t *Table) Lookup(name string) (Def, bool) {
	for i := len(t.scopes) - 1; i >= 0; i-- {
		if d, ok := t.scopes[i][name]; ok {
			return d, true
		}
	}
	return Def{}, false
}

// ListFileScope returns all names bound in the file (top-level) scope for completion.
func (t *Table) ListFileScope() []struct{ Name string; Kind DefKind } {
	if len(t.scopes) == 0 {
		return nil
	}
	out := make([]struct{ Name string; Kind DefKind }, 0, len(t.scopes[0]))
	for name, d := range t.scopes[0] {
		out = append(out, struct{ Name string; Kind DefKind }{Name: name, Kind: d.Kind})
	}
	return out
}

// Build builds the symbol table from a program.
func Build(program *ast.Program) *Table {
	if program == nil {
		return NewTable()
	}
	t := NewTable()
	for _, stmt := range program.Statements {
		walkStatement(t, stmt)
	}
	return t
}

func walkStatement(t *Table, stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.LetStatement:
		if s.Name != nil {
			kind := DefVariable
			if _, isFunc := s.Value.(*ast.FunctionLiteral); isFunc {
				kind = DefFunction
			}
			t.Bind(s.Name.Value, s.Name.Token, kind)
		}
		if s.Value != nil {
			walkExpression(t, s.Value)
		}
	case *ast.ReturnStatement:
		if s.ReturnValue != nil {
			walkExpression(t, s.ReturnValue)
		}
	case *ast.ExpressionStatement:
		if s.Expression != nil {
			walkExpression(t, s.Expression)
		}
	case *ast.BlockStatement:
		t.PushScope()
		for _, st := range s.Statements {
			walkStatement(t, st)
		}
		t.PopScope()
	default:
		// Break, Continue, etc. have no sub-statements to walk
	}
}

func walkExpression(t *Table, expr ast.Expression) {
	switch e := expr.(type) {
	case *ast.FunctionLiteral:
		t.PushScope()
		for _, p := range e.Parameters {
			if p != nil {
				t.Bind(p.Value, p.Token, DefParameter)
			}
		}
		if e.Body != nil {
			for _, st := range e.Body.Statements {
				walkStatement(t, st)
			}
		}
		t.PopScope()
	case *ast.CallExpression:
		if e.Function != nil {
			walkExpression(t, e.Function)
		}
		for _, a := range e.Arguments {
			if a != nil {
				walkExpression(t, a)
			}
		}
	case *ast.InfixExpression:
		if e.Left != nil {
			walkExpression(t, e.Left)
		}
		if e.Right != nil {
			walkExpression(t, e.Right)
		}
	case *ast.PrefixExpression:
		if e.Right != nil {
			walkExpression(t, e.Right)
		}
	case *ast.IfExpression:
		if e.Condition != nil {
			walkExpression(t, e.Condition)
		}
		if e.Consequence != nil {
			t.PushScope()
			for _, st := range e.Consequence.Statements {
				walkStatement(t, st)
			}
			t.PopScope()
		}
		if e.Alternative != nil {
			t.PushScope()
			for _, st := range e.Alternative.Statements {
				walkStatement(t, st)
			}
			t.PopScope()
		}
	case *ast.WhileExpression:
		if e.Condition != nil {
			walkExpression(t, e.Condition)
		}
		if e.Consequence != nil {
			t.PushScope()
			for _, st := range e.Consequence.Statements {
				walkStatement(t, st)
			}
			t.PopScope()
		}
	case *ast.For:
		if e.Block != nil {
			t.PushScope()
			for _, st := range e.Block.Statements {
				walkStatement(t, st)
			}
			t.PopScope()
		}
		if e.StarterValue != nil {
			walkExpression(t, e.StarterValue)
		}
		if e.Closer != nil {
			walkExpression(t, e.Closer)
		}
		if e.Condition != nil {
			walkExpression(t, e.Condition)
		}
	case *ast.ForIn:
		if e.Block != nil {
			t.PushScope()
			for _, st := range e.Block.Statements {
				walkStatement(t, st)
			}
			t.PopScope()
		}
		if e.Iterable != nil {
			walkExpression(t, e.Iterable)
		}
	case *ast.IndexExpression:
		if e.Left != nil {
			walkExpression(t, e.Left)
		}
		if e.Index != nil {
			walkExpression(t, e.Index)
		}
	case *ast.AssignmentExpression:
		if e.Left != nil {
			walkExpression(t, e.Left)
		}
		if e.Value != nil {
			walkExpression(t, e.Value)
		}
	default:
		walkExpressionRec(t, expr)
	}
}

func walkExpressionRec(t *Table, expr ast.Expression) {
	// Recursively walk composite expressions we didn't handle above
	switch e := expr.(type) {
	case *ast.CallExpression:
		if e.Function != nil {
			walkExpression(t, e.Function)
		}
		for _, a := range e.Arguments {
			if a != nil {
				walkExpression(t, a)
			}
		}
	case *ast.InfixExpression:
		if e.Left != nil {
			walkExpression(t, e.Left)
		}
		if e.Right != nil {
			walkExpression(t, e.Right)
		}
	case *ast.IndexExpression:
		if e.Left != nil {
			walkExpression(t, e.Left)
		}
		if e.Index != nil {
			walkExpression(t, e.Index)
		}
	}
}
