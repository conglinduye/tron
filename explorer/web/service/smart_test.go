package service

import (
	"testing"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/web/entity"
)

func TestQueryContracts(t *testing.T) {
	Init()
	req := &entity.Contracts{}
	req.Sort = "-number"
	req.Limit = 5
	req.Start = 0
	//req.Number = "2287351"
	resp, err := QueryContracts(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}
