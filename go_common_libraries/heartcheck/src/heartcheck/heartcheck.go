package heartcheck

import (
	"reflect"
	"time"
)

type HeartChecker struct {
	ticker   *time.Ticker
	Dur      time.Duration
	cycleDur time.Duration
	// Times   int // 超时 和 超次数 选择实现
	Server  chan int
	Offline chan struct{}
	handler interface{}
}

func InitHeartChecker(dur time.Duration, handler interface{}) *HeartChecker {
	cycleDur := time.Second
	if dur < cycleDur {
		dur = cycleDur
	}
	return &HeartChecker{
		ticker:   time.NewTicker(1 * cycleDur),
		Dur:      dur,
		cycleDur: cycleDur,
		Server:   make(chan int, 10),
		Offline:  make(chan struct{}, 1),
		handler:  handler,
	}
}

func (hc *HeartChecker) Check() {
	begin := time.Now()
	tickTime := 1
	for range hc.ticker.C {
		if time.Now().After(begin.Add(hc.Dur)) { // 超时
			hc.Offline <- struct{}{}
			close(hc.Offline)
			hc.ticker.Stop()
			go reflect.ValueOf(hc.handler).Call([]reflect.Value{})
			close(hc.Server)
			return
		}
		// 未超时，向服务器发送
		hc.Server <- tickTime
		tickTime++
	}
}
