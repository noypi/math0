package kiwi_test

import (
	"strings"
	"testing"

	"github.com/noypi/math0/expr"
	. "github.com/noypi/math0/kiwi"
	assertpkg "github.com/stretchr/testify/assert"
)

func init() {
	SetLogLevel(LogDebug)
}

func TestSimple1(t *testing.T) {
	assert := assertpkg.New(t)

	x := Var("x")
	y := Var("y")
	xterm := expr.NewTerm(167.0, x)
	yterm := expr.NewTerm(2.0, y)

	solver := Solver()

	eqn := expr.Equation(expr.NewExpr(yterm), expr.EQ, expr.NewExpr(xterm))
	solver.AddConstraint(NewConstraint(eqn, Required()))

	solver.UpdateVariables()
	assert.Equal(0.0, x.Value())
	assert.Equal(x.Value(), y.Value())
}

func TestJustStay1(t *testing.T) {
	assert := assertpkg.New(t)

	x := Var("x")
	y := Var("y")
	xterm := expr.NewTerm(1, x)
	yterm := expr.NewTerm(1, y)

	solver := Solver()

	solver.AddConstraint(NewConstraint(expr.Equation(expr.NewExpr(xterm, expr.NewTerm(-5.0)), expr.EQ, nil), Weak()))
	solver.AddConstraint(NewConstraint(expr.Equation(expr.NewExpr(yterm, expr.NewTerm(-10.0)), expr.EQ, nil), Weak()))

	solver.UpdateVariables()
	assert.Equal(5.0, x.Value())
	assert.Equal(10.0, y.Value())

}

func TestAddDelete1(t *testing.T) {
	assert := assertpkg.New(t)

	expr.EqnBuilder_VarConstructor = func(name string, power float64) expr.IVariable {
		return Var(name)
	}

	solver := Solver()

	solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("x"))(expr.EQ)(expr.Terms("100")), Weak()))

	c10expr := expr.Eqn(expr.Terms("x"))(expr.LEQ)(expr.Terms("10"))
	c20expr := expr.Eqn(expr.Terms("x"))(expr.LEQ)(expr.Terms("20"))

	c10 := NewConstraint(c10expr, Required())
	c20 := NewConstraint(c20expr, Required())
	solver.AddConstraint(c10)
	solver.AddConstraint(c20)

	x := solver.Var("x")
	solver.UpdateVariables()
	assert.Equal(10.0, x.Value())

	solver.RemoveConstraint(c10)
	solver.UpdateVariables()
	assert.Equal(20.0, x.Value())

	solver.RemoveConstraint(c20)
	solver.UpdateVariables()
	assert.Equal(100.0, x.Value())

	c10exprAgain := expr.Eqn(expr.Terms("x"))(expr.LEQ)(expr.Terms("10"))
	c10again := NewConstraint(c10exprAgain, Required())
	solver.AddConstraint(c10)
	err := solver.AddConstraint(c10again)
	assert.NotNil(err)
	assert.True(strings.Contains(err.Error(), "DuplicateConstraint"))
	solver.UpdateVariables()
	assert.Equal(10.0, x.Value())

	solver.RemoveConstraint(c10again)
	solver.UpdateVariables()
	assert.Equal(100.0, x.Value())

	err = solver.RemoveConstraint(c10)
	assert.NotNil(err)
	assert.True(strings.Contains(err.Error(), "UnknownConstraint"))
	solver.UpdateVariables()
	assert.Equal(100.0, x.Value())
}

func TestAddDelete2(t *testing.T) {
	assert := assertpkg.New(t)

	expr.EqnBuilder_VarConstructor = func(name string, power float64) expr.IVariable {
		return Var(name)
	}

	solver := Solver()

	solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("x"))(expr.EQ)(expr.Terms("100.0")), Weak()))
	solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("y"))(expr.EQ)(expr.Terms("120.0")), Strong()))

	c10 := NewConstraint(expr.Eqn(expr.Terms("x"))(expr.LEQ)(expr.Terms("10.0")), Required())
	c20 := NewConstraint(expr.Eqn(expr.Terms("x"))(expr.LEQ)(expr.Terms("20.0")), Required())
	solver.AddConstraint(c10)
	solver.AddConstraint(c20)

	solver.UpdateVariables()
	assert.Equal(10.0, solver.Var("x").Value())
	assert.Equal(120.0, solver.Var("y").Value())

	solver.RemoveConstraint(c10)

	solver.UpdateVariables()
	assert.Equal(20.0, solver.Var("x").Value())
	assert.Equal(120.0, solver.Var("y").Value())

	cxy := NewConstraint(expr.Eqn(expr.Terms("2x"))(expr.EQ)(expr.Terms("y")), Required())
	solver.AddConstraint(cxy)

	solver.UpdateVariables()
	assert.Equal(20.0, solver.Var("x").Value())
	assert.Equal(40.0, solver.Var("y").Value())

	solver.RemoveConstraint(c20)

	solver.UpdateVariables()
	assert.Equal(60.0, solver.Var("x").Value())
	assert.Equal(120.0, solver.Var("y").Value())

	solver.RemoveConstraint(cxy)

	solver.UpdateVariables()
	assert.Equal(100.0, solver.Var("x").Value())
	assert.Equal(120.0, solver.Var("y").Value())
}

