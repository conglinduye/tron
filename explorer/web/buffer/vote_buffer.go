package buffer

import (
	"fmt"
	"sync"
	"time"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
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
		//_voteBuffer.loadQueryVoteLive()
		//_voteBuffer.getMaintenanceTimeStamp()
		//_voteBuffer.loadQueryVoteCurrentCycle()

		go voteLiveBufferLoader()
		go voteCycleBufferLoader()
	})
	return _voteBuffer
}

type voteBuffer struct {
	sync.RWMutex

	voteLive map[string]*entity.LiveInfo

	voteCurrentCycle *entity.VoteCurrentCycleResp

	nextMaintenanceTime int64
}

func voteLiveBufferLoader() {
	for {
		_voteBuffer.loadQueryVoteLive()
		time.Sleep(30 * time.Second)
	}
}
func voteCycleBufferLoader() {
	for {
		_voteBuffer.getMaintenanceTimeStamp()
		_voteBuffer.loadQueryVoteCurrentCycle()
		time.Sleep(60 * time.Second)
	}
}

func (w *voteBuffer) GetVoteLive() (voteLive map[string]*entity.LiveInfo, ok bool) {
	w.RLock()
	if len(w.voteLive) == 0 {
		log.Infof("get vote live info from buffer nil, data reload")
		w.loadQueryVoteLive()
	}
	log.Infof("get vote live info from buffer, buffer data updated ")
	voteLive = w.voteLive
	w.RUnlock()
	return
}

func (w *voteBuffer) GetVoteCurrentCycle() (voteCycle *entity.VoteCurrentCycleResp) {

	if w.voteCurrentCycle == nil {
		log.Infof("get VoteCurrentCycle info from buffer nil, data reload")
		w.loadQueryVoteCurrentCycle()
	}
	log.Infof("get VoteCurrentCycle info from buffer, buffer data updated ")
	return w.voteCurrentCycle
}

func (w *voteBuffer) GetNextMaintenanceTime() int64 {

	if w.nextMaintenanceTime == 0 {
		log.Infof("get NextMaintenanceTime info from buffer nil, data reload")
		w.getMaintenanceTimeStamp()
	}
	log.Infof("get NextMaintenanceTime info from buffer, buffer data updated ")
	return w.nextMaintenanceTime
}

func (w *voteBuffer) loadQueryVoteLive() { //QueryVoteLive()  实时投票数据
	strSQL := fmt.Sprintf(`
	SELECT acc.address as voteraddress,outvoter.votes,
	       acc.frozen,acc.account_name,wlwit.url
	FROM tron.tron_account acc 
	left join tron.wlcy_witness_create_info wlwit on wlwit.address=acc.address
	left join (
		select to_address,sum(vote) as votes from tron.account_vote_result 
		 group by to_address
	) outvoter on outvoter.to_address=acc.address
     where 1=1 and outvoter.votes>=0 
	 order by outvoter.votes desc `)

	liveInfo, err := module.QueryVoteLiveRealize(strSQL)
	if err != nil {
		log.Errorf("get vote live info from db err:[%v]", err)
		return
	}
	w.Lock()
	w.voteLive = liveInfo.Data
	w.Unlock()
}

func (w *voteBuffer) loadQueryVoteCurrentCycle() { //QueryVoteCurrentCycle()  上轮投票数据
	strSQL := fmt.Sprintf(`
	SELECT acc.address as voteraddress,outvoter.votes,
	acc.frozen,acc.account_name,wlwit.url,srcc.github_link
FROM tron.tron_account acc 
left join tron.wlcy_witness_create_info wlwit on wlwit.address=acc.address
left join tron.wlcy_sr_account srcc on srcc.address=acc.address
left join (
 select address,sum(vote_count) as votes from tron.witness 
  group by address
) outvoter on outvoter.address=acc.address
where 1=1 and outvoter.votes>=0  order by votes desc `)

	voteCurrent, err := module.QueryVoteCurrentCycleRealize(strSQL, "", "", "")
	if err != nil {
		log.Errorf("get last vote info from db err:[%v]", err)
		return
	}
	w.Lock()
	w.voteCurrentCycle = voteCurrent
	w.Unlock()
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
