package account

import (
	"sync"
	"time"

	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/main/module/rawmysql"
)

// AddressSyncInfo 地址同步信息
type AddressSyncInfo struct {
	Addr string // address in base48 encoding
	LOT  int64  // latest operation time
}

// NewAccountWorker 创建账户同步工作者
func NewAccountWorker(maxClientCnt int, queueLen int, maxDBCnt int, uniqBufferTimer int, maxRecordPerCommit int) *SyncWorker {
	ret := new(SyncWorker)
	ret.maxClientCnt = maxClientCnt
	ret.maxDBCnt = maxDBCnt
	ret.queueLen = queueLen
	ret.uniqBufferTimer = uniqBufferTimer
	ret.maxRecordPerCommit = maxRecordPerCommit

	ret.client = make(chan *grpcclient.Wallet, ret.maxClientCnt)
	for i := 0; i < ret.maxClientCnt; i++ {
		ret.client <- grpcclient.GetRandomWallet()
	}

	ret.addrs2 = make(chan []byte, queueLen)
	ret.addrsBuffer = make(chan []byte, queueLen)

	ret.addrs = make(chan *AddressSyncInfo, queueLen)

	ret.accs = make(chan *Account, queueLen)

	ret.quit = make(chan struct{})
	ret.sppedQuit = make(chan struct{})

	ret.db = make(chan struct{}, ret.maxDBCnt)
	for i := 0; i < ret.maxDBCnt; i++ {
		ret.db <- struct{}{}
	}

	go ret.uniqueWorker()
	go ret.speedAddr2()
	go ret.speedDB()

	return ret
}

// AppendTask 增加账户同步任务
func (aw *SyncWorker) AppendTask(addrList []*AddressSyncInfo) {
	if 0 == len(addrList) {
		return
	}
	for _, addr := range addrList {
		if nil != addr {
			aw.addrs <- addr
		}
	}
}

// AppendTask2 增加账户同步任务
func (aw *SyncWorker) AppendTask2(rawAddrList [][]byte) {
	if 0 == len(rawAddrList) {
		return
	}
	for _, addr := range rawAddrList {
		if 0 < len(addr) {
			// aw.addrs2 <- addr
			aw.addrsBuffer <- addr
		}
	}
}

// StartAccountWorker 启动所有未运行的账户同步工作者
func (aw *SyncWorker) StartAccountWorker() int {
	restWorker := len(aw.client)
	for i := 0; i < restWorker; i++ {
		aw.wg.Add(1)
		go aw.syncWorker(i)
	}
	return restWorker
}

// StartDBWorker 启动所有未运行的db工作者
func (aw *SyncWorker) StartDBWorker() int {
	resetWorker := len(aw.db)
	for i := 0; i < resetWorker; i++ {
		aw.wg.Add(1)
		go aw.dbWorker(i)
	}
	return resetWorker
}

// Stop 停止所有工作者
func (aw *SyncWorker) Stop() {
	select {
	case <-aw.quit:
	default:
		close(aw.quit)
	}
}

// WaitStop 等待任务结束
func (aw *SyncWorker) WaitStop() {
	aw.stopSpeed()
	for {
		if aw.GetStatus() {
			break
		}

		time.Sleep(3 * time.Second)
	}
	time.Sleep(3 * time.Second)
	aw.Stop()
	aw.wg.Wait()
}

// GetStatus 打印当前状态, true 表示当前时刻任务已经处理完
func (aw *SyncWorker) GetStatus() bool {
	resetTask := len(aw.addrs)
	resetAddr2Task := len(aw.addrs2)
	resetBuffer := len(aw.addrsBuffer)
	resetDBTask := len(aw.accs)

	log.Infof("rest addr:%v, rest addr2:%v, rest buffer:%v, rest account:%v, accout syn worker:%v, db writer:%v\n", resetTask, resetAddr2Task, resetBuffer, resetDBTask, aw.CurrentAccountSyncWorker(), aw.CurrentDBWorker())

	if 0 == resetTask && 0 == resetAddr2Task && 0 == resetDBTask && 0 == resetBuffer {
		log.Infof("all task done!")
		return true
	}
	if 0 != resetTask || 0 != resetDBTask || 0 != resetAddr2Task || 0 != resetBuffer {
		a := aw.StartAccountWorker()
		b := aw.StartDBWorker()
		log.Infof("start acc worker:%v, start db worker:%v", a, b)
	}
	return false
}

// CurrentAccountSyncWorker 当前账户同步工作者数量
func (aw *SyncWorker) CurrentAccountSyncWorker() int {
	return aw.maxClientCnt - len(aw.client)
}

// CurrentDBWorker 当前写db工作者数量
func (aw *SyncWorker) CurrentDBWorker() int {
	return aw.maxDBCnt - len(aw.db)
}

