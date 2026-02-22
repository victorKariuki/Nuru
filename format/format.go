package format

import (
	"bytes"
	"strings"

	"github.com/NuruProgramming/Nuru/ast"
)

// Options holds formatting options.
type Options struct {
	Indent      string // e.g. 4 spaces
	TabSize     int    // for LSP; used to build Indent if not set
	InsertSpaces bool  // if true use spaces for indent
}

// DefaultOptions returns default formatting options (4 spaces).
func DefaultOptions() Options {
	return Options{Indent: "    ", TabSize: 4, InsertSpaces: true}
}

// Program formats a program and returns the full document string.
func Program(program *ast.Program, opts Options) string {
	if program == nil || len(program.Statements) == 0 {
		return ""
	}
	indent := opts.Indent
	if indent == "" {
		if opts.InsertSpaces {
			indent = strings.Repeat(" ", opts.TabSize)
		} else {
			indent = "\t"
		}
	}
	var buf bytes.Buffer
	for i, stmt := range program.Statements {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(formatStatement(stmt, indent, ""))
	}
	return buf.String()
}

func formatStatement(stmt ast.Statement, indent, currentIndent string) string {
	switch s := stmt.(type) {
	case *ast.LetStatement:
		return formatLetStatement(s, indent, currentIndent)
	case *ast.ReturnStatement:
		return currentIndent + "rudisha " + exprString(s.ReturnValue) + ";"
	case *ast.ExpressionStatement:
		return currentIndent + exprString(s.Expression) + ";"
	case *ast.BlockStatement:
		return formatBlock(s, indent, currentIndent)
	default:
		return currentIndent + stmt.String()
	}
}

func formatLetStatement(s *ast.LetStatement, indent, currentIndent string) string {
	out := currentIndent + "fanya " + s.Name.String() + " = "
	if s.Value != nil {
		out += exprStringWithBlock(s.Value, indent, currentIndent+indent)
	}
	return out + ";"
}

func formatBlock(block *ast.BlockStatement, indent, currentIndent string) string {
	if block == nil || len(block.Statements) == 0 {
		return "{\n" + currentIndent + "}"
	}
	var buf bytes.Buffer
	buf.WriteString("{\n")
	inner := currentIndent + indent
	for _, stmt := range block.Statements {
		buf.WriteString(formatStatement(stmt, indent, inner))
		buf.WriteByte('\n')
	}
	buf.WriteString(currentIndent + "}")
	return buf.String()
}

func exprString(e ast.Expression) string {
	if e == nil {
		return ""
	}
	return e.String()
}

// exprStringWithBlock formats an expression; for FunctionLiteral and blocks uses proper newlines/indent.
func exprStringWithBlock(e ast.Expression, indent, blockIndent string) string {
	if e == nil {
		return ""
	}
	switch ex := e.(type) {
	case *ast.FunctionLiteral:
		return formatFunctionLiteral(ex, indent, blockIndent)
	case *ast.IfExpression:
		return formatIfExpression(ex, indent, blockIndent)
	case *ast.WhileExpression:
		return formatWhileExpression(ex, indent, blockIndent)
	case *ast.For:
		return formatFor(ex, indent, blockIndent)
	case *ast.ForIn:
		return formatForIn(ex, indent, blockIndent)
	case *ast.SwitchExpression:
		return formatSwitch(ex, indent, blockIndent)
	case *ast.TryCatchExpression:
		return formatTryCatch(ex, indent, blockIndent)
	default:
		return e.String()
	}
}

func formatFunctionLiteral(fl *ast.FunctionLiteral, indent, currentIndent string) string {
	var buf bytes.Buffer
	buf.WriteString("unda(")
	if len(fl.Parameters) > 0 {
		parts := make([]string, len(fl.Parameters))
		for i, p := range fl.Parameters {
			parts[i] = p.String()
		}
		buf.WriteString(strings.Join(parts, ", "))
	}
	buf.WriteString(") ")
	if fl.Body != nil {
		buf.WriteString(formatBlock(fl.Body, indent, currentIndent))
	} else {
		buf.WriteString("{}")
	}
	return buf.String()
}

func formatIfExpression(ie *ast.IfExpression, indent, currentIndent string) string {
	out := "kama " + exprString(ie.Condition) + " "
	if ie.Consequence != nil {
		out += formatBlock(ie.Consequence, indent, currentIndent)
	}
	if ie.Alternative != nil {
		out += " sivyo " + formatBlock(ie.Alternative, indent, currentIndent)
	}
	return out
}

func formatWhileExpression(we *ast.WhileExpression, indent, currentIndent string) string {
	out := "wakati " + exprString(we.Condition) + " "
	if we.Consequence != nil {
		out += formatBlock(we.Consequence, indent, currentIndent)
	}
	return out
}

func formatFor(f *ast.For, indent, currentIndent string) string {
	out := "kwa "
	if f.StarterValue != nil {
		out += exprString(f.StarterValue) + "; "
	}
	if f.Condition != nil {
		out += exprString(f.Condition) + "; "
	}
	if f.Closer != nil {
		out += exprString(f.Closer) + " "
	}
	if f.Block != nil {
		out += formatBlock(f.Block, indent, currentIndent)
	}
	return out
}

func formatForIn(fi *ast.ForIn, indent, currentIndent string) string {
	out := "kwa "
	if fi.Key != "" {
		out += fi.Key + ", "
	}
	out += fi.Value + " ktk " + exprString(fi.Iterable) + " "
	if fi.Block != nil {
		out += formatBlock(fi.Block, indent, currentIndent)
	}
	return out
}

func formatSwitch(se *ast.SwitchExpression, indent, currentIndent string) string {
	var buf bytes.Buffer
	buf.WriteString("badili (")
	buf.WriteString(exprString(se.Value))
	buf.WriteString(")\n")
	buf.WriteString(currentIndent + "{\n")
	inner := currentIndent + indent
	for _, c := range se.Choices {
		if c == nil {
			continue
		}
		if c.Default {
			buf.WriteString(inner + "kawaida ")
		} else {
			buf.WriteString(inner + "ikiwa ")
			for i, e := range c.Expr {
				if i > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(exprString(e))
			}
			buf.WriteString(" ")
		}
		if c.Block != nil {
			buf.WriteString(formatBlock(c.Block, indent, inner))
		}
		buf.WriteByte('\n')
	}
	buf.WriteString(currentIndent + "}")
	return buf.String()
}

func formatTryCatch(tce *ast.TryCatchExpression, indent, currentIndent string) string {
	out := "jaribu " + formatBlock(tce.TryBlock, indent, currentIndent)
	out += " jaribu_bila " + formatBlock(tce.CatchBlock, indent, currentIndent)
	return out
}
