package controller

import (
	"fmt"
	"net/http"
	"webhook/issue"
	"webhook/model"
	"webhook/msgqueue"

	"github.com/gin-gonic/gin"
)

func WxIssueSend(c *gin.Context) {
	var msg model.WxMsgModel
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.String(http.StatusOK, err.Error())
	}
	// 添加事件入库
	query := issue.DbQuery{}
	query.DB = issue.DbConn.DB
	query.Table = "issues"
	if msg.MsgType == "text" {
		issueId := query.Insert(msg.Text.Content)
		//根据issueId构造访问url，添加到告警内容中
		click := fmt.Sprintf("\n\n点击[此处](%s/%d)闭环告警", "url", issueId)
		msg.Text.Content += click
	} else {
		issueId := query.Insert(msg.Markdown.Content)
		//根据issueId构造访问url，添加到告警内容中
		click := fmt.Sprintf("\n\n点击[此处](%s/%d)闭环告警", "http://url", issueId)
		msg.Text.Content += click
	}
	// 再添加消息队列中，如果队列满了，超时返回异常，不过异常也没关系，后面还有定时任务提醒未关闭的告警
	err := msgqueue.MsgQueue.Send(&msg)
	if err != nil {
		fmt.Println("消息队列已满,添加失败！！！")
		c.String(http.StatusOK, "Add Failed!! Because the message queue is full.")
	}
	c.String(http.StatusOK, "OK.")
}
