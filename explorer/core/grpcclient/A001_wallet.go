package grpcclient

import (
	"time"

	"github.com/wlcy/tron/explorer/core/utils"
	"github.com/tronprotocol/grpc-gateway/api"
	"github.com/tronprotocol/grpc-gateway/core"
)

// Wallet grpc wallet client wrapper
type Wallet struct {
	_conn
	client api.WalletClient
}

// NewWallet create new wallet grpc client
func NewWallet(serverAddr string) *Wallet {
	ret := &Wallet{}
	ret.serverAddr = serverAddr
	return ret
}

// Connect estable connect to server
func (w *Wallet) Connect() (err error) {
	err = w._conn.Connect()
	if nil != err {
		return err
	}

	w.client = api.NewWalletClient(w.c)

	if nil == w.client {
		return utils.ErrorCreateGrpClient
	}

	return nil
}

// cat solidity_node_list.txt|grep -v "//" |grep -v "}"|awk -F '(' '{print $1}' | awk '{print "\n\n// "$1" ...\nfunc (w *Wallet) "$1" () (error) {\n\treturn nil\n}"}'

// GetAccount 获取账户信息
func (w *Wallet) GetAccount(addr string) (*core.Account, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	account := &core.Account{}
	account.Address = utils.Base58DecodeAddr(addr)
	account, err := w.client.GetAccount(ctx, account, callOpt)

	return account, err
}

// GetAccountByID 获取账户信息
//	当前Account返回的数据中AccountID都为空，所以此接口的入参现在不可得
func (w *Wallet) GetAccountByID(accountID string) (*core.Account, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	account := &core.Account{}
	// account.Address = utils.Base58DecodeAddr(accountID)
	account.AccountId = []byte(accountID)
	account, err := w.client.GetAccountById(ctx, account, callOpt)

	return account, err
}

// ListNodes 返回节点列表
func (w *Wallet) ListNodes() ([]*api.Node, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	nodeList, err := w.client.ListNodes(ctx, emptyMsg, callOpt)

	if nil == nodeList {
		return nil, err
	}
	return nodeList.Nodes, err
}

// GetAssetIssueByAccount 根据账户地址查询账户发型的通证信息
//	{"owner_address":"QU0e+Gc/kW3rt+JRWo8+yvJhEDSq","name":"U0VFRA==","abbr":"U0VFRA==","total_supply":100000000000,"trx_num":1000000,"num":1,"start_time":1529987043000,"end_time":1530342060000,"description":"U2VzYW1lc2VlZCB0b2tlbnMgZm9yIGNvbW11bml0eSByZXdhcmRzIGFuZCBTRUVEZ2VybWluYXRvciBpbnZlc3RtZW50IG9mIGNvbW11bml0eS12b3RlZCBwcm9qZWN0cy4=","url":"aHR0cDovL3d3dy5zZXNhbWVzZWVkLm9yZw=="}
func (w *Wallet) GetAssetIssueByAccount(addr string) ([]*core.AssetIssueContract, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	account := &core.Account{}
	account.Address = utils.Base58DecodeAddr(addr)
	assetIssueList, err := w.client.GetAssetIssueByAccount(ctx, account, callOpt)

	if nil == assetIssueList {
		return nil, err
	}
	return assetIssueList.AssetIssue, nil
}

// GetAccountNet ...
func (w *Wallet) GetAccountNet(addr string) (*api.AccountNetMessage, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	account := &core.Account{}
	account.Address = utils.Base58DecodeAddr(addr)
	accountNet, err := w.client.GetAccountNet(ctx, account, callOpt)

	return accountNet, err
}

/*
// GetAccountResource ...
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/GetAccountResource
func (w *Wallet) GetAccountResource(addr string) (*api.AccountResourceMessage, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	account := &core.Account{}
	account.Address = utils.Base58DecodeAddr(addr)
	accountRes, err := w.client.GetAccountResource(ctx, account, callOpt)

	return accountRes, err
}
*/

