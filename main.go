package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/d2verb/bee/generator"

	"github.com/d2verb/bee/checker"
	"github.com/d2verb/bee/lexer"
	"github.com/d2verb/bee/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE: bee <file>")
	} else {
		content, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		l := lexer.New(string(content))

		p := parser.New(l)
		program := p.ParseProgram()

		if errors := p.Errors(); len(errors) != 0 {
			for _, err := range errors {
				fmt.Println(err)
			}
			os.Exit(1)
		}

		c := checker.New(program)
		c.Check()

		if errors := c.Errors(); len(errors) != 0 {
			for _, err := range errors {
				fmt.Println(err)
			}
			os.Exit(1)
		}

		generator := generator.New(program)
		irProgram := generator.Generate()

		fmt.Print(irProgram.String())
	}
}
