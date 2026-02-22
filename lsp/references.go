package lsp

import (
	"github.com/NuruProgramming/Nuru/ast"
)

// referencesInProgram returns all LSP ranges where an identifier with the given name appears.
func referencesInProgram(program *ast.Program, name string) []Range {
	if program == nil || name == "" {
		return nil
	}
	var out []Range
	collectRefs(program, name, &out)
	return out
}

func collectRefs(node ast.Node, name string, out *[]Range) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *ast.Identifier:
		if n.Value == name {
			*out = append(*out, tokenToRange(n.Token))
		}
		return
	case *ast.Program:
		for _, stmt := range n.Statements {
			collectRefs(stmt, name, out)
		}
	case *ast.LetStatement:
		collectRefs(n.Name, name, out)
		if n.Value != nil {
			collectRefs(n.Value, name, out)
		}
	case *ast.ReturnStatement:
		if n.ReturnValue != nil {
			collectRefs(n.ReturnValue, name, out)
		}
	case *ast.ExpressionStatement:
		if n.Expression != nil {
			collectRefs(n.Expression, name, out)
		}
	case *ast.BlockStatement:
		for _, stmt := range n.Statements {
			collectRefs(stmt, name, out)
		}
	case *ast.FunctionLiteral:
		for _, p := range n.Parameters {
			collectRefs(p, name, out)
		}
		if n.Body != nil {
			collectRefs(n.Body, name, out)
		}
	case *ast.CallExpression:
		collectRefs(n.Function, name, out)
		for _, arg := range n.Arguments {
			collectRefs(arg, name, out)
		}
	case *ast.InfixExpression:
		collectRefs(n.Left, name, out)
		collectRefs(n.Right, name, out)
	case *ast.PrefixExpression:
		collectRefs(n.Right, name, out)
	case *ast.IfExpression:
		collectRefs(n.Condition, name, out)
		if n.Consequence != nil {
			collectRefs(n.Consequence, name, out)
		}
		if n.Alternative != nil {
			collectRefs(n.Alternative, name, out)
		}
	case *ast.IndexExpression:
		collectRefs(n.Left, name, out)
		collectRefs(n.Index, name, out)
	case *ast.AssignmentExpression:
		collectRefs(n.Left, name, out)
		collectRefs(n.Value, name, out)
	case *ast.MethodExpression:
		collectRefs(n.Object, name, out)
		collectRefs(n.Method, name, out)
		for _, arg := range n.Arguments {
			collectRefs(arg, name, out)
		}
	case *ast.PropertyExpression:
		collectRefs(n.Object, name, out)
		collectRefs(n.Property, name, out)
	case *ast.WhileExpression:
		collectRefs(n.Condition, name, out)
		if n.Consequence != nil {
			collectRefs(n.Consequence, name, out)
		}
	case *ast.For:
		if n.Block != nil {
			collectRefs(n.Block, name, out)
		}
		if n.StarterValue != nil {
			collectRefs(n.StarterValue, name, out)
		}
		if n.Closer != nil {
			collectRefs(n.Closer, name, out)
		}
		if n.Condition != nil {
			collectRefs(n.Condition, name, out)
		}
	case *ast.ForIn:
		if n.Block != nil {
			collectRefs(n.Block, name, out)
		}
		if n.Iterable != nil {
			collectRefs(n.Iterable, name, out)
		}
	default:
		collectRefsChildren(node, name, out)
	}
}

func collectRefsChildren(node ast.Node, name string, out *[]Range) {
	switch n := node.(type) {
	case *ast.SwitchExpression:
		if n.Value != nil {
			collectRefs(n.Value, name, out)
		}
		for _, c := range n.Choices {
			if c != nil {
				for _, e := range c.Expr {
					collectRefs(e, name, out)
				}
				if c.Block != nil {
					collectRefs(c.Block, name, out)
				}
			}
		}
	case *ast.TryCatchExpression:
		if n.TryBlock != nil {
			collectRefs(n.TryBlock, name, out)
		}
		if n.CatchBlock != nil {
			collectRefs(n.CatchBlock, name, out)
		}
	case *ast.AssignEqual:
		if n.Left != nil {
			collectRefs(n.Left, name, out)
		}
		if n.Value != nil {
			collectRefs(n.Value, name, out)
		}
	}
}
