package main

type Expression interface {
	evaluate() (Value, error)
}

type Statement interface {
	do() error
}

// Program is a the highest level node in a Jlang program AST.
type Program struct {
	Statements []Statement
}

type ExpressionStatement struct {
	Expression
}

type PrintStatement struct {
	Expression
}

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

type Operator struct{ Token }

// A Literal is a number, string, boolean, or nil
type Literal struct{ Token }

// Value is the base atom for all derived jlang types
// Different type(s) implementations are determined at run time
type Value interface{}
