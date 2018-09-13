package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"

	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/grpcclient"

	"github.com/go-redis/redis"
	"github.com/wlcy/tron/explorer/web/entity"
)

/*
use memory & redis as buffer, buffer will read db
redis rule for blocks

type: key
1. blocks --> key: "block:block_id", value: BlockInfo JSON string, TTL: 6 hour

rule:
1. get max confirmed block from db, if this is the start time of buffer, load max 100 blocks from db and save it to memory
2. get now block from full node, this block is unconfirmed; check the gap between max confirmed block id and now block (gap = now_block_id_unconfirmed - max_block_id_confirmed_db)
	1. if $gap > maxUnconfirmedBlockNum, load (max_block_id_confiremd_db +1, max_block_id_confirmed_db + 1 + maxUnconfirmedBlockNum) unconfirmed block from fullnode to memory
		we only load maxUnconfirmedBlockNum blocks from fullnode to memory, set maxBlockID to it, and the user of buffer can get max block_id is max_block_id_confirmed_db + maxUnconfirmedBlockNum
	2. if $gap < maxUnconfirmedBlockNum, load (max_block_id_confiremd_db +1, now_block_id_unconfirmed) unconfirmed blocks to memory
	3. if now_block_id_unconfirmed == 0, now_block_id_unconfirmed = max_block_id_confirmed_db, no unconfirmed block in buffer (can't connect to fullnode logic)
3. when user load blocks from buffer, if max_block_id > now_block_id_unconfirmed == maxBlockID, the block buffer can read is maxBlockID, ignore blocks > maxBlockID
4. if can't get blocks from memory, buffer first try to read them fromr redis, if can't get from redis, read them from db and buffered to redis with TTL = 6 hour
5. memory clean: keep at least maxMemoryBufferBlock in memory, which block_id should in between (max_block_id_confirmed_db-3000, max_block_id_confirmed_db), use a time to swap memory buffer to redis every 10 second

read db only use block id range as condition
if we can't find block in redis, load it from db and write to redis
*/

var _redisCli *redis.Client

// GetMaxBlockID 获取最大的可用块ID
func (b *blockBuffer) GetMaxBlockID() int64 {

	blockIDUnconfirmed := atomic.LoadInt64(&b.maxBlockID)
	blockIDConfirmed := atomic.LoadInt64(&b.maxConfirmedBlockID)
	if blockIDUnconfirmed == 0 {
		return blockIDConfirmed
	}
	return blockIDUnconfirmed
}

// GetMaxConfirmedBlockID 获取最大的确认块ID
func (b *blockBuffer) GetMaxConfirmedBlockID() int64 {

	blockID := atomic.LoadInt64(&b.maxConfirmedBlockID)

	return blockID
}

// GetBlocks 从缓存批量读取blocks
//	startID: blockID start to get, -1 mean get from maxBlockID, if startID == -1, use offset to decide which is the max block_id in the buffer
//	offset: 从最新块开始的偏移量，返回的blocks max(block_id) = 缓存的currentMaxBlockID - startNum), if startID >= 0, ignore offset
//	count: 需要返回的块的数量
func (b *blockBuffer) GetBlocks(startID int64, offset int64, count int64) (blocks []*entity.BlockInfo, err error) {
	// fmt.Printf("GetBlocks startNum:%v, offset:%v, count:%v\n", startID, offset, count)
	if count <= 0 {
		return nil, nil
	}
	maxBlockID := b.GetMaxBlockID()
	if startID >= 0 {
		maxBlockID = startID
		offset = 0
	}
	numEnd := maxBlockID - offset
	if numEnd <= 0 {
		return nil, nil
	}
	numStart := numEnd - count + 1
	if numStart <= 0 {
		numStart = 0
	}
	// fmt.Printf("GetBlocks finah startNum:%v, offset:%v, count:%v, maxBlockID:%v, numStart:%v, numEnd:%v, count:%v\n", startID, offset, count, maxBlockID, numStart, numEnd, numEnd-numStart+1)
	ret := b.readBuffer(numStart, numEnd)
	return ret, nil
}

func (b *blockBuffer) GetBlock(blockID int64) (block *entity.BlockInfo) {
	if blockID > b.GetMaxBlockID() {
		return nil
	}
	ret := b.readBuffer(blockID, blockID)
	if len(ret) > 0 {
		return ret[0]
	}
	return nil
}

func getBlockBuffer() *blockBuffer {
	_onceBlockBuffer.Do(func() {
		_blockBuffer = &blockBuffer{}

		_blockBuffer.solidityClient = grpcclient.GetRandomSolidity()
		_blockBuffer.walletClient = grpcclient.GetRandomWallet()
		_blockBuffer.maxNodeErr = 3
		_blockBuffer.maxUnconfirmedBlockRead = 100
	})
	return _blockBuffer
}

var _blockBuffer *blockBuffer
var _onceBlockBuffer sync.Once

type blockBuffer struct {
	// sync.RWMutex

	maxBlockID          int64
	maxConfirmedBlockID int64

	solidityClient *grpcclient.WalletSolidity
	solidityErrCnt int
	walletClient   *grpcclient.Wallet
	walletErrCnt   int

	buffer sync.Map // blockID, blockInfo

	maxNodeErr              int   // 3                 // 单个node连接允许的最大错误数
	maxUnconfirmedBlockRead int64 //  = int64(50) // 需要缓存的最新的unconfirmed block的数量

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
	nowBlockID := blockInfo.Number

	numEnd := nowBlockID
	numStart := b.GetMaxConfirmedBlockID() + 1

	fmt.Printf("current max block_id:%v, current confirmed block_id:%v, unconfirmed block_id:%v, may need synchronize block:%v\n", nowBlockID, b.maxConfirmedBlockID, b.maxBlockID, nowBlockID-b.maxConfirmedBlockID)

	if numStart < b.maxBlockID {
		numStart = b.maxBlockID + 1 // maxBlockID we have store in memory
	}
	if numStart+b.maxUnconfirmedBlockRead < nowBlockID { // only read maxUnconfirmedBlock block
		numEnd = numStart + b.maxUnconfirmedBlockRead
	}

	fmt.Printf("current time need buffer unconfirmed block range:%v ~ %v\n", numStart, numEnd)

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

	blocks, err := QueryBlocksRealize(strSQL, filter, orderBy, limit)
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

func (b *blockBuffer) syncBlockToRedisWitoutExpire(blocks []*entity.BlockInfo) bool {
	tmp := make([]interface{}, 0, len(blocks)*2)
	for _, block := range blocks {
		tmp = append(tmp, getRedisBlockKey(block), utils.ToJSONStr(block))
	}

	ret := _redisCli.MSet(tmp...)
	if nil == ret || nil != ret.Err() {
		log.Errorf("store blocks to redis failed:%v\n", ret)
		return false
	}
	return true
}

func (b *blockBuffer) syncBlockToRedis(blocks []*entity.BlockInfo) bool {
	tmp := make([]interface{}, 0, len(blocks)*2)
	for _, block := range blocks {
		tmp = append(tmp, getRedisBlockKey(block), utils.ToJSONStr(block))
		_redisCli.Set(getRedisBlockKey(block), utils.ToJSONStr(block), 6*time.Hour)
	}
	return true
}

func getRedisBlockKey(block *entity.BlockInfo) string {
	return fmt.Sprintf("block:%v", block.Number)
}

func (b *blockBuffer) loadBlockFromRedis(blockIds []int64) {
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

	retRaw, _ := QueryBlocksRealize(strSQL, filter, "", "")
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
