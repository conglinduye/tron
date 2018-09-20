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
	req.Limit = 5
	req.Start = 0

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
	req.Address = "TSNbzxac4WhxN91XvaUfPTKP2jNT18mP6T"

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

	resp, err := UpdateAccountSr(req, "token")
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

func TestGenWebToken(t *testing.T) {
	address := "RTFGHJK6GCVHHByui765CVBCVBVB"
	ss, err := GenWebToken(address)
	log.Printf("%v-%v", ss, err) //eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiUlRGR0hKSzZHQ1ZISEJ5dWk3NjVDVkJDVkJWQiJ9.zCZU-28PsQJ_GwLsZdBbGcT6ZrsTQmXsTpUDltdnfiM
}

func TestVerifyWebToken(t *testing.T) {
	address := "RTFGHJK6GCVHHByui765CVBCVBVB"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZGRyZXNzIjoiUlRGR0hKSzZHQ1ZISEJ5dWk3NjVDVkJDVkJWQiJ9.zCZU-28PsQJ_GwLsZdBbGcT6ZrsTQmXsTpUDltdnfiM"
	tt := VerifyWebToken(address, token)
	log.Printf("%v", tt) //true
}
