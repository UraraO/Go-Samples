package containers

import (
	"errors"
	"fmt"
)

// type Stack[T int | uint | string | byte] struct {
// 	rawData []T
// 	cap     uint
// 	size    uint
// }

// func InitStack(cap int) T {

// }

// func (self Stack) Push(elem T) {

// }

type Item interface{}

type Stack struct {
	rawData []Item
	Cap     uint
	Size    uint
}

func InitStack(cap uint) Stack {
	if cap == 0 {
		err := errors.New("cap is 0, please provide a capacity number which is not 0")
		panic(err)
	}
	return Stack{
		rawData: make([]Item, cap),
		Cap:     cap,
		Size:    0,
	}
}

func (stk *Stack) Top() (bool, Item) {
	if stk.Size == 0 {
		return false, nil
	}
	return true, stk.rawData[stk.Size-1]
}

func (stk *Stack) Push(elem Item) bool {
	if stk.Size >= stk.Cap {
		return false
	}
	stk.rawData[stk.Size] = elem
	stk.Size += 1
	return true
}

func (stk *Stack) Pop() (bool, Item) {
	if stk.Size == 0 {
		return false, nil
	}
	res := stk.rawData[stk.Size]
	stk.Size -= 1
	return true, res
}

func (stk Stack) String() string {
	res := "Stack: "
	for i := 0; i < int(stk.Size); i++ {
		res += fmt.Sprintf("%v ", stk.rawData[i])
	}
	return res
}

func StackTest() {
	stk := InitStack(10)
	fmt.Println("StackTest: ")
	fmt.Println("stk.Cap =", stk.Cap)
	fmt.Println("stk.Size =", stk.Size)
	stk.Push(1)
	fmt.Println("stk.Push 1")
	_, item := stk.Top()
	fmt.Println("stk.Top =", item)
	stk.Push(2)
	_, item = stk.Top()
	fmt.Println("stk.Push 2")
	fmt.Println("stk.Top =", item)
	fmt.Println("stk.Size =", stk.Size)
	stk.Pop()
	fmt.Println("stk.Pop")
	fmt.Println("stk.Size =", stk.Size)
	stk.Pop()
	fmt.Println("stk.Pop")
	fmt.Println("stk.Size =", stk.Size)
	stk.Pop()
	fmt.Println("stk.Pop")
	fmt.Println("stk.Size =", stk.Size)

	stk.Push(1)
	stk.Push(1)
	stk.Push(1)
	stk.Push(1)
	stk.Push(1)
	stk.Push(1)
	stk.Push(1)
	fmt.Printf("%v\n", stk)
}
