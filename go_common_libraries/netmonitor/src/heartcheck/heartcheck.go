package heartcheck

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

type HeartChecker struct {
	ticker   *time.Ticker
	Dur      time.Duration
	cycleDur time.Duration
	// Times   int // 超时 和 超次数 选择实现
	Url       string
	Server    chan int
	Offline   chan struct{}
	handler   interface{}
	tickDatas *tickDataBlock
	stop      bool
	exit      bool
}

func InitHeartChecker(dur time.Duration, url string, handler interface{}) *HeartChecker {
	cycleDur := time.Second
	if dur < cycleDur {
		dur = cycleDur
	}
	hc := &HeartChecker{
		ticker:    nil,
		Dur:       dur,
		cycleDur:  cycleDur,
		Url:       url,
		Server:    make(chan int, 10),
		Offline:   make(chan struct{}, 1),
		handler:   handler,
		tickDatas: InittickDataBlock(url),
		stop:      false,
		exit:      false,
	}
	return hc
}

func InitHeartCheckerWithTimeout(ctx context.Context, timeout, dur time.Duration, url string, handler interface{}) (*HeartChecker, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	hc := InitHeartChecker(dur, url, handler)
	go func(ctx context.Context, hc *HeartChecker) {
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println(ctx.Err().Error())
		} else {
			fmt.Println("ctx has been canceled")
		}
		hc.Stop()
	}(ctx, hc)
	return hc, cancel
}

func (hc *HeartChecker) check() {
	begin := time.Now()
	tickTime := 1
	for {
		if hc.exit || hc.stop {
			return
		}
		<-hc.ticker.C // tick and ping
		if hc.exit || hc.stop {
			return
		}
		if time.Now().After(begin.Add(hc.Dur)) { // 超时 || 退出
			if len(hc.Offline) != 0 && !hc.stop && !hc.exit {
				<-hc.Offline
			}
			hc.Offline <- struct{}{}
			hc.Stop()
			return
		}
		// 未超时，向服务器发送
		hc.Server <- tickTime
		tickTime++
	}
}

func (hc *HeartChecker) serverRequesting() {
	for i := range hc.Server {
		if hc.exit {
			return
		}
		go hc.innerServerRequesting(i)
		// fmt.Println("serverRequesting request", i)
		// // ********************** tickdata ctor
		// now := time.Now()
		// resp, err := http.Get(hc.Url)
		// if err != nil {
		// 	fmt.Println("serverRequesting http.Get error:", err.Error())
		// 	resp.Body.Close()
		// 	continue
		// }
		// resp.Body.Close()
		// // ********************** tickdata add
		// hc.tickDatas.Add(now, time.Now(), true)
		// hc.Server <- 1
	}
}

func (hc *HeartChecker) innerServerRequesting(i int) {
	fmt.Println("serverRequesting request", i)
	// ********************** tickdata ctor
	now := time.Now()
	resp, err := http.Get(hc.Url)
	if err != nil {
		fmt.Println("serverRequesting http.Get error:", err.Error())
		hc.tickDatas.Add(now, time.Now(), false)
		resp.Body.Close()
		return
	}
	resp.Body.Close()
	// ********************** tickdata add
	hc.tickDatas.Add(now, time.Now(), true)
	if hc.exit || hc.stop {
		return
	}
}

func (hc *HeartChecker) StartBackground() {
	if hc.ticker == nil {
		hc.ticker = time.NewTicker(hc.cycleDur)
	}
	go hc.check()
	go hc.serverRequesting()
	hc.stop = false
	hc.ticker.Reset(hc.cycleDur)
}

func (hc *HeartChecker) StartFrontground() {
	if hc.ticker == nil {
		hc.ticker = time.NewTicker(hc.cycleDur)
	}
	go hc.serverRequesting()
	hc.stop = false
	hc.ticker.Reset(hc.cycleDur)
	hc.check()
}

func (hc *HeartChecker) Stop() {
	hc.stop = true
	hc.ticker.Stop()
}

func (hc *HeartChecker) Record(filename string) {
	go hc.tickDatas.PersistToFile(filename)
}

func (hc *HeartChecker) Close() {
	hc.exit = true
	close(hc.Offline)
	hc.ticker.Stop()
	go reflect.ValueOf(hc.handler).Call([]reflect.Value{})
	close(hc.Server)
}
