package buffer

import (
	"fmt"

	"github.com/wlcy/tron/explorer/web/entity"
)

// redis keys 存放transaction列表的 redis key
var (
	TrxRedisDescListKey = "transaction:list:desc" // 倒叙列表 缓存
	TrxRedisAscListKey  = "transaction:list:asc"  // 正序列表 缓存
)

func (b *blockBuffer) GetTransactions(offset, count int64) []*entity.TransactionInfo {
	// uncTrxLen := int64(len(b.trxListUnconfirmed))
	uncTrxLen, uncTrxMinBlockID := b.getUnconfirmdTrxListInfo()

	fmt.Printf("get trx(offset:%v, count:%v), uncLen:%v, uncMinBlockID:%v\n", offset, count, uncTrxLen, uncTrxMinBlockID)

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

	if blockID > b.GetMaxConfirmedBlockID() {
		raw, ok := b.uncBlockTrx.Load(blockID)
		if !ok {
			return nil
		}
		ret, ok := raw.([]*entity.TransactionInfo)
		if ok {
			return ret
		}
	}

	return b.getConfirmedBlockTransaction(blockID)
}

func (b *blockBuffer) GetTransactionByHash(hash string) []*entity.TransactionInfo {
	return nil
}

func (b *blockBuffer) GetTransactionByOwnerAddr(addr string) []*entity.TransactionInfo {
	return nil
}
