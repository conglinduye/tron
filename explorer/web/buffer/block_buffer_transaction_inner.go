package buffer

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync/atomic"
	"time"

	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

func (b *blockBuffer) getConfirmedBlockTransaction(blockID int64) []*entity.TransactionInfo {

	raw, ok := b.cBlockTrx.Load(blockID)
	if ok {
		ret, ok := raw.([]*entity.TransactionInfo)
		if ok {
			return ret
		}
	}

	filter := fmt.Sprintf(` and block_id = '%v'`, blockID)
	retTrxs := b.loadTransactionFromDBFilter(filter)

	if nil != retTrxs {
		b.cBlockTrx.Store(blockID, retTrxs)
	}

	transList := make([]*entity.TransferInfo, 0, len(retTrxs))
	for _, trx := range retTrxs {
		b.trxHash.Store(trx.Hash, trx)
		if tran := b.getTransferFromTrx(trx); nil != tran {
			transList = append(transList, tran)
			b.tranHash.Store(tran.TransactionHash, tran)
		}
	}
	b.cBlockTrans.Store(blockID, transList)

	return retTrxs
}

// sweep transaction buffer size
func (b *blockBuffer) sweepTrxHash() {

}

func (b *blockBuffer) loadTransactionCountFromDB() {
	strSQL := fmt.Sprintf(`select count(1) as totalNum from tron.transactions`)
	log.Debug(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("loadTransactionCountFromDB error :[%v]\n", err)
		return
	}
	if dataPtr == nil {
		log.Errorf("loadTransactionCountFromDB dataPtr is nil ")
		return
	}
	//填充数据
	for dataPtr.NextT() {
		totalNum := mysql.ConvertDBValueToInt64(dataPtr.GetField("totalNum"))
		if totalNum > 0 {
			atomic.StoreInt64(&b.transactionCount, totalNum+int64(len(b.trxListUnconfirmed)))
		}
	}
	return
}

func (b *blockBuffer) loadTransferCountFromDB() {
	strSQL := fmt.Sprintf(`select count(1) as totalNum from tron.contract_transfer`)
	log.Debug(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf("loadTransferCountFromDB error :[%v]\n", err)
		return
	}
	if dataPtr == nil {
		log.Errorf("loadTransferCountFromDB dataPtr is nil ")
		return
	}
	//填充数据
	for dataPtr.NextT() {
		totalNum := mysql.ConvertDBValueToInt64(dataPtr.GetField("totalNum"))
		if totalNum > 0 {
			atomic.StoreInt64(&b.transferCount, totalNum+int64(len(b.tranListUnconfirmed)))
		}
	}
	return
}

func (b *blockBuffer) loadTransactionFromDBFilter(filter string) []*entity.TransactionInfo {
	strSQL := fmt.Sprintf(`
	select block_id,owner_address,to_address,
	trx_hash,contract_data,result_data,fee,
	contract_type,confirmed,create_time,expire_time
	from tron.transactions
	where 1=1 `)

	order := " order by block_id desc "
	ret, err := module.QueryTransactionsRealize(strSQL, filter, order, "")
	if nil != err {
		return nil
	}
	return ret.Data
}

func (b *blockBuffer) getUnconfirmdTrxListInfo() (int64, int64) {
	if len(b.trxListUnconfirmed) > 0 {
		return int64(len(b.trxListUnconfirmed)), b.trxListUnconfirmed[len(b.trxListUnconfirmed)-1].Block
	}
	return 0, -1
}

func (b *blockBuffer) getUnconfirmdTranListInfo() (int64, int64) {
	if len(b.tranListUnconfirmed) > 0 {
		return int64(len(b.tranListUnconfirmed)), b.tranListUnconfirmed[len(b.tranListUnconfirmed)-1].Block
	}
	return 0, -1
}

func (b *blockBuffer) getConfirmdTrxListInfo() (int64, int64) {
	if len(b.trxList) > 0 {
		return int64(len(b.trxList)), b.trxList[len(b.trxList)-1].Block
	}
	return 0, -1
}

func (b *blockBuffer) getConfirmdTranListInfo() (int64, int64) {
	if len(b.tranList) > 0 {
		return int64(len(b.tranList)), b.tranList[len(b.tranList)-1].Block
	}
	return 0, -1
}

