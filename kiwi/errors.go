package kiwi

import (
	"fmt"

	"github.com/noypi/math0/expr"
)

const (
	ErrTypeDuplicateConstraint     = "DuplicateConstraint"
	ErrTypeUnsatisfiableConstraint = "UnsatisfiableConstraint"
	ErrTypeInternalSolverError     = "InternalSolverError"
	ErrTypeUnknownConstraint       = "UnknownConstraint"
	ErrTypeDuplicateEditVariable   = "DuplicateEditVariable"
	ErrTypeBadRequiredStrength     = "BadRequiredStrength"
	ErrTypeUnknownEditVariable     = "UnknownEditVariable"
)

type ErrorData struct {
	Data interface{}
	t    string
}

func UnknownEditVariable(v expr.IVariable) error {
	o := &ErrorData{
		Data: v,
		t:    ErrTypeUnknownEditVariable,
	}
	return o
}

func BadRequiredStrength() error {
	o := &ErrorData{
		t: ErrTypeBadRequiredStrength,
	}
	return o
}

func DuplicateEditVariable(v expr.IVariable) error {
	o := &ErrorData{
		Data: v,
		t:    ErrTypeDuplicateEditVariable,
	}
	return o
}

func UnknownConstraint(cn *Constraint) error {
	o := &ErrorData{
		Data: cn,
		t:    ErrTypeUnknownConstraint,
	}
	return o
}

func DuplicateConstraint(cn *Constraint) error {
	o := &ErrorData{
		Data: cn,
		t:    ErrTypeDuplicateConstraint,
	}
	return o
}

func UnsatisfiableConstraint(cn *Constraint) error {
	o := &ErrorData{
		Data: cn,
		t:    ErrTypeUnsatisfiableConstraint,
	}
	return o
}

func InternalSolverError(s string) error {
	o := &ErrorData{
		Data: s,
		t:    ErrTypeInternalSolverError,
	}
	return o
}

func (this ErrorData) Error() string {
	return fmt.Sprintf("%s: %v", this.t, this.Data)
}
