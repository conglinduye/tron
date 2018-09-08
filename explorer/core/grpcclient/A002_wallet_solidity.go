package grpcclient

import (
	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/tronprotocol/grpc-gateway/api"
	"github.com/tronprotocol/grpc-gateway/core"
)

// WalletSolidity grpc wallet client wrapper
type WalletSolidity struct {
	_conn
	client api.WalletSolidityClient
}

// NewWalletSolidity create new wallet grpc client
func NewWalletSolidity(serverAddr string) *WalletSolidity {
	ret := &WalletSolidity{}
	ret.serverAddr = serverAddr
	return ret
}

// Connect estable connect to server
func (ws *WalletSolidity) Connect() (err error) {
	err = ws._conn.Connect()
	if nil != err {
		return err
	}

	ws.client = api.NewWalletSolidityClient(ws.c)

	if nil == ws.client {
		return utils.ErrorCreateGrpClient
	}

	return nil
}

// GetAccount 获取账户信息
func (ws *WalletSolidity) GetAccount(addr string) (*core.Account, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	account := &core.Account{}
	account.Address = utils.Base58DecodeAddr(addr)
	account, err := ws.client.GetAccount(ctx, account, callOpt)

	return account, err
}

// GetAccountByID 获取账户信息
//	当前Account返回的数据中AccountID都为空，所以此接口的入参现在不可得
func (ws *WalletSolidity) GetAccountByID(accountID string) (*core.Account, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	account := &core.Account{}
	// account.Address = utils.Base58DecodeAddr(accountID)
	account.AccountId = []byte(accountID)
	account, err := ws.client.GetAccountById(ctx, account, callOpt)

	return account, err
}

// ListWitnesses 获取见证人节点账户列表
func (ws *WalletSolidity) ListWitnesses() ([]*core.Witness, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	witnessList, err := ws.client.ListWitnesses(ctx, emptyMsg, callOpt)

	if nil == witnessList {
		return nil, err
	}
	return witnessList.Witnesses, err
}

// GetAssetIssueList 获取通证信息
func (ws *WalletSolidity) GetAssetIssueList() ([]*core.AssetIssueContract, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	assetIssueList, err := ws.client.GetAssetIssueList(ctx, emptyMsg, callOpt)

	if nil == assetIssueList {
		return nil, err
	}
	return assetIssueList.AssetIssue, err
}

// GetPaginatedAssetIssueList 分页获取通证信息
//	page: page size start from 1, default 1
//	limit: should >= 1, default 1
func (ws *WalletSolidity) GetPaginatedAssetIssueList(page, limit int64) ([]*core.AssetIssueContract, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	if limit < 1 {
		limit = 1
	}
	if page < 1 {
		page = 1
	}
	pagingMsg := &api.PaginatedMessage{}
	pagingMsg.Limit = limit
	pagingMsg.Offset = limit * (page - 1)

	assetIssueList, err := ws.client.GetPaginatedAssetIssueList(ctx, pagingMsg, callOpt)

	if nil == assetIssueList {
		return nil, err
	}
	return assetIssueList.AssetIssue, err
}

// GetNowBlock 获取区块信息
func (ws *WalletSolidity) GetNowBlock() (*core.Block, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	block, err := ws.client.GetNowBlock(ctx, emptyMsg, callOpt)

	return block, err
}

/* Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.WalletSolidity/GetNowBlock2
// GetNowBlock2 获取区块扩展信息
func (ws *WalletSolidity) GetNowBlock2() (*api.BlockExtention, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	blockExt, err := ws.client.GetNowBlock2(ctx, emptyMsg, callOpt)

	return blockExt, err
}
*/

// GetBlockByNum 获取区块信息
func (ws *WalletSolidity) GetBlockByNum(num int64) (*core.Block, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	numMsg := &api.NumberMessage{Num: num}

	block, err := ws.client.GetBlockByNum(ctx, numMsg, callOpt)
	return block, err
}

/* Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.WalletSolidity/GetBlockByNum2
// GetBlockByNum2 获取区块扩展信息
func (ws *WalletSolidity) GetBlockByNum2(num int64) (*api.BlockExtention, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	numMsg := &api.NumberMessage{Num: num}

	blockExt, err := ws.client.GetBlockByNum2(ctx, numMsg, callOpt)

	return blockExt, err
}
*/

/* Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.WalletSolidity/GetTransactionCountByBlockNum
// GetTransactionCountByBlockNum 获取区块交易数
func (ws *WalletSolidity) GetTransactionCountByBlockNum(num int64) (int64, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	numMsg := &api.NumberMessage{Num: num}

	numRet, err := ws.client.GetTransactionCountByBlockNum(ctx, numMsg, callOpt)

	if nil == numRet {
		return 0, err
	}
	return numRet.Num, err
}
*/

// GetTransactionByID 通过交易hash获取交易内容
func (ws *WalletSolidity) GetTransactionByID(id string) (*core.Transaction, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	byteMsg := &api.BytesMessage{Value: utils.HexDecode(id)}

	transaction, err := ws.client.GetTransactionById(ctx, byteMsg, callOpt)

	return transaction, err
}

// GetTransactionInfoByID 通过交易hash获取信息交易信息
//	返回交易花费，所属区块，区块时间戳
func (ws *WalletSolidity) GetTransactionInfoByID(id string) (*core.TransactionInfo, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	byteMsg := &api.BytesMessage{Value: utils.HexDecode(id)}

	transInfo, err := ws.client.GetTransactionInfoById(ctx, byteMsg, callOpt)

	return transInfo, err
}

// GenerateAddress 生成新的地址秘钥对, 返回地址和私钥
//	{"address":"TYrK6ZH9Ji8HcewcAL4mRcF1WUSU5zbNgp","privateKey":"d2b2f8e5c93e2dc884ee2a24ad5eab830e1d72a0a7de3791432fac8ae021ef05"}
func (ws *WalletSolidity) GenerateAddress() (*api.AddressPrKeyPairMessage, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	addrMsg, err := ws.client.GenerateAddress(ctx, emptyMsg, callOpt)

	return addrMsg, err
}
