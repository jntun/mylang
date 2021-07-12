package lang

import (
	"unicode"
)

// Scanner is how a Jlang input string gets scanned and tokenized
type Scanner struct {
	src     string
	start   uint
	current uint
	line    uint
	tokens  []Token
	Fatal   error
	Errors  []error
}

// Scan takes an input string and either returns a Tokenized array or an error specifying why it's
// an invalid input sequence to be scanned.
func (scan *Scanner) Scan(input string) ([]Token, error) {
	scan.src = input
	scan.flush()

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
		for scan.isNumeric() || scan.src[scan.current] == '.' {
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
	case '\t':
		break
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
	case '.':
		scan.addToken(Dot)
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
	case '*':
		scan.addToken(Star)
	case '/':
		if scan.match('/') {
			scan.comment()
		} else {
			scan.addToken(Slash)
		}
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
	case '%':
		scan.addToken(Mod)
	default:
		scan.multi(val)
	}
}

func (scan *Scanner) multi(val byte) {
	if val == 'i' && scan.match('f') {
		scan.addToken(If)
		return
	}
	if val == 'e' && scan.matchStr("lse") {
		scan.addToken(Else)
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
		if scan.matchStr("alse") {
			scan.addToken(False)
			return
		}
	}

	if val == 'p' && scan.matchStr("rint") {
		scan.addToken(Print)
		return
	}
	if val == 'v' && scan.matchStr("ar") {
		scan.addToken(Var)
		return
	}
	if val == 'r' && scan.matchStr("eturn") {
		scan.addToken(Return)
		return
	}
	if val == 'w' && scan.matchStr("hile") {
		scan.addToken(While)
		return
	}

	if val == 'n' && scan.matchStr("il") {
		scan.addToken(Nil)
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

	// Temporarily moves the addToken() consume to *inside* the quotation marks "X____________Y" X=start Y=current
	scan.start++
	scan.current--
	scan.addToken(String)
	scan.current++
}

func (scan *Scanner) comment() {
	for true {
		val := scan.advance()
		if val == '\n' {
			break
		}
	}
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
	return scan.current >= uint(len(scan.src))
}

func (scan *Scanner) isNumeric() bool {
	val := scan.src[scan.current]
	return unicode.IsDigit(rune(val))
}

func (scan *Scanner) isIdentifier() bool {
	val := scan.src[int(scan.current)]
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
	scan.Errors = make([]error, 0)
	scan.tokens = make([]Token, 0)
	scan.start = 0
	scan.current = 0
	scan.line = 1
	scan.Fatal = nil
}
