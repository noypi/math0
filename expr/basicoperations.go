package expr

type BasicOperation int

const (
	Addition BasicOperation = iota
	Subtraction
	Multiplication
	Division
)

func (this BasicOperation) Apply(expr IExpression, c float64) {
	switch this {
	case Addition:
		fallthrough
	case Subtraction:
		expr.AddTerm(NewTerm(c))

	case Multiplication:
		fallthrough
	case Division:
		expr.EachTerm(func(term ITerm) bool {
			switch this {
			case Multiplication:
				term.SetC(term.C() * c)
			case Division:
				term.SetC(term.C() / c)
			}

			// get all terms
			return true
		})

	}
	SimplifyExpression(expr)
}
