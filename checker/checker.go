package checker

import (
	"fmt"

	"github.com/d2verb/bee/ast"
)

// Context represents context of semantic checker
type Context struct {
	function  *ast.Function
	variables map[string]struct{}
}

// Checker represents a semantic checker
type Checker struct {
	program    *ast.Program
	signatures map[string]int
	context    Context
	errors     []string
}

// New returns a new Checker
func New(program *ast.Program) *Checker {
	c := &Checker{
		program:    program,
		signatures: make(map[string]int),
		errors:     []string{},
	}
	return c
}

// Check does some semantic checking
func (c *Checker) Check() {
	c.checkFunctionSignature()

	if len(c.errors) != 0 {
		return
	}

	for _, function := range c.program.Functions {
		c.checkFunction(function)

		if len(c.errors) != 0 {
			return
		}
	}
}

func (c *Checker) checkFunction(function *ast.Function) {
	c.newContext(function)

	for _, parameter := range function.Parameters {
		c.context.variables[parameter.Name] = struct{}{}
	}

	c.checkBlockStatement(function.Body)
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
		count, _ := c.signatures[node.Function]
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

// Errors return errors of checker
func (c *Checker) Errors() []string {
	return c.errors
}

func (c *Checker) checkFunctionSignature() {
	for _, function := range c.program.Functions {
		if c.checkDuplicatedParameterExists(function) {
			return
		}
		c.signatures[function.Name] = len(function.Parameters)
	}
}

func (c *Checker) checkDuplicatedParameterExists(function *ast.Function) bool {
	parameters := map[string]struct{}{}
	for _, parameter := range function.Parameters {
		if _, ok := parameters[parameter.Name]; ok {
			msg := fmt.Sprintf("duplicated parameter '%s' in function '%s'", parameter.Name, function.Name)
			c.errors = append(c.errors, msg)
			return true
		}
		parameters[parameter.Name] = struct{}{}
	}
	return false
}

func (c *Checker) newContext(function *ast.Function) {
	c.context = Context{
		function:  function,
		variables: make(map[string]struct{}),
	}

}

func (c *Checker) isFunctionExists(functionName string) bool {
	_, ok := c.signatures[functionName]
	return ok
}

func (c *Checker) isVariableExists(variableName string) bool {
	_, ok := c.context.variables[variableName]
	return ok
}

func (c *Checker) parameterCount(functionName string) int {
	count, _ := c.signatures[functionName]
	return count
}

func (c *Checker) registerVariable(variableName string) {
	c.context.variables[variableName] = struct{}{}
	c.context.function.Variables = append(c.context.function.Variables,
		&ast.Variable{Name: variableName})
}
