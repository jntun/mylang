package lang

import (
	"fmt"
	"reflect"
)

type Len struct {
	v Value
}

func (l Len) evaluate(intptr *Interpreter) (Value, error) {
	v, err := intptr.VariableResolver(Variable{identifier: Token{Lexeme: "v", Type: Identifier, Line: 0}})
	if err != nil {
		return nil, err
	}

	if kind := reflect.TypeOf(v).Kind(); kind == reflect.Array || kind == reflect.Slice {
		return len(v.([]*Value)) - 1, nil
	} else {
		return 0, fmt.Errorf("type '%s' doesn't have len() implementation", kind)
	}
}

func globals() []Statement {
	globals := make([]Statement, 0)

	globals = append(globals, VariableStatement{Identifier: Token{Lexeme: "pi", Type: Identifier, Line: 0}, Expr: Literal{Token{Lexeme: "3.1415926535", Type: Number, Line: 0}}})

	globals = append(globals, makeBuiltinFunc("len", []string{"v"}, []Statement{
		ReturnStatement{Len{}, nil},
	}))

	return globals
}

func makeBuiltinFunc(identifier string, args []string, block []Statement) FunctionDeclarationStatement {
	tokenArgs := make([]Token, len(args))
	for i, arg := range args {
		tokenArgs[i] = Token{arg, Identifier, 0}
	}

	return FunctionDeclarationStatement{
		Identifier: Token{identifier, Identifier, 0},
		args:       &tokenArgs,
		block:      block,
	}
}
