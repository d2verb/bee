package parser

import (
	"fmt"
	"testing"

	"github.com/d2verb/bee/lexer"
)

func TestFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"fn main(){}", "fn main(){}"},
		{"fn main(){} fn foo(){}", "fn main(){}fn foo(){}"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestAssignExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"fn main(){ x = 5; }", "fn main(){(x=5);}"},
		{"fn main(){ x = 551; }", "fn main(){(x=551);}"},
		{"fn main(){ x = y; }", "fn main(){(x=y);}"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"fn main(){ return 5; }", "fn main(){return 5;}"},
		{"fn main(){ return x; }", "fn main(){return x;}"},
		{"fn main(){ return (x + 6); }", "fn main(){return (x+6);}"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestPutsStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"fn main(){ puts 5; }", "fn main(){puts 5;}"},
		{"fn main(){ puts x; }", "fn main(){puts x;}"},
		{"fn main(){ puts (x + 6); }", "fn main(){puts (x+6);}"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"fn main(){ 5 + 5; }", "fn main(){(5+5);}"},
		{"fn main(){ 5 - 5; }", "fn main(){(5-5);}"},
		{"fn main(){ 5 * 5; }", "fn main(){(5*5);}"},
		{"fn main(){ 5 / 5; }", "fn main(){(5/5);}"},
		{"fn main(){ 5 == 5; }", "fn main(){(5==5);}"},
		{"fn main(){ 5 < 5; }", "fn main(){(5<5);}"},
		{"fn main(){ 5 && 5; }", "fn main(){(5&&5);}"},
		{"fn main(){ 5 || 5; }", "fn main(){(5||5);}"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"fn main(){ a + b + c; }", "fn main(){((a+b)+c);}"},
		{"fn main(){ a - b + c; }", "fn main(){((a-b)+c);}"},
		{"fn main(){ a + b - c; }", "fn main(){((a+b)-c);}"},
		{"fn main(){ a * b + c; }", "fn main(){((a*b)+c);}"},
		{"fn main(){ a + b * c; }", "fn main(){(a+(b*c));}"},
		{"fn main(){ a * b / c; }", "fn main(){((a*b)/c);}"},
		{"fn main(){ a / b * c; }", "fn main(){((a/b)*c);}"},
		{"fn main(){ !5 + 5; }", "fn main(){(!(5)+5);}"},
		{"fn main(){ !(5 + 6) + 5; }", "fn main(){(!((5+6))+5);}"},
		{"fn main(){ 1 == 3 < 4; }", "fn main(){(1==(3<4));}"},
		{"fn main(){ 1 + 0 == 3 < 4; }", "fn main(){((1+0)==(3<4));}"},
		{"fn main(){ 1 + 0 && 3 < 4; }", "fn main(){((1+0)&&(3<4));}"},
		{"fn main(){ 1 + 0 || 3 < 4; }", "fn main(){((1+0)||(3<4));}"},
		{"fn main(){ x < a && x == y; }", "fn main(){((x<a)&&(x==y));}"},
		{"fn main(){ x = a < 5 && x == y; }", "fn main(){(x=((a<5)&&(x==y)));}"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestWhileStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"fn main(){ while (1) {} }", "fn main(){while(1){}}"},
		{"fn main(){ while (x<y) {x=x+1;} }", "fn main(){while((x<y)){(x=(x+1));}}"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		fmt.Println("[DEBUG] " + p.curToken.Literal)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"fn main(){ if (1) {} }", "fn main(){if(1){}}"},
		{"fn main(){ if (x<y) {x=x+1;} }", "fn main(){if((x<y)){(x=(x+1));}}"},
		{"fn main(){ if (x<y) {x=x+1;} else {puts 1;} }", "fn main(){if((x<y)){(x=(x+1));}else{puts 1;}}"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestTrailingSemicolon(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"fn main(){ puts 1 }", "fn main(){puts 1;}"},
		{"fn main(){ puts 1; }", "fn main(){puts 1;}"},
		{"fn main(){ puts 1 2 3 }", "fn main(){puts 1;2;3;}"},
		{"fn main(){ puts 1; main(1, 2) 3 }", "fn main(){puts 1;main(1,2);3;}"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
