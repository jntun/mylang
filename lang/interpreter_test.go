package lang

import (
	"reflect"
	"testing"
)

func TestInterpretArithmetic(t *testing.T) {
	if err := genFile("arithmetic"); err != nil {
		t.Error(err)
	}
}

func TestInterpretNumeric(t *testing.T) {
	if err := genFile("numeric"); err != nil {
		t.Error(err)
	}
}

func TestInterpretComment(t *testing.T) {
	if err := genFile("comment"); err != nil {
		t.Error(err)
	}
}

func TestInterpretVar(t *testing.T) {
	if err := genFile("var"); err != nil {
		if reflect.TypeOf(err).Name() != "NilReference" {
			t.Error(err)
		}
	}
}

func TestInterpretArray(t *testing.T) {
	if err := genFile("array"); err != nil {
		t.Error(err)
	}
}

func TestInterpretWhile(t *testing.T) {
	if err := genFile("while"); err != nil {
		t.Error(err)
	}
}

func TestInterpretIf(t *testing.T) {
	if err := genFile("if"); err != nil {
		t.Error(err)
	}
}

func TestInterpretIfElse(t *testing.T) {
	if err := genFile("ifelse"); err != nil {
		t.Error(err)
	}
}

func TestInterpretFunc(t *testing.T) {
	if err := genFile("function"); err != nil {
		if reflect.TypeOf(err).Name() != "BadCall" {
			t.Error(err)
		}
	}
}

func TestInterpretFor(t *testing.T) {
	if err := genFile("for"); err != nil {
		t.Error(err)
	}
}

func TestInterpretClass(t *testing.T) {
	if err := genFile("class"); err != nil {
		t.Error(err)
	}
}

func TestPropertyAssignment(t *testing.T) {
	if err := genFile("propertyassignment"); err != nil {
		t.Error(err)
	}
}

func TestInterpretConstructor(t *testing.T) {
	if err := genFile("constructor"); err != nil {
		t.Error(err)
	}
}

func TestInterpretRectangle(t *testing.T) {
	if err := genFile("rect_class"); err != nil {
		t.Error(err)
	}
}

func genFile(filename string) error {
	intptr := NewInterpreter()
	if err := intptr.File("../tests/" + filename + ".jlang"); err != nil {
		return err
	}
	return nil
}
