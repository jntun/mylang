package lang

import (
	"testing"
)

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
		Token{"func", 29, 1},
		Token{"test", 22, 1},
		Token{"(", 0, 1},
		Token{")", 1, 1},
		Token{"{", 2, 1},
		Token{"return", 35, 2},
		Token{"test", 23, 2},
		Token{";", 6, 2},
		Token{"}", 3, 3},
		Token{"func", 29, 5},
		Token{"negate", 22, 5},
		Token{"(", 0, 5},
		Token{"x", 22, 5},
		Token{")", 1, 5},
		Token{"{", 2, 5},
		Token{"return", 35, 6},
		Token{"-", 12, 6},
		Token{"x", 22, 6},
		Token{";", 6, 6},
		Token{"}", 3, 7},
		Token{"print", 34, 9},
		Token{"test", 22, 9},
		Token{"(", 0, 9},
		Token{")", 1, 9},
		Token{";", 6, 9},
		Token{"print", 34, 10},
		Token{"test2: ", 23, 10},
		Token{"+", 10, 10},
		Token{"test", 22, 10},
		Token{"(", 0, 10},
		Token{")", 1, 10},
		Token{";", 6, 10},
		Token{"print", 34, 11},
		Token{"negate", 22, 11},
		Token{"(", 0, 11},
		Token{"5", 24, 11},
		Token{")", 1, 11},
		Token{";", 6, 11},
		Token{"print", 34, 13},
		Token{"negate", 22, 13},
		Token{"(", 0, 13},
		Token{"5", 24, 13},
		Token{",", 4, 13},
		Token{"3", 24, 13},
		Token{")", 1, 13},
		Token{";", 6, 13},
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
		Token{"for", 30, 1},
		Token{"var", 39, 1},
		Token{"i", 22, 1},
		Token{"=", 16, 1},
		Token{"0", 24, 1},
		Token{";", 6, 1},
		Token{"i", 22, 1},
		Token{"<", 20, 1},
		Token{"5", 24, 1},
		Token{";", 6, 1},
		Token{"i", 22, 1},
		Token{"=", 16, 1},
		Token{"i", 22, 1},
		Token{"+", 10, 1},
		Token{"1", 24, 1},
		Token{"{", 2, 1},
		Token{"print", 34, 2},
		Token{"i", 22, 2},
		Token{";", 6, 2},
		Token{"}", 3, 3},
		Token{"for", 30, 5},
		Token{"i", 22, 5},
		Token{">", 18, 5},
		Token{"0", 24, 5},
		Token{";", 6, 5},
		Token{"i", 22, 5},
		Token{"=", 16, 5},
		Token{"i", 22, 5},
		Token{"-", 12, 5},
		Token{"1", 24, 5},
		Token{"{", 2, 5},
		Token{"print", 34, 6},
		Token{"going back: ", 23, 6},
		Token{"+", 10, 6},
		Token{"i", 22, 6},
		Token{";", 6, 6},
		Token{"}", 3, 7},
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
		scanner.Scan(input)
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
