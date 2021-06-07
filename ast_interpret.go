package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func (program Program) evaluate() (Value, error) {
	return program.Expr.evaluate()
}

func (grouping Grouping) evaluate() (Value, error) {
	return grouping.Expr.evaluate()
}

func (unary Unary) evaluate() (Value, error) {
	expr, err := unary.Expr.evaluate()
	if err != nil {
		return nil, err
	}
	switch unary.Op.Type {
	case Minus:
		switch reflect.TypeOf(expr).Kind() {
		case reflect.Int:
			return -expr.(int), nil
		case reflect.Float64:
			return -expr.(float64), nil
		}
	case Bang:
		switch reflect.TypeOf(expr).Kind() {
		case reflect.Bool:
			return !expr.(bool), nil
		}
	}

	return InternalError{21, fmt.Sprintf("Unable to evaluate unary expression: %s", unary)}, nil
}

func (binary Binary) evaluate() (Value, error) {
	fmt.Printf("%v %s %v = ", binary.Left, binary.Op.Lexeme, binary.Right)
	fmt.Println(reflect.TypeOf(binary.Left))
	left, err := binary.Left.evaluate()
	if err != nil {
		return nil, err
	}
	right, err := binary.Right.evaluate()
	if err != nil {
		return nil, err
	}

	switch binary.Op.Type {
	case Plus:
		return binary.plus(left, right)
	case Minus:
		return binary.minus(left, right)
	case Star:
		return binary.multiply(left, right)
	case Slash:
		return binary.divide(left, right)
	case EqualEqual:
		return binary.equality(left, right)
	case BangEqual:
		return !(left == right), nil
	}

	return nil, InternalError{20, fmt.Sprintf("Unable to determine operator for binary expr: %s", binary)}
}

func (binary Binary) plus(left Value, right Value) (Value, error) {
	lKind, rKind := getLeftRightKinds(left, right)
	switch lKind {
	case reflect.Int:
		if lKind == rKind {
			return left.(int) + right.(int), nil
		} else if rKind == reflect.Float64 {
			return float64(left.(int)) + right.(float64), nil
		}
	case reflect.String:
		return fmt.Sprint(left, right), nil
	case reflect.Bool:
		if lKind == rKind {
			if left.(bool) && right.(bool) {
				return true, nil
			}
			return false, nil
		}
	case reflect.Float64:
		if lKind == rKind {
			return left.(float64) + right.(float64), nil
		} else if rKind == reflect.Int {
			return left.(float64) + float64(right.(int)), nil
		}
	}

	return nil, InvalidTypeCombination{"addition", lKind, rKind}
}

func (binary Binary) minus(left Value, right Value) (Value, error) {
	lKind, rKind := getLeftRightKinds(left, right)
	switch lKind {
	case reflect.Int:
		if lKind == rKind {
			return left.(int) - right.(int), nil
		} else if rKind == reflect.Float64 {
			return float64(left.(int)) - right.(float64), nil
		}
	case reflect.Float64:
		if lKind == rKind {
			return left.(float64) - right.(float64), nil
		} else if rKind == reflect.Int {
			return left.(float64) - float64(right.(int)), nil
		}
	}
	return nil, InvalidTypeCombination{"subtraction", lKind, rKind}
}

func (binary Binary) multiply(left Value, right Value) (Value, error) {
	lKind, rKind := getLeftRightKinds(left, right)
	switch lKind {
	case reflect.Int:
		if rKind == lKind {
			return left.(int) * right.(int), nil
		} else if rKind == reflect.Float64 {
			return float64(left.(int)) * right.(float64), nil
		}
	case reflect.Float64:
		if rKind == lKind {
			return left.(float64) * right.(float64), nil
		} else if rKind == reflect.Int {
			return left.(float64) * float64(right.(int)), nil
		}
	}
	return nil, InvalidTypeCombination{"multiplication", lKind, rKind}
}

func (binary Binary) divide(left Value, right Value) (Value, error) {
	lKind, rKind := getLeftRightKinds(left, right)
	switch lKind {
	case reflect.Int:
		if rKind == reflect.Int {
			return left.(int) / right.(int), nil
		} else if rKind == reflect.Float64 {
			return float64(left.(int)) / right.(float64), nil
		}
	case reflect.Float64:
		if rKind == lKind {
			return left.(float64) / right.(float64), nil
		} else if rKind == reflect.Int {
			return left.(float64) / float64(right.(int)), nil
		}
	}

	return nil, InvalidTypeCombination{"division", lKind, rKind}
}

func (binary Binary) equality(left Value, right Value) (Value, error) {
	return equal(left, right), nil
}

func (binary Binary) inequality(left Value, right Value) (Value, error) {
	return !equal(left, right), nil
}

func (literal Literal) evaluate() (Value, error) {
	switch literal.Type {
	case Number:
		if strings.Contains(literal.Lexeme, ".") {
			fmt.Println("Number is rational:", literal.Lexeme)
			return strconv.ParseFloat(literal.Lexeme, 64)
		}
		fmt.Println("Number is irrational:", literal.Lexeme)
		return strconv.Atoi(literal.Lexeme)
	case String:
		return literal.Lexeme, nil
	case True:
		return true, nil
	case False:
		return false, nil
	case EOF:
		return "EOF", nil
	}

	return nil, fmt.Errorf("Unable to match literal %s, with a known value", literal.Token.Lexeme)
}

func getLeftRightKinds(left Value, right Value) (reflect.Kind, reflect.Kind) {
	return reflect.TypeOf(left).Kind(), reflect.TypeOf(right).Kind()
}

func equal(left Value, right Value) bool {
	// This is kind of cool :^) let's see how long it can hold...
	return left == right
}
