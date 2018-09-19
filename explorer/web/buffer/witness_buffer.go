package buffer

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

/*
store all witness in memory
load from db every 30 seconds
*/

var _witnessBuffer *witnessBuffer
var onceWitnessBuffer sync.Once

//GetWitnessBuffer ...
func GetWitnessBuffer() *witnessBuffer {
	return getWitnessBuffer()
}

// getWitnessBuffer
func getWitnessBuffer() *witnessBuffer {
	onceWitnessBuffer.Do(func() {
		_witnessBuffer = &witnessBuffer{}
		_witnessBuffer.load()
		_witnessBuffer.loadStatistic()

		go witnessBufferLoader()
	})
	return _witnessBuffer
}

func witnessBufferLoader() {
	for {
		_witnessBuffer.load()
		_witnessBuffer.loadStatistic()
		time.Sleep(30 * time.Second)

	}
}

type witnessBuffer struct {
	sync.RWMutex

	addrMap map[string]*entity.WitnessInfo

	sortList []*entity.WitnessInfo

	statisticList []*entity.WitnessStatisticInfo
}

func (w *witnessBuffer) GetWitnessNameByAddr(addr string) (name string, ok bool) {
	var witness *entity.WitnessInfo
	w.RLock()
	witness, ok = w.addrMap[addr]
	if ok && nil != witness {
		name = witness.Name
	}
	w.RUnlock()
	return
}

func (w *witnessBuffer) GetWitnessByAddr(addr string) (witness *entity.WitnessInfo, ok bool) {
	w.RLock()
	witness, ok = w.addrMap[addr]
	w.RUnlock()
	return
}

func (w *witnessBuffer) GetWitness() (witness []*entity.WitnessInfo) {
	if len(w.sortList) == 0 {
		log.Infof("get Witness info from buffer nil, data reload")
		w.load()
	}
	log.Infof("get Witness info from buffer, buffer data updated ")
	return w.sortList
}

func (w *witnessBuffer) GetWitnessStatistic() (witness []*entity.WitnessStatisticInfo) {
	if len(w.statisticList) == 0 {
		log.Infof("get WitnessStatistic info from buffer nil, data reload")
		w.loadStatistic()
	}
	log.Infof("get WitnessStatistic info from buffer, buffer data updated ")
	return w.statisticList
}

func (w *witnessBuffer) load() { //QueryWitness()
	strSQL := fmt.Sprintf(`
			select witt.address,witt.vote_count,witt.public_key,witt.url,
			witt.total_produced,witt.total_missed,acc.account_name,
			witt.latest_block_num,witt.latest_slot_num,witt.is_job
			from witness witt
			left join tron_account acc on acc.address=witt.address
			where 1=1 order by witt.vote_count desc`)

	witnessList, err := module.QueryWitnessRealize(strSQL)
	if nil != err {
		log.Errorf("load witness from db failed:%v\n", err)
		return
	}

	totalVotes := module.QueryTotalVotes()
	for _, witnessInfo := range witnessList {
		//log.Debugf("get witness list :[%#v]", witnessInfo)
		witnessInfo.ProducePercentage = 0
		witnessInfo.VotesPercentage = 0
		if witnessInfo.ProducedTotal > 0 {
			witnessInfo.ProducePercentage = float64(witnessInfo.ProducedTotal-witnessInfo.MissedTotal) / float64(witnessInfo.ProducedTotal) * 100
		}
		if totalVotes > 0 {
			witnessInfo.VotesPercentage = float64(witnessInfo.Votes) / float64(totalVotes) * 100
		}

	}

	addrMap := make(map[string]*entity.WitnessInfo, len(witnessList))
	sortList := make([]*entity.WitnessInfo, 0, len(witnessList))
	for _, witness := range witnessList {
		//log.Debugf("after calc rate for  witness list :[%#v]", witness)
		addrMap[witness.Address] = witness
		sortList = append(sortList, witness)
	}
	// votes 大的排在前面
	sort.SliceStable(sortList, func(i, j int) bool { return sortList[i].Votes > sortList[j].Votes })

	w.Lock()
	w.addrMap = addrMap
	w.sortList = sortList
	log.Debugf("set buffer data done.")
	w.Unlock()
}

func (w *witnessBuffer) loadStatistic() { //QueryWitnessStatistic()
	var blocks int64
	curMaintenanceTime, err := getMaintenanceTimeStamp()
	if err != nil {
		log.Error(err)
		return
	}
	totalBlocks, err := module.QueryTotalBlocks(curMaintenanceTime)
	if err != nil {
		log.Error(err)
		return
	}
	if totalBlocks == 0 {
		log.Errorf("total blocks is 0, replace it with large number")
		blocks = 10000000
	} else {
		blocks = totalBlocks
	}
	strSQL := fmt.Sprintf(`
	select acc.address, acc.account_name,witt.url
		   ,ifnull(blocks.blockproduce,0) as blockproduce , 
		   ifnull(blocks.blockproduce,0)/%v as blockRate
    from  tron.tron_account acc
    left join tron.witness witt on witt.address=acc.address 
    left join (
	    select witness_address,count(block_id) as blockproduce
        from tron.blocks blk
        where 1=1 and blk.create_time>%v 
        group by witness_address
    ) blocks on blocks.witness_address=acc.address
    where 1=1 and acc.is_witness=1`, blocks, curMaintenanceTime)

	statistic, err := module.QueryWitnessStatisticRealize(strSQL, totalBlocks)
	if err != nil {
		log.Error(err)
		return
	}
	w.Lock()
	w.statisticList = statistic
	w.Unlock()
}

//获取当前轮开始时间戳
func getMaintenanceTimeStamp() (int64, error) {
	nextMaintenanceTime := GetVoteBuffer().GetNextMaintenanceTime()

	curMaintenanceTime := nextMaintenanceTime - 6*60*60*1000 //6小时
	return curMaintenanceTime, nil
}
