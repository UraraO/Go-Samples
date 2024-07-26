package concurrency

import (
	"reflect"
	"sync"
)

type StartPistol struct {
	mut   *sync.Mutex
	cond  *sync.Cond
	Fired bool
}

func InitStartPistol() *StartPistol {
	return &StartPistol{
		mut:   &sync.Mutex{},
		cond:  sync.NewCond(mut),
		Fired: false,
	}
}

func (sp *StartPistol) Wait(condition interface{}, args []reflect.Value) {
	if sp.Fired {
		return
	}
	for !reflect.ValueOf(condition).Call(args)[0].Bool() {
		sp.cond.Wait()
	}
}

func (sp *StartPistol) Start() {
	sp.cond.Broadcast()
	sp.Fired = true
}

func (sp *StartPistol) Reload() {
	sp.Fired = false
}
