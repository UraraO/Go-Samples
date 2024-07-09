package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	// fmt.Println("hello world")
	// Print53()
	// fmt.Println(utf8.RuneCountInString("hello, 123爱我国防"))
	str := "byDSWEW ggg inxudwmn xy"
	// fmt.Println(CalcNumOfChar(&str))
	// fmt.Println(ReverseString("foobar"))
	// fmt.Println(ReverseStringbyPointer(&str))
	fmt.Println(ReplaceCharsInString(str, "abc", 4))
}

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
