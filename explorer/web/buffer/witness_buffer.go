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

// getWitnessBuffer
func getWitnessBuffer() *witnessBuffer {
	onceWitnessBuffer.Do(func() {
		_witnessBuffer = &witnessBuffer{}
		_witnessBuffer.load()

		go func() {
			time.Sleep(30 * time.Second)
			_witnessBuffer.load()
		}()
	})
	return _witnessBuffer
}

type witnessBuffer struct {
	sync.RWMutex

	addrMap map[string]*entity.WitnessInfo

	sortList []*entity.WitnessInfo
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

func (w *witnessBuffer) load() {
	strSQL := fmt.Sprintf(`
			select witt.address,witt.vote_count,witt.public_key,witt.url,
			witt.total_produced,witt.total_missed,acc.account_name,
			witt.latest_block_num,witt.latest_slot_num,witt.is_job
			from witness witt
			left join tron_account acc on acc.address=witt.address
			where 1=1 `)

	witnessList, err := module.QueryWitnessRealize(strSQL)
	if nil != err {
		log.Errorf("load witness from db failed:%v\n", err)
		return
	}

	addrMap := make(map[string]*entity.WitnessInfo, len(witnessList))
	sortList := make([]*entity.WitnessInfo, 0, len(witnessList))
	for _, witness := range witnessList {
		addrMap[witness.Address] = witness
		sortList = append(sortList, witness)
	}
	// votes 大的排在前面
	sort.SliceStable(sortList, func(i, j int) bool { return sortList[i].Votes > sortList[j].Votes })

	w.Lock()
	w.addrMap = addrMap
	w.sortList = sortList
	w.Unlock()
}
