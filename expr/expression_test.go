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
		NewTerm(10, NewVarN("x", 2)),
		"x^2",
		"10(x^2)"},
	//1
	_testTerm{
		NewTerm(11, NewVarN("x", 2)),
		"x^2",
		"21(x^2)",
	},

	//2
	_testTerm{
		NewTerm(-3, NewVarN("x", 2)),
		"x^2",
		"18(x^2)",
	},

	//3
	_testTerm{
		NewTerm(-20, NewVarN("x", 2)),
		"x^2",
		"-2(x^2)",
	},

	//4
	_testTerm{
		NewTerm(-20, NewVar("y")),
		"x^2,y",
		"-2(x^2) + -20(y)",
	},

	//5
	_testTerm{
		NewTerm(-123),
		"x^2,y",
		"-123 + -2(x^2) + -20(y)",
	},

	//6
	_testTerm{
		NewTerm(5, NewVarN("x", 2)),
		"x^2,y",
		"-123 + 3(x^2) + -20(y)",
	},

	//7
	_testTerm{
		NewTerm(5, NewVarN("x", 2), NewVar("y")),
		"x^2,x^2*y,y",
		"-123 + 3(x^2) + 5(x^2*y) + -20(y)",
	},

	//8
	_testTerm{
		NewTerm(5, NewVarN("x", 2), NewVar("y")),
		"x^2,x^2*y,y",
		"-123 + 3(x^2) + 10(x^2*y) + -20(y)",
	},

	//9
	_testTerm{
		NewTerm(125),
		"x^2,x^2*y,y",
		"2 + 3(x^2) + 10(x^2*y) + -20(y)",
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
	NewTerm(1, NewVar("x")),
	NewTerm(2, NewVar("y")),
	NewTerm(3, NewVar("z")),
	NewTerm(4, NewVar("x")),
	NewTerm(5, NewVar("y")),
	NewTerm(6, NewVar("y")),
	NewTerm(7, NewVar("z")),
	NewTerm(8, NewVar("x")),
	NewTerm(9, NewVar("y")),
	NewTerm(10, NewVar("z")),
}

var ttTermsRemoveTen = []_testTerm{
	//0
	_testTerm{
		NewTerm(10, NewVarN("x", 2)),
		"x^2",
		"10(x^2)"},
	//1
	_testTerm{
		NewTerm(-10, NewVarN("x", 2)),
		"",
		"0",
	},

	//2
	_testTerm{
		NewTerm(-3, NewVarN("x", 2)),
		"x^2",
		"-3(x^2)",
	},

	//3
	_testTerm{
		NewTerm(3, NewVarN("x", 2)),
		"",
		"0",
	},

	//4
	_testTerm{
		NewTerm(3, NewVarN("x", 2)),
		"x^2",
		"3(x^2)",
	},

	//5
	_testTerm{
		NewTerm(3, NewVar("y")),
		"x^2,y",
		"3(x^2) + 3(y)",
	},

	//6
	_testTerm{
		NewTerm(-3),
		"x^2,y",
		"-3 + 3(x^2) + 3(y)",
	},

	//7
	_testTerm{
		NewTerm(-3, NewVar("y")),
		"x^2",
		"-3 + 3(x^2)",
	},

	//8
	_testTerm{
		NewTerm(3),
		"x^2",
		"3(x^2)",
	},

	//9
	_testTerm{
		NewTerm(-3, NewVarN("x", 2)),
		"",
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

func TestConstant(t *testing.T) {
	assert := assertpkg.New(t)

	expr := NewExpr()
	expr.AddTerm(
		NewTerm(-3, NewVar("y")),
		NewTerm(4),
		NewTerm(4, NewVar("x")),
		NewTerm(6),
		NewTerm(5, NewVar("z")),
	)

	assert.Equal("10 + 4(x) + -3(y) + 5(z)", expr.String())
	assert.Equal(10.0, expr.Constant())

	terms := expr.Terms()
	assert.Equal(4, len(terms))
}

func TestWithVars(t *testing.T) {
	assert := assertpkg.New(t)

	eqn := Eqn(Terms("5(x*y^2)", "2(x*y)", "3(z)"))(EQ)(nil)
	assert.Equal("3(z)", eqn.Left().WithVars("z").String())
	assert.Equal("5(x*y^2)", eqn.Left().WithVars("x*y^2").String())
	assert.Equal("2(x*y)", eqn.Left().WithVars("x*y").String())

}
