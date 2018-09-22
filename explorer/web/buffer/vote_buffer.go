package buffer

import (
	"fmt"
	"sync"
	"time"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
	"sort"
)

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
		go voteWitnessBufferLoader()
		go maintenanceTimeStampLoader()
	})
	return _voteBuffer
}

type voteBuffer struct {
	sync.RWMutex

	voteWitness *entity.VoteWitnessResp

	nextMaintenanceTime int64
}

func voteWitnessBufferLoader() {
	for {
		_voteBuffer.loadVoteWitness()
		time.Sleep(30 * time.Second)
	}
}

func maintenanceTimeStampLoader() {
	for {
		_voteBuffer.getMaintenanceTimeStamp()
		time.Sleep(60 * time.Second)
	}
}


func (w *voteBuffer) GetNextMaintenanceTime() int64 {
	if w.nextMaintenanceTime == 0 {
		log.Infof("get NextMaintenanceTime info from buffer nil, data reload")
		w.getMaintenanceTimeStamp()
	}
	log.Infof("get NextMaintenanceTime info from buffer, buffer data updated ")
	return w.nextMaintenanceTime
}

func (w *voteBuffer) GetVoteWitness() (voteWitness *entity.VoteWitnessResp) {
	w.RLock()
	voteWitness = w.voteWitness
	w.RUnlock()
	return
}

func (w *voteBuffer) loadVoteWitness() {
	var filterSQL, sortSQL, pageSQL string
	strSQL := fmt.Sprintf(`
		select witt.address, witt.vote_count, srac.github_link, acc.account_name,votes.realTimeVotes
		from witness witt
		left join tron_account acc on acc.address=witt.address
		left join wlcy_sr_account srac on witt.address=srac.address
		left join (
			select to_address,sum(vote) as realTimeVotes from account_vote_result  group by to_address 
		) votes on votes.to_address=witt.address
		where 1=1 `)

	sortSQL = "order by votes.realTimeVotes desc"

	voteWitnessResp, err := module.QueryVoteWitness(strSQL, filterSQL, sortSQL, pageSQL)
	if err != nil {
		log.Errorf("QueryVoteWitness strSQL:%v, err:[%v]",strSQL, err)
		return
	}

	totalVotes := module.QueryTotalVotes()
	voteWitnessResp.TotalVotes = totalVotes

	voteWitnessList:= voteWitnessResp.Data
	for index, voteWitness := range voteWitnessList {
		voteWitness.ChangeVotes = voteWitness.RealTimeVotes - voteWitness.LastCycleVotes
		if voteWitness.URL != "" {
			voteWitness.HasPage = true
		}
		if totalVotes != 0 {
			voteWitness.VotesPercentage = float64(voteWitness.LastCycleVotes) / float64(totalVotes) * 100
		}
		voteWitness.RealTimeRanking = int32(index + 1)
	}

	// getVoteWitnessRankingChange
	getVoteWitnessRankingChange(voteWitnessList)

	sortList := make([]*entity.VoteWitness, 0, len(voteWitnessList))
	for _, temp := range voteWitnessList {
		voteWitness := new(entity.VoteWitness)
		*voteWitness = *temp
		sortList = append(sortList, voteWitness)
	}

	if len(sortList) > 0 {
		sort.SliceStable(sortList, func(i, j int) bool { return sortList[i].ChangeCycle > sortList[j].ChangeCycle })
		voteWitnessResp.FastestRise = sortList[0]
	}

	w.Lock()
	w.voteWitness = voteWitnessResp
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

// getVoteWitnessRankingChange
func getVoteWitnessRankingChange(voteWitnessList []*entity.VoteWitness) {
	lastCycleSortList := make([]*entity.VoteWitness, 0, len(voteWitnessList))
	for _, temp := range voteWitnessList {
		voteWitness := new(entity.VoteWitness)
		*voteWitness = *temp
		lastCycleSortList = append(lastCycleSortList, voteWitness)
	}

	if len(lastCycleSortList) > 0 {
		sort.SliceStable(lastCycleSortList, func(i, j int) bool { return lastCycleSortList[i].LastCycleVotes > lastCycleSortList[j].LastCycleVotes })
	}

	for _, temp1 := range voteWitnessList {
		for index := range lastCycleSortList {
			temp2 := lastCycleSortList[index]
			if temp1.Address == temp2.Address {
				temp1.ChangeCycle = int32(index+1)-temp1.RealTimeRanking
			}
		}
	}
}

