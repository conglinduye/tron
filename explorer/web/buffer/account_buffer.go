package buffer

import (
	"fmt"
	"sync"
	"time"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
)

/*
store all tokenInfo for account  in memory
load from db every 30 seconds
*/

var _accountTokenBuffer *accountTokenBuffer
var onceAccountTokentOnce sync.Once

//GetAccountTokenBuffer ...
func GetAccountTokenBuffer() *accountTokenBuffer {
	return getAccountTokenBuffer()
}

// getAccountTokenBuffer
func getAccountTokenBuffer() *accountTokenBuffer {
	onceAccountTokentOnce.Do(func() {
		_accountTokenBuffer = &accountTokenBuffer{}
		_accountTokenBuffer.getAccountTokenBuffer()

		go accountBufferLoader()
	})
	return _accountTokenBuffer
}
func accountBufferLoader() {
	for {
		_accountTokenBuffer.getAccountTokenBuffer()
		time.Sleep(5 * time.Second)
	}
}

type accountTokenBuffer struct {
	sync.RWMutex

	accountTokenInfoList map[string]map[string]int64
}

func (w *accountTokenBuffer) GetAccountTokenBuffer(address string) (tokenBalance map[string]int64) {
	if len(w.accountTokenInfoList) == 0 {
		log.Infof("GetAccountTokenBuffer info from buffer nil, data reload")
		w.getAccountTokenBuffer()
	}
	log.Infof("GetAccountTokenBuffer info from buffer, buffer data updated ")
	return w.accountTokenInfoList[address]
}

func (w *accountTokenBuffer) getAccountTokenBuffer() {
	strSQL := fmt.Sprintf(`
	select acc.address,acc.asset_name as token_name,acc.creator_address,acc.balance
	from tron.account_asset_balance acc
	where 1=1 `)
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("getAccountTokenBuffer error :[%v]\n", err)
		return
	}
	if dataPtr == nil {
		log.Errorf("getAccountTokenBuffer dataPtr is nil ")
		return
	}
	accountInfoMap := make(map[string]map[string]int64, 0)

	for dataPtr.NextT() {
		address := dataPtr.GetField("address")
		tokenName := dataPtr.GetField("token_name")
		balance := mysql.ConvertDBValueToInt64(dataPtr.GetField("balance"))
		if address != "" {
			if tokenBlanceInfo, ok := accountInfoMap[address]; ok {
				tokenBlanceInfo[tokenName] = balance
				accountInfoMap[address] = tokenBlanceInfo
			} else {
				accountTokenMap := make(map[string]int64, 0)
				accountTokenMap[tokenName] = balance
				accountInfoMap[address] = accountTokenMap
			}
		}
	}

	log.Infof("account token info in buffer :data done.")
	w.Lock()
	w.accountTokenInfoList = accountInfoMap
	w.Unlock()
	return
}