func (b *blockBuffer) bufferUnconfirmTransactions(blockID int64, trxList []*entity.TransactionInfo) {
	sort.SliceStable(trxList, func(i, j int) bool { return trxList[i].Block > trxList[j].Block })

	// buffer to block id map
	b.uncBlockTrx.Store(blockID, trxList)
	// log.Debugf("store uncBlock [%v] trx:%v\n", blockID, len(trxList))

	// buffer trx hash
	for _, trx := range trxList {
		b.trxHash.Store(trx.Hash, trx)
	}

	// buffer to trx list
	b.trxListUnconfirmed = append(trxList, b.trxListUnconfirmed...)
	// log.Debugf("### buffer uncTrx, len:%v, total len:%v\n", len(trxList), len(b.trxListUnconfirmed))

}

func (b *blockBuffer) bufferUnconfirmTransfers(blockID int64, trans []*entity.TransferInfo) {
	sort.SliceStable(trans, func(i, j int) bool { return trans[i].Block > trans[j].Block })

	// buffer to block id map
	b.uncBlockTrans.Store(blockID, trans)
	// log.Debugf("store uncBlock [%v] trx:%v\n", blockID, len(trxList))

	// buffer trx hash
	for _, tran := range trans {
		b.tranHash.Store(tran.TransactionHash, tran)
	}

	// buffer to trx list
	b.tranListUnconfirmed = append(trans, b.tranListUnconfirmed...)
	// log.Debugf("### buffer uncTrx, len:%v, total len:%v\n", len(trxList), len(b.trxListUnconfirmed))

}

// 清除unconfirmed缓存中已经被确认的transaction
func (b *blockBuffer) cleanConfirmedTrxBufferFromUncTrxList() {
	// clean up confirmed transaction
	//	the min block unconfirmed is GetMaxConfirmedBlockID + 1, remove transactions which block id small than it
	validUnconfirmedBlockID := b.GetMaxConfirmedBlockID() + 1
	uncTrxIdx := len(b.trxListUnconfirmed) - 1
	for uncTrxIdx >= 0 {
		if b.trxListUnconfirmed[uncTrxIdx].Block < validUnconfirmedBlockID {
			uncTrxIdx--
		} else {
			break
		}
	}
	b.trxListUnconfirmed = b.trxListUnconfirmed[0 : uncTrxIdx+1] // +1 mean include index of uncTrxLen
	//	log.Debugf("### clean uncTrx, uncTrxIdx:%v, uncTrx len:%v\n", uncTrxIdx, len(b.trxListUnconfirmed))
	b.cleanConfirmedTranBufferFromUncTranList()
}

func (b *blockBuffer) cleanConfirmedTranBufferFromUncTranList() {
	validUnconfirmedBlockID := b.GetMaxConfirmedBlockID() + 1
	uncTranIdx := len(b.tranListUnconfirmed) - 1
	for uncTranIdx >= 0 {
		if b.tranListUnconfirmed[uncTranIdx].Block < validUnconfirmedBlockID {
			uncTranIdx--
		} else {
			break
		}
	}
	b.tranListUnconfirmed = b.tranListUnconfirmed[0 : uncTranIdx+1] // +1 mean include index of uncTrxLen

}

func (b *blockBuffer) bufferConfiremdTransaction(filter string, limit string) {
	data := b.loadTransactionFromDB(filter, "", limit)

	sort.SliceStable(data, func(i, j int) bool { return data[i].Block > data[j].Block })
	b.trxList = append(data, b.trxList...)
	if len(b.trxList) > b.maxConfirmedTrx {
		b.trxList = b.trxList[0:b.maxConfirmedTrx]
	}

	tranList := make([]*entity.TransferInfo, 0, len(data))
	// store blockID -> trxList, transList and transHash -> trans
	if len(data) > 0 {
		blockID := data[0].Block
		blockTrx := make([]*entity.TransactionInfo, 0, 30)
		blockTran := make([]*entity.TransferInfo, 0, 30)
		for _, trx := range data {
			b.trxHash.Store(trx.Hash, trx) // trx hash index

			if tran := b.getTransferFromTrx(trx); nil != tran {
				b.tranHash.Store(tran.TransactionHash, tran)
				blockTran = append(blockTran, tran)
				tranList = append(tranList, tran)
			}

			if blockID != trx.Block {
				trxs := make([]*entity.TransactionInfo, len(blockTrx))
				copy(trxs, blockTrx[:])
				b.cBlockTrx.Store(blockID, trxs)

				trans := make([]*entity.TransferInfo, len(blockTran))
				copy(trans, blockTran[:])
				b.cBlockTrans.Store(blockID, trans)

				blockID = trx.Block
				blockTrx = blockTrx[:0]
				blockTran = blockTran[:0]
			}
			blockTrx = append(blockTrx, trx)
		}
		b.cBlockTrx.Store(blockID, blockTrx)
		b.cBlockTrans.Store(blockID, blockTran)
	}

	// transList
	b.tranList = append(tranList, b.tranList...)
	if len(b.tranList) > b.maxConfirmedTrx {
		b.tranList = b.tranList[0:b.maxConfirmedTrx]
	}
}

