package expr

import (
	"testing"

	assertpkg "github.com/stretchr/testify/assert"
)

func init() {
	SetLogLevel(LogDebug)
}

type _testTerm struct {
	term    ITerm
	expectk string
	expectS string
}

var ttTermsTen = []_testTerm{
	//0
	_testTerm{
		NewTerm(10, VariableN("x", 2)),
		"x^2",
		"10*x^2"},
	//1
	_testTerm{
		NewTerm(11, VariableN("x", 2)),
		"x^2",
		"21*x^2",
	},

	//2
	_testTerm{
		NewTerm(-3, VariableN("x", 2)),
		"x^2",
		"18*x^2",
	},

	//3
	_testTerm{
		NewTerm(-20, VariableN("x", 2)),
		"x^2",
		"-2*x^2",
	},

	//4
	_testTerm{
		NewTerm(-20, Variable("y")),
		"x^2,y",
		"-2*x^2 + -20*y",
	},

	//5
	_testTerm{
		NewTerm(-123),
		"x^2,y",
		"-2*x^2 + -20*y + -123",
	},

	//6
	_testTerm{
		NewTerm(5, VariableN("x", 2)),
		"x^2,y",
		"3*x^2 + -20*y + -123",
	},

	//7
	_testTerm{
		NewTerm(5, VariableN("x", 2), Variable("y")),
		"x^2*y,x^2,y",
		"5*x^2*y + 3*x^2 + -20*y + -123",
	},

	//8
	_testTerm{
		NewTerm(5, VariableN("x", 2), Variable("y")),
		"x^2*y,x^2,y",
		"10*x^2*y + 3*x^2 + -20*y + -123",
	},

	//9
	_testTerm{
		NewTerm(125),
		"x^2*y,x^2,y",
		"10*x^2*y + 3*x^2 + -20*y + 2",
	},
}

func TestExpression(t *testing.T) {
	assert := assertpkg.New(t)

	expr := NewExpr()
	for i, t := range ttTermsTen {
		expr.AddTerm(t.term)
		assert.Equal(t.expectk, expr.Key(), "i=%d", i)
		assert.Equal(t.expectS, expr.String(), "i=%d", i)
	}
}

var ttTermsLinearTen = []ITerm{
	NewTerm(1, Variable("x")),
	NewTerm(2, Variable("y")),
	NewTerm(3, Variable("z")),
	NewTerm(4, Variable("x")),
	NewTerm(5, Variable("y")),
	NewTerm(6, Variable("y")),
	NewTerm(7, Variable("z")),
	NewTerm(8, Variable("x")),
	NewTerm(9, Variable("y")),
	NewTerm(10, Variable("z")),
}

var ttTermsRemoveTen = []_testTerm{
	//0
	_testTerm{
		NewTerm(10, VariableN("x", 2)),
		"x^2",
		"10*x^2"},
	//1
	_testTerm{
		NewTerm(-10, VariableN("x", 2)),
		"0",
		"0",
	},

	//2
	_testTerm{
		NewTerm(-3, VariableN("x", 2)),
		"x^2",
		"-3*x^2",
	},

	//3
	_testTerm{
		NewTerm(3, VariableN("x", 2)),
		"0",
		"0",
	},

	//4
	_testTerm{
		NewTerm(3, VariableN("x", 2)),
		"x^2",
		"3*x^2",
	},

	//5
	_testTerm{
		NewTerm(3, Variable("y")),
		"x^2,y",
		"3*x^2 + 3*y",
	},

	//6
	_testTerm{
		NewTerm(-3),
		"x^2,y",
		"3*x^2 + 3*y + -3",
	},

	//7
	_testTerm{
		NewTerm(-3, Variable("y")),
		"x^2",
		"3*x^2 + -3",
	},

	//8
	_testTerm{
		NewTerm(3),
		"x^2",
		"3*x^2",
	},

	//9
	_testTerm{
		NewTerm(-3, VariableN("x", 2)),
		"0",
		"0",
	},
}

func BenchmarkExpression_AddTerm_Ten(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr := NewExpr()
		for _, t := range ttTermsTen {
			expr.AddTerm(t.term)
		}
	}
}

func BenchmarkExpression_AddTerm_TenLinear(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr := NewExpr()
		for _, t := range ttTermsLinearTen {
			expr.AddTerm(t)
		}
	}
}

func TestExpression_RemovedTermTen(t *testing.T) {
	assert := assertpkg.New(t)

	expr := NewExpr()
	for i, t := range ttTermsRemoveTen {
		expr.AddTerm(t.term)
		assert.Equal(t.expectk, expr.Key(), "key @ i=%d", i)
		assert.Equal(t.expectS, expr.String(), "string @ i=%d", i)
	}
}

func BenchmarkExpression_RemovedTermTen(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr := NewExpr()
		for _, t := range ttTermsRemoveTen {
			expr.AddTerm(t.term)
		}
	}
}
