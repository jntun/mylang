package main

import (
	"bufio"
	"fmt"
	"os"
)

var interpreter = NewInterpreter()

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		popInterpreter()
	}

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

func repl() {
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		InvalidInput(err)
	}
	val := input[len(input)-2]

	var append string
	switch val {
	case '{':
		append = multiline('}')
	case '(':
		append = multiline(')')
	}

	input += append
	err = interpreter.Interpret(input)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
}

func multiline(delim byte) string {
	fmt.Print(">\t")
	reader := bufio.NewReader(os.Stdin)
	test, err := reader.ReadString(delim)
	InvalidInput(err)
	return test
}

func InvalidInput(err error) {
	if err != nil {
		fmt.Println("Invalid input: ", err)
	}
}

// TODO: actual error system (maybe in interpreter?)
func RuntimeError(err error) {
	fmt.Println(err)
}

// todo is a drop in code holder for future features
func todo() {
	fmt.Println("TODO.")
}
