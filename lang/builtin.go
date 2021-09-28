package lang

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"time"
)

type Len struct {
	v Value
}

func (l Len) evaluate(intptr *Interpreter) (Value, error) {
	v, err := intptr.VariableResolver(Variable{Token{"v", Identifier, 0}})
	if err != nil {
		return nil, err
	}

	if kind := reflect.TypeOf(v).Kind(); kind == reflect.Array || kind == reflect.Slice {
		return len(v.([]*Value)) - 1, nil
	} else if kind == reflect.String {
		return len(v.(string)) - 1, nil
	} else {
		return 0, fmt.Errorf("type '%s' doesn't have len() implementation", kind)
	}
}

type Time struct{}

func (t Time) evaluate(intptr *Interpreter) (Value, error) {
	return int(time.Now().UnixNano() / 1000000), nil
}

type Pow struct{}

func (p Pow) evaluate(intptr *Interpreter) (Value, error) {
	x, err := intptr.VariableResolver(Variable{Token{"x", Identifier, 0}})
	if err != nil {
		return nil, err
	}
	y, err := intptr.VariableResolver(Variable{Token{"y", Identifier, 0}})
	if err != nil {
		return nil, err
	}

	switch reflect.TypeOf(x).Kind() {
	case reflect.Int:
		if yKind := reflect.TypeOf(y).Kind(); yKind == reflect.Int {
			return math.Pow(float64(x.(int)), float64(y.(int))), nil
		} else if yKind == reflect.Float64 {
			return math.Pow(float64(x.(int)), y.(float64)), nil
		} else {
			return nil, fmt.Errorf("invalid type '%s' in 'pow' call.", yKind)
		}
	case reflect.Float64:
		if yKind := reflect.TypeOf(y).Kind(); yKind == reflect.Int {
			return math.Pow(x.(float64), float64(y.(int))), nil
		} else if yKind == reflect.Float64 {
			return math.Pow(x.(float64), y.(float64)), nil
		} else {
			return nil, fmt.Errorf("invalid type '%s' in 'pow' call.", yKind)
		}
	}

	return nil, fmt.Errorf("invalid type '%s' in 'pow' call", reflect.TypeOf(x).Kind())
}

type Quit struct{}

func (q Quit) evaluate(intptr *Interpreter) (Value, error) {
	os.Exit(0)
	return nil, nil // Unreachable
}

func globals() []Statement {
	globals := make([]Statement, 0)

	globals = append(globals, VariableStatement{Token{"pi", Identifier, 0}, Literal{Token{"3.1415926535", Number, 0}}})

	globals = append(globals, makeBuiltinFunc("len", []string{"v"}, []Statement{
		ReturnStatement{Len{}, nil},
	}))
	globals = append(globals, makeBuiltinFunc("time", []string{}, []Statement{
		ReturnStatement{Time{}, nil},
	}))
	globals = append(globals, makeBuiltinFunc("pow", []string{"x", "y"}, []Statement{
		ReturnStatement{Pow{}, nil},
	}))
	globals = append(globals, makeBuiltinFunc("quit", []string{}, []Statement{
		ReturnStatement{Quit{}, nil},
	}))

	return globals
}

func makeBuiltinFunc(identifier string, args []string, block []Statement) FunctionDeclarationStatement {
	tokenArgs := make([]Token, len(args))
	for i, arg := range args {
		tokenArgs[i] = Token{arg, Identifier, 0}
	}

	return FunctionDeclarationStatement{Token{identifier, Identifier, 0}, &tokenArgs, block}
}
