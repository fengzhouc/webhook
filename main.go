package main

import (
	"log"
	"webhook/config"
	"webhook/controller"
	"webhook/engine"
	"webhook/listeners"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

func main() {
	// 1.创建路由
	r := gin.Default()
	// 另起线程进行消息的发送
	go engine.Start()
	// 未关闭告警的监控
	aspec := config.Config.CronSetting.ListenCron
	if aspec != "" {
		a := cron.New()
		aerr := a.AddFunc(aspec, func() {
			listeners.IssueListen()
		})
		if aerr != nil {
			log.Println("AddFunc error :", aerr)
			return
		}
		a.Start()
		defer a.Stop()
	}

	// 2.绑定路由规则，执行的函数
	// 完全复刻企微webhook接口的基本参数
	r.POST("/webhook/wx/send", controller.WxSend)
	// 完全复刻企微webhook接口的基本参数,但这个接口会记录为时间，需要闭环
	r.POST("/webhook/wx/issues/send", controller.WxIssueSend)
	// 事件详情页面（展示页面/修改接口）
	// -加载所有模版页面
	r.LoadHTMLGlob("template/*")
	// -配置路由规则
	r.GET("/issues/:id", controller.IssueView)
	r.POST("/issues/:id/handle", controller.IssueHandler)
	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	r.Run(config.Config.ServerSetting.Addr)
}
