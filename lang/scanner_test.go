package lang

import (
	"testing"
)

func TestScanClass(t *testing.T) {
	input, err := openFile("../tests/class.jlang")
	if err != nil {
		return
	}
	scan := Scanner{}
	tokens, err := scan.Scan(*input)
	if err != nil {
		t.Error(err)
	}

	expectedTokens := []Token{
		Token{"class", Class, 1},
		Token{"Test", Identifier, 1},
		Token{"{", LeftBrace, 1},
		Token{"var", Var, 2},
		Token{"empty", Identifier, 2},
		Token{";", Semicolon, 2},
		Token{"var", Var, 3},
		Token{"name", Identifier, 3},
		Token{"=", Equal, 3},
		Token{"test_class", String, 3},
		Token{";", Semicolon, 3},
		Token{"var", Var, 4},
		Token{"id", Identifier, 4},
		Token{"=", Equal, 4},
		Token{"1", Number, 4},
		Token{";", Semicolon, 4},
		Token{"func", Function, 6},
		Token{"Test", Identifier, 6},
		Token{"(", LeftParen, 6},
		Token{")", RightParen, 6},
		Token{"{", LeftBrace, 6},
		Token{"print", Print, 7},
		Token{"init", String, 7},
		Token{";", Semicolon, 7},
		Token{"}", RightBrace, 8},
		Token{"func", Function, 10},
		Token{"getName", Identifier, 10},
		Token{"(", LeftParen, 10},
		Token{")", RightParen, 10},
		Token{"{", LeftBrace, 10},
		Token{"return", Return, 11},
		Token{"hello world!", String, 11},
		Token{";", Semicolon, 11},
		Token{"}", RightBrace, 12},
		Token{"func", Function, 14},
		Token{"nothing", Identifier, 14},
		Token{"(", LeftParen, 14},
		Token{"x", Identifier, 14},
		Token{",", Comma, 14},
		Token{"y", Identifier, 14},
		Token{")", RightParen, 14},
		Token{"{", LeftBrace, 14},
		Token{"}", RightBrace, 14},
		Token{"}", RightBrace, 15},
		Token{"var", Var, 17},
		Token{"test", Identifier, 17},
		Token{"=", Equal, 17},
		Token{"Test", Identifier, 17},
		Token{"(", LeftParen, 17},
		Token{")", RightParen, 17},
		Token{";", Semicolon, 17},
		Token{"print", Print, 18},
		Token{"test", Identifier, 18},
		Token{";", Semicolon, 18},
		Token{"print", Print, 19},
		Token{"test", Identifier, 19},
		Token{".", Dot, 19},
		Token{"id", Identifier, 19},
		Token{";", Semicolon, 19},
		Token{"print", Print, 20},
		Token{"test", Identifier, 20},
		Token{".", Dot, 20},
		Token{"getName", Identifier, 20},
		Token{"(", LeftParen, 20},
		Token{")", RightParen, 20},
		Token{";", Semicolon, 20},
	}

	if matched, got, expect := tokenMatch(t, tokens, expectedTokens); !matched {
		gotExpectError(t, got, expect)
	}
}

