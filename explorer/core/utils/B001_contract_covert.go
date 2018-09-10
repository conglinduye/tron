package utils

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/tronprotocol/grpc-gateway/core"
)

// contractTypeMap 类型字典
// grep -E "^type .* struct" Contract.pb.go | awk -v q=\" '{print q$2q":reflect.TypeOf(core."$2"{}),"}'
var contractTypeMap = map[string]reflect.Type{
	"AccountCreateContract":           reflect.TypeOf(core.AccountCreateContract{}),
	"AccountUpdateContract":           reflect.TypeOf(core.AccountUpdateContract{}),
	"SetAccountIdContract":            reflect.TypeOf(core.SetAccountIdContract{}),
	"TransferContract":                reflect.TypeOf(core.TransferContract{}),
	"TransferAssetContract":           reflect.TypeOf(core.TransferAssetContract{}),
	"VoteAssetContract":               reflect.TypeOf(core.VoteAssetContract{}),
	"VoteWitnessContract":             reflect.TypeOf(core.VoteWitnessContract{}),
	"VoteWitnessContract_Vote":        reflect.TypeOf(core.VoteWitnessContract_Vote{}),
	"UpdateSettingContract":           reflect.TypeOf(core.UpdateSettingContract{}),
	"WitnessCreateContract":           reflect.TypeOf(core.WitnessCreateContract{}),
	"WitnessUpdateContract":           reflect.TypeOf(core.WitnessUpdateContract{}),
	"AssetIssueContract":              reflect.TypeOf(core.AssetIssueContract{}),
	"AssetIssueContract_FrozenSupply": reflect.TypeOf(core.AssetIssueContract_FrozenSupply{}),
	"ParticipateAssetIssueContract":   reflect.TypeOf(core.ParticipateAssetIssueContract{}),
	"FreezeBalanceContract":           reflect.TypeOf(core.FreezeBalanceContract{}),
	"UnfreezeBalanceContract":         reflect.TypeOf(core.UnfreezeBalanceContract{}),
	"UnfreezeAssetContract":           reflect.TypeOf(core.UnfreezeAssetContract{}),
	"WithdrawBalanceContract":         reflect.TypeOf(core.WithdrawBalanceContract{}),
	"UpdateAssetContract":             reflect.TypeOf(core.UpdateAssetContract{}),
	"ProposalCreateContract":          reflect.TypeOf(core.ProposalCreateContract{}),
	"ProposalApproveContract":         reflect.TypeOf(core.ProposalApproveContract{}),
	"ProposalDeleteContract":          reflect.TypeOf(core.ProposalDeleteContract{}),
	"CreateSmartContract":             reflect.TypeOf(core.CreateSmartContract{}),
	"TriggerSmartContract":            reflect.TypeOf(core.TriggerSmartContract{}),
	"BuyStorageContract":              reflect.TypeOf(core.BuyStorageContract{}),
	"BuyStorageBytesContract":         reflect.TypeOf(core.BuyStorageBytesContract{}),
	"SellStorageContract":             reflect.TypeOf(core.SellStorageContract{}),
	"ExchangeCreateContract":          reflect.TypeOf(core.ExchangeCreateContract{}),
	"ExchangeInjectContract":          reflect.TypeOf(core.ExchangeInjectContract{}),
	"ExchangeWithdrawContract":        reflect.TypeOf(core.ExchangeWithdrawContract{}),
	"ExchangeTransactionContract":     reflect.TypeOf(core.ExchangeTransactionContract{}),
}

// GetContract 根据协议内通用contract解析出实际的类型数据
// in:
//	contract: core.Transaction的contract对象
// out:
// 	reflect.Type: 实际协议类型
//	interface{}: 实际协议对象
func GetContract(contract *core.Transaction_Contract) (reflect.Type, interface{}) {
	// ctxTypeInfo := strings.Split(contract.Parameter.TypeUrl, ".")
	// if len(ctxTypeInfo) > 0 {
	// 	ctxTypeName := ctxTypeInfo[len(ctxTypeInfo)-1] // . 分割的类型的最后一个字段
	// 	ctxType, ok := contractTypeMap[ctxTypeName]
	// 	if ok {
	// 		ret := reflect.New(ctxType).Interface().(proto.Message)
	// 		proto.Unmarshal(contract.Parameter.Value, ret)

	// 		return ctxType, ret
	// 	}
	// }
	// return nil, nil
	return GetContractByParamVal(contract.Type, contract.Parameter.Value)
}

// GetContractByParamVal 获取实际的协议内容
func GetContractByParamVal(ctxType core.Transaction_Contract_ContractType, val []byte) (reflect.Type, interface{}) {
	ctxRealType, ok := contractTypeMap[ctxType.String()]
	if ok {
		ret := reflect.New(ctxRealType).Interface().(proto.Message)
		err := proto.Unmarshal(val, ret)
		if nil != err {
			fmt.Printf("convert contract failed:%v, type:%s, value:%v", err, ctxType, HexEncode(val))
			return nil, nil
		}
		return ctxRealType, ret
	}
	return nil, nil
}
