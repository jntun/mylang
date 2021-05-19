package main

import "fmt"

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
	line int
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

// UnknownFile is when trying to veryify that a file exists (before opening for r/w),
// if we can't find it we give an UnkownFile error
type UnknownFile struct {
	Filepath string
	Err      error
}

func (err UnknownFile) Error() string {
	return fmt.Sprintf("unable to find file %s\n\t\"%s\"", err.Filepath, err.Err)
}

// FileReadFailure is when trying to read a file and we fail
type FileReadFailure struct {
	Filepath string
	Err      error
}

func (err FileReadFailure) Error() string {
	return fmt.Sprintf("failure to read file %s: %s", err.Filepath, err.Err)
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