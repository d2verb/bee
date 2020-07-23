package token

// Type represents the type of token
type Type string

// Token holds a single token type and its literal value
type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	NEWLINE = "NEWLINE"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"
	LT       = "<"
	EQ       = "=="
	AND      = "&&"
	OR       = "||"

	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	COMMA     = ","
	SEMICOLON = ";"

	FN     = "FN"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	WHILE  = "WHILE"
	PUTS   = "PUTS"
)

var keywords = map[string]Type{
	"fn":     FN,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
	"puts":   PUTS,
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
