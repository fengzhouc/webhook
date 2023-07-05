package controller

import (
	"log"
	"net/http"
	"webhook/model"
	"webhook/msgqueue"

	"github.com/gin-gonic/gin"
)

func WxSend(c *gin.Context) {
	var msg model.WxMsgModel
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.String(http.StatusOK, err.Error())
	}
	// 添加消息队列中，如果队列满了，超时返回异常
	err := msgqueue.MsgQueue.Send(&msg)
	if err != nil {
		log.Println("消息队列已满,添加失败！！！")
		c.String(http.StatusOK, "Add Failed!! Because the message queue is full.")
	}
}
