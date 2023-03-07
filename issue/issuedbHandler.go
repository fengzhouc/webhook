package issue

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
	"webhook/config"

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
		fmt.Printf("Open sqlite3 failed,err:%v\n", err)
	}
	db.SetConnMaxLifetime(100 * time.Second)
	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(16)
	DbConnection.DB = db
}

//构建查询器，每次查询都创建对象，这样可以避免并发产生线程安全问题，使用构建器模式
type DbQuery struct {
	Wherestring  string
	Table        string
	DB           *sql.DB
	Rows         []RowModel
	sync.RWMutex // 读写锁
}

// 查询结果的结构体，包含全字段
type RowModel struct {
	Id         int64
	Desc       string
	Handle     string
	HandleDesc string
	Status     string
}

// 查询数据
func (query *DbQuery) Search() {
	sql := fmt.Sprintf("SELECT id,desc,handle,handleDesc,status FROM %s WHERE %s;", query.Table, query.Wherestring)
	rows, err := query.DB.Query(sql)
	if err != nil {
		fmt.Println("[SELECT error] ", err)
	} else {
		for rows.Next() {
			var row RowModel
			err := rows.Scan(&row.Id, &row.Desc, &row.Handle, &row.HandleDesc, &row.Status)
			if err == nil {
				query.Rows = append(query.Rows, row)
			} else {
				fmt.Println("[Scan error] ", err)
			}
		}
	}
}

// 插入数据，不需要检查是否已有,返回id
func (query *DbQuery) Insert(msg string) (id int64) {
	sql := fmt.Sprintf("INSERT INTO %s (\"desc\",\"status\",\"handle\",\"handleDesc\") VALUES (?,?,?,?);", query.Table)
	res, err := query.DB.Exec(sql, msg, "进行中", "", "")
	if err != nil {
		fmt.Println("[insert error] ", err)
	} else {
		id, err := res.LastInsertId()
		if err == nil {
			return id
		}
	}
	return -1
}

// 更新数据,用于
func (query *DbQuery) Update(id string, handle string, handleDesc string, status string) error {
	sql := fmt.Sprintf("UPDATE %s SET \"handle\"=?,\"handleDesc\"=?,\"status\"=?, \"update\"=CURRENT_TIMESTAMP WHERE id=%s", query.Table, id)
	_, err := query.DB.Exec(sql, handle, handleDesc, status)
	if err != nil {
		fmt.Println("[update error] ", err)
		return err
	}
	return nil
}
