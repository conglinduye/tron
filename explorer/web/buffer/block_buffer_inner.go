package buffer

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

type blockBuffer struct {
	// sync.RWMutex
	realMaxBlockID          int64 // fullnode的最大区块ID
	realMaxConfirmedBlockID int64 //solidity node的最大区块ID
	maxBlockID              int64
	maxConfirmedBlockID     int64

	solidityClient *grpcclient.WalletSolidity
	solidityErrCnt int
	walletClient   *grpcclient.Wallet
	walletErrCnt   int

	buffer sync.Map // blockID, blockInfo

	maxNodeErr              int   // 3                 // 单个node连接允许的最大错误数
	maxUnconfirmedBlockRead int64 //  = int64(50) // 需要缓存的最新的unconfirmed block的数量
	maxBlockInMemory        int64 // max number of confirmed block in memory
	maxBlockTimeStamp       int64 //max timestamp for confirmed block
	maxConfirmedTrx         int   // 内存中最大的confirmed transaction 数量

	// transactions
	trxListUnconfirmed []*entity.TransactionInfo // tranx 列表，最新开始, unconfirmed
	trxList            []*entity.TransactionInfo // tranx 列表，最新开始, confirmed
	uncBlockTrx        sync.Map                  // 未确认block的 trx 缓存 blockID -> trxList: *entity.TransactionInfo type
	cBlockTrx          sync.Map                  // 确认block的 trx 缓存 blockID -> trxList: *entity.TransactionInfo type
	trxHash            sync.Map                  // trx_hash -> entity.TransactionInfo

	// transfers, the same as tansactions
	tranListUnconfirmed []*entity.TransferInfo
	tranList            []*entity.TransferInfo
	uncBlockTrans       sync.Map
	cBlockTrans         sync.Map
	tranHash            sync.Map

	// not implement yet
	ownerAddrTrx sync.Map // owner_addr -> entity.TransactionInfo, not use yet

	transactionCount int64 //total transaction record
	transferCount    int64 //total transaction record
}

func (b *blockBuffer) getSolidityNodeMaxBlockID() bool {
	if nil == b.solidityClient {
		b.solidityClient = grpcclient.GetRandomSolidity()
	}
	block, err := b.solidityClient.GetNowBlock()
	if nil != err || nil == block || nil == block.BlockHeader || nil == block.BlockHeader.RawData {
		b.solidityErrCnt++
		if b.solidityErrCnt > b.maxNodeErr {
			b.solidityClient = grpcclient.GetRandomSolidity()
			b.solidityErrCnt = 0
			log.Debugf("reset solidity connection, new client:%v!!!\n", b.solidityClient.Target())
		}
		return false
	}
	blockInfo := coreBlockConvert(block)
	atomic.StoreInt64(&b.realMaxConfirmedBlockID, blockInfo.Number)

	return true
}

// getNowBlock 获取最新的未确认块并存入redis，更新 maxBlockID 字段
func (b *blockBuffer) getNowBlock() bool {
	if nil == b.walletClient {
		b.walletClient = grpcclient.GetRandomWallet()
	}
	block, err := b.walletClient.GetNowBlock()
	if nil != err || nil == block || nil == block.BlockHeader || nil == block.BlockHeader.RawData {
		b.walletErrCnt++
		if b.walletErrCnt > b.maxNodeErr {
			b.walletClient = grpcclient.GetRandomWallet()
			b.walletErrCnt = 0
			log.Debugf("reset wallet connection, new client:%v!!!\n", b.walletClient.Target())
		}
		return false
	}

	blockInfo := coreBlockConvert(block)
	atomic.StoreInt64(&b.realMaxBlockID, blockInfo.Number)

	nowBlockID := blockInfo.Number
	b.maxBlockTimeStamp = blockInfo.CreateTime
	numEnd := nowBlockID
	numStart := b.GetMaxConfirmedBlockID() + 1

	if numStart < b.maxBlockID {
		numStart = b.maxBlockID + 1 // maxBlockID we have store in memory
	}
	if numStart+b.maxUnconfirmedBlockRead < nowBlockID { // only read maxUnconfirmedBlock block
		numEnd = numStart + b.maxUnconfirmedBlockRead
	}

	log.Debugf("fullnode maxBlockID:%v, buffer maxConfirmedBlockID:%v, buffer maxBlockID:%v, (unconfirmed) blocks need synchronize:%v, buffer unconfirmed block range:(%v, %v) (count:%v)\n",
		nowBlockID, b.maxConfirmedBlockID, b.maxBlockID, nowBlockID-numStart, numStart, numEnd, numEnd-numStart+1)

	ts := time.Now()
	rawBlocks := b.getBlocksStable(numStart, numEnd)
	log.Debugf("get blockStable cost:%v, get unconfirmed block count:%v, need load:%v, gap:%v\n", time.Since(ts), len(rawBlocks), blockInfo.Number-b.maxConfirmedBlockID, blockInfo.Number-b.maxConfirmedBlockID-int64(len(rawBlocks)))

	blocks := make([]*entity.BlockInfo, 0, len(rawBlocks)+1)
	for _, rawBlock := range rawBlocks {
		bi := coreBlockConvert(rawBlock)
		if nil != bi {
			blocks = append(blocks, bi)
		}
	}

	if b.bufferBlock(blocks) {
		atomic.StoreInt64(&b.maxBlockID, numEnd)
	}

	return true
}

