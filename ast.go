package main

type Expression interface {
	evaluate() (Value, error)
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
