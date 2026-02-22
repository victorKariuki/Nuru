// This will convert the sequence of characters into a sequence of tokens

package lexer

import (
	"github.com/NuruProgramming/Nuru/token"
)

type Lexer struct {
	input        []rune
	position     int
	readPosition int
	ch           rune
	line         int
	column       int // 1-based column of the current character
}

func New(input string) *Lexer {
	l := &Lexer{input: []rune(input), line: 1, column: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	// Update line/column when advancing (not on first read)
	leavingPos := l.position
	l.position = l.readPosition
	l.readPosition += 1
	if l.readPosition > 1 && leavingPos >= 0 && leavingPos < len(l.input) {
		if l.input[leavingPos] == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	if l.ch == rune('/') && l.peekChar() == rune('/') {
		l.skipSingleLineComment()
		return l.NextToken()
	}
	if l.ch == rune('/') && l.peekChar() == rune('*') {
		l.skipMultiLineComment()
		return l.NextToken()
	}

	switch l.ch {
	case rune('='):
		if l.peekChar() == rune('=') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: col}
		} else {
			tok = newToken(token.ASSIGN, l.line, l.column, l.ch)
		}
	case rune(';'):
		tok = newToken(token.SEMICOLON, l.line, l.column, l.ch)
	case rune('('):
		tok = newToken(token.LPAREN, l.line, l.column, l.ch)
	case rune(')'):
		tok = newToken(token.RPAREN, l.line, l.column, l.ch)
	case rune('{'):
		tok = newToken(token.LBRACE, l.line, l.column, l.ch)
	case rune('}'):
		tok = newToken(token.RBRACE, l.line, l.column, l.ch)
	case rune(','):
		tok = newToken(token.COMMA, l.line, l.column, l.ch)
	case rune('+'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.PLUS_ASSIGN, Literal: string(ch) + string(l.ch), Line: l.line, Column: col}
		} else if l.peekChar() == rune('+') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.PLUS_PLUS, Literal: string(ch) + string(l.ch), Line: l.line, Column: col}
		} else {
			tok = newToken(token.PLUS, l.line, l.column, l.ch)
		}
	case rune('-'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.MINUS_ASSIGN, Literal: string(ch) + string(l.ch), Line: l.line, Column: col}
		} else if l.peekChar() == rune('-') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.MINUS_MINUS, Literal: string(ch) + string(l.ch), Line: l.line, Column: col}
		} else {
			tok = newToken(token.MINUS, l.line, l.column, l.ch)
		}
	case rune('!'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: l.line, Column: col}
		} else {
			tok = newToken(token.BANG, l.line, l.column, l.ch)
		}
	case rune('/'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.SLASH_ASSIGN, Line: l.line, Column: col, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.SLASH, l.line, l.column, l.ch)
		}
	case rune('*'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.ASTERISK_ASSIGN, Literal: string(ch) + string(l.ch), Line: l.line, Column: col}
		} else if l.peekChar() == rune('*') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.POW, Literal: string(ch) + string(l.ch), Line: l.line, Column: col}
		} else {
			tok = newToken(token.ASTERISK, l.line, l.column, l.ch)
		}
	case rune('<'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: col}
		} else {
			tok = newToken(token.LT, l.line, l.column, l.ch)
		}
	case rune('>'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Line: l.line, Column: col}
		} else {
			tok = newToken(token.GT, l.line, l.column, l.ch)
		}
	case rune('"'):
		tok.Type = token.STRING
		tok.Line = l.line
		tok.Column = l.column
		tok.Literal = l.readString()
	case rune('\''):
		tok = token.Token{Type: token.STRING, Literal: l.readSingleQuoteString(), Line: l.line, Column: l.column}
	case rune('['):
		tok = newToken(token.LBRACKET, l.line, l.column, l.ch)
	case rune(']'):
		tok = newToken(token.RBRACKET, l.line, l.column, l.ch)
	case rune(':'):
		tok = newToken(token.COLON, l.line, l.column, l.ch)
	case rune('@'):
		tok = newToken(token.AT, l.line, l.column, l.ch)
	case rune('.'):
		tok = newToken(token.DOT, l.line, l.column, l.ch)
	case rune('&'):
		if l.peekChar() == rune('&') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: l.line, Column: col}
		}
	case rune('|'):
		if l.peekChar() == rune('|') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: l.line, Column: col}
		}
	case rune('%'):
		if l.peekChar() == rune('=') {
			ch := l.ch
			col := l.column
			l.readChar()
			tok = token.Token{Type: token.MODULUS_ASSIGN, Line: l.line, Column: col, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.MODULUS, l.line, l.column, l.ch)
		}
	case rune('#'):
		if l.peekChar() == rune('!') && l.line == 1 {
			l.skipSingleLineComment()
			return l.NextToken()
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLetter(l.ch) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) && isLetter(l.peekChar()) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok = l.readDecimal()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.line, l.column, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, line, column int, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '@'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) readDecimal() token.Token {
	startLine, startCol := l.line, l.column
	integer := l.readNumber()
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		fraction := l.readNumber()
		return token.Token{Type: token.FLOAT, Literal: integer + "." + fraction, Line: startLine, Column: startCol}
	}
	return token.Token{Type: token.INT, Literal: integer, Line: startLine, Column: startCol}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return rune(0)
	} else {
		return l.input[l.readPosition]
	}
}

// func (l *Lexer) peekTwoChar() rune {
// 	if l.readPosition+1 >= len(l.input) {
// 		return rune(0)
// 	} else {
// 		return l.input[l.readPosition+1]
// 	}
// }

func (l *Lexer) skipSingleLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	l.skipWhitespace()
}

func (l *Lexer) skipMultiLineComment() {
	endFound := false

	for !endFound {
		if l.ch == 0 {
			endFound = true
		}

		if l.ch == '*' && l.peekChar() == '/' {
			endFound = true
			l.readChar()
		}

		l.readChar()
		l.skipWhitespace()
	}

}

func (l *Lexer) readString() string {
	var str string
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		} else if l.ch == '\\' {
			switch l.peekChar() {
			case 'n':
				l.readChar()
				l.ch = '\n'
			case 'r':
				l.readChar()
				l.ch = '\r'
			case 't':
				l.readChar()
				l.ch = '\t'
			case '"':
				l.readChar()
				l.ch = '"'
			case '\\':
				l.readChar()
				l.ch = '\\'
			}
		}
		str += string(l.ch)
	}
	return str
}

func (l *Lexer) readSingleQuoteString() string {
	var str string
	for {
		l.readChar()
		if l.ch == '\'' || l.ch == 0 {
			break
		} else if l.ch == '\\' {
			switch l.peekChar() {
			case 'n':
				l.readChar()
				l.ch = '\n'
			case 'r':
				l.readChar()
				l.ch = '\r'
			case 't':
				l.readChar()
				l.ch = '\t'
			case '"':
				l.readChar()
				l.ch = '"'
			case '\\':
				l.readChar()
				l.ch = '\\'
			}
		}
		str += string(l.ch)
	}
	return str
}