// SyncWorker 从主网获取用户信息
type SyncWorker struct {
	maxClientCnt int                     // init info
	client       chan *grpcclient.Wallet // worker connection to main net

	queueLen        int                   // chan长度
	addrs2          chan []byte           // address without latest operation time
	addrsBuffer     chan []byte           // buffer as set to unique addr
	addrs           chan *AddressSyncInfo // address need to sync
	accs            chan *Account         // address sync result account info
	uniqBufferTimer int                   // 去重缓存的内容时间范围, 单位秒

	maxDBCnt           int           // init info
	db                 chan struct{} // db worker limit
	maxRecordPerCommit int           // 单次SQL最大记录数
	// latestOperTime sync.Map      // addr->latest_operation_time, db worker result record, not used yet

	quit chan struct{}  // quit flag
	wg   sync.WaitGroup //

	sppedQuit chan struct{}
}

//
// you may not need to see
//

// uniqueWorker 将短时间内的addr去重
func (aw *SyncWorker) uniqueWorker() {
	list := make([][]byte, 0, aw.queueLen)
	ticker := time.NewTicker(time.Duration(aw.uniqBufferTimer) * time.Second)
	for {
		select {
		case addr := <-aw.addrsBuffer:
			list = append(list, addr)
		case <-ticker.C:
			tmp := removeDup(list)
			aw.appendRealTask(tmp)
			list = list[:0]
			// default: // addr is block and
			// 	tmp := removeDup(list)
			// 	aw.AppendTask2(tmp)
			// 	list = list[:0]
		}
	}
}

func (aw *SyncWorker) stopSpeed() {
	select {
	case <-aw.sppedQuit:
		return
	default:
		close(aw.sppedQuit)
	}
}

func (aw *SyncWorker) speedAddr2() {
	tmpList := make([][]byte, 0, aw.queueLen/10)
	ticker := time.NewTicker(time.Duration(aw.uniqBufferTimer) * time.Second)
	lastResult := 0
	for {
		select {
		case <-ticker.C:
			restTask := len(aw.addrs2)
			if restTask > aw.queueLen/20 && restTask-lastResult > aw.queueLen/20/2 {
			dataLoop:
				for {
					select {
					case addr := <-aw.addrs2:
						tmpList = append(tmpList, addr)
					default:
						break dataLoop
					}
				}
				tmpList = removeDup(tmpList)
				log.Infof("sppedAddr2 work start at task count:%v, current task count:%v, push back task:%v, last round spped result:%v", restTask, len(aw.addrs2), len(tmpList), lastResult)
				lastResult = len(tmpList)
				aw.appendRealTask(tmpList)
				tmpList = tmpList[:0]
			}

		case <-aw.sppedQuit:
			return
		case <-aw.quit:
			return
		}
	}
}

func (aw *SyncWorker) speedDB() {
	tmpList := make(map[string]*Account)
	ticker := time.NewTicker(time.Duration(aw.uniqBufferTimer) * time.Second)
	lastResult := 0
	for {
		select {
		case <-ticker.C:
			restTask := len(aw.accs)
			if restTask > aw.queueLen/40 && restTask-lastResult > aw.queueLen/40/2 {
			dataLoop:
				for {
					select {
					case acc := <-aw.accs:
						tmpList[acc.Addr] = acc // new one will cover old one
					default:
						break dataLoop
					}
				}

				log.Infof("sppedDB work start at task count:%v, current task count:%v, push back task:%v, last round spped result:%v", restTask, len(aw.accs), len(tmpList), lastResult)
				lastResult = len(tmpList)
				for _, val := range tmpList {
					aw.accs <- val
					delete(tmpList, val.Addr)
				}
			}
		case <-aw.sppedQuit:
			return
		case <-aw.quit:
			return
		}
	}
}

func (aw *SyncWorker) appendRealTask(in [][]byte) {
	for _, raw := range in {
		aw.addrs2 <- raw
	}
}

func removeDup(in [][]byte) [][]byte {
	if 0 == len(in) {
		return in
	}
	tmpMap := make(map[string]struct{})

	for _, raw := range in {
		tmpMap[utils.Base64Encode(raw)] = struct{}{}
	}

	log.Infof("origin len:%v, uniq len:%v", len(in), len(tmpMap))

	i := 0
	for key := range tmpMap {
		in[i] = utils.Base64Decode(key)
		i++
	}
	if i > 1 {
		log.Infof("address need to sync:%v", i)
	}
	return in[:i]
}

