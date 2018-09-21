package buffer

import (
	"sync"
	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/lib/log"
)

/*
store all vote data in memory
load from db every 30 seconds
*/

var _voteBuffer *voteBuffer
var onceVoteBuffer sync.Once

//GetVoteBuffer ...
func GetVoteBuffer() *voteBuffer {
	return getVoteBuffer()
}

// getVoteBuffer
func getVoteBuffer() *voteBuffer {
	onceVoteBuffer.Do(func() {
		_voteBuffer = &voteBuffer{}
		_voteBuffer.getMaintenanceTimeStamp()
	})
	return _voteBuffer
}

type voteBuffer struct {
	sync.RWMutex
	nextMaintenanceTime int64
}


func (w *voteBuffer) GetNextMaintenanceTime() int64 {

	if w.nextMaintenanceTime == 0 {
		log.Infof("get NextMaintenanceTime info from buffer nil, data reload")
		w.getMaintenanceTimeStamp()
	}
	log.Infof("get NextMaintenanceTime info from buffer, buffer data updated ")
	return w.nextMaintenanceTime
}

//获取下轮开始时间戳
func (w *voteBuffer) getMaintenanceTimeStamp() {

	client := grpcclient.GetRandomWallet()

	nextMaintenanceTime, err := client.GetNextMaintenanceTime()
	if err != nil {
		log.Errorf("get maintenance timestamp from db err:[%v]", err)
		return
	}
	w.Lock()
	w.nextMaintenanceTime = nextMaintenanceTime
	w.Unlock()
}
