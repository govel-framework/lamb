package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Col     int
	Line    int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	EOC     = "EOC"

	// Identifiers
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"
	HTML   = "HTML"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	DOT       = "."

	LPAREN = "("
	RPAREN = ")"

	LBRACE = "{"
	RBRACE = "}"

	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	VAR        = "var"
	TRUE       = "true"
	FALSE      = "false"
	IF         = "if"
	ELSE       = "else"
	ENDIF      = "endif"
	FOR        = "for"
	ENDFOR     = "endfor"
	IN         = "in"
	EXTENDS    = "extends"
	SECTION    = "section"
	ENDSECTION = "endsection"
	DEFINE     = "define"
	END        = "end"
	INCLUDE    = "include"
	AND        = "and"
)

var keywords = map[string]TokenType{
	"var":        VAR,
	"true":       TRUE,
	"false":      FALSE,
	"if":         IF,
	"else":       ELSE,
	"endif":      ENDIF,
	"for":        FOR,
	"endfor":     ENDFOR,
	"in":         IN,
	"extends":    EXTENDS,
	"section":    SECTION,
	"endsection": ENDSECTION,
	"define":     DEFINE,
	"end":        END,
	"include":    INCLUDE,
	"and":        AND,
}

func LookUpIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