func TestCasso1(t *testing.T) {
	assert := assertpkg.New(t)

	expr.EqnBuilder_VarConstructor = func(name string, power float64) expr.IVariable {
		return Var(name)
	}

	solver := Solver()
	solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("x"))(expr.LEQ)(expr.Terms("y")), Required()))
	solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("y"))(expr.EQ)(expr.Terms("x + 3")), Required()))
	solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("x"))(expr.EQ)(expr.Terms("10")), Weak()))
	solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("y"))(expr.EQ)(expr.Terms("10")), Weak()))

	solver.UpdateVariables()
	assert.Equal(10.0, solver.Var("x").Value())
	bTrue := 13.0 == solver.Var("y").Value() ||
		7.0 == solver.Var("x").Value() ||
		10.0 == solver.Var("y").Value()
	assert.True(bTrue)
}

func TestInconsistent1(t *testing.T) {
	assert := assertpkg.New(t)

	expr.EqnBuilder_VarConstructor = func(name string, power float64) expr.IVariable {
		return Var(name)
	}

	solver := Solver()
	err := solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("x"))(expr.EQ)(expr.Terms("10.0")), Required()))
	assert.Nil(err)
	err = solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("x"))(expr.EQ)(expr.Terms("5.0")), Required()))
	assert.NotNil(err)
	assert.True(strings.Contains(err.Error(), "UnsatisfiableConstraint"))
}

func TestInconsistent2(t *testing.T) {
	assert := assertpkg.New(t)

	expr.EqnBuilder_VarConstructor = func(name string, power float64) expr.IVariable {
		return Var(name)
	}

	solver := Solver()
	err := solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("x"))(expr.GEQ)(expr.Terms("10.0")), Required()))
	assert.Nil(err)
	err = solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("x"))(expr.LEQ)(expr.Terms("5.0")), Required()))
	assert.NotNil(err)
	assert.True(strings.Contains(err.Error(), "UnsatisfiableConstraint"))
}

func TestMultiedit(t *testing.T) {
	assert := assertpkg.New(t)

	expr.EqnBuilder_VarConstructor = func(name string, power float64) expr.IVariable {
		return Var(name)
	}

	solver := Solver()
	x, y, w, h := Var("x"), Var("y"), Var("w"), Var("h")
	solver.AddEditVariable(x, Strong())
	solver.AddEditVariable(y, Strong())
	solver.AddEditVariable(w, Strong())
	solver.AddEditVariable(h, Strong())

	solver.SuggestValue(x, 10.0)
	solver.SuggestValue(y, 20.0)

	solver.UpdateVariables()
	assert.Equal(10.0, x.Value())
	assert.Equal(20.0, y.Value())
	assert.Equal(0.0, w.Value())
	assert.Equal(0.0, h.Value())

	solver.SuggestValue(w, 30.0)
	solver.SuggestValue(h, 40.0)

	solver.UpdateVariables()
	assert.Equal(10.0, x.Value())
	assert.Equal(20.0, y.Value())
	assert.Equal(30.0, w.Value())
	assert.Equal(40.0, h.Value())

	solver.SuggestValue(x, 50.0)
	solver.SuggestValue(y, 60.0)

	solver.UpdateVariables()
	assert.Equal(50.0, x.Value())
	assert.Equal(60.0, y.Value())
	assert.Equal(30.0, w.Value())
	assert.Equal(40.0, h.Value())
}

func TestInconsistent3(t *testing.T) {
	assert := assertpkg.New(t)

	expr.EqnBuilder_VarConstructor = func(name string, power float64) expr.IVariable {
		return Var(name)
	}

	solver := Solver()
	err := solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("w"))(expr.GEQ)(expr.Terms("10.0")), Required()))
	assert.Nil(err)
	err = solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("x"))(expr.GEQ)(expr.Terms("w")), Required()))
	assert.Nil(err)
	err = solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("y"))(expr.GEQ)(expr.Terms("x")), Required()))
	assert.Nil(err)
	err = solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("z"))(expr.GEQ)(expr.Terms("y")), Required()))
	assert.Nil(err)
	err = solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("z"))(expr.GEQ)(expr.Terms("8.0")), Required()))
	assert.Nil(err)
	err = solver.AddConstraint(NewConstraint(expr.Eqn(expr.Terms("z"))(expr.LEQ)(expr.Terms("4.0")), Required()))
	assert.NotNil(err)
	assert.True(strings.Contains(err.Error(), "UnsatisfiableConstraint"))
}
