package main

import (
	"fmt"
	"sync"

	"github.com/tronprotocol/grpc-gateway/core"
)

var contractBufferMap sync.Map // contract_type -> trans

func anaylzeTransaction(trx *transaction) {
	if nil != trx && nil != trx.contract {

		switch v := trx.contract.(type) {
		case *core.AccountCreateContract:
			handleAccountCreateContract(trx, v)

		case *core.AccountUpdateContract:
			handleAccountUpdateContract(trx, v)

		case *core.SetAccountIdContract:
			handleSetAccountIDContract(trx, v)

		case *core.TransferContract:
			handleTransferContract(trx, v)

		case *core.TransferAssetContract:
			handleTransferAssetContract(trx, v)

		case *core.VoteAssetContract:

		case *core.VoteWitnessContract:
			handleVoteWitnessContract(trx, v)

		case *core.VoteWitnessContract_Vote: // no api use this type, no OwnerAddress
		case *core.UpdateSettingContract:
		case *core.WitnessCreateContract:
			handleWitnessCreateContract(trx, v)

		case *core.WitnessUpdateContract:
		case *core.AssetIssueContract:
			handleAssetIssueContract(trx, v)

		case *core.AssetIssueContract_FrozenSupply: // no api use this type, no OwnerAddress
		case *core.ParticipateAssetIssueContract:
			handleParticipateAssetIssueContract(trx, v)

		case *core.FreezeBalanceContract:
			handleFreezeBalanceContract(trx, v)

		case *core.UnfreezeBalanceContract:
			handleUnfreezeBalanceContract(trx, v)

		case *core.UnfreezeAssetContract:
			handleUnfreezeAssetContract(trx, v)

		case *core.WithdrawBalanceContract:
			handleWithdrawBalanceContract(trx, v)

		case *core.UpdateAssetContract:
		case *core.ProposalCreateContract:
		case *core.ProposalApproveContract:
		case *core.ProposalDeleteContract:
		case *core.CreateSmartContract:
		case *core.TriggerSmartContract:
		case *core.BuyStorageContract:
		case *core.BuyStorageBytesContract:
		case *core.SellStorageContract:
		case *core.ExchangeCreateContract:
		case *core.ExchangeInjectContract:
		case *core.ExchangeWithdrawContract:
		case *core.ExchangeTransactionContract:
		default:
			fmt.Printf("new type:%T-->%v\n", v, v)
		}
	}

}

func handleAccountCreateContract(trx *transaction, ctx *core.AccountCreateContract) {
	var buff []*core.AccountCreateContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.AccountCreateContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.AccountCreateContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress(), ctx.GetAccountAddress())

	contractBufferMap.Store(trx.ctxType, buff)
}

func handleAccountUpdateContract(trx *transaction, ctx *core.AccountUpdateContract) {
	var buff []*core.AccountUpdateContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.AccountUpdateContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.AccountUpdateContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress())

	contractBufferMap.Store(trx.ctxType, buff)

}

func handleSetAccountIDContract(trx *transaction, ctx *core.SetAccountIdContract) {
	var buff []*core.SetAccountIdContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.SetAccountIdContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.SetAccountIdContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress())

	contractBufferMap.Store(trx.ctxType, buff)
}

func handleVoteWitnessContract(trx *transaction, ctx *core.VoteWitnessContract) {
	var buff []*core.VoteWitnessContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.VoteWitnessContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.VoteWitnessContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress())

	contractBufferMap.Store(trx.ctxType, buff)
}

func handleWitnessCreateContract(trx *transaction, ctx *core.WitnessCreateContract) {
	var buff []*core.WitnessCreateContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.WitnessCreateContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.WitnessCreateContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress())

	contractBufferMap.Store(trx.ctxType, buff)
}

func handleAssetIssueContract(trx *transaction, ctx *core.AssetIssueContract) {
	var buff []*core.AssetIssueContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.AssetIssueContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.AssetIssueContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress())

	contractBufferMap.Store(trx.ctxType, buff)
}

func handleParticipateAssetIssueContract(trx *transaction, ctx *core.ParticipateAssetIssueContract) {
	var buff []*core.ParticipateAssetIssueContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.ParticipateAssetIssueContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.ParticipateAssetIssueContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress(), ctx.GetToAddress())

	contractBufferMap.Store(trx.ctxType, buff)
}

func handleFreezeBalanceContract(trx *transaction, ctx *core.FreezeBalanceContract) {
	var buff []*core.FreezeBalanceContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.FreezeBalanceContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.FreezeBalanceContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress())

	contractBufferMap.Store(trx.ctxType, buff)
}

func handleUnfreezeBalanceContract(trx *transaction, ctx *core.UnfreezeBalanceContract) {
	var buff []*core.UnfreezeBalanceContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.UnfreezeBalanceContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.UnfreezeBalanceContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress())

	contractBufferMap.Store(trx.ctxType, buff)
}

func handleUnfreezeAssetContract(trx *transaction, ctx *core.UnfreezeAssetContract) {
	var buff []*core.UnfreezeAssetContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.UnfreezeAssetContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.UnfreezeAssetContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress())

	contractBufferMap.Store(trx.ctxType, buff)
}

func handleWithdrawBalanceContract(trx *transaction, ctx *core.WithdrawBalanceContract) {
	var buff []*core.WithdrawBalanceContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.WithdrawBalanceContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.WithdrawBalanceContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress())

	contractBufferMap.Store(trx.ctxType, buff)
}

func handleTransferContract(trx *transaction, ctx *core.TransferContract) {
	var buff []*core.TransferContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.TransferContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.TransferContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress(), ctx.GetToAddress())

	contractBufferMap.Store(trx.ctxType, buff)
}

func handleTransferAssetContract(trx *transaction, ctx *core.TransferAssetContract) {
	var buff []*core.TransferAssetContract
	buffRaw, ok := contractBufferMap.Load(trx.ctxType)
	if ok {
		buff, ok = buffRaw.([]*core.TransferAssetContract)
		if ok && nil != buff {
			buff = append(buff, ctx)
		}
	} else {
		buff = make([]*core.TransferAssetContract, 0, 2000)
	}

	buff = append(buff, ctx)

	AddRefreshAddress(ctx.GetOwnerAddress(), ctx.GetToAddress())

	contractBufferMap.Store(trx.ctxType, buff)
}
