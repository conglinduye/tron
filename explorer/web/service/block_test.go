package service

import (
	"testing"

	"github.com/wlcy/tron/explorer/lib/mysql"

	"github.com/wlcy/tron/explorer/lib/log"

	"github.com/wlcy/tron/explorer/web/entity"
)

func Init() {
	//mysql.InitializeReader("127.0.0.1", "3306", "tron", "budev", "tron**1")
	/*	blockBuffer := module.GetBlockBuffer()
		go blockBuffer.BackgroundWorker()*/
}

func TestBlockList(t *testing.T) {
	Init()
	req := &entity.Blocks{}
	req.Sort = "-number"
	req.Limit = 5
	req.Start = 0
	//req.Number = "2287351"

	//resp, err := QueryBlocks(req)
	resp, err := QueryBlocksBuffer(req)
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
