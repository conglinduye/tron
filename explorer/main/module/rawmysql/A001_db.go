package rawmysql

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql" // import mysql
)

// GetMysqlDB 获取全局唯一的mysql db连接，使用rawmysql的DSN设置msyql的连接信息
func GetMysqlDB() *sql.DB {
	_once.Do(func() {
		var err error
		_dsn = DSN
		// _db, err = sql.Open("mysql", _dsn)
		// if nil != err {
		// 	panic(fmt.Sprintf("can't open mysql database using dsn:[%v], err:[%v]", _dsn, err))
		// }
		// err = _db.Ping()
		// if nil != err {
		// 	panic(fmt.Sprintf("ping db with dsn:[%v] failed:%v", _dsn, err))
		// }
		_db, err = initDB("mysql", _dsn)
		if nil != err {
			panic(fmt.Sprintf("can't open mysql database using dsn:[%v], err:[%v]", _dsn, err))
		}
	})
	return _db
}

// DSN ...
var DSN string
var _dsn string
var _db *sql.DB
var _once sync.Once

var _dbMap sync.Map

// GetDB ...
func GetDB(dsn string) *sql.DB {

	dbConnRaw, ok := _dbMap.Load(dsn)

	if ok && nil != dbConnRaw {
		dbConn, ok := dbConnRaw.(*sql.DB)
		if ok && nil != dbConn {
			return dbConn
		}
	}

	db, err := initDB("mysql", dsn)
	if nil != err && nil != db {
		_dbMap.Store(dsn, db)
	}
	return db
}

func initDB(driver string, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if nil != err || nil == db {
		return db, err
	}
	err = db.Ping()
	return db, err
}
