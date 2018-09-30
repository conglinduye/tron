package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wlcy/tron/explorer/lib/log"
	"github.com/wlcy/tron/explorer/lib/mysql"
	"gopkg.in/mgo.v2/bson"
)

//Contracts 查询确认后的智能合约列表的请求参数
type Contracts struct {
	Sort    string `json:"sort,omitempty"`    // 按时间戳倒序
	Limit   int64  `json:"limit,omitempty"`   // 每页记录数
	Count   string `json:"count,omitempty"`   // 是否返回总数
	Total   int64  `json:"total,omitempty"`   // 上一次分页查询的总数，当且仅当分页查询接口使用
	Start   int64  `json:"start,omitempty"`   // 记录的起始序号
	Address string `json:"address,omitempty"` // 合约地址
	Type    string `json:"type,omitempty"`    // 合约类型  internal  token
}

//ContractsResp 查询确认后的智能合约列表的结果
type ContractsResp struct {
	Total  int64               `json:"total"`  // 总记录数
	Status *State              `json:"status"` // 状态
	Data   []*ContractListInfo `json:"data"`   // 记录详情
}

//State 返回状态信息
type State struct {
	Code    int    `json:"code"`    // 接口返回状态码
	Message string `json:"message"` // 状态描述
}

//ContractListInfo 智能合约信息
type ContractListInfo struct {
	Address      string `json:"address"`      //:"TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBi", //合约地址
	Name         string `json:"name"`         //:"xxxxxxx",//合约名称
	Compiler     string `json:"compiler"`     //:"v0.0.1",//编译器版本
	Balance      int64  `json:"balance"`      //:100,//余额
	TrxCount     int64  `json:"trxCount"`     //:6000, //交易数量
	IsSetting    bool   `json:"isSetting"`    //:true,//是否优化
	DateVerified int64  `json:"dateVerified"` //:1531711638107,//合约验证时间
}

//###############合约详情##################

//ContractBaseResp 查询确认后的智能合约
type ContractBaseResp struct {
	Status *State              `json:"status"` // 状态
	Data   []*ContractBaseInfo `json:"data"`   // 记录详情
}

//ContractBaseInfo 合约详情基础信息
type ContractBaseInfo struct {
	Address      string   `json:"address"`      //:"TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBi", //合约地址
	Balance      int64    `json:"balance"`      //余额
	BalanceInUsd int64    `json:"balanceInUsd"` //usd 余额
	TrxCount     int64    `json:"trxCount"`     //交易数量
	TokenTracke  string   `json:"tokenTracke"`  //
	Creator      *Creator `json:"creator"`      //
}

//Creator 创建人信息
type Creator struct {
	Address      string `json:"address"`       //:"TAUN6FwrnwwmaEqYcckffC7wYmbaS6cBi",
	TxHash       string `json:"txHash"`        //:"******************",
	TokenBalance int64  `json:"token_balance"` //:1000
}

//###############合约交易列表##################

//ContractTransactionResp 合约交易列表
type ContractTransactionResp struct {
	Total  int64                      `json:"total"`  // 总记录数
	Status *State                     `json:"status"` // 状态
	Data   []*ContractTransactionInfo `json:"data"`   // 记录详情
}

//ContractTransactionInfo 合约交易详情
type ContractTransactionInfo struct {
	TxHash     string `json:"txHash"`     //:"******************", //交易hash
	ParentHash string `json:"parentHash"` //:"******************", //交易父hash only for internal交易
	Block      int64  `json:"block"`      //:66000,//区块高度
	Timestamp  int64  `json:"timestamp"`  //:1531711638107,//交易时间
	OwnAddress string `json:"ownAddress"` //:"*************",//发起人地址
	ToAddress  string `json:"toAddress"`  //:"************", //接收人地址
	Value      int64  `json:"value"`      //:10,//
	TxFee      int64  `json:"txFee"`      //:1,//交易费
	Token      string `json:"token"`      //:"ATRON",//token only for token 交易
}

//###############合约code列表##################

//ContractCodeResp 合约Code列表
type ContractCodeResp struct {
	Status *State            `json:"status"` // 状态
	Data   *ContractCodeInfo `json:"data"`   // 记录详情
}