// worker 工作者主逻辑
func (aw *SyncWorker) syncWorker(id int) {
	client := <-aw.client
	// var addrInfo *accountSyncInfo
	idleCnt := 0
	finished := 0
	errCnt := 0

	clientList := make(map[string][]int, 3)
	curClientSuccCnt := 0

workLoop:
	for {
		select {
		case addrInfo := <-aw.addrs:
			if nil == addrInfo {
				log.Errorf("worker[%v] get addrInfo nil", id)
				break workLoop
			}

			accInfo := aw.getAccountInfo(client, utils.Base58DecodeAddr(addrInfo.Addr), addrInfo.LOT)
			if nil == accInfo {
				aw.addrs <- addrInfo

				clientList[client.Target()] = append(clientList[client.Target()], curClientSuccCnt)

				errCnt++
				client.Close()
				client = grpcclient.GetRandomWallet()
				curClientSuccCnt = 0
			} else {
				aw.accs <- accInfo
				finished++
				curClientSuccCnt++
			}

		case address := <-aw.addrs2:
			if 0 == len(address) {
				log.Errorf("worker[%v] get empty addr", id)
				break workLoop
			}

			accInfo := aw.getAccountInfo(client, address, 0)
			if nil == accInfo {
				aw.addrs2 <- address

				clientList[client.Target()] = append(clientList[client.Target()], curClientSuccCnt)

				errCnt++
				client.Close()
				client = grpcclient.GetRandomWallet()
				curClientSuccCnt = 0
			} else {
				aw.accs <- accInfo
				finished++
				curClientSuccCnt++
			}

		case <-aw.quit:
			break workLoop
		default:
			if idleCnt%6 == 0 {
				client.Close()
				// log.Infof("account worker [%v] idle.... total handle account:%v, idle cnt:%v, total error:%v, client list:%v", id, finished, idleCnt, errCnt, clientList)
			}
			idleCnt++
			time.Sleep(5 * time.Second)
		}
	}

	clientList[client.Target()] = append(clientList[client.Target()], curClientSuccCnt)
	log.Infof("account worker [%v] quit, total handle account:%v, idle cnt:%v, total error:%v, client list:%v", id, finished, idleCnt, errCnt, clientList)
	aw.client <- client
	aw.wg.Done()
}

// getAccount 获取账户信息
func (aw *SyncWorker) getAccountInfo(client *grpcclient.Wallet, rawAddr []byte, lot int64) *Account {
	acc, err := client.GetAccountRawAddr(rawAddr, 3)

	if nil != err || nil == acc {
		// fmt.Printf("get account error:%v\n", err)
		return nil
	}
	if lot > acc.LatestOprationTime {
		return nil
	}

	accNet, err := client.GetAccountNetRawAddr(rawAddr, 3)
	if nil != err || nil == accNet {
		// fmt.Printf("get accountNet error:%v\n", err)
		return nil
	}

	accInfo := new(Account)
	accInfo.SetRaw(acc)
	accInfo.SetNetRaw(accNet)

	return accInfo
}

func (aw *SyncWorker) dbWorker(id int) {
	<-aw.db
	buffer := make([]*Account, 0, aw.maxRecordPerCommit)
	bufferMap := make(map[string]*Account)
	finished := 0

workLoop:
	for {
		select {
		case accInfo := <-aw.accs:
			if nil == accInfo {
				break workLoop
			}
			// buffer = append(buffer, accInfo)
			bufferMap[accInfo.Addr] = accInfo
			if len(buffer) >= aw.maxRecordPerCommit {
				finished += len(buffer)
				// buffer = aw.storeAcc(buffer)
				for _, val := range bufferMap {
					buffer = append(buffer, val)
					delete(bufferMap, val.Addr)
				}
				oriLen := len(buffer)
				buffer = aw.storeAcc(buffer)
				log.Infof("store account (full) write to db origin len:%v, reset len:%v", oriLen, len(buffer))
			}

		case <-aw.quit:
			break workLoop
		default:
			if len(bufferMap) > 0 {
				finished += len(bufferMap)
				for _, val := range bufferMap {
					buffer = append(buffer, val)
					delete(bufferMap, val.Addr)
				}
				oriLen := len(buffer)
				buffer = aw.storeAcc(buffer)
				log.Infof("store account (gap) write to db origin len:%v, reset len:%v", oriLen, len(buffer))
			} else {
				time.Sleep(3 * time.Second)
			}
		}
	}

	finished += len(buffer)
	buffer = aw.storeAcc(buffer)

	log.Infof("db worker [%v] quit, total handle account:%v", id, finished)
	aw.db <- struct{}{}
	aw.wg.Done()
}

func (aw *SyncWorker) storeAcc(accList []*Account) []*Account {
	if 0 == len(accList) {
		return accList
	}

	dbc := rawmysql.GetMysqlDB()
	ret := StoreAccount(accList, dbc)
	if !ret {
		log.Errorf("storeAcc failed, account info count:%v", len(accList))
		return accList
	}
	return accList[:0]
}
