package lang

import (
	"fmt"
)

// Parser is how a Jlang token sequence gets parsed and turned into Expressions
type Parser struct {
	src     []Token
	current uint
	error   error
	Errors  []error
}

// Parse takes a sequence of scanned Tokens and turns them into a corresponding Jlang Program statement
// If the parser is unable to form a valid Program, it returns a ParseError specifying why it couldn't
func (p *Parser) Parse(tokens []Token) (*Program, error) {
	p.src = tokens
	p.flush()
	statements := make([]Statement, 0)
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		p.consume(Semicolon, "Want ';' to close statement.")

		statements = append(statements, stmt)
	}

	return &Program{statements}, nil
}

func (p *Parser) statement() (Statement, error) {
	switch p.advance().Type {
	case Print:
		return p.PrintStatement()
	case Var:
		return p.variableStatement()
	case Identifier:
		if p.match(Equal) {
			p.reverse().reverse()
			return p.assignmentStatement()
		}
	case If:
		return p.IfStatement()
	case While:
		return p.WhileStatement()
	case For:
		return p.ForStatement()
	case Function:
		return p.FunctionDeclaration()
	case Return:
		if expr := p.expression(); expr != nil {
			return ReturnStatement{expr, nil}, nil
		}
	}
	p.reverse()

	return p.ExpressionStatement()
}

func (p *Parser) FunctionDeclaration() (Statement, error) {
	identifier := p.consume(Identifier, "Expect identifier after 'func' keyword.")
	args := make([]Token, 0)
	var block []Statement
	var err error

	p.consume(LeftParen, "Expect '(' after function identifier.")

	for p.match(Identifier) {
		args = append(args, p.previous())
		if p.peek().is(RightParen) {
			break
		}
		if !p.peek().is(Comma) {
			break
		}
		p.consume(Comma, "Expect ',' to separate parameter names.")
	}

	p.consume(RightParen, "Want ')' to close function parameter(s).")

	if block, err = p.blockStatement("func"); err != nil {
		return nil, err
	}

	if len(args) == 0 {
		return FunctionDeclarationStatement{*identifier, nil, block}, nil
	}

	return FunctionDeclarationStatement{*identifier, &args, block}, nil
}

func (p *Parser) blockStatement(stmtType string) ([]Statement, error) {
	p.consume(LeftBrace, fmt.Sprintf("Want '{' after %s statement.", stmtType))
	block := make([]Statement, 0)
	for true {
		if p.peek().is(RightBrace) {
			break
		}
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		p.consume(Semicolon, "Want ';' in block statement.")
		block = append(block, stmt)
		if p.isAtEnd() {
			return nil, ParseError{token: p.src[p.current], msg: fmt.Sprintf("Couldn't find '}' to close %s statement before end of file.", stmtType)}
		}
	}
	p.consume(RightBrace, fmt.Sprintf("Want '}' to close %s statement.", stmtType))
	return block, nil
}

func (p *Parser) variableStatement() (Statement, error) {
	identifier := p.consume(Identifier, "Expect identifier after var keyword.")
	if identifier == nil {
		return nil, ParseError{p.src[p.current], "Invalid identifier."}
	}

	var expr Expression
	if p.match(Equal) {
		expr = p.expression()
	}

	stmt := VariableStatement{*identifier, expr}
	return stmt, nil
}

func (p *Parser) assignmentStatement() (Statement, error) {
	identifier := p.advance()

	if ok := p.expect(Equal); ok != nil {
		return nil, ok
	}
	expr := p.expression()

	return AssignmentStatement{VariableStatement{identifier, expr}}, nil
}

func (p *Parser) IfStatement() (Statement, error) {
	expr := p.expression()
	stmts := make([]Statement, 0)

	p.consume(LeftBrace, "Expect '{' after if statement expression.")
	for true {
		if p.peek().is(RightBrace) {
			p.consume(RightBrace, "Expect '}' to close if statement.")

			// Build else block
			if p.match(Else) {
				elseBlock, err := p.blockStatement("else")
				if err != nil {
					return nil, err
				}
				return IfStatement{expr, stmts, &elseBlock}, nil

			}

			return IfStatement{expr, stmts, nil}, nil
		}
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		p.consume(Semicolon, "Want ';' to end statement in if block.")
		stmts = append(stmts, stmt)
		if len(stmts) > 255 {
			return nil, InternalError{42, "Maximum statements in if block reached."}
		}
	}

	return IfStatement{expr, stmts, nil}, nil
}

func (p *Parser) WhileStatement() (Statement, error) {
	expr := p.expression()
	/*
		stmts := make([]Statement, 0)
		p.consume(LeftBrace, "Expect '{' after while statement expression.")
		for true {
			val := p.peek()
			if val.is(RightBrace) {
				p.consume(RightBrace, "Expect '}' after while statement.")
				return WhileStatement{expr, stmts}, nil
			}
			stmt, err := p.statement()
			if err != nil {
				return nil, err
			}
			stmts = append(stmts, stmt)
			if p.isAtEnd() {
				return nil, ParseError{p.src[p.current], "Want '}' to close while statement"}
			}

		}
	*/
	stmts, err := p.blockStatement("while")
	if err != nil {
		return nil, err
	}

	return WhileStatement{
		test:  expr,
		block: stmts,
	}, nil
}

