package main

import (
	"fmt"
	"gotraining/internal/calculator"
)

func main() {
	// testCalcAdd(1, 2)
	// testCalcAdd(-1, -1)
	// testCalcAdd(9999999999999, 9999999999999999)
	// ctx := context.Background()
	// ctx.Done()
	// ctx, cancel := context.WithCancel(context.Background())
	// cancel()

}

func testCalcAdd(lhs, rhs int) {
	c := calculator.Calc{}
	fmt.Printf("%v + %v = %v\n", lhs, rhs, c.Add(lhs, rhs))
}
