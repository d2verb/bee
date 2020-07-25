package checker

import (
	"testing"

	"github.com/d2verb/bee/parser"

	"github.com/d2verb/bee/lexer"
)

func TestChecker(t *testing.T) {
	tests := []struct {
		input  string
		errors []string
	}{
		{
			"fn main() {}",
			[]string{},
		},
		{

			"fn main() { main(1); main(1, 2) }",
			[]string{
				"the number of arguments for 'main' is not correct. expect=0, got=1",
				"the number of arguments for 'main' is not correct. expect=0, got=2",
			},
		},
		{

			"fn main() { puts z; main(z); }",
			[]string{
				"variable 'z' is not defined",
				"the number of arguments for 'main' is not correct. expect=0, got=1",
			},
		},
		{

			"fn main() { foo(1, 2); bar(1, 2); } fn foo(x, y) { return x + y; }",
			[]string{
				"function 'bar' is not defined",
			},
		},
		{

			// Check stopped at first function
			"fn main(x, x) {} fn foo(x, x) { main(1); }",
			[]string{
				"duplicated parameter 'x' in function 'main'",
			},
		},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()

		c := New(program)
		c.Check()

		errors := c.Errors()
		if len(errors) != len(tt.errors) {
			t.Errorf("[test-%d] the number of error message is not correct. expected=%d, got=%d",
				i, len(tt.errors), len(errors))
			break
		}

		for i := 0; i < len(errors); i++ {
			if errors[i] != tt.errors[i] {
				t.Errorf("[test-%d] error message is not correct. expected=%q, got=%q",
					i, tt.errors[i], errors[i])
			}
		}
	}
}

func TestVariableCreation(t *testing.T) {
	input := "fn main(x, y) { x = y; z = 6; x = z; s = z; }"
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()

	c := New(program)
	c.Check()

	function := program.Functions[0]

	if len(function.Variables) != 4 {
		t.Errorf("the number of local variable should be 4, but got %d", len(function.Variables))
	}

	for i, expectedName := range []string{"x", "y", "z", "s"} {
		if function.Variables[i].Name != expectedName {
			t.Errorf("name of variables[%d] should be '%s', but got '%s",
				i, expectedName, function.Variables[i].Name)
		}
	}

	for i := 0; i < 2; i++ {
		if function.Variables[i] != function.Parameters[i] {
			t.Errorf("variables[%d] = %s should equal to parameters[%d] = %s",
				i, function.Variables[i].Name, i, function.Parameters[i].Name)
		}
	}
}
