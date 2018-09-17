package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

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

// GetContractInfoStr ...
func GetContractInfoStr(contract *core.Transaction_Contract) (ownerAddr string, contractDetail string) {
	_, ctx := GetContractByParamVal(contract.Type, contract.Parameter.Value)
	if nil != ctx {
		contractDetail = ToJSONStr(formatContractJSONStr(ToJSONStr(ctx)))
		if ownerIF, ok := ctx.(OwnerAddressIF); ok {
			ownerAddr = Base58EncodeAddr(ownerIF.GetOwnerAddress())
		}
	}
	return
}

// GetContractInfoObj ...
func GetContractInfoObj(contract *core.Transaction_Contract) (ownerAddr string, contractDetail interface{}) {
	_, ctx := GetContractByParamVal(contract.Type, contract.Parameter.Value)
	if nil != ctx {
		contractDetail = formatContractJSONStr(ToJSONStr(ctx))
		if ownerIF, ok := ctx.(OwnerAddressIF); ok {
			ownerAddr = Base58EncodeAddr(ownerIF.GetOwnerAddress())
		}
	}
	return
}

// GetContractInfoStr2 ...
func GetContractInfoStr2(ctxType int32, val []byte) (ownerAddr string, contractDetail string) {

	_, ctx := GetContractByParamVal(core.Transaction_Contract_ContractType(ctxType), val)
	if nil != ctx {
		contractDetail = ToJSONStr(formatContractJSONStr(ToJSONStr(ctx)))
		if ownerIF, ok := ctx.(OwnerAddressIF); ok {
			ownerAddr = Base58EncodeAddr(ownerIF.GetOwnerAddress())
		}
	}
	return
}

// GetContractInfoStr3 ...
func GetContractInfoStr3(ctxType int32, val []byte) (ownerAddr string, contractDetail interface{}) {

	_, ctx := GetContractByParamVal(core.Transaction_Contract_ContractType(ctxType), val)
	if nil != ctx {
		contractDetail = formatContractJSONStr(ToJSONStr(ctx))
		if ownerIF, ok := ctx.(OwnerAddressIF); ok {
			ownerAddr = Base58EncodeAddr(ownerIF.GetOwnerAddress())
		}
	}
	return
}

func formatContractJSONStr(jsonStr string) interface{} {
	var b interface{}

	err := json.Unmarshal([]byte(jsonStr), &b)
	_ = err

	m := b.(map[string]interface{})

	for key, val := range m {
		// fmt.Printf("%v-->%v\n", key, m[key])

		switch v := val.(type) {
		case string:
			m[key] = convertStringVal(key, v)
		case []interface{}:
			m[key] = convertListVal(v)
		}
	}

	return m
}

func convertMapVal(val map[string]interface{}) interface{} {
	for k, v := range val {
		if s, ok := v.(string); ok {
			val[k] = convertStringVal(k, s)
		}
	}
	return val
}

func convertStringVal(key string, val string) string {
	if strings.HasSuffix(key, "address") {
		return Base58EncodeAddr(Base64Decode(val))
	}
	// if strings.HasSuffix(key, "name") || strings.HasSuffix(key, "id") {
	return string(Base64Decode(val))
}

func convertListVal(val []interface{}) interface{} {
	ret := make([]interface{}, 0, len(val))
	for _, subVal := range val {
		switch t := subVal.(type) {
		case map[string]interface{}:
			ret = append(ret, convertMapVal(t))
		case string:
			ret = append(ret, string(Base64Decode(t)))
		default:
			ret = append(ret, t)
		}
	}
	return ret
}

// OwnerAddressIF ...
type OwnerAddressIF interface {
	GetOwnerAddress() []byte
}

// ToAddressIF ...
type ToAddressIF interface {
	GetToAddress() []byte
}

// AmountIF ...
type AmountIF interface {
	GetAmount() int64
}
