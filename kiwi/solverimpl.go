/*
Package kiwi, implements a Cassowary solver for golang based on kiwi (github.com/nucleic/kiwi).
*/
package kiwi

import (
	"bytes"
	"math"

	"github.com/noypi/math0"
	"github.com/noypi/math0/expr"
)

type ISolver interface {
	AddConstraint(cn *Constraint) error
	RemoveConstraint(cn *Constraint) error
	HasConstraint(cn *Constraint) bool
	AddEditVariable(v expr.IVariable, strength StrengthType)
	RemoveEditVariable(v expr.IVariable)
	HasEditVariable(v expr.IVariable) bool
	SuggestValue(variable expr.IVariable, value float64)
	UpdateVariables()
	Symbol(t SymbolType) _Symbol
	DualOptimize()
	Var(name string) *Variable
	Dump() string
}

type _SolverImpl struct {
	id_tick         int64
	infeasible_rows _SymbolList
	objective       *_Row
	artificial      *_Row

	cns   _CnMap
	vars  _VarMap
	rows  _RowMap
	edits _EditMap
}

type _Tag struct {
	marker, other _Symbol
}

type _EditInfo struct {
	tag        *_Tag
	constraint *Constraint
	constant   float64
}

func Solver() ISolver {
	o := new(_SolverImpl)
	o.objective = Row(0.0)
	o.cns = _CnMap{}
	o.vars = _VarMap{}
	o.rows = _RowMap{}
	o.edits = _EditMap{}
	o.id_tick = 1
	return o
}

// Add a constraint to the solver.
//
// returns
// -------
// DuplicateConstraint
//  	The given constraint has already been added to the solver.

// UnsatisfiableConstraint
//  	The given constraint is required and cannot be satisfied.
func (this *_SolverImpl) AddConstraint(cn *Constraint) error {
	if _, has := this.cns.Get(cn); has {
		return DuplicateConstraint(cn)
	}

	// Creating a row causes symbols to reserved for the variables
	// in the constraint. If this method exits with an exception,
	// then its possible those variables will linger in the var map.
	// Since its likely that those variables will be used in other
	// constraints and since exceptional conditions are uncommon,
	// i'm not too worried about aggressive cleanup of the var map.
	var tag _Tag
	row := this.createRow(cn, &tag)
	subject := this.chooseSubject(row, &tag)

	// If chooseSubject could find a valid entering symbol, one
	// last option is available if the entire row is composed of
	// dummy variables. If the constant of the row is zero, then
	// this represents redundant constraints and the new dummy
	// marker can enter the basis. If the constant is non-zero,
	// then it represents an unsatisfiable constraint.
	if Invalid == subject.Type && this.allDummies(row) {
		if !math0.IsApproxEqual(row.constant, 0.0) {
			return UnsatisfiableConstraint(cn)
		} else {
			subject = tag.marker
		}
	}

	// If an entering symbol still isn't found, then the row must
	// be added using an artificial variable. If that fails, then
	// the row represents an unsatisfiable constraint.
	if Invalid == subject.Type {
		if !this.addWithArtificialVariable(row) {
			return UnsatisfiableConstraint(cn)
		}

	} else {
		row.solveFor(subject)
		this.substitute(subject, row)
		this.rows.Put(subject, row)
	}

	this.cns.Put(cn, &tag)

	// Optimizing after each constraint is added performs less
	// aggregate work due to a smaller average system size. It
	// also ensures the solver remains in a consistent state.
	this.optimize(this.objective)

	return nil
}

