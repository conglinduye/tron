package mongo

import (
	"fmt"
	"strings"
	"time"

	"github.com/wlcy/tron/explorer/lib/log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var MongoHost = "47.90.203.178"    //主机
var MongoPort = "18890"            //端口
var MongoSchema = "EventLogCenter" //db schema
var MongoName = "root"             //用户名
var MongoPass = "root"             //密码

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

	MongoHost = strings.TrimSpace(host)
	MongoPort = strings.TrimSpace(port)
	MongoSchema = strings.TrimSpace(schema)
	MongoName = strings.TrimSpace(user)
	MongoPass = strings.TrimSpace(passwd)
	return true
}

//Mongodb ...
type Mongodb struct {
	Standard   string
	Userid     string
	IsJob      bool
	IsGroupJob bool
	//	Job           *Job
	RstNum           int
	RstKeyArray      []string
	IsMailJobRst     bool
	MaxExecutionTime time.Duration
	//	Ctx           *ripple.Context
	Head   bson.M
	Output interface{}
	URL    string
}

const (
	TDD_URL                         = ""
	FDD_URL                         = ""
	VIP_USER_MAX_EXECUTION_TIME     = time.Second * 55
	GENERAL_USER_MAX_EXECUTION_TIME = time.Second * 30
	MAX_BULK_WRITE_FILE_LEN         = 1024
)

var (
	mongodb    *Mongodb
	MgoSession *mgo.Session
)

//GetMongodbInstance 获取Mongodb实例
func GetMongodbInstance() *Mongodb {
	if nil == mongodb {
		mongodb = new(Mongodb)
		mongodb.URL = fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", MongoName, MongoPass, MongoHost, MongoPort, MongoSchema)
		log.Debugf("conStr:[%v]", mongodb.URL)
	}
	return mongodb
}

//HandleError 错误处理
func (mgodb *Mongodb) HandleError(err error, parentInfo string, returnVal bson.M) {
	errInfo := err.Error()
	if strings.Contains(errInfo, "Sort operation used more than the maximum 33554432 bytes of RAM. Add an index, or specify a smaller limit.") {
		returnVal["statusCode"] = "500"
		returnVal["errorInfo"] = "本次操作结果大小超过32M，系统排序失败，请缩小 筛选范围 或 增加筛选条件 后再尝试此操作"
		mgodb.Output = returnVal
		log.Fatal("panic in location:", parentInfo)
		panic("Fatal error in " + parentInfo)
	} else if strings.Contains(errInfo, "operation exceeded time limit") || strings.Contains(errInfo, "i/o timeout") {
		returnVal["statusCode"] = "失败"
		returnVal["errorInfo"] = "目前系统繁忙，导致您本次操作耗时太长而失败，请稍后再尝试此操作"
		mgodb.Output = returnVal
		log.Fatal("panic in location:", parentInfo)
		panic("Fatal error in " + parentInfo)
	} else if strings.Contains(errInfo, "not found") {
		log.Fatal("error in location:", parentInfo)
		return
	} else {
		returnVal["statusCode"] = "失败"
		returnVal["errorInfo"] = "未知错误"
		mgodb.Output = returnVal
		log.Fatal("panic in location:", parentInfo)
		panic("Fatal error in " + parentInfo)
	}
}

//HandleMutiRoutineError ...
func (mgodb *Mongodb) HandleMutiRoutineError(err error, parentInfo string) {
	errInfo := err.Error()
	if strings.Contains(errInfo, "Sort operation used more than the maximum 33554432 bytes of RAM. Add an index, or specify a smaller limit.") {
		log.Fatal("panic in location:", parentInfo)
		panic(err)
	} else if strings.Contains(errInfo, "operation exceeded time limit") || strings.Contains(errInfo, "i/o timeout") {
		log.Fatal("panic in location:", parentInfo)
		panic(err)
	} else if strings.Contains(errInfo, "not found") {
		log.Fatal("error in location:", parentInfo)
		return
	} else {
		log.Fatal("panic in location:", parentInfo)
		panic(err)
	}
}

//Catch ...
func (mgodb *Mongodb) Catch(methodName string) {
	if r := recover(); nil != r {
		log.Println("At method:", methodName,
			", Runtime error caught:", r)
	}
}

//SetMaxExecutionTime ...
func (mgodb *Mongodb) SetMaxExecutionTime(userid string) {
	mgodb.Userid = userid
	if true == mgodb.IsJob {
		mgodb.MaxExecutionTime = time.Minute * 20
		return
	}
	mgodb.MaxExecutionTime = GENERAL_USER_MAX_EXECUTION_TIME
}