// GetAssetIssueByName 根据通证名称查询通证信息
func (w *Wallet) GetAssetIssueByName(name string) (*core.AssetIssueContract, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	msg := &api.BytesMessage{}
	msg.Value = []byte(name)
	assetIssue, err := w.client.GetAssetIssueByName(ctx, msg, callOpt)

	return assetIssue, err
}

// GetNowBlock 获取最新区块信息
func (w *Wallet) GetNowBlock() (*core.Block, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	block, err := w.client.GetNowBlock(ctx, emptyMsg, callOpt)

	return block, err
}

// // GetNowBlock2 ...
// func (w *Wallet) GetNowBlock2() error {
// 	return nil
// }

// GetBlockByNum 获取区块信息
func (w *Wallet) GetBlockByNum(num int64) (*core.Block, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	numMsg := &api.NumberMessage{Num: num}

	block, err := w.client.GetBlockByNum(ctx, numMsg, callOpt)
	return block, err
}

// // GetBlockByNum2 ...
// func (w *Wallet) GetBlockByNum2() error {
// 	return nil
// }

// // GetTransactionCountByBlockNum ...
// func (w *Wallet) GetTransactionCountByBlockNum() error {
// 	return nil
// }

// GetBlockByID 根据区块hash获取区块信息 id 为hex encoding编码
func (w *Wallet) GetBlockByID(id string) (*core.Block, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	msg := &api.BytesMessage{}
	msg.Value = utils.HexDecode(id)
	block, err := w.client.GetBlockById(ctx, msg, callOpt)

	return block, err
}

// GetBlockByLimitNext 批量范围区块，numStart为最小区块号，numEnd为最大区块号，返回的为最大-1
func (w *Wallet) GetBlockByLimitNext(numStart, numEnd int64) ([]*core.Block, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	msg := &api.BlockLimit{}
	msg.StartNum = numStart
	msg.EndNum = numEnd
	blockList, err := w.client.GetBlockByLimitNext(ctx, msg, callOpt)
	if nil == blockList {
		return nil, err
	}

	return blockList.Block, err
}

/*
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/GetBlockByLimitNext2
// GetBlockByLimitNext2 ...
func (w *Wallet) GetBlockByLimitNext2(numStart, numEnd int64) ([]*api.BlockExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	msg := &api.BlockLimit{}
	msg.StartNum = numStart
	msg.EndNum = numEnd
	blockListExt, err := w.client.GetBlockByLimitNext2(ctx, msg, callOpt)

	if nil == blockListExt {
		return nil, err
	}
	return blockListExt.Block, err
}
*/

// GetBlockByLatestNum 获取指定个数的最新块
func (w *Wallet) GetBlockByLatestNum(num int64) ([]*core.Block, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	msg := &api.NumberMessage{}
	msg.Num = num
	blockList, err := w.client.GetBlockByLatestNum(ctx, msg, callOpt)

	if nil == blockList {
		return nil, err
	}

	return blockList.Block, err
}

/*
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/GetBlockByLatestNum2
// GetBlockByLatestNum2 获取指定个数的最新块扩展信息
func (w *Wallet) GetBlockByLatestNum2(num int64) ([]*api.BlockExtention, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	msg := &api.NumberMessage{}
	msg.Num = num
	blockExtList, err := w.client.GetBlockByLatestNum2(ctx, msg, callOpt)

	if nil == blockExtList {
		return nil, err
	}

	return blockExtList.Block, err
}
*/

// GetTransactionByID 通过交易hash获取交易内容
func (w *Wallet) GetTransactionByID(id string) (*core.Transaction, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	byteMsg := &api.BytesMessage{Value: utils.HexDecode(id)}

	transaction, err := w.client.GetTransactionById(ctx, byteMsg, callOpt)

	return transaction, err
}