// Remove a constraint from the solver.
//
// returns
// -------
// UnknownConstraint
// The given constraint has not been added to the solver.
func (this *_SolverImpl) RemoveConstraint(cn *Constraint) error {
	tag, has := this.cns.Get(cn)
	if !has {
		return UnknownConstraint(cn)
	}

	this.cns.Delete(cn)

	// Remove the error effects from the objective function
	// *before* pivoting, or substitutions into the objective
	// will lead to incorrect solver results.
	this.removeConstraintEffects(cn, tag)

	// If the marker is basic, simply drop the row. Otherwise,
	// pivot the marker into the basis and then drop the row.
	itrow, has := this.rows.Get(tag.marker)
	if has {
		this.rows.Delete(tag.marker)
	} else {
		var leaving _Symbol
		leaving, itrow = this.getMarkerLeavingRow(tag.marker)
		if nil == itrow {
			return InternalSolverError("failed to find leaving row")
		}
		this.rows.Delete(leaving)
		itrow.solveForLhs(leaving, tag.marker)
		this.substitute(tag.marker, itrow)
	}

	// Optimizing after each constraint is removed ensures that the
	// solver remains consistent. It makes the solver api easier to
	// use at a small tradeoff for speed.
	this.optimize(this.objective)
	return nil
}

func (this _SolverImpl) HasConstraint(cn *Constraint) bool {
	_, has := this.cns.Get(cn)
	return has
}

// Add an edit variable to the solver.
//
// This method should be called before the `suggestValue` method is
// used to supply a suggested value for the given edit variable.
//
// Throws
// ------
// DuplicateEditVariable
// 	The given edit variable has already been added to the solver.
//
// BadRequiredStrength
// 	The given strength is >= required.
func (this *_SolverImpl) AddEditVariable(v expr.IVariable, strength StrengthType) {
	if _, has := this.edits.Get(v); has {
		panic(DuplicateEditVariable(v))
	}

	strength = clipStrength(strength)
	if math0.IsApproxEqual(float64(strength), float64(_Required)) {
		panic(BadRequiredStrength())
	}

	cn := NewConstraint(expr.Equation(expr.NewExpr(expr.NewTerm(1.0, v)), expr.EQ, nil), strength)
	this.AddConstraint(cn)
	var info _EditInfo
	info.tag, _ = this.cns.Get(cn)
	info.constraint = cn
	info.constant = 0.0
	this.edits.Put(v, &info)
}

// Remove an edit variable from the solver.
//
// Throws
// ------
// UnknownEditVariable
// 	The given edit variable has not been added to the solver.
func (this *_SolverImpl) RemoveEditVariable(v expr.IVariable) {
	info, has := this.edits.Get(v)
	if !has {
		panic(UnknownEditVariable(v))
	}

	this.RemoveConstraint(info.constraint)
	this.edits.Delete(v)
}

func (this _SolverImpl) HasEditVariable(v expr.IVariable) bool {
	_, has := this.edits.Get(v)
	return has
}

// Suggest a value for the given edit variable.
//
// This method should be used after an edit variable as been added to
// the solver in order to suggest the value for that variable.
//
// Throws
// ------
// UnknownEditVariable
// 	The given edit variable has not been added to the solver.
//
func (this *_SolverImpl) SuggestValue(variable expr.IVariable, value float64) {
	info, has := this.edits.Get(variable)
	if !has {
		panic(UnknownEditVariable(variable))
	}

	defer this.DualOptimize()
	delta := value - info.constant
	info.constant = value

	// Check first if the positive error variable is basic.
	itrow, has := this.rows.Get(info.tag.marker)
	if has {
		if itrow.Add(-delta) < 0.0 {
			this.infeasible_rows = append(this.infeasible_rows, info.tag.marker)
		}
		return
	}

	// Check next if the negative error variable is basic.
	itrow, has = this.rows.Get(info.tag.other)
	if has {
		if itrow.Add(delta) < 0.0 {
			this.infeasible_rows = append(this.infeasible_rows, info.tag.other)
		}
		return
	}

	// Otherwise update each row where the error variables exist.
	for k, itrow := range this.rows {
		coeff := itrow.coefficientFor(info.tag.marker)
		if !math0.IsApproxEqual(coeff, 0.0) &&
			itrow.Add(delta*coeff) < 0.0 &&
			External != k.Type {
			this.infeasible_rows = append(this.infeasible_rows, k)
		}
	}
}

// Update the values of the external solver variables.
func (this *_SolverImpl) UpdateVariables() {
	this.vars.Each(func(variable expr.IVariable, symbol _Symbol) bool {
		if row, has := this.rows.Get(symbol); !has {
			variable.(*Variable).value = 0.0
		} else {
			variable.(*Variable).value = row.constant
		}
		return true
	})
}

