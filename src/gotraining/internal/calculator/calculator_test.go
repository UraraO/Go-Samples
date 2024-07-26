package calculator

import (
	"context"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	cases := []struct {
		lhs      int
		rhs      int
		expected int
	}{
		{1, 1, 2},
		{0, 0, 0},
		{1, -1, 0},
		{99999999, 1, 100000000},
		{-1, -1, -2},
	}

	cclt := Calc{}
	t.Helper()
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if ans := cclt.Add(c.lhs, c.rhs); ans != c.expected {
				// t.Errorf("%v + %v = %v, but not, return %v\n", c.lhs, c.rhs, c.expected, ans)
				assert.Equal(t, ans, c.expected, "error, ans should be equal to expected")
			}
		})
	}
}

func TestSubs(t *testing.T) {
	cases := []struct {
		lhs      int
		rhs      int
		expected int
	}{
		{1, 1, 0},
		{0, 0, 0},
		{1, -1, 2},
		{99999999, 1, 99999998},
		{-1, -1, 0},
	}

	cclt := Calc{}
	t.Helper()
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if ans := cclt.Subs(c.lhs, c.rhs); ans != c.expected {
				// t.Errorf("%v + %v = %v, but not, return %v\n", c.lhs, c.rhs, c.expected, ans)
				assert.Equal(t, ans, c.expected, "error, ans should be equal to expected")
			}
		})
	}
}

func TestAddInDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockDB(ctrl)

	ctx := context.Background()
	lhs, rhs := 1, 1

	mockDB.EXPECT().CreateCalc(ctx, lhs, rhs).Return(nil)

	err := AddInDB(ctx, mockDB, lhs, rhs)

	assert.NoError(t, err)
}
