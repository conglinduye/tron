package buffer

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/wlcy/tron/explorer/lib/log"

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
var blockBF *blockBuffer
var once sync.Once

// Buffer ...
type Buffer interface {
	// blocks
	GetMaxBlockID() int64
	GetMaxConfirmedBlockID() int64
	GetFullNodeMaxBlockID() int64
	GetSolidityNodeMaxBlockID() int64
	GetMaxBlockTimestamp() int64

	GetBlocks(startID int64, offset int64, count int64) (blocks []*entity.BlockInfo, err error)
	GetBlock(blockID int64) (block *entity.BlockInfo)

	// transaction
	GetTransactions(offset, count, total int64) []*entity.TransactionInfo
	GetTransactionByBlockID(blockID int64) []*entity.TransactionInfo
	GetTransactionByHash(hash string) *entity.TransactionInfo
	GetTotalTransactions() int64

	// transfer
	GetTransfers(offset, count, total int64) []*entity.TransferInfo
	GetTransferByBlockID(blockID int64) []*entity.TransferInfo
	GetTransferByHash(hash string) *entity.TransferInfo
	GetTotalTransfers() int64
}

// GetMaxBlockID 获取DB最大的可用块ID （从fullnode获取，在缓存中可用的最大blockID）
func (b *blockBuffer) GetMaxBlockID() int64 {

	blockIDUnconfirmed := atomic.LoadInt64(&b.maxBlockID)
	blockIDConfirmed := atomic.LoadInt64(&b.maxConfirmedBlockID)
	if blockIDUnconfirmed == 0 {
		return blockIDConfirmed
	}
	return blockIDUnconfirmed
}

// GetMaxConfirmedBlockID 从db获取，存储在缓存中最大的确认块ID
func (b *blockBuffer) GetMaxConfirmedBlockID() int64 {

	blockID := atomic.LoadInt64(&b.maxConfirmedBlockID)

	return blockID
}

// GetFullNodeMaxBlockID 从fullnode获取最大的块ID
func (b *blockBuffer) GetFullNodeMaxBlockID() int64 {
	return atomic.LoadInt64(&b.realMaxBlockID)
}

// GetMaxBlockTimestamp 获取最大块的时间戳
func (b *blockBuffer) GetMaxBlockTimestamp() int64 {
	return atomic.LoadInt64(&b.maxBlockTimeStamp)
}

// GetSolidityNodeMaxBlockID 获取solidity最大确认块
func (b *blockBuffer) GetSolidityNodeMaxBlockID() int64 {
	return atomic.LoadInt64(&b.realMaxConfirmedBlockID)
}

// GetBlocks 从缓存批量读取blocks
//	startID: blockID start to get, -1 mean get from maxBlockID, if startID == -1, use offset to decide which is the max block_id in the buffer
//	offset: 从最新块开始的偏移量，返回的blocks max(block_id) = 缓存的currentMaxBlockID - startNum), if startID >= 0, ignore offset
//	count: 需要返回的块的数量
func (b *blockBuffer) GetBlocks(startID int64, offset int64, count int64) (blocks []*entity.BlockInfo, err error) {
	log.Debugf("get blocks data in buffer...")
	// log.Debugf("GetBlocks startNum:%v, offset:%v, count:%v\n", startID, offset, count)
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
	// log.Debugf("GetBlocks finah startNum:%v, offset:%v, count:%v, maxBlockID:%v, numStart:%v, numEnd:%v, count:%v\n", startID, offset, count, maxBlockID, numStart, numEnd, numEnd-numStart+1)
	ret := b.readBuffer(numStart, numEnd)
	return ret, nil
}

func (b *blockBuffer) GetBlock(blockID int64) (block *entity.BlockInfo) {
	log.Debugf("get block data in buffer...")
	if blockID > b.GetMaxBlockID() {
		return nil
	}
	ret := b.readBuffer(blockID, blockID)
	if len(ret) > 0 {
		return ret[0]
	}
	return nil
}

var BufferRedisDSN = "oline-tronscan-redis.jw6c3k.ng.0001.use2.cache.amazonaws.com:6379"

// GetBlockBuffer ...
func GetBlockBuffer() Buffer {
	_onceBlockBuffer.Do(func() {
		log.Infof("buffer redis server:%v", BufferRedisDSN)
		initRedis([]string{BufferRedisDSN})

		_blockBuffer = &blockBuffer{}

		_blockBuffer.solidityClient = grpcclient.GetRandomSolidity()
		_blockBuffer.walletClient = grpcclient.GetRandomWallet()
		_blockBuffer.maxNodeErr = 3
		_blockBuffer.maxUnconfirmedBlockRead = 50
		_blockBuffer.maxBlockInMemory = 5000
		_blockBuffer.maxConfirmedTrx = 30000

		go _blockBuffer.backgroundWorker()
		go _blockBuffer.backgroundSwaper()

	})
	return _blockBuffer
}

func initRedis(redisAddr []string) {
	redisOpt := &redis.Options{
		Addr:     redisAddr[0],
		Password: "",
		DB:       0,
	}
	_redisCli = redis.NewClient(redisOpt)

	pong, err := _redisCli.Ping().Result()
	fmt.Println(pong, err)
}
