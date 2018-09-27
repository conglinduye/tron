package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/tronprotocol/grpc-gateway/api"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/grpcclient"
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"github.com/wlcy/tron/explorer/lib/util"

	"github.com/wlcy/tron/explorer/web/buffer"

	"github.com/wlcy/tron/explorer/web/entity"
	"github.com/wlcy/tron/explorer/web/module"
)

// WalletClient ...
type WalletClient struct {
	*grpcclient.Wallet
	sync.Mutex
}

// Refresh ...
func (wc *WalletClient) Refresh() {
	wc.Lock()
	if nil != wc.Wallet {
		wc.Wallet.Close()
	}
	wc.Wallet = grpcclient.GetRandomWallet()
	wc.Unlock()
}

// BroadcastTransaction ...
func (wc *WalletClient) BroadcastTransaction(trx *core.Transaction) (*api.Return, error) {
	tryCnt := 3
	for tryCnt > 0 {
		tryCnt--
		ret, err := wc.Wallet.BroadcastTransaction(trx)
		if nil != err {
			wc.Refresh()
			continue
		}
		return ret, err
	}
	return nil, nil
}

var _wallet *WalletClient
var _walletOnce sync.Once

//QueryTransactionsBuffer ...
func QueryTransactionsBuffer(req *entity.Transactions) (*entity.TransactionsResp, error) {
	transactions := &entity.TransactionsResp{}
	if req.Number != "" { //按blockID查询
		transactions.Data = buffer.GetBlockBuffer().GetTransactionByBlockID(mysql.ConvertStringToInt64(req.Number, 0))
		transactions.Total = int64(len(transactions.Data))
	} else if req.Hash != "" { //按照交易hash查询
		transact := buffer.GetBlockBuffer().GetTransactionByHash(req.Hash)
		if transact == nil {
			transact, _ = QueryTransaction(req)
		}
		transacts := make([]*entity.TransactionInfo, 0)
		transacts = append(transacts, transact)
		transactions.Data = transacts
		transactions.Total = int64(len(transactions.Data))
	} else if req.Address != "" { //按照交易所属人查询，包含转出的交易，和转入的交易
		transactions, _ = QueryTransactionsByAddress(req)
	} else { //分页查询
		transactions.Data = buffer.GetBlockBuffer().GetTransactions(req.Start, req.Limit, req.Total)
		transactions.Total = buffer.GetBlockBuffer().GetTotalTransactions()
	}

	return transactions, nil
}

//QueryTransactionsByAddress  根据地址查询其下所有相关的交易列表
func QueryTransactionsByAddress(req *entity.Transactions) (*entity.TransactionsResp, error) {
	var filterSQL, sortSQL, pageSQL string
	mutiFilter := false
	strSQL := fmt.Sprintf(`
	select oo.block_id,oo.owner_address,oo.to_address,oo.contract_type,oo.trx_hash,oo.create_time from (
	SELECT block_id,owner_address,to_address,contract_type,trx_hash,create_time 
	FROM contract_transfer 
	where to_address='%v' 
	union 
	SELECT block_id,owner_address,to_address,contract_type,trx_hash,create_time 
	FROM transactions 
	where owner_address='%v') oo
	where 1=1 `, req.Address, req.Address)

	for _, v := range strings.Split(req.Sort, ",") {
		if strings.Index(v, "timestamp") > 0 {
			if mutiFilter {
				sortSQL = fmt.Sprintf("%v ,", sortSQL)
			}
			sortSQL = fmt.Sprintf("%v create_time", sortSQL)
			if strings.Index(v, "-") == 0 {
				sortSQL = fmt.Sprintf("%v desc", sortSQL)
			}
			mutiFilter = true
		}

		if strings.Index(v, "number") > 0 {
			if mutiFilter {
				sortSQL = fmt.Sprintf("%v ,", sortSQL)
			}
			sortSQL = fmt.Sprintf("%v block_id", sortSQL)
			if strings.Index(v, "-") == 0 {
				sortSQL = fmt.Sprintf("%v desc", sortSQL)
			}
			mutiFilter = true
		}
	}
	if sortSQL != "" {
		if strings.Index(sortSQL, ",") == 0 {
			sortSQL = sortSQL[1:]
		}
		sortSQL = fmt.Sprintf("order by %v", sortSQL)
	}

	pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)

	return module.QueryTransactionsRealize(strSQL, filterSQL, sortSQL, pageSQL, true)
}