func (b *blockBuffer) getTransferFromTrx(trx *entity.TransactionInfo) *entity.TransferInfo {
	if trx.ContractType == int64(core.Transaction_Contract_TransferContract) {
		transfer := new(entity.TransferInfo)
		transfer.Block = trx.Block
		transfer.TransactionHash = trx.Hash
		transfer.CreateTime = trx.CreateTime
		transfer.Confirmed = trx.Confirmed

		_, tranRaw := utils.GetContractByParamVal(core.Transaction_Contract_ContractType(int32(trx.ContractType)), utils.HexDecode(trx.ContractDataRaw))
		if tran, ok := tranRaw.(*core.TransferContract); ok && nil != tran {
			transfer.Amount = tran.Amount
			transfer.TokenName = "TRX"
			transfer.TransferFromAddress = utils.Base58EncodeAddr(tran.OwnerAddress)
			transfer.TransferToAddress = utils.Base58EncodeAddr(tran.ToAddress)
		}
		return transfer

	} else if trx.ContractType == int64(core.Transaction_Contract_TransferAssetContract) {
		transfer := new(entity.TransferInfo)
		transfer.Block = trx.Block
		transfer.TransactionHash = trx.Hash
		transfer.CreateTime = trx.CreateTime
		transfer.Confirmed = trx.Confirmed

		_, tranRaw := utils.GetContractByParamVal(core.Transaction_Contract_ContractType(int32(trx.ContractType)), utils.HexDecode(trx.ContractDataRaw))
		if tran, ok := tranRaw.(*core.TransferAssetContract); ok && nil != tran {
			transfer.Amount = tran.Amount
			transfer.TokenName = string(tran.AssetName)
			transfer.TransferFromAddress = utils.Base58EncodeAddr(tran.OwnerAddress)
			transfer.TransferToAddress = utils.Base58EncodeAddr(tran.ToAddress)
		}
		return transfer
	}
	return nil
}

func (b *blockBuffer) loadTransactionFromDB(filter string, order string, limit string) []*entity.TransactionInfo {
	strSQL := fmt.Sprintf(`
			select block_id,owner_address,to_address,
			trx_hash,contract_data,result_data,fee,
			contract_type,confirmed,create_time,expire_time
			from tron.transactions
			where 1=1 `)

	if len(order) == 0 {
		order = "order by block_id desc"
	}
	ret, err := module.QueryTransactionsRealize(strSQL, filter, order, limit)
	if nil != err || nil == ret && 0 == len(ret.Data) {
		log.Debugf("query trx failed:%v\n", err)
		return nil
	}

	sort.SliceStable(ret.Data, func(i, j int) bool { return ret.Data[i].Block > ret.Data[j].Block })
	return ret.Data
}