// Test whether an edit variable has been added to the solver.

// Create a new Row object for the given constraint.
//
// The terms in the constraint will be converted to cells in the row.
// Any term in the constraint with a coefficient of zero is ignored.
// This method uses the `getVarSymbol` method to get the symbol for
// the variables added to the row. If the symbol for a given cell
// variable is basic, the cell variable will be substituted with the
// basic row.
//
// The necessary slack and error variables will be added to the row.
// If the constant for the row is negative, the sign for the row
// will be inverted so the constant becomes positive.
//
// The tag will be updated with the marker and error symbols to use
// for tracking the movement of the constraint in the tableau.

func (this *_SolverImpl) createRow(cn *Constraint, tag *_Tag) *_Row {
	expression := cn.expression.Clone()
	row := Row(expression.Constant())

	// Substitute the current basic variables into the row.
	expression.EachTerm(func(term expr.ITerm) bool {
		if 0 == len(term.Vars()) {
			return true
		}
		symbol := this.getVarSymbol(term.VarAt(0))
		if o, has := this.rows.Get(symbol); has {
			row.insertRow(o, term.C())
		} else {
			row.insert(symbol, term.C())
		}
		return true
	})

	// Add the necessary slack, error, and dummy variables.
	switch cn.relation {
	case expr.LEQ:
		fallthrough
	case expr.GEQ:
		var coeff float64
		if expr.LEQ == cn.relation {
			coeff = 1.0
		} else {
			coeff = -1.0
		}
		slack := this.Symbol(Slack)
		tag.marker = slack
		row.insert(slack, coeff)
		if cn.strength < _Required {
			symerr := this.Symbol(Error)
			tag.other = symerr
			row.insert(symerr, -coeff)
			this.objective.insert(symerr, float64(cn.strength))
		}
		break

	case expr.EQ:
		if cn.strength < _Required {
			errplus := this.Symbol(Error)
			errminus := this.Symbol(Error)
			tag.marker = errplus
			tag.other = errminus
			row.insert(errplus, -1.0) // v = eplus - eminus
			row.insert(errminus, 1.0) // v - eplus + eminus = 0
			this.objective.insert(errplus, float64(cn.strength))
			this.objective.insert(errminus, float64(cn.strength))

		} else {
			dummy := this.Symbol(Dummy)
			tag.marker = dummy
			row.insert(dummy, 1.0)
		}
		break
	}

	if row.constant < 0.0 {
		row.reverseSign()
	}

	return row
}

func (this *_SolverImpl) Symbol(t SymbolType) _Symbol {
	sym := _Symbol{Id: this.id_tick, Type: t}
	this.id_tick++
	return sym
}

// Get the symbol for the given variable.
//
// If a symbol does not exist for the variable, one will be created.
func (this *_SolverImpl) getVarSymbol(v expr.IVariable) _Symbol {
	sym, has := this.vars.Get(v)
	if has {
		return sym
	}
	sym = _Symbol{
		Id:   this.id_tick,
		Type: External,
	}
	this.vars.Put(v, sym)
	return sym
}

// Choose the subject for solving for the row.
//
// This method will choose the best subject for using as the solve
// target for the row. An invalid symbol will be returned if there
// is no valid target.
//
// The symbols are chosen according to the following precedence:
//
// 1) The first symbol representing an external variable.
// 2) A negative slack or error tag variable.
//
// If a subject cannot be found, an invalid symbol will be returned.
func (this *_SolverImpl) chooseSubject(row *_Row, tag *_Tag) _Symbol {
	for k, _ := range row.cells {
		if External == k.Type {
			return k
		}
	}

	if Slack == tag.marker.Type || Error == tag.marker.Type {
		if row.coefficientFor(tag.marker) < 0.0 {
			return tag.marker
		}
	}

	if Slack == tag.other.Type || Error == tag.other.Type {
		if row.coefficientFor(tag.other) < 0.0 {
			return tag.other
		}
	}

	return _Symbol{Type: Invalid}
}

// Test whether a row is composed of all dummy variables.
func (this _SolverImpl) allDummies(row *_Row) bool {
	for k, _ := range row.cells {
		if Dummy != k.Type {
			return false
		}
	}
	return true
}

