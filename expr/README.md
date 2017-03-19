```code
func ExampleEquation() {
	left := Terms("3(x)", "-4(y)", "100")
	right := Terms("55z", "6(x^3*y^4)")

	eqn := Eqn(left)(OpLEQ)(right)

	fmt.Printf("%v", eqn)
	// Output: 3(x) + -4(y) + 100 <= 6(x^3*y^4) + 55(z)
}
```