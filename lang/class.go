package lang

import "fmt"

type JlangClass struct {
	identifier Token
	Stmt       struct {
		constructor *FunctionDeclarationStatement
		varDecls    *[]VariableStatement
		funcDecls   *[]FunctionDeclarationStatement
	}
}

// evaluate is when a class is being _called_ i.e 'MyClass()'.
// This is for creating a new JlangClassInstance using a JlangClass
func (class JlangClass) evaluate(intptr *Interpreter) (Value, error) {
	scope := NewEnvironment(class.identifier.Lexeme)
	instance := JlangClassInstance{&class, scope}

	if class.Stmt.constructor != nil {
		// TODO: implement constructor execution
	}
	if varDecls := class.Stmt.varDecls; varDecls != nil {
		for _, varDecl := range *varDecls {
			err := instance.updateMember(intptr, varDecl)
			if err != nil {
				return nil, err
			}
		}
	}
	if funcDecls := class.Stmt.funcDecls; funcDecls != nil {
		for _, funcDecl := range *funcDecls {
			instance.storeFunc(funcDecl)
		}
	}

	return instance, nil
}

// execute is when a JlangClass is being declared and has already been parsed.
// This is when the class 'type' gets put into the interpreter's
// environment to reference in the future i.e 'evaluate()'.
func (class JlangClass) execute(intptr *Interpreter) error {
	intptr.env.classStore(class)
	return nil
}

type JlangClassInstance struct {
	parent *JlangClass
	scope  Environment
}

func (this JlangClassInstance) invoke(intptr *Interpreter, call FunctionCall) (Value, error) {
	intptr.env.push(fmt.Sprintf("%s#%s", this.parent.identifier.Lexeme, call.identifier.Lexeme))
	defer intptr.env.pop()

	if funk, found := this.scope.funcResolve(call); found {
		args := append([]Expression{this}, *call.args...)
		funk.FillArgs(&args)
		return funk.evaluate(intptr)
	}

	reason := fmt.Errorf("unresolved method '%s' for class of type '%s'.", call.identifier.Lexeme, this.parent.identifier.Lexeme)
	if call.identifier.Lexeme == this.parent.identifier.Lexeme {
		reason = fmt.Errorf("cannot call constructor of '%s' directly.", call.identifier.Lexeme)
	}

	return nil, BadMethodInvocation{call.identifier, reason}
}

func (this JlangClassInstance) evaluate(intptr *Interpreter) (Value, error) {
	return this, nil
}

func (this JlangClassInstance) propertyAccess(identifier Token) (Value, error) {
	var val Value
	var err error

	if val, err = this.scope.varResolve(Variable{identifier}); err != nil {
		return nil, err
	}

	return val, nil
}

func (this JlangClassInstance) updateMember(intptr *Interpreter, vari VariableStatement) error {
	var val Value
	var err error

	if vari.Expr != nil {
		val, err = vari.Expr.evaluate(intptr)
		if err != nil {
			return err
		}
	} else {
		val = nil
	}
	this.scope.varStore(vari.Identifier.Lexeme, &val)
	return nil
}

func (this *JlangClassInstance) storeFunc(fun FunctionDeclarationStatement) {
	funInv := FunctionInvocation{fun, nil}
	this.scope.funcStore(funInv)
}