// GetSession 公共方法，获取session 如果存在则拷贝一份
func (mgodb *Mongodb) GetSession() *mgo.Session {
	if MgoSession == nil {
		var err error
		MgoSession, err = mgo.Dial(mgodb.URL)
		if err != nil {
			panic(err)
		}
		MgoSession.SetSocketTimeout(time.Minute * 20)
	}
	// 最大连接池默认4096
	return MgoSession.Clone()
}

//WithCollection 获取collection对象执行操作
func (mgodb *Mongodb) WithCollection(dataBase string, collection string, s func(*mgo.Collection) error) error {
	session := mgodb.GetSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}

//WithGridFS 获取GridFS对象执行操作
func (mgodb *Mongodb) WithGridFS(dataBase string, s func(*mgo.GridFS) error) error {
	session := mgodb.GetSession()
	defer session.Close()
	gridsFS := session.DB(dataBase).GridFS("fs")
	return s(gridsFS)
}

//GetOneRecord ...
func (mgodb *Mongodb) GetOneRecord(dataBase, collection string, queryCondition, selector bson.M) (result bson.M, err error) {
	query := func(c *mgo.Collection) error {
		return c.Find(queryCondition).SetMaxTime(mgodb.MaxExecutionTime).Select(selector).One(&result)
	}
	err = mgodb.WithCollection(dataBase, collection, query)
	if err != nil {
		log.Println("mongodb GetOneRecord, query failed! db:", dataBase,
			", Collection:", collection,
			", queryCondition:", queryCondition,
			", selector:", selector,
			", ErrInfo:", err)
	}
	return
}

//GetOneRecordWithSort ...
func (mgodb *Mongodb) GetOneRecordWithSort(dataBase, collection string, queryCondition, selector bson.M, fields ...string) (result bson.M, err error) {
	query := func(c *mgo.Collection) error {
		return c.Find(queryCondition).SetMaxTime(mgodb.MaxExecutionTime).Select(selector).Sort(fields...).One(&result)
	}
	err = mgodb.WithCollection(dataBase, collection, query)
	if err != nil {
		log.Println("mongodb GetOneRecordWithSort, query failed! db:", dataBase,
			", Collection:", collection,
			", queryCondition:", queryCondition,
			", selector:", selector,
			", sort:", fields,
			", ErrInfo:", err)
	}
	return
}

//GetMultiRecord ...
func (mgodb *Mongodb) GetMultiRecord(dataBase, collection string, queryCondition, selector bson.M) (result []bson.M, err error) {
	query := func(c *mgo.Collection) error {
		return c.Find(queryCondition).SetMaxTime(mgodb.MaxExecutionTime).Select(selector).All(&result)
	}
	err = mgodb.WithCollection(dataBase, collection, query)
	if err != nil {
		log.Println("mongodb GetMultiRecord, query failed! db:", dataBase,
			", Collection:", collection,
			", queryCondition:", queryCondition,
			", selector:", selector,
			", ErrInfo:", err)
	}
	return
}

//GetMultiRecordWithSort ...
func (mgodb *Mongodb) GetMultiRecordWithSort(dataBase, collection string, queryCondition, selector bson.M, fields ...string) (result []bson.M, err error) {
	query := func(c *mgo.Collection) error {
		return c.Find(queryCondition).SetMaxTime(mgodb.MaxExecutionTime).Select(selector).Sort(fields...).All(&result)
	}
	err = mgodb.WithCollection(dataBase, collection, query)
	if err != nil {
		log.Println("mongodb GetMultiRecordWithSort, query failed! db:", dataBase,
			", Collection:", collection,
			", queryCondition:", queryCondition,
			", selector:", selector,
			", sort:", fields,
			", ErrInfo:", err)
	}
	return
}

//DeleteCollection ...
func (mgodb *Mongodb) DeleteCollection(dataBase, collection string) (err error) {
	deleteOperetion := func(c *mgo.Collection) error {
		return c.DropCollection()
	}

	err = mgodb.WithCollection(dataBase, collection, deleteOperetion)
	if err != nil {
		log.Println("mongodb DeleteCollection, deleteOperation failed! db:", dataBase,
			", Collection:", collection,
			", ErrInfo:", err)
		return
	}
	return
}

//Count ...
func (mgodb *Mongodb) Count(dataBase, collection string, queryCondition bson.M) (total int, err error) {
	query := func(c *mgo.Collection) (err error) {
		total, err = c.Find(queryCondition).SetMaxTime(mgodb.MaxExecutionTime).Count()
		return
	}
	err = mgodb.WithCollection(dataBase, collection, query)
	if err != nil {
		log.Println("mongodb Count, count failed, db:", dataBase,
			", Collection:", collection,
			", QueryCondition:", queryCondition,
			", ErrInfo:", err)
	}
	return
}

