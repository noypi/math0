package kiwi

import (
	"bytes"
	"fmt"

	"github.com/noypi/math0"
	"github.com/noypi/math0/expr"
)

type CellMap map[_Symbol]float64

type _Row struct {
	constant float64
	cells    CellMap
}

func Row(constant float64) *_Row {
	o := &_Row{
		constant: constant,
		cells:    CellMap{},
	}
	return o
}

func (this _Row) Cells() CellMap {
	return this.cells
}

func (this _Row) Constant() float64 {
	return this.constant
}

func (this *_Row) Add(a float64) float64 {
	this.constant += a
	return this.constant
}

// Insert a symbol into the row with a given coefficient.
//
// If the symbol already exists in the row, the coefficient will be
// added to the existing coefficient. If the resulting coefficient
// is zero, the symbol will be removed from the row.
func (this *_Row) insert(symbol _Symbol, coefficient float64 /*= 1.0*/) {
	f, has := this.cells[symbol]
	f += coefficient
	if math0.IsApproxEqual(f, 0.0) {
		// delete when zero
		if has {
			delete(this.cells, symbol)
		}

	} else {
		this.cells[symbol] = f
	}

}

// Insert a row into this row with a given coefficient.
//
// The constant and the cells of the other row will be multiplied by
// the coefficient and added to this row. Any cell with a resulting
// coefficient of zero will be removed from the row.
func (this *_Row) insertRow(other *_Row, coefficient float64) {
	this.constant += (other.constant * coefficient)
	for k, v := range other.cells {
		coeff := v * coefficient
		f := this.cells[k] + coeff
		if math0.IsApproxEqual(0.0, f) {
			delete(this.cells, k)
		} else {
			this.cells[k] = f
		}
	}
}

// Remove the given symbol from the row.
func (this *_Row) remove(symbol _Symbol) {
	delete(this.cells, symbol)
}

// Reverse the sign of the constant and all cells in the row.
func (this *_Row) reverseSign() {
	this.constant = -this.constant
	for k, v := range this.cells {
		this.cells[k] = -v
	}
}

// Solve the row for the given symbol.
//
// This method assumes the row is of the form a * x + b * y + c = 0
// and (assuming solve for x) will modify the row to represent the
// right hand side of x = -b/a * y - c / a. The target symbol will
// be removed from the row, and the constant and other cells will
// be multiplied by the negative inverse of the target coefficient.
//
// The given symbol *must* exist in the row.
func (this *_Row) solveFor(symbol _Symbol) {
	DBG("solveFor this.cells[symbol]=%v", this.cells[symbol])
	coeff := -1.0 / this.cells[symbol]
	delete(this.cells, symbol)

	this.constant *= coeff
	for k, v := range this.cells {
		DBG("cells <= v * coeff=%v, v=%v, coeff=%v", v*coeff, v, coeff)
		this.cells[k] = v * coeff
	}
}

// Solve the row for the given symbols.
//
// This method assumes the row is of the form x = b * y + c and will
// solve the row such that y = x / b - c / b. The rhs symbol will be
// removed from the row, the lhs added, and the result divided by the
// negative inverse of the rhs coefficient.
//
// The lhs symbol *must not* exist in the row, and the rhs symbol
// *must* exist in the row.
func (this *_Row) solveForLhs(lhs, rhs _Symbol) {
	this.insert(lhs, -1.0)
	this.solveFor(rhs)
}

// Get the coefficient for the given symbol.
//
// If the symbol does not exist in the row, zero will be returned.
func (this _Row) coefficientFor(symbol _Symbol) float64 {
	return this.cells[symbol]
}

// Substitute a symbol with the data from another row.
//
// Given a row of the form a * x + b and a substitution of the
// form x = 3 * y + c the row will be updated to reflect the
// expression 3 * a * y + a * c + b.
//
// If the symbol does not exist in the row, this is a no-op.
func (this *_Row) substitute(symbol _Symbol, row *_Row) {
	if coefficient, has := this.cells[symbol]; has {
		delete(this.cells, symbol)
		this.insertRow(row, coefficient)
	}
}

func (this _Row) Clone() *_Row {
	o := new(_Row)
	o.cells = CellMap{}
	for k, v := range this.cells {
		o.cells[k] = v
	}
	o.constant = this.constant
	return o
}

func (this _Row) Dump() string {
	buf := bytes.NewBufferString(expr.ToTrimZero(this.constant))
	for k, v := range this.cells {
		buf.WriteString(fmt.Sprintf(" + %v*", v))
		buf.WriteString(k.Dump())
	}
	buf.WriteString("\n")

	return buf.String()
}
