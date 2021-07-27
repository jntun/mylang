package lang

import (
	"fmt"
	"reflect"
	"strings"
)

type InvalidClassStatement struct {
	class Token
	src   Token
}

func (err InvalidClassStatement) Error() string {
	return fmt.Sprintf("Not a valid class statement at '%s' in class '%s' on %d.", err.src.Lexeme, err.class.Lexeme, err.src.Line)
}

type OutOfBounds struct {
	arrLex string
	index  int
	actual int
}

func (err OutOfBounds) Error() string {
	return fmt.Sprintf("Out of bounds access on array '%s'\n\t'%s' is length %d while the index was %d", err.arrLex, err.arrLex, err.actual, err.index)
}

type DivisionByZero struct {
	offender Expression
}

func (err DivisionByZero) Error() string {
	return fmt.Sprintf("Divide by zero error at '%s'.", err.offender)
}

type ArgumentMismatch struct {
	identifier Token
	expected   uint
	got        uint
}

func (err ArgumentMismatch) Error() string {
	return fmt.Sprintf("argument length mismatch for '%s' call: want %d, got %d.", err.identifier.Lexeme, err.expected, err.got)
}

type BadCall struct {
	id   Token
	more error
}

func (err BadCall) Error() string {
	if err.more != nil {
		return fmt.Sprintf("bad call to '%s'\n\tmore: %s", err.id.Lexeme, err.more)
	}
	return fmt.Sprintf("unable to reference unknown '%s'", err.id.Lexeme)
}

type InvalidOperation struct {
	op Operator
}

func (err InvalidOperation) Error() string {
	return fmt.Sprintf("Invalid type operation '%s' on line %d", err.op.Lexeme, err.op.Line)
}

type NilReference struct {
	reference Token
}

func (err NilReference) Error() string {
	return fmt.Sprintf("Reference to nil value on line %d at '%s'", err.reference.Line, err.reference.Lexeme)
}

type UnknownIdentifier struct {
	Token
}

func (err UnknownIdentifier) Error() string {
	return fmt.Sprintf("Unable to reference unknown variable '%s' on line %d.", err.Lexeme, err.Line)
}

type InvalidTypeCombination struct {
	Operation string
	Left      reflect.Kind
	Rite      reflect.Kind
}

func (err InvalidTypeCombination) Error() string {
	return fmt.Sprintf("Invalid %s between type %s and %s.", err.Operation, err.Left, err.Rite)
}

type ScanError struct {
	err error
}

func (err ScanError) Error() string {
	return fmt.Sprintf("Scan error: %s", err.err)
}

type ParseError struct {
	token Token
	msg   string
}

func (err ParseError) Error() string {
	if err.token.is(EOF) {
		return fmt.Sprintf("%d at end %s", err.token.Line, err.msg)
	} else {
		return fmt.Sprintf("%d at '%s' %s", err.token.Line, err.token.Lexeme, err.msg)
	}
}

// UnclosedStringError is when the scanner is attempting to scan a string lexeme but never reaches a closing (right) closing '"'
type UnclosedString struct {
	line uint
}

func (err UnclosedString) Error() string {
	return fmt.Sprintf("[UnclosedString] expected \" for string on line %d\n", err.line)
}

// UnknownToken is when we encounter a lexeme we don't have a matching token for
type UnknownToken struct {
	lexeme string
	line   int
}

func (err UnknownToken) Error() string {
	return fmt.Sprintf("unknown token %s on line %d", err.lexeme, err.line)
}

// UnknownFile is when trying to verify that a file exists (before opening for r/w),
// if we can't find it we give an UnknownFile error
type UnknownFile struct {
	Filepath string
	Err      error
}

func (err UnknownFile) Error() string {
	return fmt.Sprintf("Unable to find file %s\n\t\"%s\"", err.Filepath, err.Err)
}

// FileReadFailure is when trying to read a file and we fail
type FileReadFailure struct {
	Filepath string
	Err      error
}

func (err FileReadFailure) Error() string {
	return fmt.Sprintf("failure to read file %s: %s", err.Filepath, err.Err)
}

type ErrorList struct {
	errs []error
}

func (list ErrorList) Error() string {
	str := strings.Builder{}
	str.WriteString("Errors: ")
	for _, err := range list.errs {
		str.WriteString(err.Error() + "\n")
	}
	return str.String()
}

// InternalError is a generic error type for any subsystem to return on a failure of some
// generic functionality.
/*
+----------------------------------------------+
| 	    	    Code   Table 		           |
+----------------------------------------------+
|  1  | nil                                    |
| ... |      				  				   |
| 10  | scanner.seek() reached end of file      |
+----------------------------------------------+
*/
type InternalError struct {
	code int
	msg  string
}

func (err InternalError) Error() string {
	return fmt.Sprintf("internal error %d: %s", err.code, err.msg)
}
