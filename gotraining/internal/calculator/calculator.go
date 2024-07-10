package calculator

import "context"

type Calc struct{}

type DB interface {
	CreateCalc(ctx context.Context, lhs, rhs int) error
}

func (c Calc) Add(lhs, rhs int) int {
	return lhs + rhs
}

func (c Calc) Subs(lhs, rhs int) int {
	return lhs - rhs
}

func AddInDB(ctx context.Context, db DB, lhs, rhs int) error {
	return db.CreateCalc(ctx, lhs, rhs)
}
