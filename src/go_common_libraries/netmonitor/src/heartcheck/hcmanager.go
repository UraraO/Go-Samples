package heartcheck

import (
	"fmt"
	"sync"
	"time"
)

type HCManager struct {
	hcs map[int]*HeartChecker
}

var HCID int = 1
var HCIDMut sync.Mutex

func InitHCManager() *HCManager {
	return &HCManager{
		hcs: make(map[int]*HeartChecker, 10),
	}
}

func (hcm *HCManager) AddHC(dur time.Duration, url string, handler interface{}) {
	HCIDMut.Lock()
	defer HCIDMut.Unlock()
	hc := InitHeartChecker(dur, url, handler)
	fmt.Printf("new heart checker ID is: %v\nUsing the ID to stop and close the heart checker\n", HCID)
	hcm.hcs[HCID] = hc
	HCID++
}

func (hcm *HCManager) RemoveHC(id int) {
	HCIDMut.Lock()
	defer HCIDMut.Unlock()
	// hc := InitHeartChecker(dur, url, handler)
	hc, ok := hcm.hcs[id]
	if !ok {
		fmt.Println("the hc is not exist, please check")
		return
	}
	hc.Close()
	fmt.Printf("heart checker ID: %v has beed closed\n", id)
	delete(hcm.hcs, id)
}

func (hcm *HCManager) StartAll() {
	HCIDMut.Lock()
	defer HCIDMut.Unlock()
	for k := range hcm.hcs {
		hcm.hcs[k].StartBackground()
	}
}

func (hcm *HCManager) StopAll() {
	HCIDMut.Lock()
	defer HCIDMut.Unlock()
	for k := range hcm.hcs {
		hcm.hcs[k].Stop()
	}
}

func (hcm *HCManager) ListAll() {
	HCIDMut.Lock()
	defer HCIDMut.Unlock()
	for k := range hcm.hcs {
		fmt.Println(hcm.hcs[k].Url)
	}
}

func (hcm *HCManager) ReportAll() {
	HCIDMut.Lock()
	defer HCIDMut.Unlock()
	for k := range hcm.hcs {
		hcm.hcs[k].Record("")
	}
}

func (hcm *HCManager) Clear() {
	HCIDMut.Lock()
	defer HCIDMut.Unlock()
	for k := range hcm.hcs {
		hcm.hcs[k].Close()
		delete(hcm.hcs, k)
	}
}

func (hcm *HCManager) Quit() {
	hcm.Clear()
	hcm.hcs = nil
}

var dur time.Duration = 5 * time.Second
var url = "https://www.baidu.com"
var handler = func() {
	fmt.Println("handler handling------")
}

func HCMTest() {
	hcm := InitHCManager()
	hcm.AddHC(dur, url, handler)
	hcm.ListAll()
	hcm.StartAll()
	time.Sleep(dur)
	hcm.StopAll()
	time.Sleep(time.Second)
	hcm.ReportAll()
	time.Sleep(time.Second)
	hcm.Quit()
}
