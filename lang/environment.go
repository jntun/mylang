package lang

type Environment struct {
	blocks []Block
}

type varMap map[string]*Value
type Block struct {
	id   string
	vars varMap
}

func (env Environment) resolve(variable Variable) (Value, error) {
	for i := len(env.blocks) - 1; i >= 0; i-- {
		blockEnv := env.blocks[i]
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

func (env Environment) store(identifier string, val *Value) {
	env.blocks[len(env.blocks)-1].vars[identifier] = val
}

func (env *Environment) pop() {
	env.blocks = env.blocks[:len(env.blocks)-1]
}

func (env *Environment) push(id string) {
	env.blocks = append(env.blocks, Block{id, make(varMap)})
}

func NewEnvironment(id string) Environment {
	env := Environment{make([]Block, 1)}
	env.blocks[0] = Block{id, make(varMap)}
	return env
}
