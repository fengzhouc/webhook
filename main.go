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
	r.POST("/webhook/send/text", func(c *gin.Context) {
		content := c.PostForm("content")
		body := model.WxMsgTextModel{}
		body.MsgType = "text"
		body.Text.Content = content
		body.Text.MentionedList = append(body.Text.MentionedList, "周峰")
		err := body.Send()
		if err != nil {
			c.String(http.StatusBadRequest, "bad!!bad!!bad!!")
		}
		c.String(http.StatusOK, "OK!!")
	})
	r.POST("/webhook/send/md", func(c *gin.Context) {
		content := c.PostForm("content")
		body := model.WxMsgMdModel{}
		body.MsgType = "markdown"
		body.Markdown.Content = content
		body.Markdown.MentionedList = append(body.Markdown.MentionedList, "周峰")
		err := body.Send()
		if err != nil {
			c.String(http.StatusBadRequest, "bad!!bad!!bad!!")
		}
		c.String(http.StatusOK, "OK!!")
	})
	// 完全复刻企微webhook接口的基本参数
	r.POST("/webhook/send", func(c *gin.Context) {
		var msg model.WxMsgModel
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		err := msg.Send()
		if err != nil {
			c.String(http.StatusBadRequest, "bad!!bad!!bad!!")
		}
		c.String(http.StatusOK, "OK!!")
	})
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(":8000")
}
