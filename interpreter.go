package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// TODO Logging system for the interpreter

type Interpreter struct {
	s       *Scanner
	p       *Parser
	varEnv  Environment
	funcRet *Value
	funcEnv map[string]FunctionInvocation
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
		errList := make([]error, 0)
		for _, err2 := range intptr.p.Errors {
			errList = append(errList, err2)
		}
	}

	err = intptr.interpret(*ast)
	if err != nil {
		return err
	}

	return nil
}

func (intptr *Interpreter) interpret(program Program) error {
	return program.execute(intptr)
}

// VariableMap is for hooking into the scanner when encountering a VariableStatement.
// This allows the interpreter to handle state and higher-order operations.
// It is assumed that when called, the Scanner has already determined it to be a lexically _valid_
// variable statement and now it's up to the interpreter to breathe life into it.
func (intptr *Interpreter) VariableMap(stmt VariableStatement) {
	if stmt.Expr == nil {
		//intptr.varEnv[stmt.Identifier.Lexeme] = nil
		intptr.varEnv.store(stmt.Identifier.Lexeme, nil)
		return
	}
	val, err := stmt.Expr.evaluate(intptr)
	if err != nil {
		fmt.Printf("%s\n", InternalError{30, fmt.Sprintf("Invalid variable binding: %s", err)})
	}

	intptr.varEnv.store(stmt.Identifier.Lexeme, &val)
}

// VariableResolver is how an Identifier gets resolved to a real Value.
// If it is invalid for any reason, an error is returned instead.
func (intptr *Interpreter) VariableResolver(variable Variable) (Value, error) {
	return intptr.varEnv.resolve(variable)
}

// FunctionMap assumes the parser has *correctly* parsed a
// FunctionDeclarationStatement and is now ready to be breathed life into from the interpreter.
func (intptr *Interpreter) FunctionMap(stmt FunctionDeclarationStatement) {
	intptr.funcEnv[stmt.Identifier.Lexeme] = FunctionInvocation{stmt, nil}
}

func (intptr *Interpreter) FunctionResolve(caller FunctionCall) (Value, error) {
	if fun, found := intptr.funcEnv[caller.identifier.Lexeme]; found == true {
		if err := fun.FillArgs(intptr, caller.args); err != nil {
			return nil, err
		}

		val, err := fun.evaluate(intptr)
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	return nil, BadCall{caller.identifier}
}

func (intptr *Interpreter) FunctionReturn(val Value) {
	intptr.funcRet = &val
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
	intptr := &Interpreter{varEnv: NewEnvironment("var-global"), funcEnv: make(map[string]FunctionInvocation)}
	intptr.s = &Scanner{}
	intptr.p = &Parser{}
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
