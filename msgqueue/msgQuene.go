package msgqueue

import (
	"errors"
	"sync"
	"time"
	"webhook/model"
)

var (
	MsgQueue *MyMessageQueue
	once     = &sync.Once{} //保障线程安全
)

type MyMessageQueue struct {
	queue    chan model.MsgModel
	capacity int
}

// 默认会产生一个全局的消息队列，如果需要新的，通过NewMsgQueue创建
func init() {
	getInstance()
}

// 获取消息队列对象，单例模式，默认容量100
func getInstance() *MyMessageQueue {
	once.Do(func() {
		MsgQueue = NewMsgQueue(100)
	})
	return MsgQueue
}

// 创建新的消息队列
func NewMsgQueue(capacity int) *MyMessageQueue {
	mq := &MyMessageQueue{
		queue:    make(chan model.MsgModel, capacity),
		capacity: capacity,
	}
	return mq
}

// 向队列中添加消息
func (mq *MyMessageQueue) Send(message model.MsgModel) error {
	// 添加消息队列中，如果队列满了，超时返回异常
	select {
	case mq.queue <- message:
	case <-time.After(3 * time.Second):
		return errors.New("the queue is full,")
	}
	return nil
}

// 从队列中获取消息
func (mq *MyMessageQueue) Pull(timeout time.Duration) model.MsgModel {
	select {
	case msg := <-mq.queue:
		return msg
	case <-time.After(timeout):
		return nil
	}
}

// 消息队列当前的大小
func (mq *MyMessageQueue) Size() int {
	return len(mq.queue)
}

// 消息队列的容量
func (mq *MyMessageQueue) Capacity() int {
	return mq.capacity
}
