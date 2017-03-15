package expr

import (
	"testing"

	assertpkg "github.com/stretchr/testify/assert"
)

func TestVariable(t *testing.T) {
	assert := assertpkg.New(t)

	varlist := VariableList{
		Variable("z"),
		Variable("x"),
		Variable("v"),
		Variable("y"),
		Variable("z"),
		Variable("x"),
		Variable("x"),
		Variable("y"),
		Variable("y"),
		Variable("za"),
	}
	assert.Equal("v*x^3*y^3*z^2*za", varlist.String())

	varlist = VariableList{
		Variable("z"),
	}
	assert.Equal("z", varlist.String())

}

func BenchmarkVariableStringOneVars(b *testing.B) {
	varlist := VariableList{
		Variable("z"),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		varlist.String()
	}
}

func BenchmarkVariableStringFiveVars(b *testing.B) {
	varlist := VariableList{
		Variable("z"),
		Variable("x"),
		Variable("v"),
		Variable("y"),
		Variable("z"),
		Variable("x"),
		Variable("x"),
		Variable("y"),
		Variable("y"),
		Variable("za"),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		varlist.String()
	}
}
