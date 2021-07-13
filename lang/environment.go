package lang

type Environment struct {
	vars  []VarBlock
	funcs []FuncBlock
}

type varMap map[string]*Value
type funcMap map[string]FunctionInvocation

type VarBlock struct {
	id   string
	vars varMap
}

type FuncBlock struct {
	id    string
	items funcMap
}

func (env Environment) resolve(variable Variable) (Value, error) {
	for i := len(env.vars) - 1; i >= 0; i-- {
		blockEnv := env.vars[i]
		val, found := blockEnv.vars[variable.identifier.Lexeme]
		if !found {
			continue
		}
		if val == nil {
			return nil, nil
		}
		return *val, nil
	}
	return nil, UnknownIdentifier{variable}
}

func (env Environment) funcResolve(call FunctionCall) (FunctionInvocation, bool) {
	fun, found := env.funcs[len(env.funcs)-1].items[call.identifier.Lexeme]
	return fun, found
}

func (env Environment) varStore(identifier string, val *Value) {
	env.vars[len(env.vars)-1].vars[identifier] = val
}

func (env Environment) funcStore(fun FunctionInvocation) {
	env.funcs[len(env.funcs)-1].items[fun.stmt.Identifier.Lexeme] = fun
}

func (env *Environment) pop() {
	env.vars = env.vars[:len(env.vars)-1]
	env.funcs = env.funcs[:len(env.funcs)-1]
}

func (env *Environment) push(blockID string) {
	env.vars = append(env.vars, VarBlock{blockID, make(varMap)})
	env.funcs = append(env.funcs, FuncBlock{blockID, make(funcMap)})
}

func NewEnvironment(id string) Environment {
	env := Environment{make([]VarBlock, 1), make([]FuncBlock, 1)}
	env.vars[0] = VarBlock{id, make(varMap)}
	env.funcs[0] = FuncBlock{id, make(funcMap)}

	return env
}
