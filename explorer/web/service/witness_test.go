package service

import (
	"testing"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
)

//true
func TestQueryWitness(t *testing.T) {
	Init()
	/*
		resp, err := QueryWitness()
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
func TestQueryWitnessStatistic(t *testing.T) {
	Init()
	/*
		resp, err := QueryWitnessStatistic()
		if err != nil {
			log.Error(err)
		}
		ss, _ := mysql.JSONObjectToString(resp)
		log.Printf("total:%v", ss)
	*/
}

//false
func TestgetMaintenanceTimeStamp(t *testing.T) {
	Init()

	resp, err := getMaintenanceTimeStamp()
	if err != nil {
		log.Error(err)
	}
	ss, _ := mysql.JSONObjectToString(resp)
	log.Printf("total:%v", ss)

}
