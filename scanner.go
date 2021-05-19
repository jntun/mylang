package main

import "log"

// Scanner
type Scanner struct {
	src     string
	start   int
	current int
	line    int
	tokens  []Token
	Fatal   error
	Errors  []error
}

func (scan *Scanner) Scan(input string) ([]Token, error) {
	scan.tokens = make([]Token, 0)
	scan.src = input
	scan.start = 0
	scan.current = 0
	scan.line = 1

	for !scan.isAtEnd() {
		scan.start = scan.current
		scan.scanToken()

		if scan.Fatal != nil {
			return nil, scan.Fatal
		}
	}

	return scan.tokens, nil
}

func (scan *Scanner) scanToken() {
	if scan.isNumeric() {
		for scan.isNumeric() {
			scan.current++
		}
		scan.addToken(Number)
		return
	}
	val := scan.advance()
	switch val {
	case ' ':
		break
	case '\n':
		scan.line++
	case ';':
		scan.addToken(Semicolon)
	case '(':
		scan.addToken(LeftParen)
	case ')':
		scan.addToken(RightParen)
	case '{':
		scan.addToken(LeftBrace)
	case '}':
		scan.addToken(RightBrace)
	case ',':
		scan.addToken(Comma)
	case '+':
		if scan.match('+') {
			scan.addToken(PlusPlus)
			break
		}
		scan.addToken(Plus)
	case '-':
		if scan.match('-') {
			scan.addToken(MinusMinus)
			break
		}
		scan.addToken(Minus)
	case '/':
		scan.addToken(Slash)
	case '=':
		if scan.match('=') {
			scan.addToken(EqualEqual)
			break
		}
		scan.addToken(Equal)
	case '!':
		if scan.match('=') {
			scan.addToken(BangEqual)
			break
		}
		scan.addToken(Bang)
	case '>':
		if scan.match('=') {
			scan.addToken(GreaterEqual)
			break
		}
		scan.addToken(Greater)
	case '<':
		if scan.match('=') {
			scan.addToken(LessEqual)
			break
		}
		scan.addToken(Less)
	case '"':
		scan.stringParse()
	default:
		scan.multi(val)
	}
}

func (scan *Scanner) multi(val byte) {
	if val == 'i' && scan.match('f') {
		scan.addToken(If)
		return
	}
	if val == 't' && scan.matchStr("rue") {
		scan.addToken(True)
		return
	}

	if val == 'f' {
		if scan.matchStr("or") {
			scan.addToken(For)
			return
		}
		if scan.matchStr("unc") {
			scan.addToken(Function)
			return
		}
	}

	if val == 'p' && scan.matchStr("rint") {
		scan.addToken(Print)
		return
	}
	if val == 'v' && scan.matchStr("ar ") {
		scan.addToken(Var)
		return
	}
	if val == 'r' && scan.matchStr("eturn") {
		scan.addToken(Return)
		return
	}

	scan.identifier()
}

func (scan *Scanner) identifier() {
	for scan.isIdentifier() {
		scan.current++
	}
	scan.addToken(Identifier)
	//scan.Fatal(UnknownToken{string(scan.src[scan.current]), scan.line})
}

func (scan *Scanner) stringParse() {
	err := scan.seek('"')
	if err != nil {
		scan.Fatal = UnclosedString{scan.line}
		return
	}
	scan.addToken(String)
}

func (scan *Scanner) addToken(tokenType int) {
	scan.tokens = append(scan.tokens, Token{scan.src[scan.start:scan.current], tokenType, scan.line})
}

func (scan *Scanner) advance() byte {
	val := scan.src[scan.current]
	scan.current++
	return val
}

func (scan *Scanner) peek(expected byte) bool {
	if scan.isAtEnd() {
		return false
	}

	return scan.src[scan.current] == expected
}

func (scan *Scanner) match(expected byte) bool {
	if scan.peek(expected) {
		scan.current++
		return true
	}
	return false
}

func (scan *Scanner) matchStr(expected string) bool {
	for _, expChar := range expected {
		if !scan.match(byte(expChar)) {
			return false
		}
	}
	return true
}

func (scan *Scanner) seek(expected byte) error {
	for !scan.isAtEnd() {
		if a := scan.advance(); a == expected {
			return nil
		}
	}

	return InternalError{10, "scanner.seek() reached end of file"}
}

func (scan *Scanner) isAtEnd() bool {
	return scan.current >= len(scan.src)
}

func (scan *Scanner) isNumeric() bool {
	val := scan.src[scan.current]
	if val == 46 || val >= 48 && val <= 57 {
		return true
	}
	return false
}

func (scan *Scanner) isIdentifier() bool {
	val := scan.src[scan.current]
	if val >= 48 && val <= 56 {
		return true
	}
	if val >= 65 && val <= 90 {
		return true
	}
	if val >= 97 && val <= 122 {
		return true
	}
	return false
}

func (scan *Scanner) flush() {
	log.Println("Flushing scanner...")
	scan.tokens = scan.tokens[:0]
}
