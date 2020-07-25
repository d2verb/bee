package parser

import (
	"fmt"
	"strconv"

	"github.com/d2verb/bee/ast"
	"github.com/d2verb/bee/lexer"
	"github.com/d2verb/bee/token"
)

const (
	_ int = iota
	LOWEST
	ASSIGN  // =
	AND     // && or ||
	EQUALS  // ==
	LESS    // <
	SUM     // + -
	PRODUCT // * /
	PREFIX  // !X
	CALL    // myFunction(X)
)

var precedences = map[token.Type]int{
	token.ASSIGN:   ASSIGN,
	token.EQ:       EQUALS,
	token.LT:       LESS,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.DIVIDE:   PRODUCT,
	token.MULTIPLY: PRODUCT,
	token.AND:      AND,
	token.OR:       AND,
	token.LPAREN:   CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser represents a parser and contains the internal state
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

// New creates a new parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.DIVIDE, p.parseInfixExpression)
	p.registerInfix(token.MULTIPLY, p.parseInfixExpression)
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	p.nextToken()
	p.nextToken()
	return p
}

// ParseProgram parses  thewhole source code and returns the AST
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Functions = []*ast.Function{}

	for p.curTokenIs(token.FN) {
		fn := p.parseFunction()
		program.Functions = append(program.Functions, fn)
		p.nextToken()
	}

	return program
}

func (p *Parser) parseFunction() *ast.Function {
	fn := &ast.Function{}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	fn.Name = p.curToken.Literal

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fn.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

func (p *Parser) parseFunctionParameters() []*ast.Variable {
	params := []*ast.Variable{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	params = append(params, &ast.Variable{Name: p.curToken.Literal})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if p.expectPeek(token.IDENT) {
			params = append(params, &ast.Variable{Name: p.curToken.Literal})
		} else {
			return nil
		}
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{}

	// skip `{`
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	if !p.expect(token.RBRACE) {
		return nil
	}

	return block
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.RETURN:
		return p.parseReturnStatement()
	case token.PUTS:
		return p.parsePutsStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.IF:
		return p.parseIfStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{}

	// skip `if`
	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		p.nextToken()
		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{}

	// skip `while`
	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parsePutsStatement() *ast.PutsStatement {
	stmt := &ast.PutsStatement{}

	// skip `puts`
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{}

	// skip `return`
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			p.noInfixParseFnError(p.peekToken.Type)
			return nil
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Name: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Operator: p.curToken.Literal,
		Left:     left,
	}

	if p.curToken.Type == token.ASSIGN {
		switch left.(type) {
		case *ast.Identifier:
			break
		default:
			msg := fmt.Sprintf("the left hand side of '=' must be identifier")
			p.errors = append(p.errors, msg)
			return nil
		}
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{}

	switch node := function.(type) {
	case *ast.Identifier:
		exp.Function = node.Name
		break
	default:
		msg := fmt.Sprintf("only identifier is allowed to call")
		p.errors = append(p.errors, msg)
		return nil
	}

	exp.Arguments = p.parseExpressionList(token.RPAREN)

	return exp
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()

	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noInfixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no infix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expect(t token.Type) bool {
	if p.curTokenIs(t) {
		return true
	}
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.curToken.Type)
	p.errors = append(p.errors, msg)
	return false
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