func parseBlockTransaction(block *core.Block, confirmed bool) (ret []*entity.TransactionInfo, transfers []*entity.TransferInfo) {
	if nil == block || nil == block.BlockHeader || nil == block.BlockHeader.RawData || 0 == len(block.Transactions) {
		return nil, nil
	}

	///	log.Debugf("### raw block:%v, trans count:%v\n", block.BlockHeader.RawData.Number, len(block.Transactions))

	blockID := block.BlockHeader.RawData.Number
	ret = make([]*entity.TransactionInfo, 0, len(block.Transactions))
	for _, rawTrx := range block.Transactions {
		if nil == rawTrx || nil == rawTrx.RawData || 0 == len(rawTrx.RawData.Contract) {
			continue
		}
		ctx := rawTrx.RawData.Contract[0]
		_, realCtx := utils.GetContract(ctx)

		trx := new(entity.TransactionInfo)

		trx.Block = blockID
		trx.Hash = utils.HexEncode(utils.CalcTransactionHash(rawTrx))
		// trx.CreateTime = rawTrx.RawData.Timestamp
		trx.CreateTime = block.BlockHeader.RawData.Timestamp // use block timestamp
		if ownerCtxIF, ok := realCtx.(utils.OwnerAddressIF); ok {
			trx.OwnerAddress = utils.Base58EncodeAddr(ownerCtxIF.GetOwnerAddress())
		}
		if toCtxIF, ok := realCtx.(utils.ToAddressIF); ok {
			trx.ToAddress = utils.Base58EncodeAddr(toCtxIF.GetToAddress())
		}
		trx.ContractType = int64(ctx.Type)
		trx.Confirmed = confirmed
		_, trx.ContractData = utils.GetContractInfoStr3(int32(ctx.Type), ctx.Parameter.Value)

		ret = append(ret, trx)

		// parse transfer
		if ctx.Type == core.Transaction_Contract_TransferContract {
			transfer := new(entity.TransferInfo)
			rawTransfer := realCtx.(*core.TransferContract)
			transfer.Amount = rawTransfer.Amount
			transfer.TransferToAddress = utils.Base58EncodeAddr(rawTransfer.ToAddress)
			transfer.TransferFromAddress = utils.Base58EncodeAddr(rawTransfer.OwnerAddress)
			transfer.Block = trx.Block
			transfer.TransactionHash = trx.Hash
			transfer.CreateTime = trx.CreateTime
			if 0 == transfer.CreateTime {
				transfer.CreateTime = block.BlockHeader.RawData.Timestamp
			}
			transfer.Confirmed = trx.Confirmed
			transfer.TokenName = "TRX"

			transfers = append(transfers, transfer)
		} else if ctx.Type == core.Transaction_Contract_TransferAssetContract {
			transfer := new(entity.TransferInfo)
			rawTransfer := realCtx.(*core.TransferAssetContract)
			transfer.Amount = rawTransfer.Amount
			transfer.TransferToAddress = utils.Base58EncodeAddr(rawTransfer.ToAddress)
			transfer.TransferFromAddress = utils.Base58EncodeAddr(rawTransfer.OwnerAddress)
			transfer.Block = trx.Block
			transfer.TransactionHash = trx.Hash
			transfer.CreateTime = trx.CreateTime
			if 0 == transfer.CreateTime {
				transfer.CreateTime = block.BlockHeader.RawData.Timestamp
			}
			transfer.Confirmed = trx.Confirmed
			transfer.TokenName = string(rawTransfer.AssetName)

			transfers = append(transfers, transfer)
		}
	}
	// log.Debugf("parse raw block trx, ret size:%v\n", len(ret))
	return
}

// minBlockID: -1 mean get from the very beginnin of the list, otherwise need minBlockID read transaction from db
func (b *blockBuffer) getRestTrx(minBlockID int64, offset, count int64) []*entity.TransactionInfo {
	ret := make([]*entity.TransactionInfo, count, count)
	// cTrxLen := int64(len(b.trxList))
	cTrxLen, minCTrxBlockID := b.getConfirmdTrxListInfo()
	log.Debugf("get trx confirmed(offset:%v, count:%v), cLen:%v, cMinBlockID:%v, uncMinBlockID:%v\n", offset, count, cTrxLen, minCTrxBlockID, minBlockID)

	if minCTrxBlockID == -1 {
		minCTrxBlockID = minBlockID
	}
	if offset > cTrxLen {
		offset = offset - cTrxLen
		return b.getRestTrxRedis(minCTrxBlockID, offset, count)
	}
	//else { // part in confirmed list ...
	cTrxBegin := offset
	if cTrxBegin+count > cTrxLen { // part in confirmed list, part in redis
		copy(ret, b.trxList[cTrxBegin:])
		cList := b.getRestTrxRedis(minCTrxBlockID, 0, cTrxBegin+count-cTrxLen)
		copy(ret[cTrxLen-cTrxBegin:], cList)
		// ret = append(ret, cList...)
		return ret
	}

	// else { all in confirmed list
	copy(ret, b.trxList[cTrxBegin:cTrxBegin+count])
	return ret
}

func (b *blockBuffer) getRestTrxRedis(blockID int64, offset, count int64) []*entity.TransactionInfo {

	// redisList := make([]*entity.TransactionInfo, 0, count)
	// retLen := int64(0)

	redisList := b.getTrxDescListFromRedis(offset, count)

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
	limit = fmt.Sprintf("limit %v, %v", offset, count) // load from db +100  record

	filter, order, limit := b.getTransactionIndexOffset(offset+int64(len(redisList))+int64(len(b.trxList)), count)
	fmt.Printf("filter:%v\norder:%\nlimit:%v\n", filter, order, limit)

	retList := b.loadTransactionFromDB(filter, order, limit)
	// b.storeTrxDescListToRedis(retList, true)
	if len(retList) > int(count) {
		redisList = append(redisList, retList[0:count]...)
	} else {
		redisList = append(redisList, retList...)
	}
	log.Debugf("get trx db(offset:%v, count:%v), read db Len:%v\n", offset, count, len(retList))

	return redisList
}