//ContractCodeInfo 合约code信息
type ContractCodeInfo struct {
	Address    string     `json:"address"`    //:"***********",//合约地址
	Name       string     `json:"name"`       //:"*****",//合约名称
	Compiler   string     `json:"compiler"`   //:"v4.0.0",//编译器版本
	IsSetting  bool       `json:"isSetting"`  //:true,//是否优化
	Source     string     `json:"source"`     //:"***********",//合约源代码
	ByteCode   string     `json:"byteCode"`   //:"*****",//编译生成的二进制代码
	ABI        string     `json:"abi"`        //:"*****",//编译生成的abi
	AbiEncoded string     `json:"abiEncoded"` //:"********",//编译所需参数
	Librarys   []*Library `json:"librarys"`   //:
}

//InsertContractCode ...
func (w *ContractCodeInfo) InsertContractCode() error {
	library := ""
	if len(w.Librarys) > 0 {
		value, err := json.Marshal(w.Librarys)
		if err != nil {
			log.Errorf("InsertContractCode marshal json err:[%v]", err)
		}
		library = string(value)
	}
	isSetting := 0
	if w.IsSetting {
		isSetting = 1
	}
	strSQL := fmt.Sprintf(`
		insert into wlcy_smart_contract 
		(address,contract_name,compiler_version,is_optimized,verify_time,source_code,
			byte_code,abi,abi_encoded,contract_library)
		values('%v','%v','%v','%v','%v','%v','%v','%v','%v','%v')`,
		w.Address, w.Name, w.Compiler, isSetting, time.Now().UnixNano(), w.Source, w.ByteCode, w.ABI,
		w.AbiEncoded, library)
	insID, _, err := mysql.ExecuteSQLCommand(strSQL, true)
	if err != nil {
		log.Errorf("insert wlcy_smart_contract  fail:[%v]  sql:%s", err, strSQL)
		return err
	}
	log.Debugf("insert wlcy_smart_contract  success, insert id: [%v]", insID)
	return nil
}

//Library 合约代码库信息
type Library struct {
	Index   int64  `json:"index"`   //:1,//库序号，最大五个library
	Name    string `json:"name"`    //:"name",//名称
	Address string `json:"address"` //:"address",//地址
}

//###############合约event列表##################

//ContractEventResp 合约event列表
type ContractEventResp struct {
	Status *State               `json:"status"` // 状态
	Data   []*ContractEventInfo `json:"data"`   // 记录详情
}

//ContractEventInfo 合约event信息
type ContractEventInfo struct {
	TxHash    string `json:"txHash"`    //:"******************", //交易hash
	Block     int64  `json:"block"`     //:66000,//区块高度
	Timestamp int64  `json:"timestamp"` //:1531711638107,//交易时间
	Method    string `json:"method"`    //:"transfer(address,uint256)",//方法名
	EventLog  string `json:"eventLog"`  //:"************"//event log 日志
}

//###############合约内部交易列表##################

//ContractInternalTxsResp 合约内部交易列表
type ContractInternalTxsResp struct {
	Total  int64                      `json:"total"`  // 总记录数
	Status *State                     `json:"status"` // 状态
	Data   []*ContractInternalTxsInfo `json:"data"`   // 记录详情
}

//ContractInternalTxsInfo 内部交易信息
type ContractInternalTxsInfo struct {
	Block        int64  `json:"block"`        //:21232, //区块id
	Timestamp    int64  `json:"timestamp"`    //:1531711638107,//交易时间
	ParentHash   string `json:"parentHash"`   //:"******************",//父hash
	TxType       string `json:"txType"`       //:"call",//交易类型
	OwnerAddress string `json:"ownerAddress"` //:"**************", //发起人地址
	ToAddress    string `json:"toAddress"`    //:"**************",//接收人地址
	Value        int64  `json:"value"`        //:100,//交易额
	TxFee        int64  `json:"txFee"`        //:1,//交易费
}

//EventLogResp eventlog 返回内容
type EventLogResp struct {
	Status *State      `json:"status"` // 状态
	Data   []*EventLog `json:"data"`   // 记录详情
}

//EventLog eventlog
type EventLog struct {
	ID              bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	BlockNum        int64         `json:"block_number,omitempty" bson:"block_number,omitempty"`
	BlockTimestamp  int64         `json:"block_timestamp,omitempty" bson:"block_timestamp,omitempty"`
	ContractAddress string        `json:"contract_address,omitempty" bson:"contract_address,omitempty"`
	EventName       string        `json:"event_name,omitempty" bson:"event_name,omitempty"`
	Raw             interface{}   `json:"raw,omitempty" bson:"raw,omitempty"`
	Class           string        `json:"_class,omitempty" bson:"_class,omitempty"`
	Result          interface{}   `json:"result,omitempty" bson:"result,omitempty"`
	TransactionID   string        `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
}
