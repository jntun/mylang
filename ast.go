package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Expression interface{}

type Equality struct{}
type Primary struct{}

type Grouping struct {
	Expr Expression
}
type Unary struct {
	Op   Operator
	Expr Expression
}
type Binary struct {
	Left  Expression
	Op    Operator
	Right Expression
}
type Operator struct{ Token Token }

// A Literal is a number, string, boolean, or nil
type Literal struct{ Token }

// Value is the base atom for all derived jlang types
// Different type(s) implementations are determined at run time
type Value interface{}

func (l Literal) get() (Value, error) {
	switch l.Type {
	case Number:
		if strings.Contains(l.Lexeme, ".") {
			fmt.Println("Number is rational:", l.Lexeme)
			return strconv.ParseFloat(l.Lexeme, 64)
		}
		fmt.Println("Number is irrational:", l.Lexeme)
		return strconv.Atoi(l.Lexeme)
	case String:
		return l.Lexeme, nil
	case True:
		return true, nil
	case False:
		return false, nil
	case EOF:
		return "EOF", nil
	}

	return nil, fmt.Errorf("Unable to match literal %s, with a known value", l.Token.Lexeme)
}
