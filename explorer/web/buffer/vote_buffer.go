package buffer

import (
	"fmt"
	"sync"
	"time"

	"sort"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
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

	voteCurrentCycle *entity.VoteCurrentCycleResp

	voteLive *entity.VoteLiveResp
}

func voteWitnessBufferLoader() {
	for {
		_voteBuffer.loadVoteWitness()
		_voteBuffer.loadVoteCurrentCycle()
		time.Sleep(30 * time.Second)
	}
}

func maintenanceTimeStampLoader() {
	for {
		_voteBuffer.getMaintenanceTimeStamp()
		time.Sleep(10 * time.Second)
	}
}

func (w *voteBuffer) GetNextMaintenanceTime() int64 {
	if w.nextMaintenanceTime == 0 {
		log.Debugf("get NextMaintenanceTime info from buffer nil, data reload")
		w.getMaintenanceTimeStamp()
	}
	log.Debugf("get NextMaintenanceTime info from buffer, buffer data updated ")
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
		log.Errorf("QueryVoteWitness strSQL:%v, err:[%v]", strSQL, err)
		return
	}

	totalVotes := module.QueryTotalVotes()
	voteWitnessResp.TotalVotes = totalVotes

	voteWitnessList := voteWitnessResp.Data
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

	// getVoteWitnessFastestRise
	voteWitnessResp.FastestRise = getVoteWitnessFastestRise(voteWitnessList)

	w.Lock()
	w.voteWitness = voteWitnessResp
	w.Unlock()

}

//获取下轮开始时间戳
func (w *voteBuffer) getMaintenanceTimeStamp() {

	client := grpcclient.GetRandomWallet()
	if nil != client {
		defer client.Close()
	}

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
				temp1.ChangeCycle = int32(index+1) - temp1.RealTimeRanking
				break
			}
		}
	}
}

// getVoteWitnessFastestRise
func getVoteWitnessFastestRise(voteWitnessList []*entity.VoteWitness) *entity.VoteWitness {
	fastestRise := &entity.VoteWitness{}
	sortList := make([]*entity.VoteWitness, 0, len(voteWitnessList))
	for _, temp := range voteWitnessList {
		voteWitness := new(entity.VoteWitness)
		*voteWitness = *temp
		sortList = append(sortList, voteWitness)
	}
	if len(sortList) > 0 {
		sort.SliceStable(sortList, func(i, j int) bool { return sortList[i].ChangeCycle > sortList[j].ChangeCycle })
		fastestRise = sortList[0]
	}
	return fastestRise
}

func (w *voteBuffer) GetVoteCurrentCycle() (voteCurrentCycle *entity.VoteCurrentCycleResp) {
	w.RLock()
	voteCurrentCycle = w.voteCurrentCycle
	w.RUnlock()
	return
}

func (w *voteBuffer) loadVoteCurrentCycle() {
	var filterSQL, sortSQL, pageSQL string
	voteCurrentCycleResp := &entity.VoteCurrentCycleResp{}
	strSQL := fmt.Sprintf(`
		select witt.address, witt.vote_count, witt.url, acc.account_name,votes.realTimeVotes
		from witness witt
		left join tron_account acc on acc.address=witt.address
		left join (
			select to_address,sum(vote) as realTimeVotes from account_vote_result  group by to_address 
		) votes on votes.to_address=witt.address
		where 1=1 `)

	sortSQL = "order by votes.realTimeVotes desc"

	voteCurrentCycleList, err := module.QueryVoteCurrentCycle(strSQL, filterSQL, sortSQL, pageSQL)
	if err != nil {
		log.Errorf("loadVoteCurrentCycle strSQL:%v, err:[%v]", strSQL, err)
		return
	}

	for _, voteCurrentCycle := range voteCurrentCycleList {
		if voteCurrentCycle.URL != "" {
			voteCurrentCycle.HasPage = true
		}
	}

	getVoteCurrentCycleRankingChange(voteCurrentCycleList)

	voteCurrentCycleResp.Candidates = voteCurrentCycleList

	totalVotes := module.QueryTotalVotes()
	voteCurrentCycleResp.TotalVotes = totalVotes

	w.Lock()
	w.voteCurrentCycle = voteCurrentCycleResp
	w.Unlock()
}

func getVoteCurrentCycleRankingChange(voteCurrentCycleList []*entity.VoteCurrentCycle) {
	lastCycleSortList := make([]*entity.VoteCurrentCycle, 0, len(voteCurrentCycleList))
	for _, temp := range voteCurrentCycleList {
		voteCurrentCycle := new(entity.VoteCurrentCycle)
		*voteCurrentCycle = *temp
		lastCycleSortList = append(lastCycleSortList, voteCurrentCycle)
	}

	if len(lastCycleSortList) > 0 {
		sort.SliceStable(lastCycleSortList, func(i, j int) bool { return lastCycleSortList[i].Votes > lastCycleSortList[j].Votes })
	}

	for index1 := range voteCurrentCycleList {
		temp1 := voteCurrentCycleList[index1]
		for index2 := range lastCycleSortList {
			temp2 := lastCycleSortList[index2]
			if temp1.Address == temp2.Address {
				temp1.ChangeCycle = int32(index2+1) - int32(index1+1)
				temp1.ChangeDay = 0
				break
			}
		}
	}
}

func (w *voteBuffer) GetVoteLive() (voteLive *entity.VoteLiveResp) {
	w.RLock()
	voteLive = w.voteLive
	w.RUnlock()
	return
}


// QueryVoteLive
func (w *voteBuffer) loadVoteLive() () {
	voteLiveResp := &entity.VoteLiveResp{}
	data := make(map[string]*entity.VoteLive, 0)
	strSQL := fmt.Sprintf(`select to_address as address, sum(vote) as totalVote from account_vote_result group by to_address`)
	VoteLiveList, err := module.QueryVoteLive(strSQL)
	if err != nil {
		log.Errorf("QueryVoteLive strSQL:%v, err:[%v]", strSQL, err)
		return
	} else {
		for _, voteLive := range VoteLiveList {
			data[voteLive.Address] = voteLive
		}
	}

	voteLiveResp.Data = data

	w.Lock()
	w.voteLive = voteLiveResp
	w.Unlock()
}
