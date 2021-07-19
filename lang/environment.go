package lang

import (
	"fmt"
	"reflect"
)

type Environment struct {
	vars   []Block
	funcs  []Block
	arrays []Block
}

type Block struct {
	id    string
	store store
}

type varMap map[string]*Value
type funcMap map[string]FunctionInvocation
type arrayMap map[string][]*Value

type store interface {
	query(id string) (interface{}, bool)
}

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

func (arrays arrayMap) query(id string) (interface{}, bool) {
	arr, found := arrays[id]
	return arr, found
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
	for i := len(env.arrays) - 1; i >= 0; i-- {
		arr, found := env.arrays[i].store.query(variable.identifier.Lexeme)
		if !found {
			continue
		}
		return arr, nil
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
	for i := len(env.arrays) - 1; i >= 0; i-- {
		queryArray, found := env.arrays[i].store.query(arr.identifier.Lexeme)
		if !found {
			continue
		}
		valArr = queryArray.([]*Value)
	}

	if valArr == nil {
		varBlock, err := env.varResolve(Variable{arr.identifier})
		if err != nil {
			return nil, err
		}
		if kind := reflect.TypeOf(varBlock).Kind(); kind != reflect.Slice {
			return nil, InternalError{50, fmt.Sprintf("wanted array type for access statement got '%s'.", kind)}
		}

		valArr = varBlock.([]*Value)
	}

	if arrLen := len(valArr) - 1; arrLen < index {
		return nil, OutOfBounds{arr.identifier.Lexeme, index, arrLen}
	}
	if valArr[index] != nil {
		return *valArr[index], nil
	} else {
		return nil, nil
	}
}

func (env Environment) varStore(identifier string, val *Value) {
	env.vars[len(env.vars)-1].store.(varMap)[identifier] = val
}

func (env Environment) funcStore(fun FunctionInvocation) {
	env.funcs[len(env.funcs)-1].store.(funcMap)[fun.stmt.Identifier.Lexeme] = fun
}

func (env Environment) arrayStore(identifier string, arr []*Value) {
	env.arrays[len(env.arrays)-1].store.(arrayMap)[identifier] = arr
}

func (env *Environment) pop() {
	env.vars = env.vars[:len(env.vars)-1]
	env.funcs = env.funcs[:len(env.funcs)-1]
}

func (env *Environment) push(blockID string) {
	env.vars = append(env.vars, Block{"var-" + blockID, make(varMap)})
	env.funcs = append(env.funcs, Block{"fun-" + blockID, make(funcMap)})
	env.arrays = append(env.arrays, Block{"arr-" + blockID, make(arrayMap)})
}

func NewEnvironment(id string) Environment {
	env := Environment{make([]Block, 1), make([]Block, 1), make([]Block, 1)}
	env.vars[0] = Block{id, make(varMap)}
	env.funcs[0] = Block{id, make(funcMap)}
	env.arrays[0] = Block{id, make(arrayMap)}
	return env
}
