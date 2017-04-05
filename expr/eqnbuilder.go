package expr

import (
	"fmt"
	"strings"

	"github.com/noypi/math0"
)

var (
	EqnBuilder_ExprConstructor = NewExpr
	EqnBuilder_TermConstructor = NewTerm
	EqnBuilder_VarConstructor  = NewVarN
)

// Eqn(Terms("3x", "-4y"))(NEQ)(Terms("5"))
func Eqn(left TermList) func(Relation) func(right TermList) IEquation {
	eqn := Equation(EqnBuilder_ExprConstructor(left...), EQ, nil)
	return func(rel Relation) func(right TermList) IEquation {
		eqn.SetRelation(rel)
		return func(right TermList) IEquation {
			if nil == right {
				return eqn
			}
			eqn.SetRight(EqnBuilder_ExprConstructor(right...))
			return eqn
		}
	}
}

// Terms("2x^2*y^3", "2y")
func Terms(terms ...string) TermList {
	var ls TermList
	for _, s := range terms {
		var f float64
		var vars string
		fmt.Sscanf(s, "%f%s", &f, &vars)
		if 0.0 == f && 0 == len(vars) && 0 < len(s) {
			f = 1.0
			vars = s
		}
		vs := []IVariable(Vars(strings.Split(strings.Trim(vars, "()"), "*")...))
		ls = append(ls, EqnBuilder_TermConstructor(f, vs...))
	}
	return ls
}

// Vars("x^2", "x", "y", "y^2")
func Vars(vars ...string) VariableList {
	var vs VariableList
	for _, s := range vars {
		if 0 == len(s) {
			continue
		}
		var x string
		var power float64 = 1.0
		fmt.Sscanf(s, "%s^%f", &x, &power)
		if math0.IsApproxEqual(0.0, power) {
			power = 0.0
		}
		vs = append(vs, EqnBuilder_VarConstructor(x, power))
	}
	return vs
}
