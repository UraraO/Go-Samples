package calculator

type Calc struct{}

func (c Calc) Add(lhs, rhs int) int {
	return lhs + rhs
}

func (c Calc) Subs(lhs, rhs int) int {
	return lhs - rhs
}
