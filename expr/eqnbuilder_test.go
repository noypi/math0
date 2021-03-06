package expr

import (
	"fmt"
	"testing"

	assertpkg "github.com/stretchr/testify/assert"
)

func init() {
	SetLogLevel(LogDebug)
}

type _eqnbuildtest struct {
	left     TermList
	right    TermList
	relation Relation
	expected string
}

var ttEqnBuilder = []_eqnbuildtest{
	//0
	_eqnbuildtest{
		left:     Terms("x"),
		right:    Terms("5y"),
		relation: EQ,
		expected: "1(x) == 5(y)",
	},

	//1
	_eqnbuildtest{
		left:     Terms("3x", "-4y"),
		right:    Terms("5"),
		relation: NEQ,
		expected: "3(x) + -4(y) != 5",
	},

	//2
	_eqnbuildtest{
		left:     Terms("3x", "-4y"),
		right:    Terms("5"),
		relation: GEQ,
		expected: "3(x) + -4(y) >= 5",
	},

	//3
	_eqnbuildtest{
		left:     Terms("3(x)", "-4(y)"),
		right:    Terms("55z", "6(x^3*y^4)"),
		relation: LEQ,
		expected: "3(x) + -4(y) <= 6(x^3*y^4) + 55(z)",
	},

	//4
	_eqnbuildtest{
		left:     Terms("3(x)", "-4(y)"),
		right:    nil,
		relation: Lesser,
		expected: "3(x) + -4(y) < 0",
	},

	//5
	_eqnbuildtest{
		left:     Terms("3(x)", "-4(y)"),
		right:    Terms(""),
		relation: Greater,
		expected: "3(x) + -4(y) > 0",
	},
}

func TestEqnBuilder(t *testing.T) {
	assert := assertpkg.New(t)
	for i, t := range ttEqnBuilder {
		eqn := Eqn(t.left)(t.relation)(t.right)
		assert.Equal(t.expected, eqn.String(), "i=%v", i)
	}
}

func ExampleEquation() {
	left := Terms("3(x)", "-4(y)", "100")
	right := Terms("55z", "6(x^3*y^4)")

	eqn := Eqn(left)(LEQ)(right)

	fmt.Printf("%v", eqn)
	// Output: 100 + 3(x) + -4(y) <= 6(x^3*y^4) + 55(z)
}
