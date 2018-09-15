package buffer

import (
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

	ret := make([]*entity.TransactionInfo, 0, count)
	if offset > uncTrxLen { // trx is in confirmed list or other
		offset = offset - uncTrxLen
		return b.getRestTrx(uncTrxMinBlockID, offset, count)
	} //else { // at least part of trx is in unconfirmed trx list

	uncTrxBegin := offset
	if uncTrxBegin+count > uncTrxLen { // part in unconfirmed, part in other
		copy(ret, b.trxListUnconfirmed[uncTrxBegin:])
		cList := b.getRestTrx(uncTrxMinBlockID, 0, uncTrxBegin+count-uncTrxLen)
		// TODO: verify the first element of cList's BlockID should be uncTRxMinBLockID -1
		ret = append(ret, cList...)
		return ret
	}
	// else { // all trx is in unconfirmed list
	copy(ret, b.trxListUnconfirmed[uncTrxBegin:uncTrxBegin+count])
	return ret
}

func (b *blockBuffer) GetTransactionByBlockID(blockID int64) []*entity.TransactionInfo {
	return nil
}

func (b *blockBuffer) GetTransactionByHash(hash string) []*entity.TransactionInfo {
	return nil
}

func (b *blockBuffer) GetTransactionByOwnerAddr(addr string) []*entity.TransactionInfo {
	return nil
}
