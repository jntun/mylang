package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// TODO Logging system for the interpreter
type Interpreter struct {
	s      *Scanner
	p      *Parser
	global Environment
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

func (intptr *Interpreter) interpret(program Program) error {
	return program.execute()
}

// VariableMap is for hooking into the scanner when encountering a VariableStatement.
// This allows the interpreter to handle state and higher-order operations.
// It is assumed that when called, the Scanner has already determined it to be a lexically _valid_
// variable statement and now it's up to the interpreter to breathe life into it.
func (intptr *Interpreter) VariableMap(stmt VariableStatement) {
	if stmt.Expr == nil {
		//intptr.global[stmt.Identifier.Lexeme] = nil
		intptr.global.store(stmt.Identifier.Lexeme, nil)
		return
	}
	val, err := stmt.Expr.evaluate()
	if err != nil {
		fmt.Printf("%s\n", InternalError{30, fmt.Sprintf("Invalid variable binding: %s", err)})
	}

	//intptr.global[stmt.Identifier.Lexeme] = &val
	intptr.global.store(stmt.Identifier.Lexeme, &val)
}

// VariableResolver is how an Identifier gets resolved to a real Value.
// If it is invalid for any reason, an error is returned instead.
func (intptr *Interpreter) VariableResolver(variable Variable) (Value, error) {
	return intptr.global.resolve(variable)
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
	intptr := &Interpreter{global: NewEnvironment()}
	intptr.s = &Scanner{}
	intptr.p = &Parser{
		variableDecl:    intptr.VariableMap,
		variableResolve: intptr.VariableResolver,
	}
	return intptr
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