func (p *Parser) ForStatement() (Statement, error) {
	var varStmt Statement
	var assign Statement
	var test Expression
	var err error

	if p.match(Var) {
		varStmt, err = p.variableStatement()
		if err != nil {
			return nil, err
		}
		p.consume(Semicolon, "Want ';' to close for var declaration.")
	}

	test = p.expression()
	if test == nil {
		return nil, p.error
	}

	p.consume(Semicolon, "Expected ';' after for statement condition.")

	if p.peek().is(Identifier) {
		assign, err = p.assignmentStatement()
		if err != nil {
			return nil, err
		}
	}

	stmts, err := p.blockStatement("for")
	if err != nil {
		return nil, err
	}

	var retStmt VariableStatement
	if varStmt != nil {
		retStmt = varStmt.(VariableStatement)
	}
	return ForStatement{&retStmt, test, assign.(AssignmentStatement), stmts}, nil
}

func (p *Parser) ExpressionStatement() (Statement, error) {
	expr := p.expression()
	if expr == nil {
		return nil, p.error
	}
	return ExpressionStatement{expr}, nil
}

func (p *Parser) PrintStatement() (Statement, error) {
	expr := p.expression()

	return PrintStatement{expr}, nil
}

/********** Recursive descent parsing **********/
func (p *Parser) expression() Expression {
	return p.equality()
}

func (p *Parser) equality() Expression {
	expr := p.comparison()
	for p.match(BangEqual, EqualEqual) {
		op := Operator{p.previous()}
		right := p.comparison()
		expr = Binary{expr, op, right}
	}
	return expr
}

func (p *Parser) comparison() Expression {
	expr := p.term()
	for p.match(Greater, GreaterEqual, Less, LessEqual, Mod) {
		op := Operator{p.previous()}
		right := p.term()
		expr = Binary{expr, op, right}
	}

	return expr
}

func (p *Parser) term() Expression {
	expr := p.factor()
	for p.match(Minus, Plus) {
		op := Operator{p.previous()}
		right := p.factor()
		expr = Binary{expr, op, right}
	}

	return expr
}

func (p *Parser) factor() Expression {
	expr := p.unary()
	for p.match(Slash, Star) {
		op := Operator{p.previous()}
		right := p.unary()
		if right == nil {
			return nil
		}
		expr = Binary{expr, op, right}
	}
	return expr
}

func (p *Parser) unary() Expression {
	if p.match(Bang, Minus) {
		op := Operator{p.previous()}
		right := p.unary()
		return Unary{op, right}
	}
	// TODO: appendable ++ or -- instead of pre-expression
	if p.match(PlusPlus, MinusMinus) {
		op := Operator{p.previous()}
		right := p.unary()
		return Unary{op, right}
	}

	return p.primary()
}

func (p *Parser) primary() Expression {
	if p.match(False) {
		return Literal{p.previous()}
	}
	if p.match(True) {
		return Literal{p.previous()}
	}
	if p.match(Nil) {
		return Literal{p.previous()}
	}

	if p.match(Number, String) {
		return Literal{p.previous()}
	}
	if p.match(Identifier) {
		if p.peek().is(LeftParen) {
			identifier := p.previous()
			var args []Expression
			p.advance()

			args = nil
			expr := p.primary()
			args = make([]Expression, 0)
			for expr != nil {

				args = append(args, expr)
				if !p.peek().is(RightParen) {
					p.consume(Comma, "Want ',' after argument.")
				} else {
					break
				}

				expr = p.primary()
			}
			p.consume(RightParen, "Want ')' to close call.")
			funcCall := FunctionCall{identifier, &args}
			return funcCall
		}
		return Variable{p.previous()}
	}

	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "Expect ')' after expression.")
		return Grouping{expr}
	}

	p.error = p.hadError(p.previous(), "Expected expression.")
	return nil
}

/******* Scanning state functions *********/
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(tokenType int) bool {
	if p.isAtEnd() {
		return false
	}
	if p.peek().is(tokenType) {
		return true
	}

	return false
}

func (p *Parser) expect(tokens ...int) error {
	for _, expType := range tokens {
		if !p.match(expType) {
			return ParseError{p.peek(), fmt.Sprintf("expected %s got '%s'", MasterTokenMap[expType], p.peek().Lexeme)}
		}
	}
	return nil
}

func (p *Parser) match(types ...int) bool {
	for _, token := range types {
		if p.check(token) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType int, msg string) *Token {
	if p.check(tokenType) {
		t := p.advance()
		return &t
	}

	p.hadError(p.peek(), msg)
	return nil
}

func (p *Parser) peek() Token {
	if p.current > uint(len(p.src)-1) {
		return p.src[len(p.src)-1]
	}
	return p.src[p.current]
}

func (p *Parser) previous() Token {
	if p.current == 0 {
		return p.src[0]
	}
	return p.src[p.current-1]
}

func (p *Parser) reverse() *Parser {
	if p.current == 0 {
		p.error = InternalError{13, "Tried to reverse while at 0."}
		return p
	}
	p.current -= 1
	return p
}

func (p *Parser) isAtEnd() bool {
	return p.peek().is(EOF)
}

func (p *Parser) hadError(token Token, msg string) ParseError {
	err := ParseError{token, msg}
	p.Errors = append(p.Errors, err)
	//RuntimeError(err)
	return err
}

func (p *Parser) flush() {
	//log.Println("Flushing parser...")
	p.current = 0
	p.Errors = make([]error, 0)
}