func TestScanFunc(t *testing.T) {
	input, err := openFile("../tests/function.jlang")
	if err != nil {
		t.Error(err)
	}
	scan := Scanner{}
	tokens, err := scan.Scan(*input)
	if err != nil {
		t.Error(err)
	}

	expectedTokens := []Token{
		Token{"func", Function, 1},
		Token{"test", Identifier, 1},
		Token{"(", LeftParen, 1},
		Token{")", RightParen, 1},
		Token{"{", LeftBrace, 1},
		Token{"return", Return, 2},
		Token{"test", String, 2},
		Token{";", Semicolon, 2},
		Token{"}", RightBrace, 3},
		Token{"func", Function, 5},
		Token{"negate", Identifier, 5},
		Token{"(", LeftParen, 5},
		Token{"x", Identifier, 5},
		Token{")", RightParen, 5},
		Token{"{", LeftBrace, 5},
		Token{"return", Return, 6},
		Token{"-", Minus, 6},
		Token{"x", Identifier, 6},
		Token{";", Semicolon, 6},
		Token{"}", RightBrace, 7},
		Token{"print", Print, 9},
		Token{"test", Identifier, 9},
		Token{"(", LeftParen, 9},
		Token{")", RightParen, 9},
		Token{";", Semicolon, 9},
		Token{"print", Print, 10},
		Token{"test2: ", String, 10},
		Token{"+", Plus, 10},
		Token{"test", Identifier, 10},
		Token{"(", LeftParen, 10},
		Token{")", RightParen, 10},
		Token{";", Semicolon, 10},
		Token{"print", Print, 11},
		Token{"negate", Identifier, 11},
		Token{"(", LeftParen, 11},
		Token{"5", Number, 11},
		Token{")", RightParen, 11},
		Token{";", Semicolon, 11},
		Token{"print", Print, 13},
		Token{"negate", Identifier, 13},
		Token{"(", LeftParen, 13},
		Token{"5", Number, 13},
		Token{",", Comma, 13},
		Token{"3", Number, 13},
		Token{")", RightParen, 13},
		Token{";", Semicolon, 13},
	}

	if matched, got, expect := tokenMatch(t, tokens, expectedTokens); !matched {
		gotExpectError(t, got, expect)
	}
}

func TestScanDoubleOperator(t *testing.T) {
	input, err := openFile("../tests/doubleop.jlang")
	if err != nil {
		t.Error(err)
	}
	scan := Scanner{}
	tokens, err := scan.Scan(*input)
	if err != nil {
		t.Error(err)
	}

	expectedTokens := []Token{
		{"var", Var, 1},
		{"val", Identifier, 1},
		{"=", Equal, 1},
		{"1", Number, 1},
		{";", Semicolon, 1},
		{"val", Identifier, 2},
		{">", Greater, 2},
		{"1", Number, 2},
		{";", Semicolon, 2},
		{"val", Identifier, 3},
		{">=", GreaterEqual, 3},
		{"1", Number, 3},
		{";", Semicolon, 3},
		{"val", Identifier, 4},
		{"!=", BangEqual, 4},
		{"1", Number, 4},
		{";", Semicolon, 4},
		{"val", Identifier, 5},
		{">", Greater, 5},
		{"1", Number, 5},
		{";", Semicolon, 5},
		{"val", Identifier, 6},
		{">=", GreaterEqual, 6},
		{"1", Number, 6},
		{";", Semicolon, 6},
		{"val", Identifier, 7},
		{"==", EqualEqual, 7},
		{"1", Number, 7},
		{";", Semicolon, 7},
	}

	if matched, got, expect := tokenMatch(t, tokens, expectedTokens); !matched {
		gotExpectError(t, got, expect)
	}

	//t.Log(tokens)
}

