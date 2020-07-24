package ast

import "testing"

func TestString(t *testing.T) {
	program := &Program{
		Functions: []*Function{
			&Function{
				Name: "foo",
				Parameters: []*Variable{
					&Variable{Name: "x"},
					&Variable{Name: "y"},
				},
				Body: &BlockStatement{
					Statements: []Statement{
						&ReturnStatement{
							Value: &InfixExpression{
								Left:     &Identifier{Value: "x"},
								Operator: "+",
								Right:    &Identifier{Value: "y"},
							},
						},
					},
				},
			},
		},
	}

	expected := "fn foo(x,y){return (x+y);}"

	if program.String() != expected {
		t.Errorf("program.String() wrong. expected=%q, got=%q", expected, program.String())
	}
}
