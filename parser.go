package main

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

		statements = append(statements, stmt)
	}

	return &Program{statements}, nil
}

func (p *Parser) statement() (Statement, error) {
	switch p.advance().Type {
	case Print:
		return p.printStatement()
	case Var:
		return p.variableStatement()
	case Identifier:
		if p.match(Equal) {
			p.reverse().reverse()
			return p.assignmentStatement()
		}
		if p.peek().is(LeftParen) {
			todo()
		}
	case If:
		return p.ifStatement()
	case While:
		return p.WhileStatement()
	case For:
		return p.ForStatement()
	case Function:
		todo()
		return nil, nil
	}
	p.reverse()

	return p.expressionStatement()
}

func (p *Parser) printStatement() (Statement, error) {
	expr := p.expression()
	if err := p.expect(Semicolon); err != nil {
		return nil, err
	}
	return PrintStatement{expr}, nil
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

	if ok := p.expect(Semicolon); ok != nil {
		return nil, ok
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
	if ok := p.expect(Semicolon); ok != nil {
		return nil, ok
	}

	return AssignmentStatement{VariableStatement{identifier, expr}}, nil
}

func (p *Parser) ifStatement() (Statement, error) {
	expr := p.expression()
	stmts := make([]Statement, 0)

	p.consume(LeftBrace, "Expect '{' after if statement expression.")
	for true {
		if p.peek().is(RightBrace) {
			p.consume(RightBrace, "Expect '}' to close if statement.")
			var elseBlock []Statement

			// Build else block
			if p.match(Else) {
				p.consume(LeftBrace, "Expect '{' after 'else' keyword.")
				elseBlock = make([]Statement, 0)
				for true {
					if p.match(RightBrace) {
						return IfStatement{expr, stmts, &elseBlock}, nil
					}
					stmt, err := p.statement()
					if err != nil {
						return nil, err
					}
					elseBlock = append(elseBlock, stmt)
					if p.isAtEnd() {
						return nil, ParseError{p.src[p.current], "Expect '}' to close else statement."}
					}
				}
			}

			return IfStatement{expr, stmts, nil}, nil
		}
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
		if len(stmts) > 255 {
			return nil, InternalError{42, "Maximum statements in if block reached."}
		}
	}

	return IfStatement{expr, stmts, nil}, nil
}

func (p *Parser) WhileStatement() (Statement, error) {
	expr := p.expression()
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

	return nil, InternalError{43, "Unable to parse while statement."}
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
	} else {

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

	p.consume(LeftBrace, "Want '{' after for statement.")
	stmts := make([]Statement, 0)
	for true {
		if p.peek().is(RightBrace) {
			p.advance()
			break
		}
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
		if p.isAtEnd() {
			return nil, ParseError{token: p.src[p.current], msg: "Want '}' to close for statement."}
		}
	}

	var retStmt VariableStatement
	if varStmt != nil {
		retStmt = varStmt.(VariableStatement)
	}
	return ForStatement{&retStmt, test, assign.(AssignmentStatement), stmts}, nil
}

func (p *Parser) expressionStatement() (Statement, error) {
	expr := p.expression()
	if expr == nil {
		return nil, p.error
	}
	if ok := p.expect(Semicolon); ok != nil {
		return nil, ok
	}
	return ExpressionStatement{expr}, nil
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
	for p.match(Greater, GreaterEqual, Less, LessEqual) {
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

func (p *Parser) sync() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().is(Semicolon) {
			return
		}

		switch p.peek().Type {
		case Class:
		case Function:
		case Var:
		case For:
		case If:
		case While:
		case Print:
		case Return:
			return
		}

		p.advance()
	}
}

func (p *Parser) flush() {
	//log.Println("Flushing parser...")
	p.current = 0
	p.Errors = make([]error, 0)
}