// getNowConfirmedBlock 从db获取当前确认块后的所有块，从db获取的块全部都是确认块
func (b *blockBuffer) getNowConfirmedBlock() []*entity.BlockInfo {

	filter := fmt.Sprintf(" and block_id > '%v'", b.maxConfirmedBlockID)
	orderBy := "order by block_id desc"
	limit := ""
	strSQL := fmt.Sprintf(`
	select block_id,block_hash,block_size,create_time,
	transaction_num,
	tx_trie_hash,parent_hash,witness_address,confirmed
	from blocks
	where 1=1`)

	if 0 == b.maxConfirmedBlockID {
		filter = ""
		limit = "limit 1000"
	}

	blocks, err := module.QueryBlocksRealize(strSQL, filter, orderBy, limit)
	if nil != err || nil == blocks || 0 == len(blocks.Data) {
		return nil
	}
	maxBlockID := int64(0)
	for _, block := range blocks.Data {
		block.WitnessName, _ = getWitnessBuffer().GetWitnessNameByAddr(block.WitnessAddress)
		if block.Number > maxBlockID {
			maxBlockID = block.Number
		}
	}

	if b.bufferBlock(blocks.Data) {
		for _, block := range blocks.Data {
			// log.Debugf("remove block from uncBlockTrx:%v\n", block.Number)
			b.uncBlockTrx.Delete(block.Number)
		}
		atomic.StoreInt64(&b.maxConfirmedBlockID, maxBlockID)

		b.bufferConfiremdTransaction(fmt.Sprintf(" and block_id >= '%v'", blocks.Data[len(blocks.Data)-1].Number), "")
		b.cleanConfirmedTrxBufferFromUncTrxList() // clean unconfirmed block transaction
	}
	//加载 并缓存 交易总数
	b.loadTransactionCountFromDB()
	b.loadTransferCountFromDB()

	return blocks.Data
}

func (b *blockBuffer) bufferBlock(blocks []*entity.BlockInfo) bool {
	if len(blocks) == 0 {
		return false
	}
	// return b.syncBlockToRedis(blocks)
	for _, block := range blocks {
		b.buffer.Store(block.Number, block)
	}
	return true
}

