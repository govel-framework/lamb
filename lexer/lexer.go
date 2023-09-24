package lexer

import (
	"github.com/govel-framework/lamb/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	Line         int
	Column       int
	ch           byte
	inCode       bool
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.Line++

	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0

		if l.Column == 0 {
			l.Column = 1
		}

		return
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition

	l.Column++
	l.readPosition += 1

	switch l.ch {
	case '\n':
		l.Line++
		l.Column = 0
	case '\t':
		l.Column += 3
	}

}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	if !l.inCode && l.ch != 0 {

		if l.ch == '{' && l.peekChar() == '?' {
			l.inCode = true
			l.readChar()
			l.readChar()

		} else {

			tok.Type = token.HTML
			tok.Literal = string(l.ch)
			tok.Line = l.Line
			tok.Col = l.Column

			l.readChar()

			return tok
		}

	}

	l.skipWhitespace()

	if l.ch == '?' && l.peekChar() == '}' {
		l.inCode = false
		l.readChar()
		l.readChar()

		tok.Col = l.Column
		tok.Line = l.Line
		tok.Literal = ""
		tok.Type = token.EOC

		return tok
	}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = l.newToken(token.ASSIGN, l.ch)
		}

	case '+':
		tok = l.newToken(token.PLUS, l.ch)

	case '-':
		tok = l.newToken(token.MINUS, l.ch)

	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}

		} else {
			tok = l.newToken(token.BANG, l.ch)
		}

	case '*':
		tok = l.newToken(token.ASTERISK, l.ch)

	case '/':
		tok = l.newToken(token.SLASH, l.ch)

	case '<':
		tok = l.newToken(token.LT, l.ch)

	case '>':
		tok = l.newToken(token.GT, l.ch)

	case ';':
		tok = l.newToken(token.SEMICOLON, l.ch)

	case '(':
		tok = l.newToken(token.LPAREN, l.ch)

	case ')':
		tok = l.newToken(token.RPAREN, l.ch)

	case ',':
		tok = l.newToken(token.COMMA, l.ch)

	case '{':
		tok = l.newToken(token.LBRACE, l.ch)

	case '}':
		tok = l.newToken(token.RBRACE, l.ch)

	case '"':
		tok = l.readString('"')

	case '\'':
		tok = l.readString('\'')

	case '[':
		tok = l.newToken(token.LBRACKET, l.ch)

	case ']':
		tok = l.newToken(token.RBRACKET, l.ch)

	case ':':
		tok = l.newToken(token.COLON, l.ch)

	case '.':
		tok = l.newToken(token.DOT, l.ch)

	case '#':
		l.readComment()

		return l.NextToken()

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.Line
		tok.Col = l.Column

	default:
		if isLetter(l.ch) {
			tok.Col = l.Column
			tok.Line = l.Line
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookUpIdent(tok.Literal)

			return tok

		} else if isDigit(l.ch) {
			tok.Col = l.Column
			tok.Line = l.Line
			tok.Type = token.INT
			tok.Literal = l.readNumber()

			return tok

		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	pos := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() string {
	pos := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Col: l.Column, Line: l.Line}
}

func (l *Lexer) readString(char byte) token.Token {
	var tok token.Token
	tok.Col = l.Column
	tok.Line = l.Line

	position := l.position + 1

	for {
		l.readChar()

		if l.ch == char || l.ch == 0 {
			break
		}
	}

	tok.Type = token.STRING
	tok.Literal = l.input[position-1 : l.position+1]

	return tok
}

func (l *Lexer) readComment() {
	l.readChar()

	for {
		if l.ch == '#' || l.ch == 0 {
			l.readChar()
			break
		}

		l.readChar()
	}
}