// Add the row to the tableau using an artificial variable.
//
// This will return false if the constraint cannot be satisfied.
func (this *_SolverImpl) addWithArtificialVariable(row *_Row) bool {
	// Create and add the artificial variable to the tableau
	art := this.Symbol(Slack)
	this.rows[art] = row.Clone()
	this.artificial = row.Clone()

	// Optimize the artificial objective. This is successful
	// only if the artificial objective is optimized to zero.
	this.optimize(this.artificial)
	success := math0.IsApproxEqual(this.artificial.constant, 0.0)
	this.artificial = nil

	itrow, has := this.rows.Get(art)
	if has {
		this.rows.Delete(art)
		if 0 == len(itrow.cells) {
			return success
		}

		entering := this.anyPivotableSymbol(itrow)
		if Invalid == entering.Type {
			return false // unsatisfiable (will this ever happen?)
		}
		itrow.solveForLhs(art, entering)
		this.substitute(entering, itrow)
		this.rows.Put(entering, itrow)
	}

	// Remove the artificial variable from the tableau.
	for _, v := range this.rows {
		v.remove(art)
	}
	this.objective.remove(art)
	return success
}

// Optimize the system for the given objective function.
//
// This method performs iterations of Phase 2 of the simplex method
// until the objective function reaches a minimum.
//
// Throws
// ------
// InternalSolverError
// 	The value of the objective function is unbounded.
func (this *_SolverImpl) optimize(objective *_Row) {
	for {
		entering := this.getEnteringSymbol(objective)
		if Invalid == entering.Type {
			return
		}
		leaving, row := this.getLeavingRow(entering)
		if nil == row {
			panic(InternalSolverError("The objective is unbounded."))
		}

		// pivot the entering symbol into the basis
		this.rows.Delete(leaving)
		row.solveForLhs(leaving, entering)
		this.substitute(entering, row)
		this.rows.Put(entering, row)

	}
}

// Compute the entering variable for a pivot operation.
//
// This method will return first symbol in the objective function which
// is non-dummy and has a coefficient less than zero. If no symbol meets
// the criteria, it means the objective function is at a minimum, and an
// invalid symbol is returned.
func (this _SolverImpl) getEnteringSymbol(objective *_Row) _Symbol {
	for k, v := range objective.cells {
		if Dummy != k.Type && v < 0.0 {
			return k
		}
	}
	return _Symbol{Type: Invalid}
}

// Compute the row which holds the exit symbol for a pivot.
//
// This method will return an iterator to the row in the row map
// which holds the exit symbol. If no appropriate exit symbol is
// found, the end() iterator will be returned. This indicates that
// the objective function is unbounded.
func (this _SolverImpl) getLeavingRow(entering _Symbol) (sym _Symbol, row *_Row) {
	ratio := math.MaxFloat64
	for k, v := range this.rows {
		if External != k.Type {
			var temp float64 = v.coefficientFor(entering)
			if temp < 0.0 {
				var temp_ratio float64 = -v.constant / temp
				if temp_ratio < ratio {
					ratio = temp_ratio
					sym = k
					row = v
				}
			}
		}
	}
	return
}

// Substitute the parametric symbol with the given row.
//
// This method will substitute all instances of the parametric symbol
// in the tableau and the objective function with the given row.
func (this *_SolverImpl) substitute(sym _Symbol, row *_Row) {
	for k, v := range this.rows {
		v.substitute(sym, row)
		if External != k.Type &&
			v.constant < 0.0 {
			this.infeasible_rows = append(this.infeasible_rows, k)
		}
	}

	this.objective.substitute(sym, row)
	if nil != this.artificial {
		this.artificial.substitute(sym, row)
	}
}

// Get the first Slack or Error symbol in the row.
//
// If no such symbol is present, and Invalid symbol will be returned.
func (this _SolverImpl) anyPivotableSymbol(row *_Row) _Symbol {
	for sym, _ := range row.cells {
		if Slack == sym.Type || Error == sym.Type {
			return sym
		}
	}
	return _Symbol{Type: Invalid}
}

