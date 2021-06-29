package main

type Environment []map[string]*Value

func (env Environment) resolve(variable Variable) (Value, error) {
	for _, blockEnv := range env {
		val, found := blockEnv[variable.name.Lexeme]
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
	env[len(env)-1][identifier] = val
}

func NewEnvironment() Environment {
	env := make(Environment, 1)
	env[0] = newMap()
	return env
}

func newMap() map[string]*Value {
	return make(map[string]*Value)
}
