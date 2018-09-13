package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/util"
)

//日期时间的格式
const DATEFORMAT = "2006-01-02"
const DATETIMEFORMAT = "2006-01-02 15:04:05"
const DATETIMESTRING = "20060102150405"
const DATEFORMATDATE = "20060102"
const DATEFORMATHOUR = "2006010215"
const DATEFORMATMINUTE = "200601021504"

var dbHost = ""       //主机
var dbPort = "3306"   //端口
var dbSchema = "tron" //db schema
var dbName = "tron"   //用户名
var dbPass = "tron"   //密码

//数据库的连接配置
type DBParam struct {
	Mode         string
	ConnSQL      string
	MaxOpenconns int
	MaxIdleConns int
}

//连接DB的实例对象
var dbInstance *TronDB

//Initialize 初始化
// appInfo spaceInfo user report appType
// centerControl
func Initialize(host, port, schema, user, passwd string) bool {
	if len(strings.TrimSpace(host)) == 0 ||
		len(strings.TrimSpace(port)) == 0 ||
		len(strings.TrimSpace(schema)) == 0 ||
		len(strings.TrimSpace(user)) == 0 ||
		len(strings.TrimSpace(passwd)) == 0 {
		return false
	}

	dbHost = strings.TrimSpace(host)
	dbPort = strings.TrimSpace(port)
	dbSchema = strings.TrimSpace(schema)
	dbName = strings.TrimSpace(user)
	dbPass = strings.TrimSpace(passwd)
	return true
}

//GetDatabase Get一个连接的数据库对象
func GetDatabase() (*TronDB, error) {
	return retrieveDatabase()
}

//GetMysqlConnectionInfo 获取连接mysql的相关信息
func GetMysqlConnectionInfo() DBParam {
	dbConfig := DBParam{
		Mode: string("mysql"),
		//ConnSQL:      string("hub:blahblah@tcp(" + dbHost + ":" + dbPort + ")/hubDB?charset=utf8"),
		ConnSQL:      fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", dbName, dbPass, dbHost, dbPort, dbSchema),
		MaxOpenconns: 10,
		MaxIdleConns: 10,
	}
	return dbConfig
}

//RefreshDatabase 刷新DB的连接
func retrieveDatabase() (*TronDB, error) {
	defer CatchError()

	if nil == dbInstance {
		//连接数据库的参数
		para := GetMysqlConnectionInfo()

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
		dbInstance = dbPtr
	}

	//测试一下是否是连接成功的
	if err := dbInstance.Ping(); err != nil {
		//dbInstance.Close()
		dbInstance = nil
		return nil, err
	}

	return dbInstance, nil
}

