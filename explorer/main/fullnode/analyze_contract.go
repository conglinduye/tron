package main

import (
	"fmt"

	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
)

func anaylzeTransaction(tran *core.Transaction) {
	if nil != tran && nil != tran.RawData {
		for _, ctxRaw := range tran.RawData.Contract {
			_, ctxIF := utils.GetContract(ctxRaw)

			switch v := ctxIF.(type) {
			case core.AccountCreateContract:
				fmt.Print("%#v\n", v)

			case core.AccountUpdateContract:

			case core.SetAccountIdContract:
			case core.TransferContract:
			case core.TransferAssetContract:
			case core.VoteAssetContract:
			case core.VoteWitnessContract:
			case core.VoteWitnessContract_Vote:
			case core.UpdateSettingContract:
			case core.WitnessCreateContract:
			case core.WitnessUpdateContract:
			case core.AssetIssueContract:
			case core.AssetIssueContract_FrozenSupply:
			case core.ParticipateAssetIssueContract:
			case core.FreezeBalanceContract:
			case core.UnfreezeBalanceContract:
			case core.UnfreezeAssetContract:
			case core.WithdrawBalanceContract:
			case core.UpdateAssetContract:
			case core.ProposalCreateContract:
			case core.ProposalApproveContract:
			case core.ProposalDeleteContract:
			case core.CreateSmartContract:
			case core.TriggerSmartContract:
			case core.BuyStorageContract:
			case core.BuyStorageBytesContract:
			case core.SellStorageContract:
			case core.ExchangeCreateContract:
			case core.ExchangeInjectContract:
			case core.ExchangeWithdrawContract:
			case core.ExchangeTransactionContract:
			default:
				fmt.Println("new type:%#T-->%v\n", v, v)
			}
		}
	}

}

func storeParticipateAssetIssueContract(tran *core.Transaction, ctx *core.ParticipateAssetIssueContract) {

}

func storeTransferAssetIssueContract(tran *core.Transaction, ctx *core.TransferAssetContract) {

}

func storeWitnessCreateContract(tran *core.Transaction, ctx *core.WitnessCreateContract) {

}

func storeTransferContract(tran *core.Transaction, ctx *core.TransferContract) {

}
