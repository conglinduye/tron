package account

import (
	"fmt"
	"time"

	"github.com/wlcy/tron/explorer/core/grpcclient"
)

// GetRawAccount 从主网获取指定地址的账户的信息
func GetRawAccount(addrs []string) (accList []*Account, err error) {
	wallet := grpcclient.GetRandomWallet()
	if nil != wallet {
		defer func() {
			wallet.Close()
		}()
	}

	errCnt := 0
	addrTaskCnt := len(addrs)
	finished := 0
	for finished < addrTaskCnt {
		ts := time.Now()
		rawAcc, err := wallet.GetAccount(addrs[finished])
		if nil != err {
			errCnt++
			if errCnt > MaxErrCnt {
				wallet.Close()
				wallet = grpcclient.GetRandomWallet()
				continue
			}
		}
		tsc := time.Since(ts)
		ts = time.Now()
		rawAccNet, err := wallet.GetAccountNet(addrs[finished])
		if nil != err {
			errCnt++
			if errCnt > MaxErrCnt {
				wallet.Close()
				wallet = grpcclient.GetRandomWallet()
				continue
			}
		}
		tsc2 := time.Since(ts)
		fmt.Printf("get [%v] cost: [%v], [%v]\n", addrs[finished], tsc, tsc2)

		acc := new(Account)
		acc.SetRaw(rawAcc)
		acc.SetNetRaw(rawAccNet)
		accList = append(accList, acc)
		finished++
	}
	return
}
