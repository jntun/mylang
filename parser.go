package main

import (
	"log"
)

type Parser struct {
	src     []Token
	current uint64
	error   error
	errors  []error
}

func (p *Parser) Parse(tokens []Token) (Expression, error) {
	p.current = 0
	p.src = tokens
	expr := p.parse()
	if expr == nil {
		return nil, p.error
	}

	return *expr, nil
}

func (p *Parser) parse() *Expression {
	expr := p.expression()
	// We don't need no try-catch :^)
	if p.error != nil {
		return nil
	}
	return &expr
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

	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "Expect ')' after expression.")
		return Grouping{expr}
	}

	p.hadError(p.previous(), "Expected expression.")
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
	if p.current > uint64(len(p.src)-1) {
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

func (p *Parser) isAtEnd() bool {
	return p.peek().is(EOF)
}

func (p *Parser) hadError(token Token, msg string) ParseError {
	err := ParseError{token, msg}
	p.errors = append(p.errors, err)
	RuntimeError(err)
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
	log.Println("Flushing parser...")
}
