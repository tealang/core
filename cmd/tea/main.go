// Tealang runtime REPL tool.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tealang/tea-go/repl"
)

const (
	welcomeText = `Tealang v0.1.0-alpha
Copyright 2017 Lennart Espe. All rights reserved.`
	replSymbol = "> "
)

func main() {
	interactive := flag.Bool("i", false, "Start in interactive mode")
	graphviz := flag.Bool("g", false, "Enable GraphViz visualization mode")
	flag.Parse()

	env := repl.New(*graphviz)
	if filename := flag.Arg(0); filename != "" {
		code, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		_, err = env.Interpret(string(code))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		*interactive = true
		fmt.Println(welcomeText)
	}

	if !*interactive {
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for env.Active {
		fmt.Print(replSymbol)
		input, err := reader.ReadString('\n')
		if err != nil {
			env.Stop()
		} else {
			output, err := env.Interpret(strings.TrimRight(input, "\n"))
			if err != nil {
				fmt.Println(err)
			} else if output != "" {
				fmt.Println(output)
			}
		}
	}
	fmt.Println()
}
