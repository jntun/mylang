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
		fmt.Println("Error occured - rs", err)
	}

	fmt.Println(string(input[len(input)-2]))

	err = interpreter.Interpret(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

}

// TODO actual error system (maybe in interpreter?)
func RuntimeError(err error) {
	fmt.Println(err)
}
