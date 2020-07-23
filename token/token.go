package token

// Type represents the type of token
type Type string

// Token holds a single token type and its literal value
type Token struct {
	Type    Type
	Literal string
}

const (
	// ILLEGAL represents an illegal token
	ILLEGAL = "ILLEGAL"

	// EOF end of file
	EOF = "EOF"

	// IDENT an identifier, e.g: add, foobar, x, y, ...
	IDENT = "IDENT"

	// INT an integer, e.g: 1234
	INT = "INT"

	//
	// Operators
	//

	// ASSIGN the assignment operator
	ASSIGN = "="

	// PLUS the addition operator
	PLUS = "+"

	// MINUS the subtraction operator
	MINUS = "-"

	// MULTIPLY the multiplication operator
	MULTIPLY = "*"

	// DEVIDE the division operator
	DIVIDE = "/"

	//
	// Logical operators
	//

	// NOT the not operator
	NOT = "!"
	// AND the logical and operator
	AND = "&&"
	// OR the logical or operator
	OR = "||"

	//
	// Comparision operators
	//

	// LT the less than comparision operator
	LT = "<"
	// EQ the equality operator
	EQ = "=="

	//
	// Delimiters
	//

	// LPAREN a left parenthesis
	LPAREN = "("
	// RPAREN a right parenthesis
	RPAREN = ")"
	// LBRACE a left brace
	LBRACE = "{"
	// RBRACE a right brace
	RBRACE = "}"
	// COMMA a comma
	COMMA = ","
	// SEMICOLON a semi-colon
	SEMICOLON = ";"

	//
	// Keywords
	//

	// FN the `fn` keyword
	FN = "FN"
	// IF the `if` keyword
	IF = "IF"
	// ELSE the `else` keyword
	ELSE = "ELSE"
	// RETURN the `return` keyword
	RETURN = "RETURN"
	// WHILE the `while` keyword
	WHILE = "WHILE"
	// PUTS the `puts` keyword
	PUTS = "PUTS"
)

var keywords = map[string]Type{
	"fn":     FN,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
	"puts":   PUTS,
}

// LookupIdent looks up the identifier in ident and returns the appropriate
// token type depending on whether the identifier is user-defined or a keyword
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
