package kiwi

import (
	"testing"

	"github.com/noypi/math0/expr"
	assertpkg "github.com/stretchr/testify/assert"
)

func init() {
	SetLogLevel(LogDebug)
}

func aTestSimple1(t *testing.T) {
	API(">>>TestSimple1")
	defer API("<<<TestSimple1")
	assert := assertpkg.New(t)

	x := expr.NewTerm(167.0, Variable("x"))
	y := expr.NewTerm(2.0, Variable("y"))

	solver := Solver()

	eqn := expr.Equation(expr.NewExpr(y), expr.OpEQ, expr.NewExpr(x))
	solver.AddConstraint(Constraint(eqn, _Required))

	solver.UpdateVariables()
	DBG("x=%v", x)
	DBG("y=%v", y)
	assert.Equal(0.0, ValueOf(x.VarAt(0)))
	assert.Equal(ValueOf(x.VarAt(0)), ValueOf(y.VarAt(0)))
}

func aTestJustStay1(t *testing.T) {
	assert := assertpkg.New(t)

	x := expr.NewTerm(1, Variable("x"))
	y := expr.NewTerm(1, Variable("y"))
	solver := Solver()

	solver.AddConstraint(Constraint(expr.Equation(expr.NewExpr(x, expr.NewTerm(-5.0)), expr.OpEQ, nil), Weak()))
	solver.AddConstraint(Constraint(expr.Equation(expr.NewExpr(y, expr.NewTerm(-10.0)), expr.OpEQ, nil), Weak()))

	solver.UpdateVariables()
	DBG("x=%v", x)
	DBG("y=%v", y)
	assert.Equal(5.0, ValueOf(x.VarAt(0)))
	assert.Equal(10.0, ValueOf(y.VarAt(0)))

}

func TestAddDelete1(t *testing.T) {
	assert := assertpkg.New(t)

	expr.EqnBuilder_VarConstructor = func(name string, power float64) expr.IVariable {
		return Variable(name)
	}

	solver := Solver()

	solver.AddConstraint(Constraint(expr.Eqn(expr.Terms("x"))(expr.OpEQ)(expr.Terms("100")), Weak()))
	DBG("dump=>\n%v", solver.Dump())

	c10 := expr.Eqn(expr.Terms("x"))(expr.OpLEQ)(expr.Terms("10"))
	c20 := expr.Eqn(expr.Terms("x"))(expr.OpLEQ)(expr.Terms("20"))

	solver.AddConstraint(Constraint(c10, Required()))
	solver.AddConstraint(Constraint(c20, Required()))

	solver.UpdateVariables()
	assert.Equal(10.0, ValueOf(c10.Left().WithVars("x").VarAt(0)))

	/*
		solver.AddConstraint(c10).AddConstraint(c20)
		assert.True(Util.IsApproxEqual(x.Value(), 10.0))

		solver.RemoveConstraint(c10)
		assert.True(Util.IsApproxEqual(x.Value(), 20.0))

		solver.RemoveConstraint(c20)
		assert.True(Util.IsApproxEqual(x.Value(), 100.0))

		c10again := LinearInequalityVV0(x, CnLEQ, 10.0)
		solver.AddConstraint(c10).AddConstraint(c10again)
		assert.True(Util.IsApproxEqual(x.Value(), 10.0))

		solver.RemoveConstraint(c10)
		assert.True(Util.IsApproxEqual(x.Value(), 10.0))

		solver.RemoveConstraint(c10again)
		assert.True(Util.IsApproxEqual(x.Value(), 100.0))
	*/
}
