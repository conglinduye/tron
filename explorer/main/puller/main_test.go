package main

import (
	"fmt"
	"testing"

	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
)

func TestA(*testing.T) {

	a := make(chan struct{}, 10)
	fmt.Println(len(a))
	for i := 0; i < 10; i++ {
		a <- struct{}{}
	}

	fmt.Println(len(a))
	<-a
	fmt.Println(len(a))
	<-a
	fmt.Println(len(a))

	a <- struct{}{}

	fmt.Println(len(a))

	b(3)

	fmt.Println(utils.HexEncode([]byte{}))
	//fmt.Println(utils.ConverTimestamp(1536348924123))
}

func b(i int) bool {
	fmt.Printf("b %v run\n", i)
	defer c(i)

	if i > 0 {
		return b(i - 1)
	}
	return true
}

func c(i int) {
	fmt.Printf("c %v call\n", i)
}

func TestGetContract(t *testing.T) {
	contracttype := core.Transaction_Contract_ContractType(int32(1))
	contractData := "0a1541cd49ffd0768fb194e39473adf5a36fdb191c3c2e1215419a9708eb1db8d0379315260226ff541119a6826e18ccfbc102"
	da := utils.HexDecode(contractData)
	_, ctx := utils.GetContractByParamVal(contracttype, da)

	/*switch v := ctx.(type) {
	case core.TransferContract:
		fmt.Printf("%#v\n", v)
	default:
		fmt.Printf("ddd%#v\n", v)
	}*/
	ownerCtx, ok := ctx.(core.TransferContract)
	if ok {
		ownerAddr := utils.Base58EncodeAddr(ownerCtx.GetOwnerAddress())
		toAddress := utils.Base58EncodeAddr(ownerCtx.GetToAddress())
		fmt.Printf("%v,%#v", ownerAddr, toAddress)
	}
	fmt.Printf("123")
}
