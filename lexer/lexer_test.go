package lexer

import (
	"testing"

	"github.com/d2verb/bee/token"
)

func TestNextToken(t *testing.T) {
	input := `
foo_bar1234 551 = + - * /
! < == && || ( ) { } , fn if else return while var puts
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.NEWLINE, "\n"},
		{token.IDENT, "foo_bar1234"},
		{token.INT, "551"},
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
		{token.NEWLINE, "\n"},
		{token.BANG, "!"},
		{token.LT, "<"},
		{token.EQ, "=="},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.FN, "fn"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.RETURN, "return"},
		{token.WHILE, "while"},
		{token.VAR, "var"},
		{token.PUTS, "puts"},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - wrong token type. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - wrong literal. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
	}
}
