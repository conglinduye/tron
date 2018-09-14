package service

import (
	"testing"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"

	"github.com/wlcy/tron/explorer/web/entity"
)

//true
func TestQueryAccounts(t *testing.T) {
	Init()
	req := &entity.Accounts{}
	req.Sort = "-number"
	req.Limit = "5"
	req.Start = "0"

	//req.Address = "2287351"

	resp, err := QueryAccounts(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}

//true//true
func TestQueryAccount(t *testing.T) {
	Init()
	req := &entity.Accounts{}
	/*req.Sort = "-number"
	req.Limit = "5"
	req.Start = "0"
	*/
	req.Address = "T9yDddzXNFeQyn3Eam1QcVzm85ekYaUkKz"

	resp, err := QueryAccount(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}

//true
func TestQueryAccountMedia(t *testing.T) {
	Init()
	req := &entity.Accounts{}
	/*req.Sort = "-number"
	req.Limit = "5"
	req.Start = "0"
	*/
	req.Address = "T9yDddzXNFeQyn3Eam1QcVzm85ekYaUkKz"

	resp, err := QueryAccountMedia(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}

func TestUpdateAccountSr(t *testing.T) {
	Init()
	req := &entity.SuperAccountInfo{}
	/*req.Sort = "-number"
	req.Limit = "5"
	req.Start = "0"
	*/
	req.Address = "2287351"
	req.GithubLink = "testurl1"

	resp, err := UpdateAccountSr(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}

func TestQueryAccountSr(t *testing.T) {
	Init()
	req := &entity.SuperAccountInfo{}
	/*req.Sort = "-number"
	req.Limit = "5"
	req.Start = "0"
	*/
	req.Address = "2287351"

	resp, err := QueryAccountSr(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}
