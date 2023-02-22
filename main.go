package main

import (
	"net/http"
	"webhook/model"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1.创建路由
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response

	// 完全复刻企微webhook接口的基本参数
	r.POST("/webhook/wx/send", func(c *gin.Context) {
		var msg model.WxMsgModel
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		err := msg.Send()
		if err != nil {
			c.String(http.StatusBadRequest, "bad!!bad!!bad!!")
			return
		}
		c.String(http.StatusOK, "OK!!")
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
