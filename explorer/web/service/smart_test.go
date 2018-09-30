package service

import (
	"testing"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mongo"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/web/entity"
)

func TestQueryContracts(t *testing.T) {
	Init()
	req := &entity.Contracts{}
	req.Sort = "-number"
	req.Limit = 5
	req.Start = 0
	resp, err := QueryContracts(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}

//fail can't get contract creation info due to contract_create_smart table
func TestQueryContractByAddress(t *testing.T) {
	Init()
	req := &entity.Contracts{}
	req.Sort = "-number"
	req.Limit = 5
	req.Start = 0
	req.Address = "TJzEofnjZ42khzk7bxZo7LyYXLQhs2ybX4"
	resp, err := QueryContractByAddress(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}

func TestQueryContractsCode(t *testing.T) {
	Init()
	req := &entity.Contracts{}
	req.Sort = "-number"
	req.Limit = 5
	req.Start = 0
	req.Address = "TJzEofnjZ42khzk7bxZo7LyYXLQhs2ybX4"
	resp, err := QueryContractsCode(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}

func TestQueryContractEvent(t *testing.T) {
	mongo.Initialize("47.90.203.178", "18890", "EventLogCenter", "root", "root")
	req := &entity.Contracts{}
	req.Sort = "-block_timestamp"
	req.Limit = 5
	req.Start = 0
	req.Address = "TKvAo1oYXKhmv5CogMQL5SMc3DTqW27FcD"
	resp, err := QueryContractEvent(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)
}
