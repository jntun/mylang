package lang

import (
	"fmt"
	"reflect"
)

type Environment struct {
	vars  []Block
	funcs []Block
}

type Block struct {
	id    string
	store interface {
		query(id string) (interface{}, bool)
		fmt.Stringer
	}
}

type varMap map[string]*Value
type funcMap map[string]FunctionInvocation

func (vars varMap) query(id string) (interface{}, bool) {
	val, found := vars[id]
	if val == nil {
		return nil, found
	}

	return *val, found
}

func (funs funcMap) query(id string) (interface{}, bool) {
	fun, found := funs[id]
	return fun, found
}

func (env Environment) varResolve(variable Variable) (Value, error) {
	for i := len(env.vars) - 1; i >= 0; i-- {
		val, found := env.vars[i].store.query(variable.identifier.Lexeme)
		if !found {
			continue
		}
		if val == nil {
			return nil, nil
		}
		return val, nil
	}
	return nil, UnknownIdentifier{variable.identifier}
}

func (env Environment) funcResolve(call FunctionCall) (FunctionInvocation, bool) {
	var fun FunctionInvocation
	for i := len(env.funcs) - 1; i >= 0; i-- {
		funIn, found := env.funcs[i].store.query(call.identifier.Lexeme)
		if found {
			fun = funIn.(FunctionInvocation)
			return fun, found
		}
	}

	return fun, false
}

func (env Environment) arrayResolve(arr ArrayAccess, index int) (Value, error) {
	var valArr []*Value
	varBlock, err := env.varResolve(Variable{arr.identifier})
	if err != nil {
		return nil, err
	}
	if kind := reflect.TypeOf(varBlock).Kind(); kind != reflect.Slice {
		return nil, InternalError{50, fmt.Sprintf("wanted array type for access statement got '%s'.", kind)}
	}

	valArr = varBlock.([]*Value)

	if arrLen := len(valArr) - 1; arrLen < index {
		return nil, OutOfBounds{arr.identifier.Lexeme, index, arrLen}
	}
	if valArr[index] != nil {
		return *valArr[index], nil
	} else {
		return nil, nil
	}
}

func (env Environment) classResolve(identifier Token) (JlangClass, error) {
	var class JlangClass

	val, err := env.varResolve(Variable{identifier})
	if err != nil {
		return JlangClass{}, err
	}
	if typeStr := reflect.TypeOf(val).String(); typeStr != "lang.JlangClass" {
		return class, fmt.Errorf("'%s' is not a Class type - is %s", identifier.Lexeme, typeStr)
	}
	class = val.(JlangClass)

	return class, nil
}

func (env Environment) varStore(identifier string, val *Value) {
	env.vars[len(env.vars)-1].store.(varMap)[identifier] = val
}

func (env Environment) funcStore(fun FunctionInvocation) {
	env.funcs[len(env.funcs)-1].store.(funcMap)[fun.stmt.Identifier.Lexeme] = fun
}

func (env Environment) classStore(class JlangClass) {
	x := magic(class)
	env.vars[len(env.vars)-1].store.(varMap)[class.identifier.Lexeme] = &x
}

func (env Environment) arrayStore(identifier string, arr []*Value) {
	/* Currently 'arr' is an array of pointers to Value interfaces. What we need to store into the varMap is a Value pointer */
	/* As far as I'm concerned, magic is the only way that's going to happen so that's exactly what'll be done */
	x := magic(arr)
	env.vars[len(env.vars)-1].store.(varMap)[identifier] = &x
}

func (env *Environment) pop() {
	env.vars = env.vars[:len(env.vars)-1]
	env.funcs = env.funcs[:len(env.funcs)-1]
}

func (env *Environment) push(blockID string) {
	env.vars = append(env.vars, Block{"var-" + blockID, make(varMap)})
	env.funcs = append(env.funcs, Block{"fun-" + blockID, make(funcMap)})
}

func NewEnvironment(id string) Environment {
	env := Environment{make([]Block, 1), make([]Block, 1)}
	env.vars[0] = Block{id, make(varMap)}
	env.funcs[0] = Block{id, make(funcMap)}
	return env
}

// I don't know what the implications of this magic are
func magic(x interface{}) Value {
	return x
}