func (b *blockBuffer) storeTrxDescListToRedis(trxList []*entity.TransactionInfo, fromDB bool) {
	if len(trxList) == 0 {
		return
	}

	if fromDB {
		redisList := make([]interface{}, 0, len(trxList))
		for _, trx := range trxList {
			redisList = append(redisList, utils.ToJSONStr(trx))
		}
		cnt, err := _redisCli.RPush(TrxRedisDescListKey, redisList...).Result()
		if nil != err {
			log.Debugf("store trx to redis failed:%v, current trx desc len:%v\n", err, cnt)
		} else {
			log.Debugf("store trx to redis ok, trx list len:%v, redis trx desc list len:%v\n", len(trxList), cnt)
		}
	} else { // from memory
		redisList := make([]interface{}, 0, len(trxList))
		for _, trx := range trxList {
			redisList[len(trxList)-1] = utils.ToJSONStr(trx)
		}
		cnt, err := _redisCli.LPush(TrxRedisDescListKey, redisList...).Result()
		if nil != err {
			log.Debugf("store trx to redis failed:%v, current trx desc len:%v\n", err, cnt)
		} else {
			log.Debugf("store trx to redis ok, trx list len:%v, redis trx desc list len:%v\n", len(trxList), cnt)
		}
	}

	// redis lpush list e1 e2 e3: push to list front side; top element is the lpush list last one, e.g: list result: (head) e3, e2, e1 (tail)
	// redis rpush list e4 e5 e6: push to list tail side; tail element is the rpush last one, e.g: list result: (head) e3, e2, e1, e4, e5, e6 (tail)
	// so trx list move out from confirmed list should use lpush with e1, e2, e3 (e1.block < e2.block < e3.block) for desc list
	// trx list read from db should use rpush with e3, e2, e1 (e3.block > e2.block > e1.block)
}

func (b *blockBuffer) getTrxDescListFromRedis(offset, count int64) (ret []*entity.TransactionInfo) {
	retList, err := _redisCli.LRange(TrxRedisDescListKey, offset, count).Result()
	if nil != err || len(retList) == 0 {
		return nil
	}
	for _, val := range retList {
		trx := new(entity.TransactionInfo)
		err := json.Unmarshal([]byte(val), trx)
		if err == nil {
			ret = append(ret, trx)
		}
	}
	return ret
}

func (b *blockBuffer) sweepTransactionRedisList() {
	minInterval := time.Duration(600) * time.Second // 10 分钟
	for {
		ts := time.Now()

		tsc := time.Since(ts)
		if tsc < minInterval {
			time.Sleep(minInterval - tsc)
		}

		_redisCli.LTrim(TrxRedisDescListKey, 0, int64(b.maxConfirmedTrx)*2) // clean transaction redis

		_redisCli.LTrim(TranRedisDescListKey, 0, int64(b.maxConfirmedTrx)*2) // clean transfer redis
	}
}

func (b *blockBuffer) getTransactionIndexOffset(offset, count int64) (filter string, order string, limit string) {
	order = " order by block_id asc "
	limit = fmt.Sprintf("limit %v, %v", 0, count)

	index := b.trxIndex.GetIndex()
	// totalTrn := b.trxIndex.GetTotal()
	totalTrn := b.GetTotalTransactions()
	step := b.trxIndex.GetStep()

	if offset >= totalTrn {
		fmt.Printf("invalid offset:%v, total count:%v, index range:[0, %v]\n", offset, totalTrn, totalTrn-1)
		return
	}

	ascOffset := totalTrn - offset - 1
	ascOffsetIdx := ascOffset / step
	ascInnerOffsetIdx := ascOffset % step

	if ascOffsetIdx >= int64(len(index)) {
		fmt.Printf("invalid offset:%v, err index:%v\n", offset, ascOffset)
		return "", "", ""
	}

	fmt.Printf("offset:%v, ascOffset:%v, ascOffsetIdx:%v, ascInnerOffsetIdx:%v\n", offset, ascOffset, ascOffsetIdx, ascInnerOffsetIdx)

	idx := index[ascOffsetIdx]
	filter = fmt.Sprintf(" and block_id >= '%v'", idx.BlockID)
	limit = fmt.Sprintf(" limit %v, %v", idx.Offset+ascInnerOffsetIdx, count)
	return
}

func (b *blockBuffer) loadTransactionIndex() {

	sqlStr := "select start_pos, block_id, inner_offset, total_record from transactions_index order by start_pos"

	rows, err := mysql.QueryTableData(sqlStr)
	if nil != err {
		log.Errorf("load transactions_index failed:%v\n", err)
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
		b.trxIndex.Lock()
		b.trxIndex.total = index[0].Count
		b.trxIndex.index = index
		b.trxIndex.step = index[1].Count
		b.trxIndex.Unlock()
	}
}
