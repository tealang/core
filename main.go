// Tealang runtime REPL tool.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tealang/tea-go/tea/repl"
)

const (
	welcomeText = `Tealang v0.1-alpha
Copyright 2017 Lennart Espe. All rights reserved.`
	replSymbol = "~> "
)

func main() {
	fmt.Println(welcomeText)

	reader := bufio.NewReader(os.Stdin)
	runtime := repl.New()

	for runtime.Active {
		fmt.Print(replSymbol)
		input, err := reader.ReadString('\n')
		if err != nil {
			runtime.Stop()
		} else {
			output := runtime.Interpret(strings.TrimRight(input, "\n"))
			fmt.Println(output)
		}
	}
	fmt.Println()
}