//Pipe ...
func (mgodb *Mongodb) Pipe(dataBase, collection string, pipelineStage []bson.M) (result []interface{}, err error) {
	aggregate := func(c *mgo.Collection) error {
		return c.Pipe(pipelineStage).All(&result)
	}
	err = mgodb.WithCollection(dataBase, collection, aggregate)
	if err != nil {
		log.Println("mongodb Pipe, pipe failed, db:", dataBase,
			", Collection:", collection,
			", PipelineStage:", pipelineStage,
			", ErrInfo:", err)
	}
	return
}

//PipeAllowDiskUse ...
func (mgodb *Mongodb) PipeAllowDiskUse(dataBase, collection string, pipelineStage []bson.M) (result []interface{}, err error) {
	aggregate := func(c *mgo.Collection) error {
		return c.Pipe(pipelineStage).AllowDiskUse().All(&result)
	}
	err = mgodb.WithCollection(dataBase, collection, aggregate)
	if err != nil {
		log.Println("mongodb Pipe, PipeAllowDiskUse failed, db:", dataBase,
			", Collection:", collection,
			", PipelineStage:", pipelineStage,
			", ErrInfo:", err)
	}
	return
}

//Distinct ...
func (mgodb *Mongodb) Distinct(dataBase, collection, fieldName string, queryCondition bson.M) (out []interface{}, err error) {
	query := func(c *mgo.Collection) error {
		return c.Find(queryCondition).SetMaxTime(mgodb.MaxExecutionTime).Distinct(fieldName, &out)
	}
	err = mgodb.WithCollection(dataBase, collection, query)
	if err != nil {
		log.Println("mongodb Distinct, distinct failed, db:", dataBase,
			", Collection:", collection,
			", FieldName:", fieldName,
			", QueryCondition:", queryCondition,
			", ErrInfo:", err)
	}
	return
}

//UpdateOne ...
func (mgodb *Mongodb) UpdateOne(dataBase, collection string, selector interface{}, updater interface{}) (err error) {
	update := func(c *mgo.Collection) error {
		return c.Update(selector, updater)
	}
	err = mgodb.WithCollection(dataBase, collection, update)
	if err != nil {
		log.Println("mongodb UpdateOne, updateOne Failed, db:", dataBase,
			", Collection:", collection,
			", Selector:", selector,
			", Updater:", updater,
			", ErrInfo:", err)
	}
	return
}

//UpdateAll ...
func (mgodb *Mongodb) UpdateAll(dataBase, collection string, selector interface{}, updater interface{}) (err error) {
	update := func(c *mgo.Collection) error {
		_, err := c.UpdateAll(selector, updater)
		return err
	}
	err = mgodb.WithCollection(dataBase, collection, update)
	if err != nil {
		log.Println("mongodb UpdateAll, UpdateAll Failed, db:", dataBase,
			", Collection:", collection,
			", Selector:", selector,
			", Updater:", updater,
			", ErrInfo:", err)
	}
	return
}

//InsertOne ...
func (mgodb *Mongodb) InsertOne(dataBase, collection string, oneRecord interface{}) (err error) {
	insert := func(c *mgo.Collection) error {
		return c.Insert(oneRecord)
	}
	err = mgodb.WithCollection(dataBase, collection, insert)
	if err != nil {
		log.Println("mongodb InsertOne, InsertOne Failed, db:", dataBase,
			", Collection:", collection,
			", OneRecord:", oneRecord,
			", ErrInfo:", err)
		return
	}
	return
}

//Upsert ...
func (mgodb *Mongodb) Upsert(dataBase, collection string, selector interface{}, updater interface{}) (err error) {
	upsert := func(c *mgo.Collection) error {
		_, err := c.Upsert(selector, updater)
		return err
	}
	err = mgodb.WithCollection(dataBase, collection, upsert)
	if err != nil {
		log.Println("mongodb Upsert, Upsert Failed, db:", dataBase,
			", Collection:", collection,
			", Selector:", selector,
			", Updater:", updater,
			", ErrInfo:", err)
	}
	return
}

//BulkInsert ...
func (mgodb *Mongodb) BulkInsert(dataBase, collection string, allRecords []interface{}) (err error) {
	if 0 == len(allRecords) {
		log.Println("bulkInsert: 0 == len(allRecords)! db:", dataBase, "collection:", collection)
		return
	}
	bulkInsert := func(c *mgo.Collection) error {
		bulk := c.Bulk()
		bulk.Unordered()
		bulk.Insert(allRecords...)
		_, err := bulk.Run()
		return err
	}
	err = mgodb.WithCollection(dataBase, collection, bulkInsert)
	if err != nil {
		log.Println("mongodb BulkInsert, BulkInsert Failed, db:", dataBase,
			", Collection:", collection,
			", ErrInfo:", err)
	}
	return
}

