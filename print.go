package main

import "fmt"

func (stmt FunctionDeclarationStatement) String() string {
	return fmt.Sprintf("func: %s | args: %v | block: %v |\n", stmt.Identifier.Lexeme, *stmt.args, stmt.block)
}

func (stmt ReturnStmt) String() string {
	return fmt.Sprintf("return %v", stmt.Expression)
}

func (un Unary) String() string {
	return fmt.Sprintf("%s%s", un.Op.Lexeme, un.Expr)
}

func (literal Literal) String() string {
	return fmt.Sprintf("literal - %s", literal.Lexeme)
}

func (v Variable) String() string {
	return v.identifier.Lexeme
}
