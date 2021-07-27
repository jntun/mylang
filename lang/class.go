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
	instance := JlangClassInstance{scope}

	if class.Stmt.constructor != nil {
		fmt.Println("constructor...")
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
	scope Environment
}

func (this JlangClassInstance) evaluate(intptr *Interpreter) (Value, error) {
	return 0, nil
}

func (this JlangClassInstance) execute(intptr *Interpreter) error {
	return nil
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

func (this JlangClassInstance) storeFunc(fun FunctionDeclarationStatement) {
	funInv := FunctionInvocation{fun, nil}
	this.scope.funcStore(funInv)
}
