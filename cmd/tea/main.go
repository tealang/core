// Tealang runtime REPL tool.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tealang/core/repl"
)

const (
	welcomeText = `Tealang v0.1.0-alpha
Copyright 2017 Lennart Espe. All rights reserved.`
	replSymbol = "> "
)

var (
	interactiveMode, graphvizMode bool
)

func main() {
	flag.BoolVar(&interactiveMode, "i", false, "Start in interactive mode")
	flag.BoolVar(&graphvizMode, "g", false, "Enable GraphViz visualization mode")
	flag.Parse()

	env := repl.New(repl.Config{
		OutputGraph: graphvizMode,
	})
	if filename := flag.Arg(0); filename != "" {
		err := env.Load(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		interactiveMode = true
		fmt.Fprintln(os.Stdout, welcomeText)
	}

	if !interactiveMode {
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for env.IsActive() {
		fmt.Fprint(os.Stdout, replSymbol)
		input, err := reader.ReadString('\n')
		if err != nil {
			env.Stop()
		} else {
			output, err := env.Interpret(strings.TrimRight(input, "\n"))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else if output != "" {
				fmt.Fprintln(os.Stdout, output)
			}
		}
	}
	fmt.Fprintln(os.Stdout)
}
