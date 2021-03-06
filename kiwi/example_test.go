package kiwi_test

import (
	"fmt"

	"github.com/noypi/math0/expr"
	"github.com/noypi/math0/kiwi"
)

func init() {
	kiwi.SetLogLevel(kiwi.LogDebug)
	expr.SetLogLevel(expr.LogDebug)
}

func aExampleSolverImpl_AddConstraint() {
	expr.EqnBuilder_VarConstructor = func(name string, power float64) expr.IVariable {
		return kiwi.Var(name)
	}

	eqn1 := expr.Eqn(expr.Terms("x"))(expr.EQ)(expr.Terms("5"))
	eqn2 := expr.Eqn(expr.Terms("y"))(expr.EQ)(expr.Terms("10"))

	solver := kiwi.Solver()

	solver.AddConstraint(kiwi.NewConstraint(eqn1, kiwi.Weak()))
	solver.AddConstraint(kiwi.NewConstraint(eqn2, kiwi.Weak()))
	solver.UpdateVariables()

	xval := kiwi.ValueOf(eqn1.Left().WithVars("x").VarAt(0))
	yval := kiwi.ValueOf(eqn2.Left().WithVars("y").VarAt(0))
	fmt.Printf("x=%v, y=%v", xval, yval)
	// Output: x=5, y=10
}
