package service

import (
	"testing"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/web/entity"
)

func TestTransfers(t *testing.T) {
	Init()
	req := &entity.Transfers{}
	req.Sort = "-timestamp"
	req.Limit = 5
	req.Start = 0

	//req.Number = "2287351"

	resp, err := QueryTransfers(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)
	/*for _, value := range resp.Data {
		log.Printf("data:%#v", value)
	}*/

}

func TestTransfer(t *testing.T) {
	Init()
	req := &entity.Transfers{}
	/*req.Sort = "-number"
	req.Limit = "5"
	req.Start = "0"
	*/
	req.Hash = "0284c1ab70afb4fc11c68f6a83e22627798fb5be3f79e264fcbd80487e2a5d8a"

	resp, err := QueryTransfer(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}
