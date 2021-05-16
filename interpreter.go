package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// TODO Logging system for the interpreter

type Interpreter struct {
	s *Scanner
	p *Parser
}

func (intptr *Interpreter) Interpret(input string) error {
	tokens, err := intptr.s.Scan(input)
	if err != nil {
		return ScanError{err}
	}

	ast, err := intptr.p.Parse(append(tokens, Token{"EOF", EOF, tokens[len(tokens)-1].Line}))
	if len(intptr.p.errors) != 0 {
		for i, err2 := range intptr.p.errors {
			fmt.Printf("Error %d: %s\n", i, err2)
		}
	}

	if err != nil {
		return err
	}

	fmt.Println(ast)
	return nil
}

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
