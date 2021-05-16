package main

import (
	"fmt"
	"log"
)

type Parser struct {
	src      []Token
	ast      []Node
	current  uint64
	tokenMap map[int][]func(t Token) error // maps a TokenType(int) to a group of monads that are interested in the given type
	fatal    error
	errors   []error
}

func (p *Parser) init() {
	p.current = 0
	p.ast = make([]Node, 0)
	p.tokenMap = make(map[int][]func(Token) error)
	fmt.Println("Setting grammar map...")
	for key, _ := range MasterTokenMap {
		p.tokenMap[key] = make([]func(Token) error, 0)
	}
	p.HookMonad(Number, p.number)
	p.HookMonad(LeftParen, p.grouping)
	p.HookMonad(LeftBrace, p.grouping)
}

func (p *Parser) Parse(tokens []Token) ([]Node, error) {
	p.init()
	p.src = tokens
	for !p.isAtEnd() {
		p.parse()
		if p.fatal != nil {
			return nil, p.fatal
		}
	}

	p.ast = append(p.ast, Literal{Token{"", EOF, p.src[len(p.src)-1].Line}})
	return p.ast, nil
}

func (p *Parser) parse() {
	token := p.consume()
	// Grab all monads registered with this token
	monads, found := p.tokenMap[token.Type]
	if !found { // FIXME: should check if tokenMap has been populated already
		return
	}
	// Go through and call all the monads
	for _, monad := range monads {
		err := monad(token)
		if err != nil {
			p.errors = append(p.errors, err)
			return
		}
	}
}

func (p *Parser) consume() Token {
	token := p.src[p.current]
	p.current++
	return token
}

func (p *Parser) peek() Token {
	if p.current > uint64(len(p.src)-1) {
		return p.src[len(p.src)-1]
	}
	return p.src[p.current+1]
}

// seek scans until it reaches the end of the src
// or finds a corresponding token's type that matches the t parameter
func (p *Parser) seek(t int) error {
	for !p.isAtEnd() {
		if a := p.consume(); a.is(t) {
			return nil
		}
	}

	return InternalError{21, "parser.seek() reached end of src input"}
}

/********** AST Node parsing implementation **********/
func (p *Parser) number(t Token) error {
	p.ast = append(p.ast, Literal{t})
	return nil
}

func (p *Parser) grouping(t Token) error {
	switch t.Type {
	case LeftParen:
		if err := p.seek(RightParen); err != nil {
			return err
		}
	case LeftBrace:
		if err := p.seek(RightBrace); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) expression(t Token) error {
	return nil
}

/********* Utility functions ************/
func (p *Parser) isAtEnd() bool {
	return p.current >= uint64(len(p.src))
}

func (p *Parser) flush() {
	log.Println("Flushing parser...")
	p.ast = p.ast[:0]
}

func (p *Parser) HookMonad(tokenType int, monad func(Token) error) int {
	/* TODO: config for statically allocating monads vs dynamically at runtime (and benchmarks one day...)
	if p.tokenMap[tokenType] == nil {
		p.tokenMap[tokenType] = make([]func(Token) error, 0)
	}*/
	p.tokenMap[tokenType] = append(p.tokenMap[tokenType], monad)
	fmt.Printf("Hooking monad for '%s' - %d listeners\n", MasterTokenMap[tokenType], len(p.tokenMap[tokenType]))
	return len(p.tokenMap) - 1
}
