package symbols

import (
	"testing"

	"github.com/NuruProgramming/Nuru/ast"
	"github.com/NuruProgramming/Nuru/lexer"
	"github.com/NuruProgramming/Nuru/parser"
	"github.com/NuruProgramming/Nuru/token"
)

func parse(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}
	return program
}

func TestBuild_LetStatement(t *testing.T) {
	program := parse(t, "fanya x = 5;")
	tab := Build(program)
	def, ok := tab.Lookup("x")
	if !ok {
		t.Fatal("expected to find x")
	}
	if def.Kind != DefVariable {
		t.Errorf("expected variable, got %s", def.Kind)
	}
	if def.Line != 1 || def.Column < 1 {
		t.Errorf("expected line 1, column >= 1, got line %d col %d", def.Line, def.Column)
	}
}

func TestBuild_FunctionLiteral(t *testing.T) {
	program := parse(t, "fanya f = unda(a, b) { rudisha a + b; };")
	tab := Build(program)
	def, ok := tab.Lookup("f")
	if !ok {
		t.Fatal("expected to find f")
	}
	if def.Kind != DefFunction {
		t.Errorf("expected function, got %s", def.Kind)
	}
	// Parameters are in inner scope
	_, okA := tab.Lookup("a")
	if okA {
		t.Error("a is parameter, should be in function scope only; lookup from top finds it in closed scope - actually we search from innermost so we wouldn't find a at file level")
	}
}

func TestBuild_NestedScopes(t *testing.T) {
	program := parse(t, `
fanya x = 1;
fanya y = unda() {
  fanya x = 2;
  x;
};
`)
	tab := Build(program)
	def, ok := tab.Lookup("x")
	if !ok {
		t.Fatal("expected to find x")
	}
	// File-level x is the one we find when looking up from the top (we don't have a "current scope" in lookup - we search all scopes from innermost to outermost, but the table is built for the whole program so the "innermost" is the last scope pushed). Actually Lookup searches from top of stack (innermost) to bottom (outermost). So the order of scopes in the table is: after building, we have [fileScope]. Then for "fanya y = unda()..." we push function scope, add x=2, pop. So we have [fileScope] with x=1, y=?. So when we Lookup("x") we only have one scope with x. So we get x=1. Good.
	if def.Line < 1 || def.Line > 5 {
		t.Errorf("expected file-level x at a reasonable line, got line %d", def.Line)
	}
}

func TestBind_Lookup(t *testing.T) {
	tab := NewTable()
	tok := token.Token{Line: 1, Column: 1, Literal: "foo"}
	tab.Bind("foo", tok, DefVariable)
	def, ok := tab.Lookup("foo")
	if !ok {
		t.Fatal("expected to find foo")
	}
	if def.Line != 1 || def.Column != 1 || def.Kind != DefVariable {
		t.Errorf("got %+v", def)
	}
}
