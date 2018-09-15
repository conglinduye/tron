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
			fmt.Printf("reset solidity connection, new client:%v!!!\n", b.solidityClient.Target())
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
			fmt.Printf("reset wallet connection, new client:%v!!!\n", b.walletClient.Target())
		}
		return false
	}

	blockInfo := coreBlockConvert(block)
	atomic.StoreInt64(&b.realMaxBlockID, blockInfo.Number)

	nowBlockID := blockInfo.Number
	b.maxBlockTimeStamp = blockInfo.CreateTime
	numEnd := nowBlockID
	numStart := b.GetMaxConfirmedBlockID() + 1

	fmt.Printf("current max block_id:%v, current confirmed block_id:%v, unconfirmed block_id:%v, may need synchronize block:%v\n", nowBlockID, b.maxConfirmedBlockID, b.maxBlockID, nowBlockID-b.maxConfirmedBlockID)

	if numStart < b.maxBlockID {
		numStart = b.maxBlockID + 1 // maxBlockID we have store in memory
	}
	if numStart+b.maxUnconfirmedBlockRead < nowBlockID { // only read maxUnconfirmedBlock block
		numEnd = numStart + b.maxUnconfirmedBlockRead
	}

	fmt.Printf("current need buffer unconfirmed block range:%v ~ %v\n", numStart, numEnd)

	ts := time.Now()
	rawBlocks := b.getBlocksStable(numStart, numEnd)
	fmt.Printf("get blockStable cost:%v, get block count:%v, need load:%v, gap:%v\n", time.Since(ts), len(rawBlocks), blockInfo.Number-b.maxConfirmedBlockID, blockInfo.Number-b.maxConfirmedBlockID-int64(len(rawBlocks)))

	blocks := make([]*entity.BlockInfo, 0, len(rawBlocks)+1)
	for _, rawBlock := range rawBlocks {
		bi := coreBlockConvert(rawBlock)
		if nil != bi {
			blocks = append(blocks, bi)
		}
	}
	blocks = append(blocks, blockInfo)

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
		limit = "limit 100"
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
		atomic.StoreInt64(&b.maxConfirmedBlockID, maxBlockID)
	}

	return blocks.Data
}

func (b *blockBuffer) bufferBlock(blocks []*entity.BlockInfo) bool {
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

	// fmt.Printf("readbuffer %v ~ %v (%v)\n", numStart, numEnd, numEnd-numStart+1)
	curMaxBlockID := b.GetMaxBlockID()
	if numEnd > curMaxBlockID { // the max block id we can get is max block id
		// fmt.Printf("read buffer change numEnd from %v to %v\n", numEnd, b.maxBlockID)
		numEnd = curMaxBlockID
	}
	// if numEnd == 0 { // GetMaxBlockID confirm curMaxBlockID >= maxConfirmedBlockID
	// 	fmt.Printf("read buffer change numEnd from %v to %v\n", numEnd, b.maxConfirmedBlockID)
	// 	numEnd = b.maxConfirmedBlockID
	// }
	// data either in buffer or in db, we do not get data from main net in readBuffer

	// fmt.Printf("readbuffer %v ~ %v (%v)\n", numStart, numEnd, numEnd-numStart+1)
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
	// fmt.Printf("readBuffer get from buffer:%v, missing:%v\n", len(ret), len(missingBlockID))

	if len(missingBlockID) > 0 {
		var redisBuf []*entity.BlockInfo
		redisBuf, missingBlockID = b.loadBlockFromRedis(missingBlockID)
		if len(redisBuf) > 0 {
			ret = append(ret, redisBuf...)
		}
	}

	if len(missingBlockID) > 0 {
		// ts := time.Now()
		blocks := b.getBlocksStableB(missingBlockID)
		// fmt.Printf("readbuffer load from db cost:%v, size:%v\n", time.Since(ts), len(blocks))
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
		b.getNowConfirmedBlock()
		for {
			if b.getSolidityNodeMaxBlockID() {
				break
			}
		}
		for {
			if b.getNowBlock() {
				break
			}
		}
		tsc := time.Since(ts)
		if tsc < minInterval {
			time.Sleep(minInterval - tsc)
		}
	}
}

func (b *blockBuffer) backgroundSwaper() {
	minInterval := time.Duration(10) * time.Second
	swapData := make([]*entity.BlockInfo, 3000)
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
		b.buffer.Range(func(key, val interface{}) bool {
			id, ok := key.(int64)
			if ok && id <= minBlockID {
				b.buffer.Delete(key)
				block := val.(*entity.BlockInfo)
				swapData = append(swapData, block)
				if maxBlockIDSwap < block.Number {
					maxBlockIDSwap = block.Number
				}
				if minBlockIDSwap > block.Number {
					minBlockIDSwap = block.Number
				}
			}
			return true
		})
		fmt.Printf("swap record count:%v, min block_id:%v, max block_id:%v\n", len(swapData), minBlockIDSwap, maxBlockIDSwap)
		b.syncBlockToRedis(swapData)
		swapData = swapData[:0]
	}
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

// numEnd do not need to get
func (b *blockBuffer) getBlocksStable(numStart int64, numEnd int64) []*core.Block {
	if numStart > numEnd {
		return nil
	}
	fmt.Printf("get block stable, start:%v, end:%v\n", numStart, numEnd)

	ret := make([]*core.Block, 0, numEnd-numStart)
	for i := numEnd - 1; i >= numStart; i-- {
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
			// fmt.Printf("success get one block:%v, total:%v\n", block.BlockHeader.RawData.Number, len(ret))
			ret = append(ret, block)
			break
		}
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
	// fmt.Printf("read buffer from db filter:[%v]", filter)

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
