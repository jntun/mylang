package main

import "fmt"

func (program Program) execute() error {
	for _, stmt := range program.Statements {
		err := stmt.execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func (stmt PrintStatement) execute() error {
	val, err := stmt.Expression.evaluate()
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func (stmt VariableStatement) execute() error {
	stmt.resolver(stmt)
	return nil
}

func (stmt AssignmentStatement) execute() error {
	_, err := Variable{stmt.VariableStatement.Identifier, stmt.resolver}.evaluate()
	if err != nil {
		return err
	}

	return stmt.VariableStatement.execute()
}

func (stmt IfStatement) execute() error {
	for _, currStmt := range stmt.block {
		var val Value
		var err error
		if val, err = stmt.Expr.evaluate(); err != nil {
			return err
		}
		if truthy(val) {
			err = currStmt.execute()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (stmt WhileStatement) execute() error {
	val, err := stmt.test.evaluate()
	if err != nil {
		return err
	}
	for truthy(val) {
		if err != nil {
			return err
		}
		for _, exec := range stmt.block {
			err := exec.execute()
			if err != nil {
				return err
			}
		}
		val, err = stmt.test.evaluate()
	}
	return nil
}

func (stmt ExpressionStatement) execute() error {
	_, err := stmt.Expression.evaluate()
	if err != nil {
		return err
	}
	//fmt.Println(val)
	return nil
}
