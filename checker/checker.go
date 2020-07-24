package checker

import (
	"fmt"

	"github.com/d2verb/bee/ast"
)

type Checker struct {
	program   *ast.Program
	functions map[string]int
	variables map[string]struct{}
	errors    []string
}

func New(program *ast.Program) *Checker {
	c := &Checker{
		program:   program,
		functions: make(map[string]int),
		errors:    []string{},
	}
	return c
}

func (c *Checker) Check() {
	c.findFunctions()
	for _, function := range c.program.Functions {
		c.variables = make(map[string]struct{})

		for _, parameter := range function.Parameters {
			c.variables[parameter.Name] = struct{}{}
		}

		c.checkBlockStatement(function.Body)
	}
}

func (c *Checker) checkBlockStatement(node *ast.BlockStatement) {
	for _, statement := range node.Statements {
		c.checkStatement(statement)
	}
}

func (c *Checker) checkStatement(node ast.Statement) {
	switch node := node.(type) {
	case *ast.BlockStatement:
		c.checkBlockStatement(node)
	case *ast.ReturnStatement:
		c.checkExpression(node.Value)
	case *ast.PutsStatement:
		c.checkExpression(node.Value)
	case *ast.IfStatement:
		c.checkIfStatement(node)
	case *ast.WhileStatement:
		c.checkWhileStatement(node)
	case *ast.ExpressionStatement:
		c.checkExpression(node.Expression)
	}
}

func (c *Checker) checkIfStatement(node *ast.IfStatement) {
	c.checkExpression(node.Condition)
	c.checkBlockStatement(node.Consequence)

	if node.Alternative != nil {
		c.checkBlockStatement(node.Alternative)
	}
}

func (c *Checker) checkWhileStatement(node *ast.WhileStatement) {
	c.checkExpression(node.Condition)
	c.checkBlockStatement(node.Body)
}

func (c *Checker) checkExpression(node ast.Expression) {
	switch node := node.(type) {
	case *ast.InfixExpression:
		if node.Operator == "=" {
			c.checkExpression(node.Right)
			c.registerVariable(node.Left.(*ast.Identifier).Value)
		} else {
			c.checkExpression(node.Left)
			c.checkExpression(node.Right)
		}
	case *ast.PrefixExpression:
		c.checkExpression(node.Right)
	case *ast.CallExpression:
		if !c.isFunctionExists(node.Function) {
			msg := fmt.Sprintf("function '%s' is not defined", node.Function)
			c.errors = append(c.errors, msg)
			return
		}
		count, _ := c.functions[node.Function]
		if len(node.Arguments) != count {
			msg := fmt.Sprintf("the number of arguments for '%s' is not correct. expect=%d, got=%d",
				node.Function,
				count,
				len(node.Arguments))
			c.errors = append(c.errors, msg)
			return
		}
		for _, argument := range node.Arguments {
			c.checkExpression(argument)
		}
	case *ast.Identifier:
		if !c.isVariableExists(node.Value) {
			msg := fmt.Sprintf("variable '%s' is not defined", node.Value)
			c.errors = append(c.errors, msg)
		}
	}
}

func (c *Checker) Errors() []string {
	return c.errors
}

func (c *Checker) findFunctions() {
	for _, function := range c.program.Functions {
		c.functions[function.Name] = len(function.Parameters)
	}
}

func (c *Checker) isFunctionExists(functionName string) bool {
	_, ok := c.functions[functionName]
	return ok
}

func (c *Checker) isVariableExists(variableName string) bool {
	_, ok := c.variables[variableName]
	return ok
}

func (c *Checker) parameterCount(functionName string) int {
	count, _ := c.functions[functionName]
	return count
}

func (c *Checker) registerVariable(variableName string) {
	c.variables[variableName] = struct{}{}
}
