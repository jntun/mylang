package main

import "fmt"

const (
	LeftParen = iota
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Semicolon
	Slash
	Star
	Mod

	Plus
	PlusPlus
	Minus
	MinusMinus
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	Identifier
	String
	Number

	And
	Class
	Else
	False
	Function
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
)

// Token is a parsed sequence of character terminal(s)
// TODO: maybe store column position as well?
type Token struct {
	Lexeme string
	Type   int
	Line   uint
}

func (t Token) is(ta int) bool {
	return t.Type == ta
}

func (t Token) String() string {
	return fmt.Sprintf("Token<'%s'|%d|%s>", t.Lexeme, t.Line, t.TypeString())
}

var MasterTokenMap = map[int]string{
	LeftParen:    "LeftParen",
	RightParen:   "RightParen",
	LeftBrace:    "LeftBrace",
	RightBrace:   "RightBrace",
	Comma:        "Comma",
	Dot:          "Dot",
	Semicolon:    "Semicolon",
	Slash:        "Slash",
	Star:         "Star",
	Plus:         "Plus",
	PlusPlus:     "PlusPlus",
	Minus:        "Minus",
	MinusMinus:   "MinusMinus",
	Bang:         "Bang",
	BangEqual:    "BangEqual",
	Equal:        "Equal",
	EqualEqual:   "EqualEqual",
	Greater:      "Greater",
	GreaterEqual: "GreaterEqual",
	Less:         "Less",
	LessEqual:    "LessEqual",
	Identifier:   "Identifier",
	String:       "String",
	Number:       "Number",
	And:          "And",
	Class:        "Class",
	Else:         "Else",
	False:        "False",
	Function:     "Function",
	For:          "For",
	If:           "If",
	Nil:          "Nil",
	Or:           "Or",
	Print:        "Print",
	Return:       "Return",
	Super:        "Super",
	This:         "This",
	True:         "True",
	Var:          "Var",
	While:        "While",
	EOF:          "EOF",
}

func (t Token) TypeString() string {
	str, found := MasterTokenMap[t.Type]
	if found {
		return str
	}

	return "no_type_str"
}

// A go fmtd Token array output, useful for making tests
func (t Token) FmtString() string {
	return fmt.Sprintf("Token{\"%s\", %d, %d},", t.Lexeme, t.Type, t.Line)
}