// DeployContract 部署智能合约
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/DeployContract
func (w *Wallet) DeployContract(addr string, smartContract *core.SmartContract) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	csc := &core.CreateSmartContract{}
	csc.OwnerAddress = utils.Base58DecodeAddr(addr)
	csc.NewContract = smartContract

	result, err := w.client.DeployContract(ctx, csc, callOpt)

	return result, err
}

// GetContract 通过智能合约地址获取智能合约信息
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/DeployContract
func (w *Wallet) GetContract(contractAddr string) (*core.SmartContract, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	msg := &api.BytesMessage{}
	msg.Value = utils.Base58DecodeAddr(contractAddr)

	smartContract, err := w.client.GetContract(ctx, msg, callOpt)

	return smartContract, err
}

// TriggerContract 触发智能合约
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/TriggerContract
func (w *Wallet) TriggerContract(addr string, contractAddr string, callValue int64, data []byte) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	msg := &core.TriggerSmartContract{}
	msg.OwnerAddress = utils.Base58DecodeAddr(addr)
	msg.ContractAddress = utils.Base58DecodeAddr(contractAddr)
	msg.CallValue = callValue
	msg.Data = data

	result, err := w.client.TriggerContract(ctx, msg, callOpt)

	return result, err
}

// ListWitnesses 获取见证人节点账户列表
func (w *Wallet) ListWitnesses() ([]*core.Witness, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	witnessList, err := w.client.ListWitnesses(ctx, emptyMsg, callOpt)

	if nil == witnessList {
		return nil, err
	}
	return witnessList.Witnesses, err
}

// ListProposals ...
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/ListProposals
func (w *Wallet) ListProposals() ([]*core.Proposal, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	proposalList, err := w.client.ListProposals(ctx, emptyMsg, callOpt)

	if nil == proposalList {
		return nil, err
	}
	return proposalList.Proposals, err
}

// GetProposalByID ...
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/GetProposalById
func (w *Wallet) GetProposalByID(id string) (*core.Proposal, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	msg := &api.BytesMessage{}
	msg.Value = utils.HexDecode(id)

	proposal, err := w.client.GetProposalById(ctx, msg, callOpt)

	return proposal, err
}

// ListExchanges ...
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/ListExchanges
func (w *Wallet) ListExchanges() ([]*core.Exchange, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	exchangeList, err := w.client.ListExchanges(ctx, emptyMsg, callOpt)

	if nil == exchangeList {
		return nil, err
	}
	return exchangeList.Exchanges, err
}

// GetExchangeByID ...
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/GetExchangeById
func (w *Wallet) GetExchangeByID(id string) (*core.Exchange, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	msg := &api.BytesMessage{}
	msg.Value = utils.HexDecode(id)

	exchange, err := w.client.GetExchangeById(ctx, msg, callOpt)

	return exchange, err
}

// GetChainParameters ...
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/GetChainParameters
func (w *Wallet) GetChainParameters() (*core.ChainParameters, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	msg := &api.EmptyMessage{}

	chainParam, err := w.client.GetChainParameters(ctx, msg, callOpt)

	return chainParam, err
}

// GetAssetIssueList ...
func (w *Wallet) GetAssetIssueList() ([]*core.AssetIssueContract, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	assetIssueList, err := w.client.GetAssetIssueList(ctx, emptyMsg, callOpt)

	if nil == assetIssueList {
		return nil, err
	}
	return assetIssueList.AssetIssue, err
}

// GetPaginatedAssetIssueList 分页获取通证信息, 偏移量 和
func (w *Wallet) GetPaginatedAssetIssueList(offset, limit int64) ([]*core.AssetIssueContract, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	// if limit < 1 {
	// 	limit = 1
	// }
	// if page < 1 {
	// 	page = 1
	// }
	pagingMsg := &api.PaginatedMessage{}
	pagingMsg.Limit = limit
	pagingMsg.Offset = offset
	// pagingMsg.Offset = limit * (page - 1)

	assetIssueList, err := w.client.GetPaginatedAssetIssueList(ctx, pagingMsg, callOpt)

	if nil == assetIssueList {
		return nil, err
	}
	return assetIssueList.AssetIssue, err
}

