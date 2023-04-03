package controller

import (
	"fmt"
	"net/http"
	"webhook/config"
	"webhook/issuedb"
	"webhook/model"
	"webhook/msgqueue"

	"github.com/gin-gonic/gin"
)

type IssueId struct {
	Id string `uri:"id"`
}

func IssueView(context *gin.Context) {
	// 获取urlpath参数
	var id IssueId
	if err := context.ShouldBindUri(&id); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 根据url的id，查询告警内容
	query := issuedb.DbQuery{}
	query.DB = issuedb.DbConn.DB
	query.Table = "issues"
	query.Wherestring = fmt.Sprintf("issueId=\"%s\"", id.Id)
	query.Search()
	if len(query.Rows) != 0 {
		for _, row := range query.Rows {
			context.HTML(http.StatusOK, "issue.tmpl", row)
		}
	}
}

type FormBody struct {
	Handle     string `form:"handle"`
	HandleDesc string `form:"handledesc"`
	Status     string `form:"status"`
	IssueType  string `form:"issueType"`
	Owner      string `form:"owner"`
}

// 用于处理issue表单请求的数据处理
func IssueHandler(context *gin.Context) {
	// 获取urlpath参数
	var id IssueId
	if err := context.ShouldBindUri(&id); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 获取form表单的数据
	var form FormBody
	if err := context.Bind(&form); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 根据url的id，更新告警内容
	query := issuedb.DbQuery{}
	query.DB = issuedb.DbConn.DB
	query.Table = "issues"
	query.Update(form.Handle, form.HandleDesc, form.Status, form.IssueType, form.Owner, id.Id)
	context.String(http.StatusOK, "OK")
	// 加一个，处置完后，展示告警处置情况
	AfterHandle(id.Id)
}

// 告警处理后，在发送消息，展示处理情况
func AfterHandle(IssueId string) {
	// 根据url的id，查询告警内容
	query := issuedb.DbQuery{}
	query.DB = issuedb.DbConn.DB
	query.Table = "issues"
	query.Wherestring = fmt.Sprintf("issueId=\"%s\"", IssueId)
	query.Search()

	//构造webhook消息并发送
	for _, row := range query.Rows {
		// 根据form获取告警模型，这样可以知道告警该发到哪
		msg := model.GetWebhookByForm(row.Form)
		if msg != nil {
			//根据issueId构造访问url，添加到告警内容中
			note := fmt.Sprintf("**处置的告警详情(id=%s)如下:**", IssueId)
			note += fmt.Sprintf("\n>%s", row.Desc)
			note += "\n\n**处置情况如下:**"
			note += fmt.Sprintf("\n>**处置动作:**%s", row.Handle)
			note += fmt.Sprintf("\n>**处置记录:**%s", row.HandleDesc)
			note += fmt.Sprintf("\n>**告警状态:**%s", row.Status)
			note += fmt.Sprintf("\n\n点击[此处](%s/issues/%s)查看告警详情", config.Config.ServerSetting.BaseUrl, row.Id)
			msg.SetContent(note, "markdown")

			// 再添加消息队列中，如果队列满了，超时返回异常，不过异常也没关系，后面还有定时任务提醒未关闭的告警
			msgqueue.MsgQueue.Send(msg)
		}
	}
}
