package buffer

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync/atomic"

	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

// offset: 降序偏移量 最新的transfer offset ==0, offset 有效范围为 0～ total-1
// count: 记录数
// total: 请求时刻的记录总数, 用户计算生序偏移量(第一个transfer == total-1)
func (b *blockBuffer) GetTransfers(offset, count, total int64) []*entity.TransferInfo {
	// uncTrxLen := int64(len(b.trxListUnconfirmed))
	uncTranLen, uncTranMinBlockID := b.getUnconfirmdTranListInfo()

	log.Debugf("get tran(offset:%v, count:%v, total:%v), uncLen:%v, uncMinBlockID:%v\n", offset, count, total, uncTranLen, uncTranMinBlockID)

	ret := make([]*entity.TransferInfo, count, count)
	if offset > uncTranLen { // trx is in confirmed list or other
		// offset = offset - uncTranLen
		return b.getRestTran(uncTranMinBlockID, offset, count, total)
	} //else { // at least part of trx is in unconfirmed trx list

	uncTranBegin := offset
	if uncTranBegin+count > uncTranLen { // part in unconfirmed, part in other
		copy(ret, b.tranListUnconfirmed[uncTranBegin:])
		cList := b.getRestTran(uncTranMinBlockID, offset, uncTranBegin+count-uncTranLen, total)
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
func (b *blockBuffer) getRestTran(minBlockID int64, offset, count, total int64) []*entity.TransferInfo {
	ret := make([]*entity.TransferInfo, count, count)
	// cTrxLen := int64(len(b.trxList))
	cTranLen, minCTranBlockID := b.getConfirmdTranListInfo()
	log.Debugf("get tran confirmed(offset:%v, count:%v), cLen:%v, cMinBlockID:%v, uncMinBlockID:%v\n", offset, count, cTranLen, minCTranBlockID, minBlockID)

	if minCTranBlockID == -1 {
		minCTranBlockID = minBlockID
	}
	if offset-int64(len(b.tranListUnconfirmed)) > cTranLen {
		// offset = offset - cTranLen
		return b.getRestTranRedis(minCTranBlockID, offset, count, total)
	}
	//else { // part in confirmed list ...
	cTranBegin := offset - int64(len(b.tranListUnconfirmed))
	if cTranBegin+count > cTranLen { // part in confirmed list, part in redis
		copy(ret, b.tranList[cTranBegin:])
		cList := b.getRestTranRedis(minCTranBlockID, offset, cTranBegin+count-cTranLen, total)
		copy(ret[cTranLen-cTranBegin:], cList)
		// ret = append(ret, cList...)
		return ret
	}

	// else { all in confirmed list
	copy(ret, b.tranList[cTranBegin:cTranBegin+count])
	return ret
}

func (b *blockBuffer) getRestTranRedis(blockID int64, offset, count, total int64) []*entity.TransferInfo {
	redisList := make([]*entity.TransferInfo, 0, count)
	retLen := int64(0)

	// redisList := b.getTranDescListFromRedis(offset, count)

	// retLen := int64(len(redisList))
	// if retLen >= count {
	// 	log.Debugf("get tran redis(offset:%v, count:%v), read redis Len:%v\n", offset, count, len(redisList))

	// 	return redisList
	// }

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
	limit = fmt.Sprintf("limit %v, %v", offset, count)

	// filter, order, limit := b.getTransferIndexOffset(offset+int64(len(redisList))+int64(len(b.trxList)), count)
	filter, order, limit := b.getTransferIndexOffset(offset, count, total)

	retList := b.loadTransferFromDB(filter, order, limit)
	// b.storeTranDescListToRedis(retList, true)
	if len(retList) > int(count) {
		redisList = append(redisList, retList[0:count]...)
	} else {
		redisList = append(redisList, retList...)
	}
	log.Debugf("get tran (offset:%v, count:%v), read db Len:%v\n", offset, count, len(retList))

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

func (b *blockBuffer) loadTransferFromDB(filter string, order string, limit string) []*entity.TransferInfo {
	strSQL := fmt.Sprintf(`
		select block_id,owner_address,to_address,amount,
		asset_name,trx_hash,
		contract_type,confirmed,create_time
		from tron.contract_transfer
		where 1=1  `)

	if len(order) == 0 {
		order = "order by block_id desc"
	}
	ret, err := module.QueryTransfersRealize(strSQL, filter, order, limit, "", false)
	if nil != err || nil == ret && 0 == len(ret.Data) {
		log.Debugf("query trx failed:%v\n", err)
		return nil
	}

	sort.SliceStable(ret.Data, func(i, j int) bool { return ret.Data[i].Block > ret.Data[j].Block })
	return ret.Data
}

func (b *blockBuffer) getTransferIndexOffset(offset, count, total int64) (filter string, order string, limit string) {
	order = " order by block_id asc "
	limit = fmt.Sprintf("limit %v, %v", 0, count)

	index := b.tranIndex.GetIndex()
	step := b.tranIndex.GetStep()
	if 0 == step {
		return
	}

	if offset > total {
		fmt.Printf("invalid offset:%v, total count:%v, index range:[0, %v]\n", offset, total, total-1)
		return
	}

	ascOffset := total - offset - 1
	ascOffsetIdx := ascOffset / step
	ascInnerOffsetIdx := ascOffset % step

	if ascOffsetIdx >= int64(len(index)) {
		fmt.Printf("invalid offset:%v, err index:%v\n", offset, ascOffset)
		return "", "", ""
	}

	fmt.Printf("transfer index: totalTrn:%v (current total:%v), step:%v, offset:%v, ascOffset:%v, ascOffsetIdx:%v, ascInnerOffsetIdx:%v\n", total, b.GetTotalTransfers(), step, offset, ascOffset, ascOffsetIdx, ascInnerOffsetIdx)

	idx := index[ascOffsetIdx]
	filter = fmt.Sprintf(" and block_id >= '%v'", idx.BlockID)
	limit = fmt.Sprintf(" limit %v, %v", idx.Offset+ascInnerOffsetIdx, count)
	return
}

func (b *blockBuffer) loadTransferIndex() {

	sqlStr := "select start_pos, block_id, inner_offset, total_record from contract_transfer_index order by start_pos"

	rows, err := mysql.QueryTableData(sqlStr)
	if nil != err {
		log.Errorf("load contract_transfer_index failed:%v\n", err)
		return
	}

	index := make([]*indexPos, 0, 10000)
	for rows.NextT() {
		idx := new(indexPos)
		idx.Position = mysql.ConvertStringToInt64(rows.GetField("start_pos"), 0)
		idx.BlockID = mysql.ConvertStringToInt64(rows.GetField("block_id"), 0)
		idx.Offset = mysql.ConvertStringToInt64(rows.GetField("inner_offset"), 0)
		idx.Count = mysql.ConvertStringToInt64(rows.GetField("total_record"), 0)
		index = append(index, idx)
	}

	if len(index) > 1 {
		b.tranIndex.Lock()
		b.tranIndex.total = index[0].Count
		b.tranIndex.index = index
		b.tranIndex.step = index[1].Count
		b.tranIndex.Unlock()
	}
}
