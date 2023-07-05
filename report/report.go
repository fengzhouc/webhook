package report

import (
	"fmt"
	"time"
	"webhook/issuedb"
)

// 报告模版
type ReportTemplete struct {
	Count      int        // 告警总数
	Trend      int        // 较上月的趋势,记录的是较上月的差值
	Proportion Proportion // 告警占比
}

// 占比结构体，格式：key 数量 百分比
type Proportion struct {
	TypeProp   []string // 告警类型top
	DepartProp []string // 部门告警top
}

// 将报告格式化为字符串格式
func (report *ReportTemplete) String() string {
	var note string
	note = "月度告警分析报告,详情如下:"
	note += fmt.Sprintf("\n**告警总数:** %d", report.Count)
	note += fmt.Sprintf("\n**较上月趋势:** %d", report.Trend)
	note += "\n**告警类型占比情况:**"
	for _, v := range report.Proportion.TypeProp {
		note += fmt.Sprintf("\n>%s", v)
	}
	note += "\n**部门告警占比情况:**"
	for _, v := range report.Proportion.DepartProp {
		note += fmt.Sprintf("\n>%s", v)
	}

	return note
}

// 按月整理报告
func (report *ReportTemplete) MonthReporter() {
	// 查询到当月的所有告警数据
	query := issuedb.DbQuery{}
	query.DB = issuedb.DbConn.DB
	query.Table = "issues"
	query.Wherestring = "STRFTIME('%m',\"update\") =" + fmt.Sprintf("'%s'", time.Now().Format("01"))
	query.Search()
	// TODO: 数据处理，会有各种查询及数据统计
	report.Count = len(query.Rows)
	report.CountHandler(&query)

}

// 计算较上月的数量情况
func (report *ReportTemplete) CountHandler(query *issuedb.DbQuery) {
	// 查询上个月的告警数量
	query.Wherestring = "\"update\"  between datetime('now','-1 month','start of month') and datetime('now','start of month')"
	query.Search()
	report.Trend = report.Count - len(query.Rows)
}

// 处理告警类型占比的数据
func (report *ReportTemplete) TypePropHandler(query *issuedb.DbQuery) {
	// 按告警类别统计数据，group by

}
