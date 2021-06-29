package main

type environment map[string]*Value

func (env environment) resolve(variable Variable) (Value, error) {
	val, found := env[variable.name.Lexeme]
	if !found {
		return nil, UnknownIdentifier{variable}
	}
	if val == nil {
		return nil, nil
	}
	return *val, nil
}

func (env environment) store(identifier string, val *Value) {
	env[identifier] = val
}
