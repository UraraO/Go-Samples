package log

import "fmt"

func Log(prefix, content string) {
	fmt.Printf("%s:%s\n", prefix, content)
}
