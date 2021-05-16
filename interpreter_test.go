package main

import (
	"testing"
)

func TestInterpretNumeric(t *testing.T) {
	if err := genFile("tests/numeric.jlang"); err != nil {
		t.Error(err)
	}
}

func TestInterpretFor(t *testing.T) {
	if err := genFile("tests/for.jlang"); err != nil {
		t.Error(err)
	}
}

func TestInterpretFunc(t *testing.T) {
	if err := genFile("tests/function.jlang"); err != nil {
		t.Error(err)
	}
}

func genFile(filepath string) error {
	intptr := NewInterpreter()
	if err := intptr.File(filepath); err != nil {
		return err
	}
	return nil
}
