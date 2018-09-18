package buffer

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/web/entity"
)

func (b *blockBuffer) GetTransfers(offset, count int64) []*entity.TransferInfo {
	// uncTrxLen := int64(len(b.trxListUnconfirmed))
	uncTranLen, uncTranMinBlockID := b.getUnconfirmdTranListInfo()

	log.Debugf("get tran(offset:%v, count:%v), uncLen:%v, uncMinBlockID:%v\n", offset, count, uncTranLen, uncTranMinBlockID)

	ret := make([]*entity.TransferInfo, count, count)
	if offset > uncTranLen { // trx is in confirmed list or other
		offset = offset - uncTranLen
		return b.getRestTran(uncTranMinBlockID, offset, count)
	} //else { // at least part of trx is in unconfirmed trx list

	uncTranBegin := offset
	if uncTranBegin+count > uncTranLen { // part in unconfirmed, part in other
		copy(ret, b.tranListUnconfirmed[uncTranBegin:])
		cList := b.getRestTran(uncTranMinBlockID, 0, uncTranBegin+count-uncTranLen)
		// TODO: verify the first element of cList's BlockID should be uncTRxMinBLockID -1
		copy(ret[uncTranLen-uncTranBegin:], cList[:])
		// ret = append(ret, cList...)
		return ret
	}
	// else { // all trx is in unconfirmed list
	copy(ret, b.tranListUnconfirmed[uncTranBegin:uncTranBegin+count])
	return ret
}

func (b *blockBuffer) GetTransferByBlockID(blockID int64) []*entity.TransferInfo {
	if blockID > b.GetMaxConfirmedBlockID() {
		raw, ok := b.uncBlockTrans.Load(blockID)
		// log.Debugf("get uncBlockTrx[%v]--->%v-->%v\n", blockID, ok, raw)
		if !ok {
			cnt := 0
			b.uncBlockTrans.Range(func(key, val interface{}) bool {
				cnt++
				// log.Debugf("%v-->%v\n", cnt, key)
				return true
			})
			return nil
		}
		ret, ok := raw.([]*entity.TransferInfo)
		if ok {
			return ret
		}
	}

	b.getConfirmedBlockTransaction(blockID) // will load transfer at the same time

	raw, ok := b.cBlockTrans.Load(blockID)
	if ok {
		trans, ok := raw.([]*entity.TransferInfo)
		if ok && nil != trans {
			return trans
		}
	}
	return nil
}

func (b *blockBuffer) GetTransferByHash(hash string) *entity.TransferInfo {
	return nil
}

func (b *blockBuffer) GetTotalTransfers() int64 {
	return atomic.LoadInt64(&b.transferCount)
}

// minBlockID: -1 mean get from the very beginnin of the list, otherwise need minBlockID read transaction from db
func (b *blockBuffer) getRestTran(minBlockID int64, offset, count int64) []*entity.TransferInfo {
	ret := make([]*entity.TransferInfo, count, count)
	// cTrxLen := int64(len(b.trxList))
	cTranLen, minCTranBlockID := b.getConfirmdTranListInfo()
	log.Debugf("get tran confirmed(offset:%v, count:%v), cLen:%v, cMinBlockID:%v, uncMinBlockID:%v\n", offset, count, cTranLen, minCTranBlockID, minBlockID)

	if minCTranBlockID == -1 {
		minCTranBlockID = minBlockID
	}
	if offset > cTranLen {
		offset = offset - cTranLen
		return b.getRestTranRedis(minCTranBlockID, offset, count)
	}
	//else { // part in confirmed list ...
	cTranBegin := offset
	if cTranBegin+count > cTranLen { // part in confirmed list, part in redis
		copy(ret, b.tranList[cTranBegin:])
		cList := b.getRestTranRedis(minCTranBlockID, 0, cTranBegin+count-cTranLen)
		copy(ret[cTranLen-cTranBegin:], cList)
		// ret = append(ret, cList...)
		return ret
	}

	// else { all in confirmed list
	copy(ret, b.tranList[cTranBegin:cTranBegin+count])
	return ret
}

