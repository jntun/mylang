package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

// TODO Logging system for the interpreter

type Interpreter struct {
	s *Scanner
	p *Parser
}

// Interpret accepts an input string and attempts to execute the given sequence
// If a fatal error is encountered at any point, the Interpreter will break out and return an error
// describing the problem
func (intptr *Interpreter) Interpret(input string) error {
	tokens, err := intptr.s.Scan(input)
	if err != nil {
		return ScanError{err}
	}

	ast, err := intptr.p.Parse(append(tokens, Token{"EOF", EOF, tokens[len(tokens)-1].Line}))
	if err != nil {
		intptr.p.error = nil
		return err
	}

	if len(intptr.p.Errors) > 1 {
		for i, err2 := range intptr.p.Errors {
			fmt.Printf("Error %d: %s\n", i, err2)
		}
	}

	err = intptr.interpret(*ast)
	if err != nil {
		return err
	}

	return nil
}

func (intptr *Interpreter) interpret(ast Program) error {
	val, err := ast.evaluate()

	if err != nil {
		return err
	}

	fmt.Println(reflect.TypeOf(val), ":", val)
	return nil
}

// File accepts a direct source file path, reads it, and then calls Interpret() with the file string
func (intptr *Interpreter) File(filepath string) error {
	//log.Printf("Scanning file %s...\n", filepath)
	src, err := openFile(filepath)
	if err != nil {
		return err
	}

	if err = intptr.Interpret(*src); err != nil {
		return err
	}

	return nil
}

func (intptr *Interpreter) flush() {
	intptr.s.flush()
	intptr.p.flush()
}

func NewInterpreter() *Interpreter {
	return &Interpreter{&Scanner{}, &Parser{}}
}

func openFile(filepath string) (*string, error) {
	if _, err := os.Stat(filepath); err != nil {
		return nil, UnknownFile{filepath, err}
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, FileReadFailure{filepath, err}
	}
	dat := string(data)

	return &dat, nil
}
