package lang

type Expression interface {
	evaluate(intptr *Interpreter) (Value, error)
}

type Statement interface {
	execute(intptr *Interpreter) error
}

// Program is a the highest level node in a Jlang program AST.
type Program struct {
	Statements []Statement
}

type VariableStatement struct {
	Identifier Token
	Expr       Expression
}

type FunctionDeclarationStatement struct {
	Identifier Token
	args       *[]Token
	block      []Statement
}

type ArrayDeclarationStatement struct {
	Identifier Token
	ExprList   []Expression
}

type AssignmentStatement struct {
	VariableStatement
}

type IfStatement struct {
	Expr      Expression
	block     []Statement
	elseBlock *[]Statement
}

type WhileStatement struct {
	test  Expression
	block []Statement
}

type ForStatement struct {
	varStmt *VariableStatement
	test    Expression
	assign  AssignmentStatement
	block   []Statement
}

type ExpressionStatement struct {
	Expression
}

type ReturnStatement struct {
	Expression
	val Value
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
	identifier Token
}

type Call struct {
	identifer Token
	args      *[]Expression
}

type FunctionCall struct {
	identifier Token
	args       *[]Expression
}

type FunctionInvocation struct {
	stmt     FunctionDeclarationStatement
	argExprs *[]Expression
}

type ArrayAccess struct {
	identifier Token
	index      Expression
}

type Operator struct{ Token }

// A Literal is a number, string, boolean, or nil
type Literal struct{ Token }

// Value is the base atom for all derived jlang types
// Different type(s) implementations are determined at run time
type Value interface{}
