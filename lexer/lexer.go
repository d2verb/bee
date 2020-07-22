package lexer

import "github.com/d2verb/bee/token"

type Lexer struct {
	input        string
	position     int // current position
	readPosition int // next position to read
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	return l
}

func (l *Lexer) NextToken() token.Token {
	return token.Token{
		Type:    token.SLASH,
		Literal: "/",
	}
}
