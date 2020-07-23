package ast

import "testing"

func TestString(t *testing.T) {
	program := &Program{
		Functions: []*Function{
			&Function{
				Name: "add",
				Parameters: []*Identifier{
					&Identifier{Value: "x"},
					&Identifier{Value: "y"},
				},
				Body: &BlockStatement{
					Statements: []Statement{
						&ReturnStatement{
							ReturnValue: &InfixExpression{
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

	expected := "fn add(x,y){return (x+y);}"

	if program.String() != expected {
		t.Errorf("program.String() wrong. expected=%q, got=%q", expected, program.String())
	}
}
