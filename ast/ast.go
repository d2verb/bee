package ast

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Node defines an interface for all nodes in the AST
type Node interface {
	String() string
}

// Statement defines the interface for all statement nodes.
type Statement interface {
	Node
	statementNode()
}

// Expression defines the interface for all expression nodes.
type Expression interface {
	Node
	expressionNode()
}

// Program is a root node and consist of a slice of Function(s)
type Program struct {
	Functions []*Function
}

// String returns a stringified version of the AST for debugging
func (p *Program) String() string {
	var out bytes.Buffer

	for _, f := range p.Functions {
		out.WriteString(f.String())
	}

	return out.String()
}

// Function is a top level node and represents a function
type Function struct {
	Name       string
	Parameters []*Identifier
	Body       *BlockStatement
}

// String returns a stringified version of the AST for debugging
func (fn *Function) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fn.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fmt.Sprintf("fn %s", fn.Name))
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fn.Body.String())

	return out.String()
}

// Identifier represents an identifier and holds the name of the identifier
type Identifier struct {
	Value string
}

func (id *Identifier) expressionNode() {}

// String returns a stringified version of the AST for debugging
func (id *Identifier) String() string { return id.Value }

// ExpressionStatement represents an expression statement and holds an expression
type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// String returns a stringified version of the AST for debugging
func (es *ExpressionStatement) String() string {
	return fmt.Sprintf("%s;", es.Expression.String())
}

// IntegerLiteral represents al literal integer and holds an integer value
type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// String returns a stringified version of the AST for debugging
func (il *IntegerLiteral) String() string {
	return strconv.FormatInt(il.Value, 10)
}

// PrefixExpression represents a prefix expression and holds the operator
// as well as the right-hand side expression
// e.g: !is_valid
type PrefixExpression struct {
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// String returns a stringified version of the AST for debugging
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("%s(%s)", pe.Operator, pe.Right.String())
}

// InfixExpression represents an infix expression and holds the left-hand
// expression, operator and right-hand expression
// e.g.: 1 + 2
type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

// String returns a stringified version of the AST for debugging
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s%s%s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

// CallExpression represents a call expression and holds the function to be
// called as well as the arguments to be passed to that function
type CallExpression struct {
	Function  string
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// String returns a stringified version of the AST for debugging
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function)
	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")

	return out.String()
}

// BlockStatement represents a block statement and holds one or more other
// statements
type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// String returns a stringified version of the AST for debugging
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}")

	return out.String()
}

// ReturnStatement represenets the `return` statement node
// e.g: return 1234;
type ReturnStatement struct {
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}

// String returns a stringified version of the AST for debugging
func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %s;", rs.Value.String())
}

// IfStatement represents an `if` statement and holds the condition,
// consequence and alternative expressions
type IfStatement struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode() {}

// String returns a stringified version of the AST for debugging
func (is *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString("if(")
	out.WriteString(is.Condition.String())
	out.WriteString(")")
	out.WriteString(is.Consequence.String())

	if is.Alternative != nil {
		out.WriteString("else")
		out.WriteString(is.Alternative.String())
	}

	return out.String()
}

// WhileStatement represents an `while` statement and holds the condition,
// and consequence expression
type WhileStatement struct {
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileStatement) statementNode() {}

// String returns a stringified version of the AST for debugging
func (we *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString("while(")
	out.WriteString(we.Condition.String())
	out.WriteString(")")
	out.WriteString(we.Body.String())

	return out.String()
}

// PutsStatement represents an `puts` statement and holds the argument
// e.g: puts(1234);
type PutsStatement struct {
	Value Expression
}

func (ps *PutsStatement) statementNode() {}

// String returns a stringified version of the AST for debugging
func (ps *PutsStatement) String() string {
	return fmt.Sprintf("puts %s;", ps.Value.String())
}
