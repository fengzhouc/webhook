package engine

import (
	"log"
	"time"
	"webhook/msgqueue"
)

// 处理发送限制的情况，另起线程进行消息的发送，比如企业微信每分钟20条的情况
// 优化：可不可以发送失败的就存到重试队列，然后发送下一个，因为不同的webhook可能限制不一样，一致等待一个就不太好
func Start() {
	retryQueue := msgqueue.NewMsgQueue(100)
	for {
		msg := msgqueue.MsgQueue.Pull(3)
		if msg != nil {
			// 正常发送webhook
			err := msg.Send()
			if err != nil {
				// 发送失败的话，传入重试队列
				log.Println("send failed!!! now send to retryqueue. msg: ", msg.String())
				err := retryQueue.Send(msg)
				// 添加失败。即重试消息队列满了，这时先取出一个处理，处理完后再添加该消息到重试消息队列
				if err != nil {
					retrymsg := retryQueue.Pull(3)
					if retrymsg != nil {
						log.Println("retryQueue was full, now to start retry. msg: ", retrymsg.String())
						err := retrymsg.Send()
						// 如果还发送失败，就再放回默认队列中，这样循环发送,这时需要睡眠10s
						if err != nil {
							// 重试还失败的话，可能限制还没过去，这时先睡眠10s
							time.Sleep(10 * time.Second)
							// 重试失败还是少数，除非webhook接口挂了，所以认为默认队列多数时间是空闲的
							log.Println("retry was failed!!! now send to default queue. msg: ", retrymsg.String())
							// 为了避免重试还失败，重试失败的话，就放到默认消息队列中，不能放回重试消息队列，不然又满了
							msgqueue.MsgQueue.Send(retrymsg)
						}
					}
					// 把失败的这次再添加进去，避免丢包
					err := retryQueue.Send(msg)
					if err != nil {
						log.Println("retryQueue was full, now send to detault. msg: ", retrymsg.String())
						msgqueue.MsgQueue.Send(retrymsg)
					} else {
						log.Println("send to retryQueue. msg: ", retrymsg.String())
					}
				}
			}
		} else {
			// 如果默认队列中取不到消息了，也就是空了，那就把重试队列中的消息处理掉
			for {
				msg := retryQueue.Pull(3)
				// 需要取到值，且默认消息队列为空，如果默认消息队列不为空了，需要立刻退出，优先处理默认消息队列的
				if msg != nil && msgqueue.MsgQueue.Size() == 0 {
					log.Println("start retry!!! ", msg.String())
					err := msg.Send()
					if err != nil {
						// 重试还失败的话，可能限制还没过去，这时先睡眠10s
						time.Sleep(10 * time.Second)
						// 发送失败就放回重试消息队列
						retryQueue.Send(msg)
					}
				} else {
					// 取不到就是空了，这时推出循环
					break
				}
			}
		}
	}
}
