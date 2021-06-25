package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// TODO: file-wide change commented out println's to debug output

func (program Program) execute() error {
	for _, stmt := range program.Statements {
		err := stmt.execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func (stmt PrintStatement) execute() error {
	val, err := stmt.Expression.evaluate()
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func (stmt VariableStatement) execute() error {
	reslv := stmt.resolver
	reslv(stmt)
	return nil
}

func (stmt IfStatement) execute() error {
	for _, currStmt := range stmt.block {
		var val Value
		var err error
		if val, err = stmt.Expr.evaluate(); err != nil {
			return err
		}
		if truthy(val) {
			err = currStmt.execute()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (stmt ExpressionStatement) execute() error {
	_, err := stmt.Expression.evaluate()
	if err != nil {
		return err
	}
	//fmt.Println(val)
	return nil
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
	case PlusPlus:
		kind := reflect.TypeOf(expr).Kind()
		if kind == reflect.Float64 {
			return expr.(float64) + 1, nil
		} else if kind == reflect.Int {
			return expr.(int) + 1, nil
		}
	case MinusMinus:
		kind := reflect.TypeOf(expr).Kind()
		if kind == reflect.Float64 {
			return expr.(float64) - 1, nil
		} else if kind == reflect.Int {
			return expr.(int) - 1, nil
		}
	}

	return nil, InvalidOperation{unary.Op}
}

func (binary Binary) evaluate() (Value, error) {
	//fmt.Printf("%v %s %v = ", binary.Left, binary.Op.Lexeme, binary.Right)
	//fmt.Println(reflect.TypeOf(binary.Left))
	left, err := binary.Left.evaluate()
	if err != nil {
		return nil, err
	}
	right, err := binary.Right.evaluate()
	if err != nil {
		return nil, err
	}

	if left == nil || right == nil {
		return nil, NilReference{}
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
	case Greater:
		return binary.Greater(left, right), nil
	case Less:
		return !binary.Greater(left, right), nil
	case GreaterEqual:
		if equal(left, right) {
			return true, nil
		}
		return binary.Greater(left, right), nil
	case LessEqual:
		if equal(left, right) {
			return true, nil
		}
		return !binary.Greater(left, right), nil
	case EqualEqual:
		return binary.Equality(left, right)
	case BangEqual:
		return binary.Inequality(left, right)
	}

	return nil, InvalidOperation{binary.Op}
}

func (binary Binary) Greater(lhs Value, rhs Value) bool {
	left, right := getLeftRightKinds(lhs, rhs)
	switch left {
	case reflect.Int:
		if right == reflect.Int {
			return lhs.(int) > rhs.(int)
		} else if right == reflect.Float64 {
			return float64(lhs.(int)) > rhs.(float64)
		}
	case reflect.Float64:
		if right == left {
			return lhs.(float64) > rhs.(float64)
		} else if right == reflect.Int {
			return lhs.(float64) > float64(rhs.(int))
		}
	}

	return false
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

func (binary Binary) Equality(left Value, right Value) (Value, error) {
	return equal(left, right), nil
}

func (binary Binary) Inequality(left Value, right Value) (Value, error) {
	return !equal(left, right), nil
}

func (literal Literal) evaluate() (Value, error) {
	switch literal.Type {
	case Number:
		if strings.Contains(literal.Lexeme, ".") {
			//fmt.Println("Number is rational:", literal.Lexeme)
			return strconv.ParseFloat(literal.Lexeme, 64)
		}
		//fmt.Println("Number is irrational:", literal.Lexeme)
		return strconv.Atoi(literal.Lexeme)
	case String:
		return literal.Lexeme, nil
	case True:
		return true, nil
	case False:
		return false, nil
	case Nil:
		return nil, nil
	case EOF:
		return "EOF", nil
	}

	return nil, fmt.Errorf("Unable to match literal %s, with a known value.", literal.Token.Lexeme)
}

func (variable Variable) evaluate() (Value, error) {
	return variable.resolver(variable)
}

func getLeftRightKinds(left Value, right Value) (reflect.Kind, reflect.Kind) {
	return reflect.TypeOf(left).Kind(), reflect.TypeOf(right).Kind()
}

func equal(left Value, right Value) bool {
	// This is kind of cool :^) let's see how long it can hold...
	return left == right
}

func truthy(val Value) bool {
	return val == true
}
