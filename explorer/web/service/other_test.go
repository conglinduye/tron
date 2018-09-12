package service

import (
	"testing"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
)

//true
func TestQuerySystemStatus(t *testing.T) {
	Init()

	resp, err := QuerySystemStatus()
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}

//true
func TestQueryMarkets(t *testing.T) {
	Init()

	resp, err := QueryMarkets()
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}
