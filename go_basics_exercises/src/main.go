package main

import (
	"fmt"
	gctracetest "go_basics_exercises/src/gctrace_test"

	// cache "go_basics_exercises/src/cache_module"

	"strings"
	"unicode/utf8"
)

func main() {
	// fmt.Println("hello world")

	// 流程控制
	// Print53()

	// 字符串处理
	// fmt.Println(utf8.RuneCountInString("hello, 123爱我国防"))
	// str := "byDSWEW ggg inxudwmn xy"
	// fmt.Println(CalcNumOfChar(&str))
	// fmt.Println(ReverseString("foobar"))
	// fmt.Println(ReverseStringbyPointer(&str))
	// fmt.Println(ReplaceCharsInString(str, "abc", 4))

	// 结构体、数组、包
	// containers.StackTest()

	// 函数
	// MapFuncTest()

	//附加题
	// hdl := bus.NormalHandler{}
	// Bus := bus.InitBusModule()
	// event := bus.InitEvent(bus.EVENT_TYPE_DEBUG, "bus test")

	// debug1 := Bus.RegisterEventListener(bus.EVENT_TYPE_DEBUG, hdl)
	// debug2 := Bus.RegisterEventListener(bus.EVENT_TYPE_DEBUG, hdl)
	// info1 := Bus.RegisterEventListener(bus.EVENT_TYPE_INFO, hdl)

	// Bus.Publish(event)
	// time.Sleep(time.Second * 3)
	// Bus.RemoveEventListener(debug1)
	// Bus.RemoveEventListener(debug2)
	// Bus.RemoveEventListener(info1)

	// 缓存测试
	// c := cache.InitCache(2)

	// c.Put("key1", "value1", 2*time.Second)
	// c.Put("key2", "value2", 2*time.Second)

	// val, found := c.Get("key1")
	// if found {
	// 	fmt.Println("Found key1:", val)
	// }

	// c.Delete("key2")
	// _, found = c.Get("key2")
	// if !found {
	// 	fmt.Println("key2 Not Found")
	// }

	// time.Sleep(3 * time.Second)

	// _, found = c.Get("key1")
	// if !found {
	// 	fmt.Println("Key1 not found or expired")
	// }

	// c.Put("key1", "value1", 2*time.Second)
	// c.Put("key2", "value2", 2*time.Second)
	// ok3 := c.Put("key3", "value3", 2*time.Second)
	// if !ok3 {
	// 	fmt.Println("Key3 put failed")
	// }

	// s1 := make([]int, 0, 2)
	// s1 = append(s1, 1)
	// s2 := append(s1, 1)
	// s3 := s1[1:2]
	// s2 = append(s2, 2, 3)
	// fmt.Println(s1, s2, s3, len(s1), cap(s3))
	// a := [5]int{1, 2, 3, 4, 5}
	// b := a[0:3]
	// fmt.Println(cap(b))
	// b = append(b, 6)
	// fmt.Println(a, b, len(b), cap(b))

	//
	// msgqueue.MQTest()

	// 超时控制
	// timeoutcontrol.SlowQuery()
	// timeoutcontrol.TimerTest()
	// timeoutcontrol.ContextTest()

	// GC Trace
	gctracetest.GCTraceTest()
}

// 流程控制
func Print53() {
	for i := 1; i <= 100; i++ {
		is5 := (i%5 == 0)
		is3 := (i%3 == 0)
		if is5 && is3 {
			fmt.Println(i, "is multiple of 5 & 3")
		} else if is5 && !is3 {
			fmt.Println(i, "is multiple of 5")
		} else if is3 && !is5 {
			fmt.Println(i, "is multiple of 3")
		} else {
			fmt.Println(i)
		}
	}

}

// 字符串处理
func CalcNumOfChar(src *string) uint {
	return uint(utf8.RuneCountInString(*src))
}

func ReplaceCharsInString(src, des string, pos int) string {
	if len(src) == 0 { // src为空，直接将src改为des
		return des
	}
	if pos > len(src)+1 {
		return src + des
	}
	sumLen := max(pos+len(des), len(src))
	var sb strings.Builder
	sb.Grow(sumLen)
	for i := 0; i < pos-1; i++ {
		sb.WriteByte(src[i])
	}
	for i := 0; i < len(des); i++ {
		sb.WriteByte(des[i])
	}
	if pos+len(des) >= len(src) { // src长度不足
		return sb.String()
	}
	for i := pos + len(des) - 1; i < sumLen; i++ {
		sb.WriteByte(src[i])
	}
	return sb.String()
}

func ReverseStringbyPointer(src *string) string {
	var sb strings.Builder
	sb.Grow(len(*src))
	for i := len(*src); i > 0; i-- {
		sb.WriteByte((*src)[i-1])
	}
	return sb.String()
}

func ReverseString(src string) string {
	var sb strings.Builder
	sb.Grow(len(src))
	for i := len(src); i > 0; i-- {
		sb.WriteByte((src)[i-1])
	}
	return sb.String()
}

// 函数
func mapFunc(f func(int) int, vec []int) []int {
	res := make([]int, len(vec))
	for i, v := range vec {
		res[i] = f(v)
	}
	return res
}

func MapFuncTest() {
	double := func(i int) int {
		return i * 2
	}
	var vec = []int{1, 2, 3, 4, 5, 6, 7, 8}
	res := mapFunc(double, vec)
	fmt.Println(res)
}

// 附加题
