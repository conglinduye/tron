package buffer

import (
	"sync/atomic"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/entity"
)

// redis keys 存放transaction列表的 redis key
var (
	TrxRedisDescListKey = "transaction:list:desc" // 倒叙列表 缓存
	TrxRedisAscListKey  = "transaction:list:asc"  // 正序列表 缓存

	TranRedisDescListKey = "transfer:list:desc" // 交易缓存 倒叙
)

func (b *blockBuffer) GetTransactions(offset, count int64) []*entity.TransactionInfo {
	// uncTrxLen := int64(len(b.trxListUnconfirmed))
	uncTrxLen, uncTrxMinBlockID := b.getUnconfirmdTrxListInfo()

	log.Debugf("get trx(offset:%v, count:%v), uncLen:%v, uncMinBlockID:%v\n", offset, count, uncTrxLen, uncTrxMinBlockID)

	ret := make([]*entity.TransactionInfo, count, count)
	if offset > uncTrxLen { // trx is in confirmed list or other
		offset = offset - uncTrxLen
		return b.getRestTrx(uncTrxMinBlockID, offset, count)
	} //else { // at least part of trx is in unconfirmed trx list

	uncTrxBegin := offset
	if uncTrxBegin+count > uncTrxLen { // part in unconfirmed, part in other
		copy(ret, b.trxListUnconfirmed[uncTrxBegin:])
		cList := b.getRestTrx(uncTrxMinBlockID, 0, uncTrxBegin+count-uncTrxLen)
		// TODO: verify the first element of cList's BlockID should be uncTRxMinBLockID -1
		copy(ret[uncTrxLen-uncTrxBegin:], cList[:])
		// ret = append(ret, cList...)
		return ret
	}
	// else { // all trx is in unconfirmed list
	copy(ret, b.trxListUnconfirmed[uncTrxBegin:uncTrxBegin+count])
	return ret
}

func (b *blockBuffer) GetTransactionByBlockID(blockID int64) []*entity.TransactionInfo {

	// log.Debugf("blockID:%v-->maxConfirmedBlockID:%v\n", blockID, b.GetMaxConfirmedBlockID())
	if blockID > b.GetMaxConfirmedBlockID() {
		raw, ok := b.uncBlockTrx.Load(blockID)
		// log.Debugf("get uncBlockTrx[%v]--->%v-->%v\n", blockID, ok, raw)
		if !ok {
			cnt := 0
			b.uncBlockTrx.Range(func(key, val interface{}) bool {
				cnt++
				// log.Debugf("%v-->%v\n", cnt, key)
				return true
			})
			return nil
		}
		ret, ok := raw.([]*entity.TransactionInfo)
		if ok {
			return ret
		}
	}

	return b.getConfirmedBlockTransaction(blockID)
}

func (b *blockBuffer) GetTransactionByHash(hash string) *entity.TransactionInfo {

	if trans, ok := b.trxHash.Load(hash); ok {
		transactionInfo := trans.(*entity.TransactionInfo)
		if transactionInfo == nil {
			//TODO loda db
		}
		return transactionInfo
	}
	return nil
}

func (b *blockBuffer) GetTransactionByOwnerAddr(addr string) []*entity.TransactionInfo {
	return nil
}

// GetTotalTransactions 获取缓存中交易总数
func (b *blockBuffer) GetTotalTransactions() int64 {
	return atomic.LoadInt64(&b.transactionCount)
}
