package model

// 模型的接口，未来其他类型的webhook实现该接口就可以了
type MsgModel interface {
	String() string
	Send() error
}
