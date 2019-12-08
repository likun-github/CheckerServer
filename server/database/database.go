package database

import (
	"database/sql"
	"fmt"
	"github.com/liangdas/mqant/log"
	"time"
)

var Db *sql.DB

func DbInit() bool {
	db, err := sql.Open("mysql", "root:yutao1994.YT@tcp(127.0.0.1:3306)/JustForFun?charset=utf8")
	if err != nil {
		log.Error("can not open db")
		return false
	}
	Db = db
	Db.SetMaxOpenConns(1000)
	Db.SetMaxIdleConns(100)
	Db.SetConnMaxLifetime(10 * time.Minute)
	Db.Ping()
	rows, err := Db.Query("select * from Test")
	if err != nil {
		log.Error("db error")
		return false
	}
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		log.Info(fmt.Sprint(record))
	}

	stmt, _ := Db.Prepare("select id,createTime,test from Test where id=?")
	row := stmt.QueryRow(1)
	var (
		id         int64
		createTime string
		test       string
	)
	row.Scan(&id, &createTime, &test)
	log.Info("id=%d, %s, %s", id, createTime, test)
	defer stmt.Close()

	defer rows.Close()
	log.Info("db success")
	return true
}
