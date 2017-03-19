package expr

import (
	"testing"

	assertpkg "github.com/stretchr/testify/assert"
)

func TestVariable(t *testing.T) {
	assert := assertpkg.New(t)

	varlist := VariableList{
		NewVar("z"),
		NewVar("x"),
		NewVar("v"),
		NewVar("y"),
		NewVar("z"),
		NewVar("x"),
		NewVar("x"),
		NewVar("y"),
		NewVar("y"),
		NewVar("za"),
	}
	assert.Equal("v*x^3*y^3*z^2*za", varlist.Key())

	varlist = VariableList{
		NewVar("z"),
	}
	assert.Equal("z", varlist.Key())

}

func BenchmarkVariableStringOneVars(b *testing.B) {
	varlist := VariableList{
		NewVar("z"),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		varlist.Key()
	}
}

func BenchmarkVariableStringFiveVars(b *testing.B) {
	varlist := VariableList{
		NewVar("z"),
		NewVar("x"),
		NewVar("v"),
		NewVar("y"),
		NewVar("z"),
		NewVar("x"),
		NewVar("x"),
		NewVar("y"),
		NewVar("y"),
		NewVar("za"),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		varlist.Key()
	}
}
