package service

import (
	"testing"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/web/entity"
)

func TestTransactions(t *testing.T) {
	Init()
	req := &entity.Transactions{}
	req.Sort = "-number"
	req.Limit = 5
	req.Start = 0

	//req.Number = "2287351"

	resp, err := QueryTransactions(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)
	/*for _, value := range resp.Data {
		log.Printf("data:%#v", value)
	}*/

}

func TestTransaction(t *testing.T) {
	Init()
	req := &entity.Transactions{}
	/*req.Sort = "-number"
	req.Limit = "5"
	req.Start = "0"
	*/
	req.Hash = "086cd2282f698c0f72b6eb4b3eb880c2eb4a2bd8249c6ae644dc82f52b82490a"

	resp, err := QueryTransaction(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}
