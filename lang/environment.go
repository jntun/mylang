package lang

type Environment struct {
	vars   []Block
	funcs  []Block
	arrays []Block
}

type store interface {
	query(id string) (interface{}, bool)
}

type Block struct {
	id    string
	store store
}

type varMap map[string]*Value

func (vars varMap) query(id string) (interface{}, bool) {
	val, found := vars[id]
	if val == nil {
		return nil, found
	}

	return *val, found
}

type funcMap map[string]FunctionInvocation

func (funs funcMap) query(id string) (interface{}, bool) {
	fun, found := funs[id]
	return fun, found
}

type arrayMap map[string][]*Value

func (arrays arrayMap) query(id string) (interface{}, bool) {
	arr, found := arrays[id]
	return arr, found
}

func (env Environment) resolve(variable Variable) (Value, error) {
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
	return nil, UnknownIdentifier{variable}
}

func (env Environment) funcResolve(call FunctionCall) (FunctionInvocation, bool) {
	fun, found := env.funcs[len(env.funcs)-1].store.query(call.identifier.Lexeme)
	return fun.(FunctionInvocation), found
}

func (env Environment) varStore(identifier string, val *Value) {
	env.vars[len(env.vars)-1].store.(varMap)[identifier] = val
}

func (env Environment) funcStore(fun FunctionInvocation) {
	env.funcs[len(env.funcs)-1].store.(funcMap)[fun.stmt.Identifier.Lexeme] = fun
}

func (env *Environment) pop() {
	env.vars = env.vars[:len(env.vars)-1]
	env.funcs = env.funcs[:len(env.funcs)-1]
}

func (env *Environment) push(blockID string) {
	env.vars = append(env.vars, Block{blockID, make(varMap)})
	env.funcs = append(env.funcs, Block{blockID, make(funcMap)})
}

func NewEnvironment(id string) Environment {
	env := Environment{make([]Block, 1), make([]Block, 1), make([]Block, 1)}
	env.vars[0] = Block{id, make(varMap)}
	env.funcs[0] = Block{id, make(funcMap)}
	env.arrays[0] = Block{id, make(arrayMap)}
	return env
}
