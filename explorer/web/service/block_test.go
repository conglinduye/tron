package service

import (
	"testing"

	"github.com/wlcy/tron/explorer/lib/mysql"

	"github.com/wlcy/tron/explorer/lib/log"

	"github.com/wlcy/tron/explorer/web/entity"
)

func Init() {
	mysql.Initialize("127.0.0.1", "3306", "tron", "budev", "tron**1")
}

func TestBlockList(t *testing.T) {
	Init()
	req := &entity.Blocks{}
	req.Sort = "-number"
	req.Limit = "5"
	req.Start = "0"

	//req.Number = "2287351"

	resp, err := QueryBlocks(req)
	if err != nil {
		log.Error(err)
	}
	log.Printf("total:%v", resp.Total)
	for _, value := range resp.Data {
		log.Printf("data:%#v", value)
	}

}

func TestBlock(t *testing.T) {
	Init()
	req := &entity.Blocks{}
	/*req.Sort = "-number"
	req.Limit = "5"
	req.Start = "0"
	*/
	req.Number = "2287351"

	resp, err := QueryBlock(req)
	if err != nil {
		log.Error(err)
	}
	log.Printf("total:%#v", resp)

}
