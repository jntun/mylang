package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Expression interface{}

// A Literal is a number, string, boolean, or nil
type Literal struct{ Token }

type Grouping struct {
	Left  Token
	Expr  Expression
	Right Token
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
type Operator struct{ Token *Token }

// Value is the base atom for all derived jlang types
// Different type(s) implementations are determined at run time
type Value interface{}

// Node is a valid AST node in a Parser
type Node interface {
	Value() (Value, error)
}

func (l Literal) Value() (Value, error) {
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

	return nil, ParseError{fmt.Errorf("Unable to match literal %s, with a known value", l.Token.Lexeme)}
}
