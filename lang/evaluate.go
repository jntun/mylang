package lang

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// TODO: file-wide change commented out println's to debug output

func (grouping Grouping) evaluate(intptr *Interpreter) (Value, error) {
	return grouping.Expr.evaluate(intptr)
}

func (unary Unary) evaluate(intptr *Interpreter) (Value, error) {
	if unary.Expr == nil {
		return nil, InvalidOperation{unary.Op}
	}
	expr, err := unary.Expr.evaluate(intptr)
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

func (binary Binary) evaluate(intptr *Interpreter) (Value, error) {
	left, err := binary.Left.evaluate(intptr)
	if err != nil {
		return nil, err
	}
	right, err := binary.Right.evaluate(intptr)
	if err != nil {
		return nil, err
	}

	if left == nil || right == nil {
		switch binary.Op.Type {
		case EqualEqual:
			if left == nil {
				return right == nil, nil
			}
			return false, nil
		case BangEqual:
			if left == nil {
				return !(right == nil), nil
			} else if right == nil {
				return true, nil
			}
			return false, nil
		}
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
		return binary.Less(left, right), nil
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
	case Mod:
		return binary.Modulo(left, right)
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

func (binary Binary) Less(lhs Value, rhs Value) bool {
	left, right := getLeftRightKinds(lhs, rhs)
	switch left {
	case reflect.Int:
		if right == reflect.Int {
			return lhs.(int) < rhs.(int)
		} else if right == reflect.Float64 {
			return float64(lhs.(int)) < rhs.(float64)
		}
	case reflect.Float64:
		if right == left {
			return lhs.(float64) < rhs.(float64)
		} else if right == reflect.Int {
			return lhs.(float64) < float64(rhs.(int))
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
		} else if rKind == reflect.String {
			return fmt.Sprintf("%v%v", left, right), nil
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
	if left == 0 {
		return nil, DivisionByZero{binary.Left}
	} else if right == 0 {
		return nil, DivisionByZero{binary.Right}
	}

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

func (binary Binary) Modulo(left Value, right Value) (Value, error) {
	lKind, rKind := getLeftRightKinds(left, right)
	switch lKind {
	case reflect.Int:
		if rKind == reflect.Int {
			return left.(int) % right.(int), nil
		}
	}

	return nil, InvalidTypeCombination{"modulo", lKind, rKind}
}

func (binary Binary) Equality(left Value, right Value) (Value, error) {
	return equal(left, right), nil
}

func (binary Binary) Inequality(left Value, right Value) (Value, error) {
	return !equal(left, right), nil
}

func (literal Literal) evaluate(intptr *Interpreter) (Value, error) {
	switch literal.Type {
	case Number:
		if strings.Contains(literal.Lexeme, ".") {
			return strconv.ParseFloat(literal.Lexeme, 64)
		}
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

func (variable Variable) evaluate(intptr *Interpreter) (Value, error) {
	return intptr.VariableResolver(variable)
}

func (fun FunctionCall) evaluate(intptr *Interpreter) (Value, error) {
	intptr.env.push(fmt.Sprintf("%s@%d", fun.identifier.Lexeme, fun.identifier.Line))
	return intptr.FunctionResolve(fun)
}

func (fun FunctionInvocation) evaluate(intptr *Interpreter) (Value, error) {
	defer intptr.env.pop()

	// If we have args, map them to the interpreter environment
	if fun.argExprs != nil {
		if len(*fun.argExprs) != len(*fun.stmt.args) {
			return nil, ArgumentMismatch{fun.stmt.Identifier, uint(len(*fun.stmt.args)), uint(len(*fun.argExprs))}
		}
		for i, expr := range *fun.argExprs {
			ids := *fun.stmt.args
			intptr.VariableMap(VariableStatement{ids[i], expr})
		}
	}

	for _, stmt := range fun.stmt.block {
		if err := stmt.execute(intptr); err != nil {
			return nil, err
		}
		if intptr.shouldBreak() {
			break
		}
	}

	if intptr.funcRet != nil {
		val := *intptr.funcRet
		intptr.funcRet = nil
		return val, nil
	}

	return nil, nil
}

func (fun *FunctionInvocation) FillArgs(argExprs *[]Expression) error {
	if argExprs == nil {
		return nil
	}
	if len(*argExprs) == 0 {
		return nil
	}

	fun.argExprs = argExprs
	return nil
}

func (array ArrayAccess) evaluate(intptr *Interpreter) (Value, error) {
	index, err := array.index.evaluate(intptr)
	if err != nil {
		return nil, err
	}
	if kind := reflect.TypeOf(index).Kind(); kind != reflect.Int {
		return nil, fmt.Errorf("invalid type '%s' for index of array '%s'", kind.String(), array.identifier.Lexeme)
	}
	return intptr.env.arrayResolve(array, index.(int))
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
