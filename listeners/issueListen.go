package listeners

import (
	"fmt"
	"webhook/config"
	"webhook/issuedb"
	"webhook/model"
	"webhook/msgqueue"
)

// 查询还没关闭的事件，重新发送告警，根据配置间隔发送
func IssueListen() {
	query := issuedb.DbQuery{}
	query.DB = issuedb.DbConn.DB
	query.Table = "issues"
	query.Wherestring = "status=\"进行中\""
	query.Search()
	for _, row := range query.Rows {
		//构造webhook消息并发送
		// 根据form获取告警模型，这样可以知道告警该发到哪
		msg := model.GetWebhookByForm(row.Form)
		//根据issueId构造访问url，添加到告警内容中
		note := fmt.Sprintf("**未处置的告警(%d),详情(id=%s)如下:**", len(query.Rows), row.Id)
		note += fmt.Sprintf("\n>%s", row.Desc)
		note += fmt.Sprintf("\n\n点击[此处](%s/issues/%s)闭环告警", config.Config.ServerSetting.BaseUrl, row.Id)
		msg.SetContent(note, row.Form)

		// 再添加消息队列中，如果队列满了，超时返回异常，不过异常也没关系，后面还有定时任务提醒未关闭的告警
		msgqueue.MsgQueue.Send(msg)
	}
}
