package main

import "fmt"

func (program Program) execute(intptr *Interpreter) error {
	for _, stmt := range program.Statements {
		err := stmt.execute(intptr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (stmt PrintStatement) execute(intptr *Interpreter) error {
	val, err := stmt.Expression.evaluate(intptr)
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func (stmt VariableStatement) execute(intptr *Interpreter) error {
	//stmt.resolver(stmt)
	intptr.VariableMap(stmt)
	return nil
}

func (stmt AssignmentStatement) execute(intptr *Interpreter) error {
	_, err := Variable{stmt.VariableStatement.Identifier}.evaluate(intptr)
	if err != nil {
		return err
	}

	return stmt.VariableStatement.execute(intptr)
}

func (stmt IfStatement) execute(intptr *Interpreter) error {
	var exec []Statement

	val, err := stmt.Expr.evaluate(intptr)
	if err != nil {
		return err
	}

	if ok := truthy(val); ok == true {
		exec = stmt.block

	} else if stmt.elseBlock != nil {
		exec = *stmt.elseBlock
	}

	for _, stmt := range exec {
		err = stmt.execute(intptr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (stmt FunctionDeclarationStatement) execute(intptr *Interpreter) error {
	intptr.FunctionMap(stmt)
	return nil
}

func (stmt ReturnStatement) execute(intptr *Interpreter) error {
	val, err := stmt.Expression.evaluate(intptr)
	if err != nil {
		return err
	}
	intptr.FunctionReturn(val)
	return nil
}

func (stmt WhileStatement) execute(intptr *Interpreter) error {
	val, err := stmt.test.evaluate(intptr)
	if err != nil {
		return err
	}
	for truthy(val) {
		if err != nil {
			return err
		}
		for _, exec := range stmt.block {
			err := exec.execute(intptr)
			if err != nil {
				return err
			}
		}
		val, err = stmt.test.evaluate(intptr)
	}
	return nil
}

func (stmt ForStatement) execute(intptr *Interpreter) error {
	if stmt.varStmt != nil {
		err := stmt.varStmt.execute(intptr)
		if err != nil {
			return err
		}
	}
	val, err := stmt.test.evaluate(intptr)
	if err != nil {
		return err
	}
	for truthy(val) {
		if err != nil {
			return err
		}
		for _, exec := range stmt.block {
			if err := exec.execute(intptr); err != nil {
				return err
			}
		}
		val, err = stmt.test.evaluate(intptr)
		if assignErr := stmt.assign.execute(intptr); assignErr != nil {
			return assignErr
		}
	}

	return nil
}

func (stmt ExpressionStatement) execute(intptr *Interpreter) error {
	_, err := stmt.Expression.evaluate(intptr)
	if err != nil {
		return err
	}
	//fmt.Println(val)
	return nil
}
