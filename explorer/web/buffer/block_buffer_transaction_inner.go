package buffer

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

func (b *blockBuffer) getConfirmedBlockTransaction(blockID int64) []*entity.TransactionInfo {

	raw, ok := b.cBlockTrx.Load(blockID)
	if !ok {
		return nil
	}
	ret, ok := raw.([]*entity.TransactionInfo)
	if ok {
		return ret
	}

	filter := fmt.Sprintf(` and block_id = '%v'`, blockID)
	ret = b.loadTransactionFromDBFilter(filter)

	if nil != ret {
		b.cBlockTrx.Store(blockID, ret)
	}

	return ret
}

func (b *blockBuffer) loadTransactionFromDBFilter(filter string) []*entity.TransactionInfo {
	strSQL := fmt.Sprintf(`
	select block_id,owner_address,to_address,
	trx_hash,contract_data,result_data,fee
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

func (b *blockBuffer) getConfirmdTrxListInfo() (int64, int64) {
	if len(b.trxList) > 0 {
		return int64(len(b.trxList)), b.trxList[len(b.trxList)-1].Block
	}
	return 0, -1
}

func (b *blockBuffer) bufferUnconfirmTransactions(blockID int64, trxList []*entity.TransactionInfo) {
	// buffer to trx list
	sort.SliceStable(trxList, func(i, j int) bool { return trxList[i].Block > trxList[j].Block })
	b.trxListUnconfirmed = append(trxList, b.trxListUnconfirmed...)
	fmt.Printf("### buffer uncTrx, len:%v, total len:%v\n", len(trxList), len(b.trxListUnconfirmed))

	// buffer to block id map
	b.uncBlockTrx.Store(blockID, trxList)

	// buffer trx hash
	for _, trx := range trxList {
		b.trxHash.Store(trx.Hash, trx)
	}

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
	fmt.Printf("### clean uncTrx, uncTrxIdx:%v, uncTrx len:%v\n", uncTrxIdx, len(b.trxListUnconfirmed))

}

func (b *blockBuffer) bufferConfiremdTransaction(filter string, limit string) {
	data := b.loadTransactionFromDB(filter, limit)

	sort.SliceStable(data, func(i, j int) bool { return data[i].Block > data[i].Block })
	b.trxList = append(data, b.trxList...)
	if len(b.trxList) > b.maxConfirmedTrx {
		b.trxList = b.trxList[0:b.maxConfirmedTrx]
	}

	for _, trx := range data {
		b.trxHash.Store(trx.Hash, trx)
	}
}

func (b *blockBuffer) loadTransactionFromDB(filter string, limit string) []*entity.TransactionInfo {
	strSQL := fmt.Sprintf(`
			select block_id,owner_address,to_address,
			trx_hash,contract_data,result_data,fee
			contract_type,confirmed,create_time,expire_time
			from tron.transactions
			where 1=1 `)

	ret, err := module.QueryTransactionsRealize(strSQL, filter, "order by block_id desc", limit)
	if nil != err || nil == ret && 0 == len(ret.Data) {
		return nil
	}

	sort.SliceStable(ret.Data, func(i, j int) bool { return ret.Data[i].Block > ret.Data[i].Block })
	return ret.Data
}

func parseBlockTransaction(block *core.Block, confirmed bool) (ret []*entity.TransactionInfo) {
	if nil == block || nil == block.BlockHeader || nil == block.BlockHeader.RawData || 0 == len(block.Transactions) {
		return nil
	}

	fmt.Printf("### raw block:%v, trans count:%v\n", block.BlockHeader.RawData.Number, len(block.Transactions))

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
		trx.CreateTime = rawTrx.RawData.Timestamp
		if ownerCtxIF, ok := realCtx.(utils.OwnerAddressIF); ok {
			trx.OwnerAddress = utils.Base58EncodeAddr(ownerCtxIF.GetOwnerAddress())
		}
		if toCtxIF, ok := realCtx.(utils.ToAddressIF); ok {
			trx.ToAddress = utils.Base58EncodeAddr(toCtxIF.GetToAddress())
		}
		trx.ContractType = int64(ctx.Type)
		trx.Confirmed = confirmed
		trx.ContractData = utils.HexEncode(ctx.Parameter.Value)

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
		}
		ret = append(ret, trx)
	}
	return
}

// minBlockID: -1 mean get from the very beginnin of the list, otherwise need minBlockID read transaction from db
func (b *blockBuffer) getRestTrx(minBlockID int64, offset, count int64) []*entity.TransactionInfo {
	ret := make([]*entity.TransactionInfo, count, count)
	// cTrxLen := int64(len(b.trxList))
	cTrxLen, minCTrxBlockID := b.getConfirmdTrxListInfo()
	fmt.Printf("get trx confirmed(offset:%v, count:%v), cLen:%v, cMinBlockID:%v, uncMinBlockID:%v\n", offset, count, cTrxLen, minCTrxBlockID, minBlockID)

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
	redisList := b.getTrxDescListFromRedis(offset, count)

	retLen := int64(len(redisList))
	if retLen >= count {
		fmt.Printf("get trx redis(offset:%v, count:%v), read redis Len:%v\n", offset, count, len(redisList))

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
	limit = fmt.Sprintf("limit %v", count)

	retList := b.loadTransactionFromDB(filter, limit)
	b.storeTrxDescListToRedis(retList, true)
	redisList = append(redisList, retList...)
	fmt.Printf("get trx db(offset:%v, count:%v), read db Len:%v\n", offset, count, len(retList))

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
			fmt.Printf("store trx to redis failed:%v, current trx desc len:%v\n", err, cnt)
		} else {
			fmt.Printf("store trx to redis ok, trx list len:%v, redis trx desc list len:%v\n", len(trxList), cnt)
		}
	} else { // from memory
		redisList := make([]interface{}, 0, len(trxList))
		for _, trx := range trxList {
			redisList[len(trxList)-1] = utils.ToJSONStr(trx)
		}
		cnt, err := _redisCli.LPush(TrxRedisDescListKey, redisList...).Result()
		if nil != err {
			fmt.Printf("store trx to redis failed:%v, current trx desc len:%v\n", err, cnt)
		} else {
			fmt.Printf("store trx to redis ok, trx list len:%v, redis trx desc list len:%v\n", len(trxList), cnt)
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

func (b *blockBuffer) loadTrxDescFromDB() []*entity.TransactionInfo {
	return nil
}

func (b *blockBuffer) sweepTransactionRedisList() {
	minInterval := time.Duration(600) * time.Second // 10 分钟
	for {
		ts := time.Now()

		tsc := time.Since(ts)
		if tsc < minInterval {
			time.Sleep(minInterval - tsc)
		}

		_redisCli.LTrim(TrxRedisDescListKey, 0, int64(b.maxConfirmedTrx)*2)
	}
}
