package fun

import "fmt"

var cover = true

func Add(a, b int) int {
	if cover {
		fmt.Println("cover")
	}
	return a + b
}
