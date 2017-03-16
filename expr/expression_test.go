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

	expr := NewExpr()

	expr.AddTerm(NewTerm(10, VariableN("x", 2)))
	assert.Equal("x^2", expr.Key())
	assert.Equal("10*x^2", expr.String())

	expr.AddTerm(NewTerm(11, VariableN("x", 2)))
	assert.Equal("x^2", expr.Key())
	assert.Equal("21*x^2", expr.String())
}
