package buffer

/*
store all tokenInfo for account  in memory
load from db every 30 seconds
*/
/*

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
		_accountTokenBuffer.load()

		go func() {
			time.Sleep(5 * time.Second)
			_accountTokenBuffer.load()
		}()
	})
	return _accountTokenBuffer
}

type accountTokenBuffer struct {
	sync.RWMutex

	accountTokenInfoList map[string]map[string]int64

	updateTime string
}

func (w *marketBuffer) GetMarket() (witness []*entity.MarketInfo) {
	if len(w.marketInfoList) == 0 {
		log.Debugf("get market info from buffer nil, data reload at :[%v]", w.updateTime)
		w.load()
	}
	log.Debugf("get market info from buffer, buffer data updated at :[%v]", w.updateTime)
	return w.marketInfoList
}

func (w *marketBuffer) getAccountTokenBuffer() {
	strSQL := fmt.Sprintf(`
	select acc.address,acc.asset_name as token_name,acc.creator_address,acc.balance
	from tron.account_asset_balance acc
	where 1=1 order by `)
	log.Debug(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("getAccountTokenBuffer error :[%v]\n", err)
		return
	}
	if dataPtr == nil {
		log.Errorf("getAccountTokenBuffer dataPtr is nil ")
		return
	}
	accountTokenMap := make(map[string]int64, 0)
	for dataPtr.NextT() {
		address := dataPtr.GetField("address")
		tokenName := dataPtr.GetField("token_name")
		balance := mysql.ConvertDBValueToInt64(dataPtr.GetField("balance"))
		if address != "" {
			if _, ok := accountTokenMap[tokenName]; !ok {
				accountTokenMap[tokenName] = balance
			}
		}
	}

	log.Debugf("market in buffer : parse page data done.")
	w.Lock()
	w.marketInfoList = marketInfos
	w.updateTime = time.Now().Local().Format(mysql.DATETIMEFORMAT)
	w.Unlock()
}
*/
