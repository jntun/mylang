package lang

import (
	"fmt"
	"strings"
)

func (stmt FunctionDeclarationStatement) String() string {
	return fmt.Sprintf("func: %s | args: %v | block: %v |\n", stmt.Identifier.Lexeme, *stmt.args, stmt.block)
}

func (stmt ReturnStatement) String() string {
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

func (block Block) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("\t-- scope %s --\n", block.id))
	b.WriteString(block.store.String())

	return b.String()
}

func (vars varMap) String() string {
	sb := strings.Builder{}
	for lexeme, value := range vars {
		if value != nil {
			sb.WriteString(fmt.Sprintf("%s=%v\n", lexeme, *value))
			continue
		}
		sb.WriteString(fmt.Sprintf("%s=%v\n", lexeme, nil))
	}
	return sb.String()
}

func (funcs funcMap) String() string {
	sb := strings.Builder{}
	for name, invocation := range funcs {
		if invocation.stmt.args != nil {
			sb.WriteString(fmt.Sprintf("%s(%v)", name, *invocation.stmt.args))
			continue
		}
		sb.WriteString(fmt.Sprintf("%s()", name))
	}
	return sb.String()
}
