package listeners

import (
	"fmt"
	"webhook/issue"
	"webhook/model"
	"webhook/msgqueue"
)

// 查询还没关闭的事件，重新发送告警，根据配置间隔发送
func IssueListen() {
	query := issue.DbQuery{}
	query.DB = issue.DbConn.DB
	query.Table = "issues"
	query.Wherestring = "status=\"进行中\""
	query.Search()
	for _, row := range query.Rows {
		//构造webhook消息并发送
		// 告警内容附带快捷访问url（根据id拼接url）
		// msgqueue.MsgQueue.Send("")
		msg := model.WxMsgModel{}
		msg.MsgType = "markdown"
		//根据issueId构造访问url，添加到告警内容中
		note := "**未处置的告警,告警详情如下~**"
		note += fmt.Sprintf("\n>%s", row.Desc)
		note += fmt.Sprintf("\n\n点击[此处](%s/%d)闭环告警", "http://url", row.Id)
		msg.Markdown.Content = note

		// 再添加消息队列中，如果队列满了，超时返回异常，不过异常也没关系，后面还有定时任务提醒未关闭的告警
		msgqueue.MsgQueue.Send(&msg)
	}
}
