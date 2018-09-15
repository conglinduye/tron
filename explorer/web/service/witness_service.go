package service

import (
	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/buffer"
	"github.com/wlcy/tron/explorer/web/entity"
)

//QueryWitnessBuffer 从QueryWitnessBuffer中获取witness信息
func QueryWitnessBuffer() ([]*entity.WitnessInfo, error) {

	witnessBuffer := buffer.GetWitnessBuffer()

	witnessList := witnessBuffer.GetWitness()

	return witnessList, nil
}

//QueryWitnessStatisticBuffer  从buffer获取
func QueryWitnessStatisticBuffer() ([]*entity.WitnessStatisticInfo, error) {
	witnessBuffer := buffer.GetWitnessBuffer()
	witnessInfo := witnessBuffer.GetWitnessStatistic()
	return witnessInfo, nil
}

//获取当前轮开始时间戳
func getMaintenanceTimeStamp() (int64, error) {

	client := grpcclient.GetRandomWallet()

	nextMaintenanceTime, err := client.GetNextMaintenanceTime()
	if err != nil {
		log.Error(err)
		return 0, err
	}
	curMaintenanceTime := nextMaintenanceTime - 6*60*60*1000 //6小时
	return curMaintenanceTime, nil
}
