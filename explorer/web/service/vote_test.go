package service

import (
	"testing"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/web/entity"
)

//true
func TestVotes(t *testing.T) {
	Init()
	req := &entity.Votes{}
	req.Sort = "-votes"
	req.Limit = 5
	req.Start = 0

	//req.Candidate = "2287351"
	//req.Voter="Voter"

	resp, err := QueryVotes(req)
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)
	/*for _, value := range resp.Data {
		log.Printf("data:%#v", value)
	}*/

}

//true
func TestVoteLivefff(t *testing.T) {
	Init()

	//resp, err := QueryVoteLive()
	/*if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)*/

}

//true
func TestQueryVoteCurrentCycle(t *testing.T) {
	Init()

	/*resp, err := QueryVoteCurrentCycle()
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)
	*/
}

//true
func TestQueryVoteNextCycle(t *testing.T) {
	Init()

	resp, err := QueryVoteNextCycle()
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}
