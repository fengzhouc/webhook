package controller

import (
	"fmt"
	"net/http"
	"webhook/issue"

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
	query := issue.DbQuery{}
	query.DB = issue.DbConn.DB
	query.Table = "issues"
	query.Wherestring = fmt.Sprintf("id=\"%s\"", id.Id)
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
	fmt.Println(form)
	// 根据url的id，更新告警内容
	query := issue.DbQuery{}
	query.DB = issue.DbConn.DB
	query.Table = "issues"
	query.Update(id.Id, form.Handle, form.HandleDesc, form.Status)
	context.String(http.StatusOK, "OK")
}
