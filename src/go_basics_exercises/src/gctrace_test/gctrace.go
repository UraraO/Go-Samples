package gctracetest

import (
	"time"
)

var sliceSize = 1 << 20
var allocFreq = 20

func makeLargeSlice() []int64 {
	s := make([]int64, sliceSize)
	for i := 0; i < sliceSize; i++ {
		s = append(s, int64(i))
	}
	return s
}

func GCTraceTest() {
	slices := make([][]int64, allocFreq)
	for i := 0; i < allocFreq; i++ {
		slices = append(slices, makeLargeSlice())
	}
	time.Sleep(1 * time.Second)
}