// include numEnd
func (b *blockBuffer) readBuffer(numStart int64, numEnd int64) []*entity.BlockInfo {
	if numStart > numEnd {
		return nil
	}

	// log.Debugf("readbuffer %v ~ %v (%v)\n", numStart, numEnd, numEnd-numStart+1)
	curMaxBlockID := b.GetMaxBlockID()
	if numEnd > curMaxBlockID { // the max block id we can get is max block id
		// log.Debugf("read buffer change numEnd from %v to %v\n", numEnd, b.maxBlockID)
		numEnd = curMaxBlockID
	}
	// if numEnd == 0 { // GetMaxBlockID confirm curMaxBlockID >= maxConfirmedBlockID
	// 	log.Debugf("read buffer change numEnd from %v to %v\n", numEnd, b.maxConfirmedBlockID)
	// 	numEnd = b.maxConfirmedBlockID
	// }
	// data either in buffer or in db, we do not get data from main net in readBuffer

	// log.Debugf("readbuffer %v ~ %v (%v)\n", numStart, numEnd, numEnd-numStart+1)
	if numStart > numEnd {
		return nil
	}

	ret := make([]*entity.BlockInfo, 0, numEnd-numStart+1)

	missingBlockID := make([]string, 0, numEnd-numStart+1)
	for i := numStart; i <= numEnd; i++ {
		tmp, ok := b.buffer.Load(i)
		if ok && nil != tmp {
			if v, ok := tmp.(*entity.BlockInfo); ok && nil != v {
				ret = append(ret, v)
			} else {
				missingBlockID = append(missingBlockID, strconv.FormatInt(i, 10))
			}
		} else {
			missingBlockID = append(missingBlockID, strconv.FormatInt(i, 10))
		}
	}
	log.Debugf("readBuffer get from buffer:%v, missing:%v\n", len(ret), len(missingBlockID))

	if len(missingBlockID) > 0 {
		ts := time.Now()
		var redisBuf []*entity.BlockInfo
		redisBuf, missingBlockID = b.loadBlockFromRedis(missingBlockID)

		if len(redisBuf) > 0 {
			ret = append(ret, redisBuf...)
		}
		log.Debugf("readBuffer load from redis cost:%v, size:%v\n", time.Since(ts), len(redisBuf))
	}

	if len(missingBlockID) > 0 {
		ts := time.Now()
		blocks := b.getBlocksStableB(missingBlockID)
		log.Debugf("readbuffer load from db cost:%v, size:%v\n", time.Since(ts), len(blocks))
		b.bufferBlock(blocks)
		ret = append(ret, blocks...)
	}

	sort.SliceStable(ret, func(i, j int) bool { return ret[i].Number > ret[j].Number })

	return ret
}

func (b *blockBuffer) backgroundWorker() {
	minInterval := time.Duration(10) * time.Second
	for {
		ts := time.Now()
		//log.Debugf("000")
		b.getNowConfirmedBlock()
		//log.Debugf("111-%v, %v, %v, %v\n", b.GetMaxBlockID(), b.GetMaxConfirmedBlockID(), b.GetSolidityNodeMaxBlockID(), b.GetFullNodeMaxBlockID())
		for {
			if b.getSolidityNodeMaxBlockID() {
				//log.Debugf("222-%v, %v, %v, %v\n", b.GetMaxBlockID(), b.GetMaxConfirmedBlockID(), b.GetSolidityNodeMaxBlockID(), b.GetFullNodeMaxBlockID())
				break
			}
			//log.Debugf("333-%v, %v, %v, %v\n", b.GetMaxBlockID(), b.GetMaxConfirmedBlockID(), b.GetSolidityNodeMaxBlockID(), b.GetFullNodeMaxBlockID())
		}
		for {
			if b.getNowBlock() {
				//log.Debugf("444-%v, %v, %v, %v\n", b.GetMaxBlockID(), b.GetMaxConfirmedBlockID(), b.GetSolidityNodeMaxBlockID(), b.GetFullNodeMaxBlockID())
				break
			}
			//log.Debugf("555-%v, %v, %v, %v\n", b.GetMaxBlockID(), b.GetMaxConfirmedBlockID(), b.GetSolidityNodeMaxBlockID(), b.GetFullNodeMaxBlockID())
		}

		tsc := time.Since(ts)
		if tsc < minInterval {
			time.Sleep(minInterval - tsc)
		}

	}
}

func (b *blockBuffer) backgroundSwaper() {

	go b.sweepBlockBuffer()
	go b.sweepTrxHash()
	go b.sweepTransactionRedisList()
}