// Remove the effects of a constraint on the objective function.
func (this *_SolverImpl) removeConstraintEffects(cn *Constraint, tag *_Tag) {
	if Error == tag.marker.Type {
		this.removeMarkerEffects(tag.marker, cn.strength)
	} else if Error == tag.other.Type {
		this.removeMarkerEffects(tag.other, cn.strength)
	}
}

// Remove the effects of an error marker on the objective function.
func (this *_SolverImpl) removeMarkerEffects(marker _Symbol, strength StrengthType) {
	if row, has := this.rows.Get(marker); has {
		this.objective.insertRow(row, -float64(strength))
	} else {
		this.objective.insert(marker, -float64(strength))
	}
}

func (this *_SolverImpl) getMarkerLeavingRow(marker _Symbol) (_Symbol, *_Row) {
	dmax := math.MaxFloat64
	r1 := dmax
	r2 := dmax
	type _pair struct {
		k _Symbol
		v *_Row
	}
	var first, second, third _pair
	for k, v := range this.rows {
		c := v.coefficientFor(marker)
		if math0.IsApproxEqual(c, 0.0) {
			continue
		}

		if External == k.Type {
			third.k, third.v = k, v

		} else if c < 0.0 {
			r := -v.constant / c
			if r < r1 {
				r1 = r
				first.k, first.v = k, v
			}

		} else {
			r := v.constant / c
			if r < r2 {
				r2 = r
				second.k, second.v = k, v
			}
		}
	}

	if nil != first.v {
		return first.k, first.v
	}
	if nil != second.v {
		return second.k, second.v
	}

	return third.k, third.v
}

// Optimize the system using the dual of the simplex method.
//
// The current state of the system should be such that the objective
// function is optimal, but not feasible. This method will perform
// an iteration of the dual simplex method to make the solution both
// optimal and feasible.
//
// Throws
// ------
// InternalSolverError
// 	The system cannot be dual optimized.
func (this *_SolverImpl) DualOptimize() {
	for 0 < len(this.infeasible_rows) {
		leaving := this.infeasible_rows[len(this.infeasible_rows)-1]
		this.infeasible_rows = this.infeasible_rows[:len(this.infeasible_rows)-1]
		row, has := this.rows.Get(leaving)
		if has && row.constant < 0.0 {
			entering := this.getDualEnteringSymbol(row)
			if Invalid == entering.Type {
				panic(InternalSolverError("Dual optimize failed."))
			}
			this.rows.Delete(leaving)
			row.solveForLhs(leaving, entering)
			this.substitute(entering, row)
			this.rows.Put(entering, row)
		}
	}
}

func (this _SolverImpl) getDualEnteringSymbol(row *_Row) _Symbol {
	entering := _Symbol{Type: Invalid}
	ratio := math.MaxFloat64
	for k, v := range row.cells {
		if v > 0.0 && Dummy != k.Type {
			coeff := this.objective.coefficientFor(k)
			r := coeff / v
			if r < ratio {
				ratio = r
				entering = k
			}
		}
	}
	return entering
}

func (this _SolverImpl) Var(name string) *Variable {
	if o, has := this.vars[name]; has {
		return o.k
	}
	return nil
}

func (this _SolverImpl) Dump() string {
	buf := bytes.NewBufferString("Objective\n")
	buf.WriteString("---------\n")
	buf.WriteString(this.objective.Dump())
	buf.WriteString("\n")

	buf.WriteString("Tableau\n")
	buf.WriteString("-------\n")
	buf.WriteString(this.rows.Dump())
	buf.WriteString("\n")

	buf.WriteString("Infeasible\n")
	buf.WriteString("----------\n")
	buf.WriteString(this.infeasible_rows.Dump())
	buf.WriteString("\n")

	buf.WriteString("Variables\n")
	buf.WriteString("--------\n")
	buf.WriteString(this.vars.Dump())
	buf.WriteString("\n")

	buf.WriteString("Constraints\n")
	buf.WriteString("-----------\n")
	buf.WriteString(this.cns.Dump())
	buf.WriteString("\n")
	buf.WriteString("\n")

	return buf.String()
}
