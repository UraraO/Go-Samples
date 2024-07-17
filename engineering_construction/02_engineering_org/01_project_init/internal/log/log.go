package log

import "fmt"

const logPrefix = "app"

func Log(content string) {
	fmt.Printf("%s:%s\n", logPrefix, content)
}