func (b *blockBuffer) sweepBlockBuffer() {
	minInterval := time.Duration(10) * time.Second
	swapData := make([]*entity.BlockInfo, b.maxBlockInMemory)
	for {
		ts := time.Now()

		tsc := time.Since(ts)
		if tsc < minInterval {
			time.Sleep(minInterval - tsc)
		}

		maxConfirmedBlockID := b.GetMaxConfirmedBlockID()

		minBlockID := maxConfirmedBlockID - b.maxBlockInMemory
		if minBlockID < 0 {
			minBlockID = 0
		}

		maxBlockIDSwap := int64(0)
		minBlockIDSwap := int64(9999999999)
		maxBlockIDBuffered := int64(0)
		minBlockIDBuffered := int64(9999999999)
		blockCnt := 0
		b.buffer.Range(func(key, val interface{}) bool {
			id, ok := key.(int64)
			blockCnt++

			if ok && id <= minBlockID {
				b.cleanBufferBlock(id)
				block := val.(*entity.BlockInfo)
				swapData = append(swapData, block)
				if maxBlockIDSwap < block.Number {
					maxBlockIDSwap = block.Number
				}
				if minBlockIDSwap > block.Number {
					minBlockIDSwap = block.Number
				}
			}

			// statistic
			block := val.(*entity.BlockInfo)
			if ok && nil != block {
				if maxBlockIDBuffered < block.Number {
					maxBlockIDBuffered = block.Number
				}
				if minBlockIDSwap > block.Number {
					minBlockIDBuffered = block.Number
				}
			}
			return true
		})
		log.Debugf("scan buffer: total bufferd block:%v (max:%v, min:%v, gap:%v), swap count:%v, min block_id:%v, max block_id:%v\n",
			blockCnt, maxBlockIDBuffered, minBlockIDBuffered, maxBlockIDBuffered-minBlockIDBuffered,
			len(swapData), minBlockIDSwap, maxBlockIDSwap)
		b.syncBlockToRedis(swapData)
		swapData = swapData[:0]
	}
}

func (b *blockBuffer) cleanBufferBlock(blockID int64) {
	b.buffer.Delete(blockID) // clean block buffer
	trxs, ok := b.cBlockTrx.Load(blockID)
	if ok && nil != trxs {
		trxList := trxs.([]*entity.TransactionInfo)
		for _, trx := range trxList {
			b.trxHash.Delete(trx.Hash)
			b.tranHash.Delete(trx.Hash)
		}
	}
	b.cBlockTrx.Delete(blockID) // clean confirmed block trx list buffer
	b.cBlockTrans.Delete(blockID)

}

// bluk store to redis, but can't control TTL
func (b *blockBuffer) syncBlockToRedisWitoutExpire(blocks []*entity.BlockInfo) bool {
	tmp := make([]interface{}, 0, len(blocks)*2)
	for _, block := range blocks {
		tmp = append(tmp, getRedisBlockKey(block.Number), utils.ToJSONStr(block))
	}

	ret := _redisCli.MSet(tmp...)
	if nil == ret || nil != ret.Err() {
		log.Errorf("store blocks to redis failed:%v\n", ret)
		return false
	}
	return true
}

// store block to redis, with ttl
func (b *blockBuffer) syncBlockToRedis(blocks []*entity.BlockInfo) bool {
	tmp := make([]interface{}, 0, len(blocks)*2)
	for _, block := range blocks {
		if block != nil {
			tmp = append(tmp, getRedisBlockKey(block.Number), utils.ToJSONStr(block))
			_redisCli.Set(getRedisBlockKey(block.Number), utils.ToJSONStr(block), 6*time.Hour)
		}
	}
	return true
}

func getRedisBlockKey(blockID interface{}) string {
	return fmt.Sprintf("block:%v", blockID)
}

// loadBlockFromRedis 从redis读取block
func (b *blockBuffer) loadBlockFromRedis(blockIDs []string) ([]*entity.BlockInfo, []string) {
	ret := make([]*entity.BlockInfo, 0, len(blockIDs))
	retIDs := make([]string, 0, len(blockIDs))
	for _, blockID := range blockIDs {
		data, err := _redisCli.Get(getRedisBlockKey(blockID)).Result()
		if nil != err || 0 == len(data) {
			retIDs = append(retIDs, blockID)
		} else {
			block := new(entity.BlockInfo)
			err := json.Unmarshal([]byte(data), block)
			if nil == err {
				ret = append(ret, block)
			} else {
				retIDs = append(retIDs, blockID)
			}
		}
	}

	return ret, retIDs
}

