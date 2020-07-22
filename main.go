package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE: bee <file>")
	} else {
		_, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	}
}
