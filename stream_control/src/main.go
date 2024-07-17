package main

import (
	timeoutcontrol "stream_control/src/timeout_control"
	"time"
)

func main() {
	timeoutcontrol.TimerTest()
	time.Sleep(1 * time.Second)
	timeoutcontrol.ContextTest()
}
