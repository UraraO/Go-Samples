package msgqueue

import (
	"fmt"
	"sync"
	"time"
)

var mqID = 1
var csmID = 1
var mqIDmut sync.Mutex
var csmIDmut sync.Mutex

type Message struct {
	SetTime time.Time
	Topic   string
	Content string
}

func NewMessage(topic, content string) Message {
	return Message{
		SetTime: time.Now(),
		Topic:   topic,
		Content: content,
	}
}

type handler func(msg *Message)

type Option func(o *MsgQueue)

func WithCapacity(cap int) Option {
	return func(o *MsgQueue) {
		o.Cap = cap
	}
}

type MsgQueue struct {
	ID        int
	Cap       int
	Consumers map[string]map[int]*Consumer // topics - consumers[id]csm
	queue     chan Message
	regMut    sync.Mutex
}

func InitMsgQueue(opts ...Option) *MsgQueue {
	mqIDmut.Lock()
	mq := &MsgQueue{
		ID:        mqID,
		Cap:       0,
		Consumers: make(map[string]map[int]*Consumer, 10),
		regMut:    sync.Mutex{},
	}
	mqID++
	mqIDmut.Unlock()
	for _, opt := range opts {
		opt(mq)
	}
	if mq.Cap == 0 {
		mq.queue = make(chan Message)
	} else {
		mq.queue = make(chan Message, mq.Cap)
	}
	go mq.Run()
	return mq
}

// mq后台协程，处理生产者方输入的消息，转发给各消费者
func (mq *MsgQueue) Run() {
	for msg := range mq.queue {
		mq.regMut.Lock()
		csms, ok := mq.Consumers[msg.Topic]
		if !ok { // 当前topic未有consumer注册
			fmt.Println("topic:", msg.Topic, " have no consumer")
			continue
		}
		for id, csm := range csms {
			// fmt.Println("id:", csm.ID, " consumer receive a msg,", msg)
			// csm.consume(msg)
			// fmt.Println(msg, "has been sendin", csm)
			// go csm.sendin(msg)
			// 延迟删除 标记为deleted 的消费者
			if csm.deleted {
				close(csms[id].msgStream)
				delete(csms, id)
				continue
			}
			csm.sendin(msg)
		}
		mq.regMut.Unlock()
	}
}

func (mq *MsgQueue) Send(msg Message) {
	select {
	case mq.queue <- msg:
		fmt.Println("Send", msg, "to", "mq.queue")
	default:
		fmt.Printf("Send, 当队列被占满, 生产者不阻塞, 丢弃消息; msg: %v\n", msg)
	}
}

// 注册消费者，实际使用中需要调用方传入handler，将consume替换掉
func (mq *MsgQueue) RegisterConsumer(topic string, handler handler) int {
	csm := NewConsumer(topic, handler)
	mq.regMut.Lock()
	defer mq.regMut.Unlock()
	csms, ok := mq.Consumers[topic]
	if ok { // 当前topic已存在
		csms[csm.ID] = csm
		go csm.consume()
		fmt.Println("Register a consumer:", csm)
		return csm.ID
	}
	mq.Consumers[topic] = make(map[int]*Consumer, 10)
	mq.Consumers[topic][csm.ID] = csm
	go csm.consume()
	fmt.Println("Register a consumer:", csm)
	return csm.ID
}

// 移除消费者
// 并发问题：创建消费者后立即删除，可能发生极短时间内的消息Send触发send to closed channel
// 修正：延迟删除，以防止并发冲突
func (mq *MsgQueue) RemoveConsumer(topic string, id int) {
	mq.regMut.Lock()
	defer mq.regMut.Unlock()
	csms, ok := mq.Consumers[topic]
	if !ok { // 当前topic不存在
		fmt.Println("in topic:", topic, "id:", id, "consumer is not exist")
		return
	}
	csms[id].deleted = true
	// close(csms[id].msgStream)
	// delete(csms, id)
}

type Consumer struct {
	ID        int
	Topic     string
	msgStream chan Message
	handler   handler
	deleted   bool
}

func NewConsumer(topic string, handler handler) *Consumer {
	csmIDmut.Lock()
	csm := &Consumer{
		ID:        csmID,
		Topic:     topic,
		msgStream: make(chan Message, 1),
		handler:   handler,
		deleted:   false,
	}
	csmID++
	csmIDmut.Unlock()
	return csm
}

func (c *Consumer) sendin(msg Message) {
	c.msgStream <- msg
}

// 消费循环，注册消费者时启动协程循环运行
func (c *Consumer) consume() {
	for msg := range c.msgStream {
		// fmt.Println("Consumer", c.ID, "received message", msg)
		c.handler(&msg)
	}
}

func MQTest() {
	fmt.Println("MQTest begin")
	mq := InitMsgQueue(WithCapacity(10))
	topic := "hello?"
	content := "world!"
	hdl := func(msg *Message) {
		fmt.Println("received message", msg)
	}
	msg := NewMessage(topic, content)
	id1 := mq.RegisterConsumer(topic, hdl)
	id2 := mq.RegisterConsumer(topic, hdl)
	fmt.Println(id1, id2)
	mq.Send(msg)
	// time.Sleep(2 * time.Second)
	mq.RemoveConsumer(topic, id1)
	mq.RemoveConsumer(topic, id2)
	mq.Send(msg)
	time.Sleep(1 * time.Second)
}