// TotalTransaction 获取总交易数, 这个接口比较慢
func (w *Wallet) TotalTransaction() (int64, error) {
	ctx, cancel := getTimeoutContext(20 * time.Second)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	numMsg, err := w.client.TotalTransaction(ctx, emptyMsg, callOpt)

	if nil == numMsg {
		return 0, err
	}
	return numMsg.Num, err
}

// GetNextMaintenanceTime 下一次维护时间，返回为毫秒
func (w *Wallet) GetNextMaintenanceTime() (int64, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	numMsg, err := w.client.GetNextMaintenanceTime(ctx, emptyMsg, callOpt)

	if nil == numMsg {
		return 0, err
	}
	return numMsg.Num, err
}

// GetTransactionSign 交易签名
func (w *Wallet) GetTransactionSign(trans *core.Transaction, privKey string) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	transSign := &core.TransactionSign{}
	transSign.Transaction = trans
	transSign.PrivateKey = utils.HexDecode(privKey)

	signedTran, err := w.client.GetTransactionSign(ctx, transSign, callOpt)
	return signedTran, err
}

// GetTransactionSign2 ...
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/GetTransactionSign2
func (w *Wallet) GetTransactionSign2(trans *core.Transaction, privKey string) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	transSign := &core.TransactionSign{}
	transSign.Transaction = trans
	transSign.PrivateKey = utils.HexDecode(privKey)

	signedTranExt, err := w.client.GetTransactionSign2(ctx, transSign, callOpt)
	return signedTranExt, err
}

// CreateAddress ...
// Faield, error:rpc error: code = Unimplemented desc = Method not found: protocol.Wallet/CreateAddress
func (w *Wallet) CreateAddress(pubKey string) (string, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	msg := &api.BytesMessage{}
	msg.Value = utils.HexDecode(pubKey)

	result, err := w.client.CreateAddress(ctx, msg, callOpt)

	if nil == result {
		return "", err
	}

	return utils.HexEncode(result.Value), err
}

// EasyTransfer 简易转账
func (w *Wallet) EasyTransfer(passPhrase, addr string, amount int64) (*api.EasyTransferResponse, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	msg := &api.EasyTransferMessage{}
	msg.PassPhrase = utils.HexDecode(passPhrase)
	msg.ToAddress = utils.Base58DecodeAddr(addr)
	msg.Amount = amount

	result, err := w.client.EasyTransfer(ctx, msg, callOpt)

	return result, err
}

// EasyTransferByPrivate 简易转账，发起人私钥，接收账户地址，金额（trx, not sun）
func (w *Wallet) EasyTransferByPrivate(privKey, addr string, amount int64) (*api.EasyTransferResponse, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	msg := &api.EasyTransferByPrivateMessage{}
	msg.PrivateKey = utils.HexDecode(privKey)
	msg.ToAddress = utils.Base58DecodeAddr(addr)
	msg.Amount = amount

	result, err := w.client.EasyTransferByPrivate(ctx, msg, callOpt)

	return result, err
}

// GenerateAddress ...
func (w *Wallet) GenerateAddress() (*api.AddressPrKeyPairMessage, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	addrMsg, err := w.client.GenerateAddress(ctx, emptyMsg, callOpt)

	return addrMsg, err
}

// GetTransactionInfoByID 通过交易hash获取信息交易信息
//	返回交易花费，所属区块，区块时间戳
func (w *Wallet) GetTransactionInfoByID(id string) (*core.TransactionInfo, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	byteMsg := &api.BytesMessage{Value: utils.HexDecode(id)}

	transInfo, err := w.client.GetTransactionInfoById(ctx, byteMsg, callOpt)

	return transInfo, err
}