func TestScanFor(t *testing.T) {
	input, err := openFile("../tests/for.jlang")
	if err != nil {
		t.Error(err)
	}
	scan := Scanner{}
	tokens, err := scan.Scan(*input)
	if err != nil {
		t.Error(err)
	}

	expectedTokens := []Token{
		Token{"for", For, 1},
		Token{"var", Var, 1},
		Token{"i", Identifier, 1},
		Token{"=", Equal, 1},
		Token{"0", Number, 1},
		Token{";", Semicolon, 1},
		Token{"i", Identifier, 1},
		Token{"<", Less, 1},
		Token{"5", Number, 1},
		Token{";", Semicolon, 1},
		Token{"i", Identifier, 1},
		Token{"=", Equal, 1},
		Token{"i", Identifier, 1},
		Token{"+", Plus, 1},
		Token{"1", Number, 1},
		Token{"{", LeftBrace, 1},
		Token{"print", Print, 2},
		Token{"i", Identifier, 2},
		Token{";", Semicolon, 2},
		Token{"}", RightBrace, 3},
		Token{"for", For, 5},
		Token{"var", Var, 5},
		Token{"i", Identifier, 5},
		Token{"=", Equal, 5},
		Token{"5", Number, 5},
		Token{";", Semicolon, 5},
		Token{"i", Identifier, 5},
		Token{">", Greater, 5},
		Token{"0", Number, 5},
		Token{";", Semicolon, 5},
		Token{"i", Identifier, 5},
		Token{"=", Equal, 5},
		Token{"i", Identifier, 5},
		Token{"-", Minus, 5},
		Token{"1", Number, 5},
		Token{"{", LeftBrace, 5},
		Token{"print", Print, 6},
		Token{"going back: ", String, 6},
		Token{"+", Plus, 6},
		Token{"i", Identifier, 6},
		Token{";", Semicolon, 6},
		Token{"}", RightBrace, 7},
	}

	if matched, got, expect := tokenMatch(t, tokens, expectedTokens); !matched {
		gotExpectError(t, got, expect)
	}
}

func TestScanArithmetic(t *testing.T) {
	input, err := openFile("../tests/arithmetic.jlang")
	if err != nil {
		t.Error(err)
	}
	scan := Scanner{}
	tokens, err := scan.Scan(*input)
	if err != nil {
		t.Error(err)
	}

	expectedTokens := []Token{
		{"5", Number, 1},
		{"+", Plus, 1},
		{"4", Number, 1},
		{"-", Minus, 1},
		{"3", Number, 1},
		{"==", EqualEqual, 1},
		{"6", Number, 1},
		{";", Semicolon, 1},
	}

	if matched, got, expect := tokenMatch(t, tokens, expectedTokens); !matched {
		gotExpectError(t, got, expect)
	}
}

func TestPeek(t *testing.T) {
	input := "xa"
	expected := 'a'
	scan := Scanner{}
	scan.src = input

	// pretending we are in the big switch statement in scanner.go (i.e: switch scan.scanToken())
	scan.advance()
	// now this is where we would be invoking a peek()
	if !(scan.peek(byte(expected))) {
		t.Error("Did not peek expected value", expected)
	}
}

func TestMatchStr(t *testing.T) {
	input := "{ true print var }"
	scan := Scanner{}
	tokens, err := scan.Scan(input)
	if err != nil {
		t.Error(err)
	}

	expectedTokens := []Token{
		{"{", LeftBrace, 1},
		{"true", True, 1},
		{"print", Print, 1},
		{"var", Var, 1},
		{"}", RightBrace, 1},
	}

	if matched, got, expect := tokenMatch(t, tokens, expectedTokens); !matched {
		gotExpectError(t, got, expect)
	}

	t.Log(tokens)
}

func BenchmarkScanner(b *testing.B) {
	input := "" +
		"for(var i=0; i < 5; i++) {\n" +
		"print(i);\n" +
		"}"
	scanner := Scanner{}
	for i := 0; i < b.N; i++ {
		_, err := scanner.Scan(input)
		if err != nil {
			b.Error(err)
		}
	}
}

func tokenMatch(t *testing.T, tokens []Token, expected []Token) (bool, *Token, *Token) {
	if len(tokens) != len(expected) {
		t.Errorf("tokens are not the same length as expected: %d - %d", len(tokens), len(expected))
		return false, nil, nil
	}

	for i, token := range expected {
		//t.Logf("matching: | expected: %s | got: %s |", token, tokens[i])
		if token == tokens[i] {
			//t.Log(" - match\n")
			continue
		}
		t.Log(" - no match\n")
		return false, &token, &tokens[i]
	}

	return true, nil, nil
}

func gotExpectError(t *testing.T, got *Token, expect *Token) {
	if got != nil && expect != nil {
		t.Errorf("tokens did not match: %s - %s", *got, *expect)
		return
	}
	t.Error("len(tokens) != len(expected) - does not match\n")
}
