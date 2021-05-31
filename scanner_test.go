package main

import (
	"testing"
)

func TestScanFunc(t *testing.T) {
	input, err := openFile("tests/function.jlang")
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
	}

	if matched, got, expect := tokenMatch(t, tokens, expectedTokens); !matched {
		gotExpectError(t, got, expect)
	}
}

func TestScanDoubleOperator(t *testing.T) {
	input, err := openFile("tests/doubleop.jlang")
	if err != nil {
		t.Error(err)
	}
	scan := Scanner{}
	tokens, err := scan.Scan(*input)
	if err != nil {
		t.Error(err)
	}

	expectedTokens := []Token{
		Token{"var", Var, 1},
		Token{"val", Identifier, 1},
		Token{"=", Equal, 1},
		Token{"1", Number, 1},
		Token{";", Semicolon, 1},
		Token{"val", Identifier, 2},
		Token{">", Greater, 2},
		Token{"1", Number, 2},
		Token{";", Semicolon, 2},
		Token{"val", Identifier, 3},
		Token{">=", GreaterEqual, 3},
		Token{"1", Number, 3},
		Token{";", Semicolon, 3},
		Token{"val", Identifier, 4},
		Token{"!=", BangEqual, 4},
		Token{"1", Number, 4},
		Token{";", Semicolon, 4},
		Token{"val", Identifier, 5},
		Token{">", Greater, 5},
		Token{"1", Number, 5},
		Token{";", Semicolon, 5},
		Token{"val", Identifier, 6},
		Token{">=", GreaterEqual, 6},
		Token{"1", Number, 6},
		Token{";", Semicolon, 6},
		Token{"val", Identifier, 7},
		Token{"==", EqualEqual, 7},
		Token{"1", Number, 7},
		Token{";", Semicolon, 7},
	}

	if matched, got, expect := tokenMatch(t, tokens, expectedTokens); !matched {
		gotExpectError(t, got, expect)
	}

	//t.Log(tokens)
}

func TestScanFor(t *testing.T) {
	input, err := openFile("tests/for.jlang")
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
		Token{"(", LeftParen, 1},
		Token{"var", Var, 1},
		Token{"i", Identifier, 1},
		Token{"=", Equal, 1},
		Token{"0", Number, 1},
		Token{";", Semicolon, 1},
		Token{"i", Identifier, 1},
		Token{"<", Less, 1},
		Token{"10", Number, 1},
		Token{";", Semicolon, 1},
		Token{"i", Identifier, 1},
		Token{"++", PlusPlus, 1},
		Token{")", RightParen, 1},
		Token{"{", LeftBrace, 1},
		Token{"print", Print, 2},
		Token{"(", LeftParen, 2},
		Token{"test", String, 2},
		Token{",", Comma, 2},
		Token{"i", Identifier, 2},
		Token{")", RightParen, 2},
		Token{";", Semicolon, 2},
		Token{"}", RightBrace, 3},
	}

	if matched, got, expect := tokenMatch(t, tokens, expectedTokens); !matched {
		gotExpectError(t, got, expect)
	}
}

func TestScanArithmetic(t *testing.T) {
	input, err := openFile("tests/arithmetic.jlang")
	if err != nil {
		t.Error(err)
	}
	scan := Scanner{}
	tokens, err := scan.Scan(*input)
	if err != nil {
		t.Error(err)
	}

	expectedTokens := []Token{
		Token{"2", 23, 1},
		Token{"*", 8, 1},
		Token{"(", 0, 1},
		Token{"3", 23, 1},
		Token{"*", 8, 1},
		Token{"10", 23, 1},
		Token{")", 1, 1},
		Token{"/", 7, 1},
		Token{"3", 23, 1},
		Token{";", 6, 1},
		Token{"2", 23, 2},
		Token{"*", 8, 2},
		Token{"(", 0, 2},
		Token{"3", 23, 2},
		Token{"*", 8, 2},
		Token{"10", 23, 2},
		Token{")", 1, 2},
		Token{"/", 7, 2},
		Token{"3", 23, 2},
		Token{"==", 16, 2},
		Token{"20", 23, 2},
		Token{";", 6, 2},
		Token{"3.14", 23, 3},
		Token{"*", 8, 3},
		Token{"3", 23, 3},
		Token{";", 6, 3},
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
		Token{"{", LeftBrace, 1},
		Token{"true", True, 1},
		Token{"print", Print, 1},
		Token{"var", Var, 1},
		Token{"}", RightBrace, 1},
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