//QueryTransactions 条件查询  	//?sort=-number&limit=1&count=true&number=2135998 TODO: cache
func QueryTransactions(req *entity.Transactions) (*entity.TransactionsResp, error) {
	var filterSQL, sortSQL, pageSQL string
	mutiFilter := false

	strSQL := fmt.Sprintf(`
			select block_id,owner_address,to_address,
			trx_hash,contract_data,result_data,fee,
			contract_type,confirmed,create_time,expire_time
			from transactions
			where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Number != "" {
		filterSQL = fmt.Sprintf(" and block_id=%v", req.Number)
	}
	if req.Hash != "" {
		filterSQL = fmt.Sprintf(" and trx_hash='%v'", req.Hash)
	}
	if req.Address != "" {
		filterSQL = fmt.Sprintf(" and (owner_address='%v' or to_address='%v')", req.Address, req.Address)
	}
	for _, v := range strings.Split(req.Sort, ",") {
		if strings.Index(v, "timestamp") > 0 {
			if mutiFilter {
				sortSQL = fmt.Sprintf("%v ,", sortSQL)
			}
			sortSQL = fmt.Sprintf("%v create_time", sortSQL)
			if strings.Index(v, "-") == 0 {
				sortSQL = fmt.Sprintf("%v desc", sortSQL)
			}
			mutiFilter = true
		}

		if strings.Index(v, "number") > 0 {
			if mutiFilter {
				sortSQL = fmt.Sprintf("%v ,", sortSQL)
			}
			sortSQL = fmt.Sprintf("%v block_id", sortSQL)
			if strings.Index(v, "-") == 0 {
				sortSQL = fmt.Sprintf("%v desc", sortSQL)
			}
			mutiFilter = true
		}
	}
	if sortSQL != "" {
		if strings.Index(sortSQL, ",") == 0 {
			sortSQL = sortSQL[1:]
		}
		sortSQL = fmt.Sprintf("order by %v", sortSQL)
	}

	pageSQL = fmt.Sprintf("limit %v, %v", req.Start, req.Limit)

	return module.QueryTransactionsRealize(strSQL, filterSQL, sortSQL, pageSQL, true)
}

//QueryTransactionByHashFromBuffer 精确查询
func QueryTransactionByHashFromBuffer(req *entity.Transactions) (*entity.TransactionInfo, error) {
	return buffer.GetBlockBuffer().GetTransactionByHash(req.Hash), nil
}

//QueryTransaction 精确查询  	//number=2135998   TODO: cache
func QueryTransaction(req *entity.Transactions) (*entity.TransactionInfo, error) {
	var filterSQL string
	strSQL := fmt.Sprintf(`
		select block_id,owner_address,to_address,
		trx_hash,contract_data,result_data,fee,
		contract_type,confirmed,create_time,expire_time
		from transactions
			where 1=1 `)

	//按传入条件拼接sql，很容易错误，需要注意
	if req.Number != "" {
		filterSQL = fmt.Sprintf(" and block_id=%v", req.Number)
	}
	if req.Hash != "" {
		filterSQL = fmt.Sprintf(" and trx_hash='%v'", req.Hash)
	}
	return module.QueryTransactionRealize(strSQL, filterSQL)
}

//PostTransaction 创建交易
func PostTransaction(req *entity.PostTransaction, dryRun string) (*entity.PostTransactionResp, error) {
	postResult := &entity.PostTransactionResp{}
	if req.Transaction == "" {
		log.Errorf("no transaction received")
		return postResult, util.NewErrorMsg(util.Error_common_request_json_no_data)
	}
	//将请求转码为transaction结构
	jsonData := req.Transaction
	tranHexData := utils.HexDecode(jsonData)
	transaction := &core.Transaction{}
	if err := proto.Unmarshal(tranHexData, transaction); err != nil || transaction.RawData == nil {
		log.Errorf("pb unmarshal err:[%v];hexData:[%v]", err, tranHexData)
		return postResult, err
	}
	if dryRun != "1" {
		//向主网发布广播
		result, err := GetWalletClient().BroadcastTransaction(transaction)
		if err != nil {
			log.Errorf("call broadcastTransaction err[%v],transaction:[%#v]", err, transaction)
			return postResult, err
		}
		//解析主网接口返回
		postResult.Code = result.Code.String()
		postResult.Message = string(result.Message)
		postResult.Success = result.Result
	}

	//计算前端需要的信息
	postData := &entity.PostTransData{}
	postData.Hash = utils.HexEncode(utils.CalcTransactionHash(transaction))
	postData.Timestamp = transaction.RawData.Timestamp
	contracts := make([]interface{}, 0)
	contractNew := &entity.TransContract{}
	for _, contractOri := range transaction.GetRawData().Contract {
		if contractOri == nil {
			continue
		}
		_, transferContract := utils.GetContractInfoStr2(1, contractOri.Parameter.Value)
		if err := json.Unmarshal([]byte(transferContract), contractNew); err != nil {
			log.Errorf("json unmarshal err:[%v];hexData:[%v]", err, transferContract)
			contracts = append(contracts, transferContract)
			//return postResult, err
		} else {
			contractNew.ContractType = "TransferContract"
			contractNew.ContractTypeID = 1
			contracts = append(contracts, contractNew)
		}
	}
	postData.Contracts = contracts
	postData.Data = string(transaction.RawData.Data)
	signs := make([]*entity.Signatures, 0)

	pubKey, _ := utils.GetSignedPublicKey(transaction)
	address, err := utils.GetTronBase58Address(utils.HexEncode(pubKey))
	if nil != err {
		return postResult, err
	}
	for _, signOri := range transaction.Signature {
		sign := &entity.Signatures{}
		sign.Bytes = utils.Base64Encode(signOri)
		//sign.Bytes1 = base64.RawStdEncoding.EncodeToString(signOri)
		//sign.Bytes2 = base64.RawURLEncoding.EncodeToString(signOri)
		//sign.Bytes3 = base64.StdEncoding.EncodeToString(signOri)
		sign.Address = address
		signs = append(signs, sign)
	}

	postData.Signatures = signs
	postResult.Transaction = postData
	ss, _ := mysql.JSONObjectToString(postResult)
	log.Debugf("Post transaction result:[%v]", ss)
	return postResult, nil
}

// GetWalletClient ...
func GetWalletClient() *WalletClient {
	_walletOnce.Do(func() {
		_wallet = new(WalletClient)
		_wallet.Wallet = grpcclient.GetRandomWallet()
	})
	return _wallet
}
