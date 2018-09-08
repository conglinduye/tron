package grpcclient

import (
	"fmt"
	"testing"

	"github.com/wlcy/tron/explorer/core/utils"
)

func TestDatabaseClient(*testing.T) {

	client := NewDatabase(fmt.Sprintf("%s:50051", utils.GetRandFullNode()))

	err := client.Connect()
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println(client.GetState(), client.Target())

	utils.VerifyCall(client.GetDynamicProperties())
	utils.VerifyCall(client.GetNowBlock())

}