//RemoveOne ...
func (mgodb *Mongodb) RemoveOne(dataBase, collection string, selector interface{}) (err error) {
	remove := func(c *mgo.Collection) error {
		return c.Remove(selector)
	}
	err = mgodb.WithCollection(dataBase, collection, remove)
	if err != nil {
		log.Println("mongodb RemoveOne, RemoveOne Failed, db:", dataBase,
			", Collection:", collection,
			", Selector:", selector,
			", ErrInfo:", err)
		return
	}
	return
}

//RemoveAll ...
func (mgodb *Mongodb) RemoveAll(dataBase, collection string, selector interface{}) (err error) {
	removeAll := func(c *mgo.Collection) (err error) {
		removeAllInfo, err := c.RemoveAll(selector)
		if nil != err {
			log.Println("mongodb RemoveAll, RemoveAll Failed, info:", *removeAllInfo)
		}
		return
	}
	err = mgodb.WithCollection(dataBase, collection, removeAll)
	if err != nil {
		log.Println("mongodb RemoveAll, RemoveAll Failed, db:", dataBase,
			", Collection:", collection,
			", Selector:", selector,
			", ErrInfo:", err)
		return
	}
	return
}

//InsertOneFile ...
func (mgodb *Mongodb) InsertOneFile(dataBase, fileName string, fileData []byte) (fileSize int, err error) {
	create := func(gridFS *mgo.GridFS) error {
		file, err := gridFS.Create(fileName)
		if nil != err {
			return err
		}
		fileSize, err = file.Write(fileData)
		if nil != err {
			return err
		}
		return file.Close()
	}
	err = mgodb.WithGridFS(dataBase, create)
	if nil != err {
		log.Println("mongodb GetOneFile, GetOneFile Failed, db:", dataBase,
			", FileName:", fileName,
			", ErrInfo:", err)
		return
	}
	return
}

//GetOneFileInfo ...
func (mgodb *Mongodb) GetOneFileInfo(dataBase string, queryCondition, selector bson.M) (result bson.M, err error) {
	query := func(gridFS *mgo.GridFS) error {
		return gridFS.Find(queryCondition).Select(selector).One(&result)
	}
	err = mgodb.WithGridFS(dataBase, query)
	if nil != err {
		log.Println("mongodb GetOneFile, GetOneFile Failed, db:", dataBase,
			", QueryCondition:", queryCondition,
			", Selector:", selector,
			", ErrInfo:", err)
		return
	}
	return
}

//GetMultiFileInfo ...
func (mgodb *Mongodb) GetMultiFileInfo(dataBase string, queryCondition, selector bson.M) (result []bson.M, err error) {
	query := func(gridFS *mgo.GridFS) error {
		return gridFS.Find(queryCondition).Select(selector).All(&result)
	}
	err = mgodb.WithGridFS(dataBase, query)
	if nil != err {
		log.Println("mongodb GetOneFile, GetOneFile Failed, db:", dataBase,
			", QueryCondition:", queryCondition,
			", Selector:", selector,
			", ErrInfo:", err)
		return
	}
	return
}

//OpenFileByName ...
func (mgodb *Mongodb) OpenFileByName(dataBase, fileName string, fileSize int) (result []byte, err error) {
	openFile := func(gridFS *mgo.GridFS) error {
		file, err := gridFS.Open(fileName)
		if nil != err {
			return err
		}
		result = make([]byte, fileSize)
		_, err = file.Read(result)
		if nil != err {
			return err
		}
		return file.Close()
	}
	err = mgodb.WithGridFS(dataBase, openFile)
	if nil != err {
		log.Println("mongodb OpenFile, OpenFile Failed, db:", dataBase,
			", FileName:", fileName,
			", FileSize:", fileSize,
			", ErrInfo:", err)
		return
	}
	return
}

//RemoveFileByName ...
func (mgodb *Mongodb) RemoveFileByName(dataBase, fileName string) (err error) {
	removeFile := func(gridFS *mgo.GridFS) error {
		return gridFS.Remove(fileName)
	}
	err = mgodb.WithGridFS(dataBase, removeFile)
	if nil != err {
		log.Println("mongodb RemoveFile, RemoveFile Failed, db:", dataBase,
			", FileName:", fileName,
			", ErrInfo:", err)
		return
	}
	return
}

//RemoveFileByID ...
func (mgodb *Mongodb) RemoveFileByID(dataBase string, id interface{}) (err error) {
	removeFile := func(gridFS *mgo.GridFS) error {
		return gridFS.RemoveId(id)
	}
	err = mgodb.WithGridFS(dataBase, removeFile)
	if nil != err {
		log.Println("mongodb RemoveFile, RemoveFile Failed, db:", dataBase,
			", Id:", id,
			", ErrInfo:", err)
		return
	}
	return
}