func (b *blockBuffer) getRestTranRedis(blockID int64, offset, count int64) []*entity.TransferInfo {
	redisList := b.getTranDescListFromRedis(offset, count)

	retLen := int64(len(redisList))
	if retLen >= count {
		log.Debugf("get trx redis(offset:%v, count:%v), read redis Len:%v\n", offset, count, len(redisList))

		return redisList
	}

	//else { load from db
	var filter, limit string
	minBlockID := int64(0)
	if retLen > 0 {
		minBlockID = redisList[retLen-1].Block
	} else {
		minBlockID = blockID
	}

	if minBlockID == -1 {
		filter = fmt.Sprintf(" and 1=1")
	} else {
		count = count - retLen
		filter = fmt.Sprintf("and block_id < '%v'", minBlockID)
	}
	limit = fmt.Sprintf("limit %v", count+100)

	retList := b.loadTransactionFromDB(filter, limit)

	transList := make([]*entity.TransferInfo, 0, len(retList))
	for _, trx := range retList {
		if tran := b.getTransferFromTrx(trx); nil != tran {
			transList = append(transList, tran)
		}
	}
	b.storeTranDescListToRedis(transList, true)
	redisList = append(redisList, transList[0:count]...)
	log.Debugf("get tran db(offset:%v, count:%v), read db Len:%v\n", offset, count, len(retList))

	return redisList
}

func (b *blockBuffer) storeTranDescListToRedis(tranList []*entity.TransferInfo, fromDB bool) {
	if len(tranList) == 0 {
		return
	}

	if fromDB {
		redisList := make([]interface{}, 0, len(tranList))
		for _, tran := range tranList {
			redisList = append(redisList, utils.ToJSONStr(tran))
		}
		cnt, err := _redisCli.RPush(TranRedisDescListKey, redisList...).Result()
		if nil != err {
			log.Debugf("store tran to redis failed:%v, current tran desc len:%v\n", err, cnt)
		} else {
			log.Debugf("store tran to redis ok, trx list len:%v, redis tran desc list len:%v\n", len(tranList), cnt)
		}
	} else { // from memory
		redisList := make([]interface{}, 0, len(tranList))
		for _, tran := range tranList {
			redisList[len(tranList)-1] = utils.ToJSONStr(tran)
		}
		cnt, err := _redisCli.LPush(TranRedisDescListKey, redisList...).Result()
		if nil != err {
			log.Debugf("store tran to redis failed:%v, current tran desc len:%v\n", err, cnt)
		} else {
			log.Debugf("store tran to redis ok, trx list len:%v, redis tran desc list len:%v\n", len(tranList), cnt)
		}
	}

	// redis lpush list e1 e2 e3: push to list front side; top element is the lpush list last one, e.g: list result: (head) e3, e2, e1 (tail)
	// redis rpush list e4 e5 e6: push to list tail side; tail element is the rpush last one, e.g: list result: (head) e3, e2, e1, e4, e5, e6 (tail)
	// so trx list move out from confirmed list should use lpush with e1, e2, e3 (e1.block < e2.block < e3.block) for desc list
	// trx list read from db should use rpush with e3, e2, e1 (e3.block > e2.block > e1.block)
}

func (b *blockBuffer) getTranDescListFromRedis(offset, count int64) (ret []*entity.TransferInfo) {
	retList, err := _redisCli.LRange(TranRedisDescListKey, offset, count).Result()
	if nil != err || len(retList) == 0 {
		return nil
	}
	for _, val := range retList {
		trx := new(entity.TransferInfo)
		err := json.Unmarshal([]byte(val), trx)
		if err == nil {
			ret = append(ret, trx)
		}
	}
	return ret
}