// OpenDataBaseTransaction 开启一个数据库事物
func OpenDataBaseTransaction() (*sql.Tx, error) {
	dataPtr, err := GetDatabase()
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
func OpenDBPrepare(query string) (*sql.Tx, *sql.Stmt, error) {
	dataPtr, err := GetDatabase()
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

//QueryTableDataPages 按分页面大小，返回表的页面数量
func QueryTableDataPages(tableName string, pageSize int) int64 {
	pagesCount := int64(0)

	//判断输入的参数是否有效，tableName在下面的函数中判断了，这里就不判断了
	if pageSize < 1 {
		return pagesCount
	}

	//计算页面数量
	if count, err := QueryTableDataCount(tableName); err == nil {
		pagesCount = (count + int64(pageSize) - 1) / int64(pageSize)
	}
	return pagesCount
}

//QueryTableData 查询数据库数据
func QueryTableData(strSQL string) (*TronDBRows, error) {
	//获取数据库对象
	var dbPtr *TronDB
	var err error
	if dbPtr, err = GetDatabase(); err != nil {
		log.Errorf("get database error :[%v]\n", err)
		return nil, util.NewErrorMsg(util.Error_common_db_not_connected)
	}

	//查询数据集
	rows, err := dbPtr.Select(strSQL)
	if err != nil {
		log.Errorf("query database using:[\n%v\n] error:[%v]", strSQL, err)
		return rows, err
	}
	return rows, err
}

//QueryTablePageData 根据传入的SQL 执行分页查询
func QueryTablePageData(strSQL string, pageColumnName string, pageIndex int, pagesize int) (*TronDBRows, error) {
	strLimitSQL := fmt.Sprintf(" %s limit %d, %d ", pageColumnName, pageIndex*pagesize, pagesize)
	strSQL = fmt.Sprintf("select * from (%s)  %s", strSQL, strLimitSQL)

	//获取数据库对象
	var dbPtr *TronDB
	var err error
	if dbPtr, err = GetDatabase(); err != nil {
		return nil, util.NewError(util.Error_common_db_not_connected,
			util.GetErrorMsgSleek(util.Error_common_db_not_connected))
	}

	log.Debug(strSQL) //log select sql
	rows, err := dbPtr.Select(strSQL)
	if err != nil {
		log.Errorf("query database using:[\n%v\n] error:[%v]", strSQL, err)
		return rows, err
	}
	return rows, err
}

//QueryTableDataCount 返回某个表的记录个数
func QueryTableDataCount(tableName string) (int64, error) {
	rowCount := int64(0) //返回的数据库表行数

	//判断输入参数
	if len(tableName) < 0 {
		return 0, util.NewError(util.Error_common_parameter_invalid, util.GetErrorMsgSleek(util.Error_common_parameter_invalid))
	}

	//获取数据库对象
	var dbPtr *TronDB
	var err error
	if dbPtr, err = GetDatabase(); err != nil {
		return 0, util.NewError(util.Error_common_db_not_connected, util.GetErrorMsgSleek(util.Error_common_db_not_connected))
	}

	strSQL := "select count(*) as rowcounts from " + tableName
	log.Debug(strSQL)
	var data *TronDBRows
	if data, err = dbPtr.Select(strSQL); err != nil {
		return 0, util.NewError(util.Error_common_internal_error, util.GetErrorMsgSleek(util.Error_common_internal_error))
	}

	if data.NextT() {
		strValue := data.GetField("rowcounts")
		if count, err := strconv.ParseInt(strValue, 10, 64); err != nil {
			return 0, util.NewError(util.Error_common_internal_error, util.GetErrorMsgSleek(util.Error_common_internal_error))
		} else {
			rowCount = count //set the count
		}
	}

	return rowCount, nil
}

//QueryTableDataCount 返回某个表的记录个数
func QuerySQLViewCount(strSQLView string) (int64, error) {
	rowCount := int64(0) //返回的数据库表行数

	//判断输入参数
	if len(strSQLView) < 0 {
		return 0, util.NewError(util.Error_common_parameter_invalid, util.GetErrorMsgSleek(util.Error_common_parameter_invalid))
	}

	//获取数据库对象
	var dbPtr *TronDB
	var err error
	if dbPtr, err = GetDatabase(); err != nil {
		return 0, util.NewError(util.Error_common_db_not_connected, util.GetErrorMsgSleek(util.Error_common_db_not_connected))
	}

	strSQL := fmt.Sprintf(`
select count(*) as rowcounts from (
				%s
) newtableName`, strSQLView)
	log.Debug(strSQL)
	var data *TronDBRows
	if data, err = dbPtr.Select(strSQL); err != nil {
		return 0, util.NewError(util.Error_common_internal_error, util.GetErrorMsgSleek(util.Error_common_internal_error))
	}

	if data.NextT() {
		strValue := data.GetField("rowcounts")
		if count, err := strconv.ParseInt(strValue, 10, 64); err != nil {
			return 0, util.NewError(util.Error_common_internal_error, util.GetErrorMsgSleek(util.Error_common_internal_error))
		} else {
			rowCount = count //set the count
		}
	}

	return rowCount, nil
}

//ExecuteSQLCommand 执行insert update操作,依次返回 插入消息的主键，影响的条数，错误对象
func ExecuteSQLCommand(strSQL string, isInsertSQL bool) (int64, int64, error) {
	var key int64
	var rows int64
	var err error
	var dbPtr *TronDB

	//获取数据库对象
	if dbPtr, err = GetDatabase(); err != nil {
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
	if dbPtr, err = GetDatabase(); err != nil {
		return util.NewError(util.Error_common_db_not_connected, util.GetErrorMsgSleek(util.Error_common_db_not_connected))
	}

	return dbPtr.TransactionDB(sqls)
}

//GenSqlPartIn 获得SQL中 in 部分的SQL  like : "10,12,13,14" or "(10,12,13,14)"
func GenSqlPartIn(keys []uint64, withBrackets bool) string {
	var strRet string
	if nil != keys && len(keys) > 0 {
		for _, item := range keys {
			strRet = strRet + strconv.FormatUint(item, 10) + ","
		}
		strRet = strRet[:len(strRet)-1]

		if withBrackets {
			strRet = "(" + strRet + ")"
		}
	}
	return strRet
}

//GenSqlPartInStr 获得SQL中 in 部分的SQL  like : "10,12,13,14" or "(10,12,13,14)"
func GenSqlPartInStr(keys []string, withBrackets bool) string {
	var strRet string
	if nil != keys && len(keys) > 0 {
		for _, item := range keys {
			strRet = strRet + item + ","
		}
		strRet = strRet[:len(strRet)-1]

		if withBrackets {
			strRet = "(" + strRet + ")"
		}
	}
	return strRet
}

//GenSQLPartInlist 获取SQL中 连续的＝号，如：a.id=1 or a.id=2 or a.id=3
func GenSQLPartInlist(colName string, keys []uint64, withBrackets bool) string {
	var strSQL string
	if len(keys) > 0 {
		for i, item := range keys {
			if i == 0 {
				strSQL = strSQL + fmt.Sprintf(" %v=%v ", colName, item)
			} else {
				strSQL = strSQL + fmt.Sprintf(" or %v=%v ", colName, item)
			}
		}
	}
	if withBrackets && len(strSQL) > 0 {
		strSQL = fmt.Sprintf("(%v)", strSQL)
	}
	return strSQL
}

//GenSQLPartInStrList 获取SQL中 连续的＝号，如：a.id='1' or a.id='2' or a.id='3'
func GenSQLPartInStrList(colName string, keys []string, withBrackets bool) string {
	var strSQL string
	if len(keys) > 0 {
		for i, item := range keys {
			if i == 0 {
				strSQL = strSQL + fmt.Sprintf(" %v='%v' ", colName, item)
			} else {
				strSQL = strSQL + fmt.Sprintf(" or %v='%v' ", colName, item)
			}
		}
	}
	if withBrackets && len(strSQL) > 0 {
		strSQL = fmt.Sprintf("(%v)", strSQL)
	}
	return strSQL
}

//GetNextKey 返回某个表的下一个主键ID，适用于AUTO_INCREMENT字段,tableName 支持schema.table结构
func GetNextKey(schema string, tableName string) (uint64, error) {
	if len(tableName) == 0 {
		return 0, util.NewError(util.Error_common_internal_error, util.GetErrorMsgSleek(util.Error_common_internal_error))
	}

	var strSQL = fmt.Sprintf("SELECT AUTO_INCREMENT FROM information_schema.tables WHERE lower(table_name)=lower('%v')", tableName)
	if len(schema) > 0 {
		strSQL += fmt.Sprintf(" and lower(table_schema) = lower('%v')", schema)
	}

	dataPtr, err := QueryTableData(strSQL)
	if err != nil {
		return 0, err
	}

	if nil == dataPtr || dataPtr.ResNum() != 1 {
		return 0, util.NewError(util.Error_common_internal_error, util.GetErrorMsgSleek(util.Error_common_internal_error))
	}

	var nextID uint64
	for dataPtr.NextT() {
		var strNextID = dataPtr.GetField("AUTO_INCREMENT")
		if len(strNextID) == 0 {
			return 0, util.NewError(util.Error_common_internal_error, util.GetErrorMsgSleek(util.Error_common_internal_error))
		}

		nextID = util.Str2Uint64(strNextID) //get the next ID
		break
	}

	//为了避免下一个ID可能被其他人占用的情况，跳过一个ID。这只是一个不安全的策略
	//nextID ＝ nextID + 1

	return nextID, nil
}

// CatchError 捕获异常
func CatchError() {
	if err := recover(); err != nil {
		log.Errorf("recover error :%v", err) //标准log输出
	}
}

//dbObjectCheck 数据库对象有效性检查接口
type dbObjectCheck interface {
	//检查对象内部值是否有效,如果有效则返回true,否则返回false
	checkObjectValid() bool
	//检查对象内部值是否有溢出,如果有溢出则返回true,否则返回false
	checkObjectDataOverflow() bool
}

//JSONObjectToString 将jason对象转换为string
func JSONObjectToString(v interface{}) (string, error) {
	//检查参数是否有效
	if nil == v {
		return "", util.NewError(util.Error_common_json_object_nil, util.GetErrorMsgSleek(util.Error_common_json_object_nil))
	}

	var strJSONString string
	buffer, err := json.Marshal(v)
	if err == nil {
		strJSONString = string(buffer)
	}
	return strJSONString, err
}

//ConvertDBValueToString 返回字符串结果
func ConvertDBValueToString(colValue string) string {
	return colValue
}

//ConvertDBValueToInt64 返加int64结果
func ConvertDBValueToInt64(colValue string) int64 {
	retInt64 := int64(0)

	if len(colValue) > 0 {
		i64, err := strconv.ParseInt(colValue, 10, 64)
		if err == nil {
			retInt64 = i64
		}
	}
	return retInt64
}

//ConvertDBValueToUint64 返回uint64结果
func ConvertDBValueToUint64(colValue string) uint64 {
	retInt64 := uint64(0)

	if len(colValue) > 0 {
		i64, err := strconv.ParseUint(colValue, 10, 64)
		if err == nil {
			retInt64 = i64
		}
	}
	return retInt64
}

//ConvertDBValueToInt 返回int结果
func ConvertDBValueToInt(colValue string) int {
	retInt := int(0)

	if len(colValue) > 0 {
		i64, err := strconv.Atoi(colValue)
		if err == nil {
			retInt = i64
		}
	}
	return retInt
}

//ConvertDBValueToUint 返回int结果
func ConvertDBValueToUint(colValue string) uint {
	retInt := uint(0)

	if len(colValue) > 0 {
		ui, err := strconv.ParseUint(colValue, 10, 0)
		if err == nil {
			retInt = uint(ui)
		}
	}
	return retInt
}

//ConvertDBValueToFloat32 返回float32结果
func ConvertDBValueToFloat32(colValue string) float32 {
	retInt := float32(0)

	if len(colValue) > 0 {
		ui, err := strconv.ParseFloat(colValue, 32)
		if err == nil {
			retInt = float32(ui)
		}
	}
	return retInt
}

//ConvertDBValueToFloat64 返回float64结果
func ConvertDBValueToFloat64(colValue string) float64 {
	retInt := float64(0)

	if len(colValue) > 0 {
		ui, err := strconv.ParseFloat(colValue, 64)
		if err == nil {
			retInt = float64(ui)
		}
	}
	return retInt
}

//ConvertDBValueToBool 返回boolean结果
func ConvertDBValueToBool(colValue string) bool {
	retBool := false

	if len(colValue) > 0 {
		b, err := strconv.ParseBool(colValue)
		if err == nil {
			retBool = b
		}
	}
	return retBool
}

//ConvertDBValueToTime 返回时间结果
func ConvertDBValueToTime(colValue string) time.Time {
	var retTime time.Time

	if len(colValue) > 0 {
		t, err := time.ParseInLocation(DATETIMEFORMAT, colValue, time.Local) // 使用系统当前时区
		if err == nil {
			retTime = t
		}
	}
	return retTime
}

//ConvertStringToTime 将时间字符串转化为时间对象
func ConvertStringToTime(strTimeString string) (time.Time, error) {
	var retTime time.Time
	if len(strTimeString) <= 0 {
		return retTime, errors.New("时间字符串无效：字符串为空值")
	}

	var err error
	retTime, err = time.Parse(DATETIMEFORMAT, strTimeString)
	if err != nil {
		return retTime, errors.New("时间字符串无效：" + strTimeString)
	}
	return retTime, nil
}

//ConvertStringToBoolean 将字符串转化为Boolean类型
func ConvertStringToBoolean(booleanString string, defaultValue bool) bool {
	retBool := defaultValue

	if len(booleanString) > 0 {
		b, err := strconv.ParseBool(booleanString)
		if err == nil {
			retBool = b
		}
	}
	return retBool
}

//ConvertStringToFloat 将字符串转化为Float类型
func ConvertStringToFloat(val string, defaultValue float64) float64 {
	ret := defaultValue

	if len(val) > 0 {
		b, err := strconv.ParseFloat(val, 64)
		if err == nil {
			ret = b
		}
	}
	return ret
}

//ConvertStringToInt 将字符串转化为int类型
func ConvertStringToInt(intString string, defaultValue int) int {
	retInt := defaultValue

	if len(intString) > 0 {
		i64, err := strconv.Atoi(intString)
		if err == nil {
			retInt = i64
		}
	}
	return retInt
}

//ConvertStringToInt 将字符串转化为uint类型
func ConvertStringToUint(intString string, defaultValue uint) uint {
	retInt := defaultValue

	if len(intString) > 0 {
		i64, err := strconv.ParseUint(intString, 10, 64)
		if err == nil {
			retInt = uint(i64)
		}
	}
	return retInt
}

//ConvertStringToInt64 将字符串转化为int64类型
func ConvertStringToInt64(intString string, defaultValue int64) int64 {
	retInt := defaultValue

	if len(intString) > 0 {
		i64, err := strconv.ParseInt(intString, 10, 64)
		if err == nil {
			retInt = i64
		}
	}
	return retInt
}

//ConvertStringToInt32 将字符串转化为int32类型
func ConvertStringToInt32(intString string, defaultValue int32) int32 {
	retInt := defaultValue

	if len(intString) > 0 {
		i64, err := strconv.ParseInt(intString, 10, 32)
		if err == nil {
			retInt = int32(i64)
		}
	}
	return retInt
}

//ConvertStringToUint64 将字符串转化为uint64类型
func ConvertStringToUint64(intString string, defaultValue uint64) uint64 {
	retInt := defaultValue

	if len(intString) > 0 {
		i64, err := strconv.ParseUint(intString, 10, 64)
		if err == nil {
			retInt = i64
		}
	}
	return retInt
}

//FormatFloat64 格式化float的小位数,f为小数的位数
func FormatFloat64(value float64, f uint) float64 {
	var strFormat = fmt.Sprintf("%%.%vf", f)
	var strFloat = fmt.Sprintf(strFormat, value)
	//fmt.Println(strFloat)
	v, _ := strconv.ParseFloat(strFloat, 64)
	return v
}

//FormatFloat64ToString 格式化float的小位数,f为小数的位数
func FormatFloat64ToString(value float64, f uint) string {
	var strFormat = fmt.Sprintf("%%.%vf", f)
	var strFloat = fmt.Sprintf(strFormat, value)
	//fmt.Println(strFloat)
	return strFloat
}

//TimeToDbTimestamp 返回插入数据库的时间字符串
func TimeToDbTimestamp(t time.Time) string {
	return "str_to_date('" + t.Format(DATETIMEFORMAT) + "','%Y-%m-%d %H:%i:%s')"
}

//GetNowString 返回当前时间的格式化后的字符串
func GetNowString() string {
	return time.Now().Format(DATETIMEFORMAT)
}

//SplitStrAsUint64 将string分隔成uint64
func SplitStrAsUint64(strValue, split string) []uint64 {
	if len(strValue) == 0 {
		return nil
	}

	var ids = make([]uint64, 0)
	keys := strings.Split(strValue, ",")
	for _, item := range keys {
		if id, err := strconv.ParseUint(item, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

//SplitStrAsInt64 将string分隔成int64
func SplitStrAsInt64(strValue, split string) []int64 {
	if len(strValue) == 0 {
		return nil
	}

	var ids = make([]int64, 0)
	keys := strings.Split(strValue, ",")
	for _, item := range keys {
		if id, err := strconv.ParseInt(item, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

//StrArrayToUint64 将字符串数组转化为数字数组
func StrArrayToUint64(keys []string) []uint64 {
	if nil == keys {
		return nil
	}
	var ids = make([]uint64, 0)
	if len(keys) == 0 {
		return ids
	}

	for _, item := range keys {
		item = strings.TrimSpace(item)
		if id, err := strconv.ParseUint(item, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

//Uint64ArrayToString 将数字数组转化为字符串数组
func Uint64ArrayToString(keys []uint64) []string {
	if nil == keys {
		return nil
	}
	var ids = make([]string, 0)
	if len(keys) == 0 {
		return ids
	}

	for _, item := range keys {
		ids = append(ids, strconv.FormatUint(item, 10))
	}
	return ids
}

//ConvertBeginEndTimeDay 转化开始和结束时间,全部转化为日期，不包括小时分钟秒
func ConvertBeginEndTimeDay(begin, end string) (time.Time, time.Time, error) {
	var invalidTime = time.Unix(0, 0)

	if len(begin) < len(DATEFORMAT) || len(end) < len(DATEFORMAT) {
		return invalidTime, invalidTime, errors.New("beginTime or endTime invalid, please confirm format as 2006-01-02")
	}

	beginDay, err := time.Parse(DATEFORMAT, begin[:len(DATEFORMAT)])
	if err != nil {
		return invalidTime, invalidTime, err
	}
	endDay, err := time.Parse(DATEFORMAT, end[:len(DATEFORMAT)])
	if err != nil {
		return invalidTime, invalidTime, err
	}

	if beginDay != endDay && false == endDay.After(beginDay) {
		return invalidTime, invalidTime, errors.New("endTime less than beginTime")
	}
	return beginDay, endDay, nil
}

//yz++ 查找数组中是否存在元素
func ArrayFinder(keys []string, key string) bool {
	if nil == keys || len(keys) == 0 {
		return false
	}

	for _, item := range keys {
		item = strings.TrimSpace(item)
		if strings.EqualFold(item, key) {
			return true
		}
	}
	return false
}

//ArrayUint64ToString 将数字数组转化为字符串
func ArrayUint64ToString(ids []uint64, spliter string) string {
	var strString string
	if len(ids) > 0 {
		for i, item := range ids {
			if i == len(ids)-1 {
				strString = strString + strconv.FormatUint(item, 10)
			} else {
				strString = strString + strconv.FormatUint(item, 10) + spliter
			}
		}
	}
	return strString
}

//IsSameDay 是否是同一天,所支持的格式为"yyyy-mm-dd" 和 "yyyy-mm-dd hh:mm:ss"
func IsSameDay(strTime1, strTime2 string) (bool, error) {
	t1, err := str2Time(strTime1)
	if err != nil {
		return false, err
	}
	t2, err := str2Time(strTime2)
	if err != nil {
		return false, err
	}
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day(), nil
}

//str2Time 日期转化　"yyyy-mm-dd" 和 "yyyy-mm-dd hh:mm:ss"
func str2Time(strTime string) (time.Time, error) {
	if len(strTime) == len(DATEFORMAT) {
		return time.Parse(DATEFORMAT, strTime)
	}
	if len(strTime) == len(DATETIMEFORMAT) {
		return time.Parse(DATETIMEFORMAT, strTime)
	}
	return time.Now(), fmt.Errorf("time string :[%v] invalid. should like :[%v] or [%v]",
		strTime, DATEFORMAT, DATETIMEFORMAT)

}

//Str2Time 日期转化　"yyyy-mm-dd" 和 "yyyy-mm-dd hh:mm:ss"
func Str2Time(strTime string) (time.Time, error) {
	return str2Time(strTime)
}

//Str2Boolean 将字符串转化为boolean类型，支持的字符串包括 “true”,"y","yes"以及非“0”字符 及其大写字符
func Str2Boolean(strValue string) bool {
	strValue = strings.ToLower(strValue)

	//判断字符串
	if "true" == strValue || "yes" == strValue || "y" == strValue {
		return true
	}

	//判断数字
	if len(strValue) > 0 {
		if v, ok := strconv.Atoi(strValue); ok == nil {
			if v != 0 {
				return true
			}
		}
	}
	return false
}

//Distinct 清除重复的信息
func Distinct(value []string) (retData []string, distinct map[string]int) {
	retData = make([]string, 0, len(value)) //返回的去重后的信息
	distinct = make(map[string]int, len(value))

	for _, item := range value {
		if v, ok := distinct[item]; ok {
			v++
			distinct[item] = v
		} else {
			distinct[item] = 1
		}
	}

	for key := range distinct {
		retData = append(retData, key)
	}
	return retData, distinct
}

//SetDefaultVal 如果src为空，则返回defaultVal
func SetDefaultVal(src, defaultVal string) string {
	if src == "" {
		src = defaultVal
	}
	return src
}
