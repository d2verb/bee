package main

import (
	"fmt"
	"io/ioutil"
	"os"

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

		if len(p.Errors()) != 0 {
			for _, err := range p.Errors() {
				fmt.Println(err)
			}
		} else {
			for _, function := range program.Functions {
				fmt.Println(function.String())
			}
		}
	}
}
