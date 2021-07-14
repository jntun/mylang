package main

import (
	"bufio"
	"fmt"
	"github.com/jntun/mylang/lang"
	"os"
)

var interpreter = lang.NewInterpreter()

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		popInterpreter()
	}

	processArgs(args)
	filename := args[0]
	if err := interpreter.File(filename); err != nil {
		RuntimeError(err)
	}
}

func popInterpreter() {
	for true {
		repl()
	}
}

var indent = 0

func repl() {
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		InvalidInput(err)
	}
	lastChar := input[len(input)-2]

	var appnd string
	switch lastChar {
	case '{':
		indent++
		appnd = multiline('}')
	case '(':
		indent++
		appnd = multiline(')')
	case '[':
		indent++
		appnd = multiline(']')
	}
	input += appnd

	err = interpreter.Interpret(input)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func space(i int) string {
	str := ""
	for x := 0; x < i; x++ {
		str += "   "
	}
	return str
}

func multiline(delim byte) string {
	fmt.Printf(">%s", space(indent))

	reader := bufio.NewReader(os.Stdin)
	appnd, err := reader.ReadString(delim)
	InvalidInput(err)

	indent--
	return appnd
}

func InvalidInput(err error) {
	if err != nil {
		fmt.Println("Invalid input: ", err)
	}
}

func processArgs(args []string) {
	for _, arg := range args {
		if arg == "-server" {
			// No return
			httpServer()
		}
	}
}

// TODO: actual error system (maybe in interpreter?)
func RuntimeError(err error) {
	fmt.Println(err)
}
