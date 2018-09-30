package service

import (
	"testing"

	"github.com/wlcy/tron/explorer/lib/mysql"

	"github.com/wlcy/tron/explorer/lib/log"

	"github.com/wlcy/tron/explorer/web/entity"
)

func Init() {
	var dbParams = make(map[string]map[string]string, 0)
	var dbConns = []string{"primary", "secondary"}
	//init read database
	for _, db := range dbConns {
		params := make(map[string]string, 0)
		params[mysql.DBHost] = "127.0.0.1"
		params[mysql.DBPort] = "3306"
		params[mysql.DBSchema] = "tron_test_net"
		params[mysql.DBName] = "budev"
		params[mysql.DBPass] = "tron**1"
		dbParams[db] = params
		log.Debugf("read database init:db:[%v],param:[%v]", db, params)
	}
	mysql.InitializeReader(dbParams, dbConns)
	mysql.InitializeWriter("127.0.0.1", "3306", "tron_test_net", "budev", "tron**1")
}

func TestBlockList(t *testing.T) {
	Init()
	req := &entity.Blocks{}
	req.Sort = "-number"
	req.Limit = 5
	req.Start = 0
	//req.Number = "2287351"

	resp, err := QueryBlocks(req)
	//resp, err := QueryBlocksBuffer(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}

func TestBlock(t *testing.T) {
	Init()
	req := &entity.Blocks{}
	/*req.Sort = "-number"
	req.Limit = "5"
	req.Start = "0"
	*/
	req.Number = "2287351"

	resp, _ := QueryBlock(req)
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)
}
