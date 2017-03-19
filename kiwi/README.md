Cassowary solver for golang base on kiwi ([github.com/nucleic/kiwi]()).

## Example

```code
func ExampleSolverImpl_AddConstraint() {
	expr.EqnBuilder_VarConstructor = func(name string, power float64) expr.IVariable {
		return kiwi.Variable(name)
	}

	eqn1 := expr.Eqn(expr.Terms("x"))(expr.OpEQ)(expr.Terms("5"))
	eqn2 := expr.Eqn(expr.Terms("y"))(expr.OpEQ)(expr.Terms("10"))

	solver := kiwi.Solver()

	solver.AddConstraint(kiwi.Constraint(eqn1, kiwi.Weak()))
	solver.AddConstraint(kiwi.Constraint(eqn2, kiwi.Weak()))
	solver.UpdateVariables()

	xval := kiwi.ValueOf(eqn1.Left().WithVars("x").VarAt(0))
	yval := kiwi.ValueOf(eqn2.Left().WithVars("y").VarAt(0))
	fmt.Printf("x=%v, y=%v", xval, yval)
	// Output: x=5, y=10
}
```
