package mq

import (
	"errors"
	"sync"
)

//实现一个基于内存的mq

var once *sync.Once
var memoryMqInstance *memoryMq //全局维护的消息队列

type memoryMq struct {
	queues sync.Map //key 为topic,每个topic单独一个队列,队列用chan来实现,chan *Message
}

func (m memoryMq) Consume(topic Topic) (*Message, error) {
	//从对应的topic中获取消息
	record, ok := m.queues.Load(topic)
	if !ok { //当前的topic不存在,则创建
		q := make(chan *Message, 10000)
		m.queues.Store(topic, q)
		record = q
	}
	queue, ok := record.(chan *Message)
	if !ok {
		return nil, errors.New("model's type is not chan *Message")
	}

	return <-queue, nil //取出一个值返回
}

func (m memoryMq) Produce(message *Message) error {
	record, ok := m.queues.Load(message.Topic())
	if !ok {
		q := make(chan *Message, 10000)
		m.queues.Store(message.Topic(), q)
		record = q
	}

	queue, ok := record.(chan *Message)
	if !ok {
		return errors.New("model's type is not chan *Message")
	}
	queue <- message
	return nil

}

//懒汉单例模式

func MemoryMqInstance() *memoryMq {
	once.Do(func() {
		memoryMqInstance = &memoryMq{queues: sync.Map{}}
	})
	return memoryMqInstance
}
