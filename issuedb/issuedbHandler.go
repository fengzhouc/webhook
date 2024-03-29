package issuedb

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
	"webhook/config"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DbConn *DbConnection
	once   = &sync.Once{} //保障线程安全
)

// 'import config'的时候就会调用,所以用来做初始化,所以可以不用调用GetInstance去获取config对象
func init() {
	DbConn = getInstance()
}

// 获取globalSetting对象，单例模式
func getInstance() *DbConnection {
	once.Do(func() {
		DbConn = &DbConnection{}
		DbConn.ConntSqlite3()
	})
	return DbConn
}

// 实现数据库连接及操作
// 1.获取db对象，单例模式，防止连接超过上限
// 2.db初始化，包含创建db文件，表初始化，重初始化
// 3.增删改查工作

//创建连接器，单例模式，提供一个操作句柄
type DbConnection struct {
	DB *sql.DB
}

//构建连接-sqlite3
func (DbConnection *DbConnection) ConntSqlite3() {
	db, err := sql.Open("sqlite3", config.Config.DbSetting.Sqitepath)
	if err != nil {
		log.Printf("Open sqlite3 failed,err:%v\n", err)
	}
	db.SetConnMaxLifetime(100 * time.Second)
	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(16)
	DbConnection.DB = db
}

//构建查询器，每次查询都创建对象，这样可以避免并发产生线程安全问题，使用构建器模式
type DbQuery struct {
	Select       string
	Wherestring  string
	Table        string
	DB           *sql.DB
	Rows         []RowModel
	sync.RWMutex // 读写锁
}

// 查询结果的结构体，包含全字段
type RowModel struct {
	Id         string
	Desc       string
	Handle     string
	HandleDesc string
	Status     string
	Form       string
	IssueType  string
	Owner      string
	Orgmsg     string
	Count      int
}

// 查询数据
func (query *DbQuery) Search() {
	sql := fmt.Sprintf("SELECT issueId,desc,handle,handleDesc,status,form,issueType,owner,orgmsg,count FROM %s WHERE %s;", query.Table, query.Wherestring)
	rows, err := query.DB.Query(sql)
	if err != nil {
		log.Println("[SELECT error] ", err)
	} else {
		for rows.Next() {
			var row RowModel
			err := rows.Scan(&row.Id, &row.Desc, &row.Handle, &row.HandleDesc, &row.Status, &row.Form, &row.IssueType, &row.Owner, &row.Orgmsg, &row.Count)
			if err == nil {
				query.Rows = append(query.Rows, row)
			} else {
				log.Println("[Scan error] ", err)
			}
		}
	}
}

// 插入数据，不需要检查是否已有,返回id
// msg: 内容
// form: 来自那个接口的，这个值会映射到配置文件中适配的webhook接口，也就是机器人列表
func (query *DbQuery) Insert(msg string, form string, orgmsg string) (issueId string) {
	query.Wherestring = fmt.Sprintf("desc=\"%s\"", msg)
	query.Search()
	// 没有添加过就直接添加
	if len(query.Rows) == 0 {
		sql := fmt.Sprintf("INSERT INTO %s (\"issueId\",\"desc\",\"status\",\"handle\",\"handleDesc\",\"form\",\"issueType\",\"owner\",\"orgmsg\",\"count\") VALUES (?,?,?,?,?,?,?,?,?,?);", query.Table)
		issueId = uuid.Must(uuid.NewV1()).String()
		_, err := query.DB.Exec(sql, issueId, msg, "进行中", "", "", form, "", "", orgmsg, 0)
		if err != nil {
			log.Println("[insert error] ", err)
		} else {
			return issueId
		}
	} else {
		// 添加过的就设置count= +1,然后返回已添加的id
		query.Update("count", query.Rows[0].Count+1)
		// 如果状态是关闭的，则重新打开
		if query.Rows[0].Status == "关闭" {
			query.Update("status", "进行中")
		}
		return query.Rows[0].Id
	}
	return "-1"
}

// 更新指定字段的数据
func (query *DbQuery) Update(key string, value interface{}) error {
	sql := fmt.Sprintf("UPDATE %s SET \"%s\"=?, \"update\"=DATETIME('now','localtime') WHERE %s", query.Table, key, query.Wherestring)
	_, err := query.DB.Exec(sql, value)
	if err != nil {
		log.Println("[update error] ", err)
		return err
	}
	return nil
}

// 更新数据,用于处置告警后的数据更新,传参要求：最后一个必须是issueId
func (query *DbQuery) IssueHandlerUpdate(msg ...interface{}) error {
	sql := fmt.Sprintf("UPDATE %s SET \"handle\"=?,\"handleDesc\"=?,\"status\"=?,\"issueType\"=?,\"owner\"=?, \"update\"=CURRENT_TIMESTAMP WHERE issueId=?", query.Table)
	_, err := query.DB.Exec(sql, msg...)
	if err != nil {
		log.Println("[update error] ", err)
		return err
	}
	return nil
}
