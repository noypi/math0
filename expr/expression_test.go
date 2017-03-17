package expr

import (
	"testing"

	assertpkg "github.com/stretchr/testify/assert"
)

func init() {
	SetLogLevel(LogDebug)
}

func TestExpression(t *testing.T) {
	assert := assertpkg.New(t)

	type _test struct {
		term    ITerm
		expectk string
		expectS string
	}

	tt := []_test{
		//0
		_test{
			NewTerm(10, VariableN("x", 2)),
			"x^2",
			"10*x^2"},
		//1
		_test{
			NewTerm(11, VariableN("x", 2)),
			"x^2",
			"21*x^2",
		},

		//2
		_test{
			NewTerm(-3, VariableN("x", 2)),
			"x^2",
			"18*x^2",
		},

		//3
		_test{
			NewTerm(-20, VariableN("x", 2)),
			"x^2",
			"-2*x^2",
		},

		//4
		_test{
			NewTerm(-20, Variable("y")),
			"x^2,y",
			"-2*x^2 + -20*y",
		},

		//5
		_test{
			NewTerm(-123),
			"x^2,y",
			"-2*x^2 + -20*y + -123",
		},

		//6
		_test{
			NewTerm(5, VariableN("x", 2)),
			"x^2,y",
			"3*x^2 + -20*y + -123",
		},

		//7
		_test{
			NewTerm(5, VariableN("x", 2), Variable("y")),
			"x^2*y,x^2,y",
			"5*x^2*y + 3*x^2 + -20*y + -123",
		},
	}

	expr := NewExpr()
	for i, t := range tt {
		expr.AddTerm(t.term)
		assert.Equal(t.expectk, expr.Key(), "i=%d", i)
		assert.Equal(t.expectS, expr.String(), "i=%d", i)
	}
}
