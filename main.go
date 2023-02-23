package main

import (
	"fmt"
	"net/http"
	"time"
	"webhook/model"

	"github.com/gin-gonic/gin"
)

// 处理发送限制的情况，另起线程进行消息的发送，比如企业微信每分钟20条的情况
// 这里会发送失败后，等待一定时间后重新发送
func engine(c chan model.MsgModel) {
	for msg := range c {
		for { // 死循环，知道发成功才进行下一条的发送
			err := msg.Send()
			if err != nil {
				// 发送失败，等待10s再重新尝试发送
				time.Sleep(10 * time.Second)
			} else {
				// 发送成功则跳出死循环，进行下一个消息的发送
				break
			}

		}
	}
}

func main() {
	// 1.创建路由
	r := gin.Default()
	// 消息队列
	channel := make(chan model.MsgModel, 50)
	// 处理企微每分钟20条的情况，另起线程进行消息的发送
	go engine(channel)
	// 2.绑定路由规则，执行的函数
	// 完全复刻企微webhook接口的基本参数
	r.POST("/webhook/wx/send", func(c *gin.Context) {
		var msg model.WxMsgModel
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		select {
		case channel <- &msg:
			fmt.Println("消息队列空闲,添加成功")
			c.String(http.StatusOK, "OK!!")
		case <-time.After(3 * time.Second):
			fmt.Println("消息队列已满,添加失败！！！")
			c.String(http.StatusOK, "Add Failed!! Because the message queue is full.")
		}
	})
	// 完全复刻企微webhook接口的基本参数,但这个接口会记录为时间，需要闭环
	r.POST("/webhook/wx/issues/send", func(c *gin.Context) {
		var msg model.WxMsgModel
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		// TODO: 将数据落地到本地数据库，而且需要起个线程，定时监测事件状态及时间
		err := msg.Send()
		if err != nil {
			c.String(http.StatusBadRequest, "bad!!bad!!bad!!")
			return
		}
		c.String(http.StatusOK, "OK!!")
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")
}
