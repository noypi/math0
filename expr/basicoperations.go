package expr

type BasicOperation int

const (
	Addition BasicOperation = iota
	Subtraction
	Multiplication
	Division
)

func (this BasicOperation) Apply(expr IExpression, c float64) {
	expr.EachTerm(func(term ITerm) bool {
		switch this {
		case Addition:
			term.SetC(term.C() + c)
		case Subtraction:
			term.SetC(term.C() - c)
		case Multiplication:
			term.SetC(term.C() * c)
		case Division:
			term.SetC(term.C() / c)
		}
		return true
	})
	SimplifyExpression(expr)
}
