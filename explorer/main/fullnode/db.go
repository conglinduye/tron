package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func getMysqlDB() *sql.DB {
	return dbb
}

var dbb *sql.DB

func initDB(dsn string) {
	var err error
	// dbb, err = sql.Open("mysql", "tron:tron@tcp(172.16.21.224:3306)/tron")
	dbb, err = sql.Open("mysql", dsn)
	if nil != err {
		panic(err)
	}
	err = dbb.Ping()
	if nil != err {
		panic(err)
	}
}
