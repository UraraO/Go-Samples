/*===========
 Author: UraraO Haru_UraraO@outlook.com
 Date: 2024-08-06 19:47:01
 LastEditors: UraraO Haru_UraraO@outlook.com
 LastEditTime: 2024-08-06 22:34:53
 FilePath: /Golang-Samples/src/_Module/bus_module/bus.go
 Description:

 消息总线模块，支持多接收者，多消息类型和消息内容，实现异步发送消息和异步接收消息

 Copyright (c) 2024 by ${git_name_email}, All Rights Reserved.
===========*/

package bus

import (
	"fmt"
	"log"
	"sync"
)

// 全局ID设置
// 此处使用锁进行并发控制，实际使用中可以改为snowflake等算法防止id冲突
var EventID = 0
var EventListenerID = 1

var eventIdMutex sync.Mutex
var listenerIdMutex sync.Mutex

// 事件类型预定义，实际使用中需注册到项目配置文件中
type EventType = int

const (
	EVENT_TYPE_ERROR = iota + 1
	EVENT_TYPE_FATAL
	EVENT_TYPE_WARN
	EVENT_TYPE_INFO
	EVENT_TYPE_DEBUG
)

// 事件类定义，事件为发送者向接收者（总线上其他方）发送的消息对象，包含消息类型，id，消息内容
type Event struct {
	Id      int
	Type    EventType
	Content string
}

func InitEvent(eventType EventType, content string) Event {
	eventIdMutex.Lock()
	defer eventIdMutex.Unlock()
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

// 事件监听器
// 即总线上的接收者，可以监听某一类型的消息，并注册handler
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

// 总线模块
// 维护一个监听器列表（接收者列表）
// 以及一个输入缓存（消息推送进该消息缓存，类似实际应用中的输入用消息队列）
type BusModule struct {
	Listeners   map[int]EventListener
	Cin         chan Event
	publishLock sync.Mutex
}

func InitBusModule() BusModule {
	return BusModule{
		Listeners: make(map[int]EventListener),
		Cin:       make(chan Event),
	}
}

// 发布消息，由发送方调用
func (bus *BusModule) Publish(event Event) {
	bus.publishLock.Lock()
	defer bus.publishLock.Unlock()
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
		bus.publishLock.Lock()
		for i, v := range bus.Listeners {
			if v.FollowingEvent == newEvent.Type {
				v.EventCh <- newEvent
				log.Default().Println("Publish, new Event ", newEvent.Id, " send to Listener ", i)
			}
		}
		bus.publishLock.Unlock()
	}
}
*/

// 注册监听者（接收者），提供需要监听的事件类型和handler
// 返回的监听器id可用于移除监听器，调用该函数的一方可以保存该id
func (bus *BusModule) RegisterEventListener(eventType EventType, handler Handler) int {
	listenerIdMutex.Lock()
	listener := EventListener{
		ID:             EventListenerID,
		FollowingEvent: eventType,
		Handler:        handler,
		EventCh:        make(chan Event),
	}
	EventListenerID++
	listenerIdMutex.Unlock()
	bus.Listeners[listener.ID] = listener
	go func(listener *EventListener) {
		listener.HandleLoop()
	}(&listener)

	log.Default().Println("RegisterEventListener, new Listener ID is: ", listener.ID)
	return listener.ID
}

// 移除监听器，id参数在注册时返回，由注册者保存
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

// 移除所有监听器
func (bus *BusModule) ClearEventListener() {
	for i, v := range bus.Listeners {
		close(v.EventCh)
		delete(bus.Listeners, i)
		log.Default().Println("RemoveEventListener, remove Listener ID is: ", v.ID)
	}
	log.Default().Println("ClearEventListener")
}
