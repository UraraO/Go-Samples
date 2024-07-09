package Bus

import (
	"fmt"
	"log"
)

var EventID = 0
var EventListenerID = 1

type EventType = int

const (
	EVENT_TYPE_ERROR = iota + 1
	EVENT_TYPE_FATAL
	EVENT_TYPE_WARN
	EVENT_TYPE_INFO
	EVENT_TYPE_DEBUG
)

type Event struct {
	Id      int
	Type    EventType
	Content string
}

func InitEvent(eventType EventType, content string) Event {
	EventID++
	return Event{
		Id:      EventID,
		Type:    eventType,
		Content: content,
	}
}

type Handler interface {
	Handle(event Event)
}

type NormalHandler struct{}

func (hdl NormalHandler) Handle(event Event) {
	fmt.Println("Handle the event: ", event.Id, " -- ", event.Content)
}

type EventListener struct {
	ID             int
	FollowingEvent EventType
	Handler        Handler
	EventCh        chan Event
}

func (lsn *EventListener) HandleLoop() {
	for event := range lsn.EventCh {
		lsn.Handler.Handle(event)
	}
}

type BusInterface interface {
	Publish(event Event)                                        // 发布消息
	RegisterEventListener(eventType EventType, handler Handler) // 注册监听器
}

type BusModule struct {
	Listeners map[int]EventListener
	Cin       chan Event
}

func InitBusModule() BusModule {
	return BusModule{
		Listeners: make(map[int]EventListener),
		Cin:       make(chan Event),
	}
}

func (bus *BusModule) Publish(event Event) {
	for i, v := range bus.Listeners {
		if v.FollowingEvent == event.Type {
			v.EventCh <- event
			log.Default().Println("Publish, new Event ", event.Id, " send to Listener ", i)
		}
	}
}

/*
// 若后台运行主进程，比如go bus.Run(),则通过PublishBackground()发布消息,Cin通道即是为此准备
// 前台运行方式则无需使用Cin
func (bus *BusModule) PublishBackground(event Event) {
	bus.Cin <- event
}

func (bus *BusModule) Run() {
	for newEvent := range bus.Cin {
		for i, v := range bus.Listeners {
			if v.FollowingEvent == newEvent.Type {
				v.EventCh <- newEvent
				log.Default().Println("Publish, new Event ", newEvent.Id, " send to Listener ", i)
			}
		}
	}
}
*/

func (bus *BusModule) RegisterEventListener(eventType EventType, handler Handler) int {
	listener := EventListener{
		ID:             EventListenerID,
		FollowingEvent: eventType,
		Handler:        handler,
		EventCh:        make(chan Event),
	}
	EventListenerID++
	bus.Listeners[listener.ID] = listener
	go func(listener *EventListener) {
		listener.HandleLoop()
	}(&listener)

	log.Default().Println("RegisterEventListener, new Listener ID is: ", listener.ID)
	return listener.ID
}

func (bus *BusModule) RemoveEventListener(listenerID int) bool {
	v, ok := bus.Listeners[listenerID]
	if !ok {
		return false
	}
	close(v.EventCh)
	delete(bus.Listeners, listenerID)
	log.Default().Println("RemoveEventListener, remove Listener ID is: ", v.ID)
	return true
}
