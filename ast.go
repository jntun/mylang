package main

type Expression interface {
	evaluate() (Value, error)
}

type Statement interface {
	execute() error
}

type Declaration interface{}

// Program is a the highest level node in a Jlang program AST.
type Program struct {
	Statements []Statement
}

type VariableStatement struct {
	Identifier Token
	Expr       Expression
	resolver   func(VariableStatement)
}

type AssignmentStatement struct {
	VariableStatement
	resolver func(Variable) (Value, error)
}

type IfStatement struct {
	Expr  Expression
	block []Statement
}

type WhileStatement struct {
	test  Expression
	block []Statement
}

type ExpressionStatement struct {
	Expression
}

type PrintStatement struct {
	Expression
}
type Binary struct {
	Left  Expression
	Op    Operator
	Right Expression
}

type Grouping struct {
	Expr Expression
}
type Unary struct {
	Op   Operator
	Expr Expression
}

type Variable struct {
	name     Token
	resolver func(Variable) (Value, error)
}

type Operator struct{ Token }

// A Literal is a number, string, boolean, or nil
type Literal struct{ Token }

// Value is the base atom for all derived jlang types
// Different type(s) implementations are determined at run time
type Value interface{}
