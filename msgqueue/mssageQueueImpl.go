package msgqueue

import "time"

type MessageQueue interface {
	// 添加消息到队列中
	Send(message interface{}) error
	// 从队列中获取指定数量的消息
	// Pull(size int, timeout time.Duration) []interface{}
	// 从队列中获取消息
	Pull(timeout time.Duration) interface{}
	// 返回队列的大小
	Size() int
	// 返回队列的最大容量
	Capacity() int
}