// include blockID == numEnd
func (b *blockBuffer) getBlocksStable(numStart int64, numEnd int64) []*core.Block {
	if numStart > numEnd {
		return nil
	}
	// log.Debugf("get block stable, start:%v, end:%v\n", numStart, numEnd)

	ret := make([]*core.Block, 0, numEnd-numStart+1)
	for i := numEnd; i >= numStart; i-- {
		for {
			block, err := b.walletClient.GetBlockByNum(i)
			if nil != err || nil == block || nil == block.BlockHeader || nil == block.BlockHeader.RawData || i != block.BlockHeader.RawData.Number {
				b.walletErrCnt++
				if b.walletErrCnt > b.maxNodeErr {
					b.walletClient = grpcclient.GetRandomWallet()
					b.walletErrCnt = 0
				}
				continue
			}
			ret = append(ret, block)
			break
		}
	}

	for idx := len(ret) - 1; idx >= 0; idx-- {
		trxs, trans := parseBlockTransaction(ret[idx], false)
		b.bufferUnconfirmTransactions(ret[idx].BlockHeader.RawData.Number, trxs)
		b.bufferUnconfirmTransfers(ret[idx].BlockHeader.RawData.Number, trans)
	}
	return ret
}

// numEnd do not need to get
func (b *blockBuffer) getBlocksStableB(blockIDs []string) []*entity.BlockInfo {
	if len(blockIDs) == 0 {
		return nil
	}

	filter := strings.Join(blockIDs, "', '")
	filter = fmt.Sprintf("and block_id in ('%v')", filter)

	strSQL := fmt.Sprintf(`
	select block_id,block_hash,block_size,create_time,
	transaction_num,
	tx_trie_hash,parent_hash,witness_address,confirmed
	from blocks
	where 1=1`)
	// log.Debugf("read buffer from db filter:[%v]", filter)

	retRaw, _ := module.QueryBlocksRealize(strSQL, filter, "", "")
	return retRaw.Data

	// ret := make([]*entity.BlockInfo, 0, len(blockIDs))
	// for _, i := range blockIDs {
	// 	for {
	// 		block, err := b.walletClient.GetBlockByNum(i)
	// 		if nil != err || nil == block || nil == block.BlockHeader || nil == block.BlockHeader.RawData || i != block.BlockHeader.RawData.Number {
	// 			b.walletErrCnt++
	// 			if b.walletErrCnt > gMaxNodeErr {
	// 				b.walletClient = grpcclient.GetRandomWallet()
	// 				b.walletErrCnt = 0
	// 			}
	// 			continue
	// 		}
	// 		ret = append(ret, coreBlockConvert(block))
	// 		break
	// 	}
	// }
	// return ret
}

var _blockBuffer *blockBuffer
var _onceBlockBuffer sync.Once

func coreBlockConvert(inblock *core.Block) *entity.BlockInfo {
	if nil == inblock || nil == inblock.BlockHeader || nil == inblock.BlockHeader.RawData {
		return nil
	}

	ret := &entity.BlockInfo{
		Number:         inblock.BlockHeader.RawData.Number,
		Hash:           utils.HexEncode(utils.CalcBlockHash(inblock)),
		Size:           utils.CalcBlockSize(inblock),
		CreateTime:     inblock.BlockHeader.RawData.Timestamp,
		TxTrieRoot:     utils.HexEncode(inblock.BlockHeader.RawData.TxTrieRoot),
		ParentHash:     utils.HexEncode(inblock.BlockHeader.RawData.ParentHash),
		WitnessID:      int32(inblock.BlockHeader.RawData.WitnessId),
		WitnessAddress: utils.Base58EncodeAddr(inblock.BlockHeader.RawData.WitnessAddress),
		NrOfTrx:        int64(len(inblock.Transactions)),
		Confirmed:      false,
	}
	ret.WitnessName, _ = getWitnessBuffer().GetWitnessNameByAddr(ret.WitnessAddress)

	return ret
}
