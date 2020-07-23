package lexer

import (
	"testing"

	"github.com/d2verb/bee/token"
)

func TestNextToken(t *testing.T) {
	input := `
foo_bar 551 = + - * /
! < == && || ( ) { } , ; fn if else return while puts`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "foo_bar"},
		{token.INT, "551"},
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.ASTERISK, "*"},
		{token.SLASH, "/"},
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
		{token.SEMICOLON, ";"},
		{token.FN, "fn"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.RETURN, "return"},
		{token.WHILE, "while"},
		{token.PUTS, "puts"},
		{token.EOF, " "},
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
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
