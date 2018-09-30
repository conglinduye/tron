package mysql

/**
 * @author [yanzheng]
 * @email [yan_zheng2018@163.com@mail.com]
 * @create date 2018-09-22 12:11:50
 * @modify date 2018-09-22 12:11:50
 * @desc [写数据库实例，封装数据库写操作]
 */
import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
)

var dbWriteHost = ""       //主机
var dbWritePort = "3306"   //端口
var dbWriteSchema = "tron" //db schema
var dbWriteName = "tron"   //用户名
var dbWritePass = "tron"   //密码

/*//数据库的连接配置
type DBWrieteParam struct {
	Mode         string
	ConnSQL      string
	MaxOpenconns int
	MaxIdleConns int
}*/

//连接DB的实例对象
var dbWriteInstance *TronDB

//InitializeWriter 初始化
// appInfo spaceInfo user report appType
// centerControl
func InitializeWriter(host, port, schema, user, passwd string) bool {
	if len(strings.TrimSpace(host)) == 0 ||
		len(strings.TrimSpace(port)) == 0 ||
		len(strings.TrimSpace(schema)) == 0 ||
		len(strings.TrimSpace(user)) == 0 ||
		len(strings.TrimSpace(passwd)) == 0 {
		return false
	}

	dbWriteHost = strings.TrimSpace(host)
	dbWritePort = strings.TrimSpace(port)
	dbWriteSchema = strings.TrimSpace(schema)
	dbWriteName = strings.TrimSpace(user)
	dbWritePass = strings.TrimSpace(passwd)
	return true
}

//GetWriteDatabase Get一个连接的数据库对象
func GetWriteDatabase() (*TronDB, error) {
	return retrieveWriteDatabase()
}

//GetMysqlWriteConnectionInfo 获取连接mysql的相关信息
func GetMysqlWriteConnectionInfo() DBParam {
	dbConfig := DBParam{
		Mode:         string("mysql"),
		ConnSQL:      fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", dbWriteName, dbWritePass, dbWriteHost, dbWritePort, dbWriteSchema),
		MaxOpenconns: 10,
		MaxIdleConns: 10,
	}
	return dbConfig
}

//retrieveWriteDatabase 刷新DB的连接
func retrieveWriteDatabase() (*TronDB, error) {
	defer CatchError()

	if nil == dbWriteInstance {
		//连接数据库的参数
		para := GetMysqlWriteConnectionInfo()

		//打开这个DB对象
		dbPtr, err := OpenDB(para.Mode, para.ConnSQL)
		if err != nil {
			return nil, err
		}
		if dbPtr == nil {
			return nil, util.NewError(util.Error_common_db_not_connected, util.GetErrorMsgSleek(util.Error_common_db_not_connected))
		}

		//设置连接池信息
		dbPtr.SetConnsParam(para.MaxOpenconns, para.MaxIdleConns)
		dbWriteInstance = dbPtr
	}

	//测试一下是否是连接成功的
	if err := dbWriteInstance.Ping(); err != nil {
		dbWriteInstance = nil
		return nil, err
	}

	return dbWriteInstance, nil
}

// OpenDataBaseWriteTransaction 开启一个数据库事物
func OpenDataBaseWriteTransaction() (*sql.Tx, error) {
	dataPtr, err := GetWriteDatabase()
	if err != nil {
		return nil, err
	}

	sqlTx, err := dataPtr.Begin()
	if err != nil {
		return nil, err
	}
	return sqlTx, nil
}

// OpenDBPrepare statement
func OpenDBWritePrepare(query string) (*sql.Tx, *sql.Stmt, error) {
	dataPtr, err := GetWriteDatabase()
	if err != nil {
		return nil, nil, err
	}

	sqlTx, err := dataPtr.Begin()
	if err != nil {
		return nil, nil, err
	}

	stmt, err := sqlTx.Prepare(query)
	if err != nil {
		return nil, nil, err
	}

	return sqlTx, stmt, nil
}

//ExecuteSQLCommand 执行insert update操作,依次返回 插入消息的主键，影响的条数，错误对象
func ExecuteSQLCommand(strSQL string, isInsertSQL bool) (int64, int64, error) {
	var key int64
	var rows int64
	var err error
	var dbPtr *TronDB

	//获取数据库对象
	if dbPtr, err = GetWriteDatabase(); err != nil {
		return 0, 0, util.NewError(util.Error_common_db_not_connected, util.GetErrorMsgSleek(util.Error_common_db_not_connected))
	}

	//执行语句
	if isInsertSQL {
		if key, rows, err = dbPtr.Insert(strSQL); err != nil {
			log.Errorf("execute [%s] error [%s] ", strSQL, err)
			return key, rows, util.NewError(util.Error_common_internal_error, util.GetErrorMsgSleek(util.Error_common_internal_error)) //返回一个逻辑错误
		}
	} else {
		if key, rows, err = dbPtr.Update(strSQL); err != nil {
			log.Errorf("execute [%s] error [%s] ", strSQL, err)
			return key, rows, util.NewError(util.Error_common_internal_error, util.GetErrorMsgSleek(util.Error_common_internal_error)) //返回一个逻辑错误
		}
	}

	return key, rows, err
}

//ExecuteSQLCommands 批量执行SQL语句
func ExecuteSQLCommands(sqls []string) error {
	var err error
	var dbPtr *TronDB

	//获取数据库对象
	if dbPtr, err = GetWriteDatabase(); err != nil {
		return util.NewError(util.Error_common_db_not_connected, util.GetErrorMsgSleek(util.Error_common_db_not_connected))
	}

	return dbPtr.TransactionDB(sqls)
}

// ExecuteDeleteSQLCommand 执行delete语句
func ExecuteDeleteSQLCommand(strSQL string) (int64, int64, error) {
	var key int64
	var rows int64
	var err error
	var dbPtr *TronDB

	//获取数据库对象
	if dbPtr, err = GetWriteDatabase(); err != nil {
		return 0, 0, util.NewError(util.Error_common_db_not_connected, util.GetErrorMsgSleek(util.Error_common_db_not_connected))
	}

	//执行语句
	if key, rows, err = dbPtr.Delete(strSQL); err != nil {
		log.Errorf("execute [%s] error [%s] ", strSQL, err)
		return key, rows, util.NewError(util.Error_common_internal_error, util.GetErrorMsgSleek(util.Error_common_internal_error)) //返回一个逻辑错误
	}

	return key, rows, err
}
