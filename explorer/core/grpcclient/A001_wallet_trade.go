package grpcclient

import (
	"github.com/tronprotocol/grpc-gateway/api"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
)

// CreateTransaction ...
func (w *Wallet) CreateTransaction() (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	transfer := &core.TransferContract{}
	trx, err := w.client.CreateTransaction(ctx, transfer, callOpt)

	return trx, err
}

// CreateTransaction2 ...
func (w *Wallet) CreateTransaction2() (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	transfer := &core.TransferContract{}
	trxExt, err := w.client.CreateTransaction2(ctx, transfer, callOpt)

	return trxExt, err
}

// BroadcastTransaction ...
func (w *Wallet) BroadcastTransaction(trx *core.Transaction) (*api.Return, error) {
	ctx, cancel := getTimeoutContext(1000000000000)
	defer cancel()
	callOpt := getDefaultCallOptions()

	ret, err := w.client.BroadcastTransaction(ctx, trx, callOpt)

	return ret, err
}

// UpdateAccount ...
func (w *Wallet) UpdateAccount(addr string, name string) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.AccountUpdateContract{}
	contract.AccountName = []byte(name)
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	trx, err := w.client.UpdateAccount(ctx, contract, callOpt)

	return trx, err
}

// SetAccountID ...
func (w *Wallet) SetAccountID(addr string, accountID string) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.SetAccountIdContract{}
	contract.AccountId = []byte(accountID)
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	trx, err := w.client.SetAccountId(ctx, contract, callOpt)

	return trx, err
}

// UpdateAccount2 ...
func (w *Wallet) UpdateAccount2(addr string, name string) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.AccountUpdateContract{}
	contract.AccountName = []byte(name)
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	trxExt, err := w.client.UpdateAccount2(ctx, contract, callOpt)

	return trxExt, err
}

// VoteWitnessAccount ...
func (w *Wallet) VoteWitnessAccount(addr string, votes map[string]int64) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.VoteWitnessContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)

	for witness, vote := range votes {
		contract.Votes = append(contract.Votes, &core.VoteWitnessContract_Vote{
			VoteAddress: utils.Base58DecodeAddr(witness),
			VoteCount:   vote,
		})
	}
	trx, err := w.client.VoteWitnessAccount(ctx, contract, callOpt)

	return trx, err
}

// UpdateSetting ... 智能合约相关
func (w *Wallet) UpdateSetting(addr string, contractAddr string) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.UpdateSettingContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.ContractAddress = utils.Base58DecodeAddr(contractAddr)

	trxExt, err := w.client.UpdateSetting(ctx, contract, callOpt)

	return trxExt, err
}

// VoteWitnessAccount2 ...
func (w *Wallet) VoteWitnessAccount2(addr string, votes map[string]int64) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.VoteWitnessContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)

	for witness, vote := range votes {
		contract.Votes = append(contract.Votes, &core.VoteWitnessContract_Vote{
			VoteAddress: utils.Base58DecodeAddr(witness),
			VoteCount:   vote,
		})
	}
	trxExt, err := w.client.VoteWitnessAccount2(ctx, contract, callOpt)

	return trxExt, err
}

// CreateAssetIssue ...
func (w *Wallet) CreateAssetIssue(addr string, assetIssue *core.AssetIssueContract) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.AssetIssueContract{}
	*contract = *assetIssue
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)

	trx, err := w.client.CreateAssetIssue(ctx, contract, callOpt)

	return trx, err
}

// CreateAssetIssue2 ...
func (w *Wallet) CreateAssetIssue2() error {
	return nil
}

// UpdateWitness ...
func (w *Wallet) UpdateWitness(addr string, url string) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.WitnessUpdateContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.UpdateUrl = []byte(url)

	trx, err := w.client.UpdateWitness(ctx, contract, callOpt)

	return trx, err
}

// UpdateWitness2 ...
func (w *Wallet) UpdateWitness2() error {
	return nil
}

// CreateAccount ...
func (w *Wallet) CreateAccount(addr, newAddr string) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.AccountCreateContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.AccountAddress = utils.Base58DecodeAddr(newAddr)

	trx, err := w.client.CreateAccount(ctx, contract, callOpt)

	return trx, err
}

// CreateAccount2 ...
func (w *Wallet) CreateAccount2() error {
	return nil
}

// CreateWitness ...
func (w *Wallet) CreateWitness(addr, url string) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.WitnessCreateContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.Url = []byte(url)

	trx, err := w.client.CreateWitness(ctx, contract, callOpt)

	return trx, err
}

// CreateWitness2 ...
func (w *Wallet) CreateWitness2() error {
	return nil
}

// TransferAsset ...
func (w *Wallet) TransferAsset(addr, toAddr, assetName string, amount int64) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.TransferAssetContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.ToAddress = []byte(toAddr)
	contract.Amount = amount
	contract.AssetName = []byte(assetName)

	trx, err := w.client.TransferAsset(ctx, contract, callOpt)

	return trx, err
}

// TransferAsset2 ...
func (w *Wallet) TransferAsset2() error {
	return nil
}

// ParticipateAssetIssue ...
func (w *Wallet) ParticipateAssetIssue(addr, toAddr, assetName string, amount int64) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.ParticipateAssetIssueContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.ToAddress = []byte(toAddr)
	contract.Amount = amount
	contract.AssetName = []byte(assetName)

	trx, err := w.client.ParticipateAssetIssue(ctx, contract, callOpt)

	return trx, err
}

// ParticipateAssetIssue2 ...
func (w *Wallet) ParticipateAssetIssue2() error {
	return nil
}

// FreezeBalance ...
func (w *Wallet) FreezeBalance(addr string, amount, duration int64) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.FreezeBalanceContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.FrozenBalance = amount
	contract.FrozenDuration = duration

	trx, err := w.client.FreezeBalance(ctx, contract, callOpt)

	return trx, err
}

// FreezeBalance2 ...
func (w *Wallet) FreezeBalance2() error {
	return nil
}

// UnfreezeBalance ...
func (w *Wallet) UnfreezeBalance(addr string) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.UnfreezeBalanceContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)

	trx, err := w.client.UnfreezeBalance(ctx, contract, callOpt)

	return trx, err
}

// UnfreezeBalance2 ...
func (w *Wallet) UnfreezeBalance2() error {
	return nil
}

// UnfreezeAsset ...
func (w *Wallet) UnfreezeAsset(addr string) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.UnfreezeAssetContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)

	trx, err := w.client.UnfreezeAsset(ctx, contract, callOpt)

	return trx, err
}

// UnfreezeAsset2 ...
func (w *Wallet) UnfreezeAsset2() error {
	return nil
}

// WithdrawBalance ...
func (w *Wallet) WithdrawBalance(addr string) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.WithdrawBalanceContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)

	trx, err := w.client.WithdrawBalance(ctx, contract, callOpt)

	return trx, err
}

// WithdrawBalance2 ...
func (w *Wallet) WithdrawBalance2() error {
	return nil
}

// UpdateAsset ...
func (w *Wallet) UpdateAsset(addr string, desc, url string, newLimit, newPublicLimit int64) (*core.Transaction, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.UpdateAssetContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.Description = []byte(desc)
	contract.Url = []byte(url)
	contract.NewLimit = newLimit
	contract.NewPublicLimit = newPublicLimit

	trx, err := w.client.UpdateAsset(ctx, contract, callOpt)

	return trx, err
}

// UpdateAsset2 ...
func (w *Wallet) UpdateAsset2() error {
	return nil
}

// ProposalCreate ...
func (w *Wallet) ProposalCreate(addr string, params map[int64]int64) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.ProposalCreateContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.Parameters = params

	trxExt, err := w.client.ProposalCreate(ctx, contract, callOpt)

	return trxExt, err
}

// ProposalApprove ...
func (w *Wallet) ProposalApprove(addr string, proposalID int64, isAdd bool) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.ProposalApproveContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.ProposalId = proposalID
	contract.IsAddApproval = isAdd

	trxExt, err := w.client.ProposalApprove(ctx, contract, callOpt)

	return trxExt, err
}

// ProposalDelete ...
func (w *Wallet) ProposalDelete(addr string, proposalID int64) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.ProposalDeleteContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.ProposalId = proposalID

	trxExt, err := w.client.ProposalDelete(ctx, contract, callOpt)

	return trxExt, err
}

// BuyStorage ...
func (w *Wallet) BuyStorage(addr string, quant int64) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.BuyStorageContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.Quant = quant

	trxExt, err := w.client.BuyStorage(ctx, contract, callOpt)

	return trxExt, err
}

// BuyStorageBytes ...
func (w *Wallet) BuyStorageBytes(addr string, byteCount int64) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.BuyStorageBytesContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.Bytes = byteCount

	trxExt, err := w.client.BuyStorageBytes(ctx, contract, callOpt)

	return trxExt, err
}

// SellStorage ...
func (w *Wallet) SellStorage(addr string, byteCount int64) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.SellStorageContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.StorageBytes = byteCount

	trx, err := w.client.SellStorage(ctx, contract, callOpt)

	return trx, err
}

// ExchangeCreate ...
func (w *Wallet) ExchangeCreate(addr string, fromTokenID, toTokenID string, fromTokenAmount, toTokenAmount int64) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.ExchangeCreateContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.FirstTokenId = []byte(fromTokenID)
	contract.SecondTokenId = []byte(toTokenID)
	contract.FirstTokenBalance = fromTokenAmount
	contract.SecondTokenBalance = toTokenAmount

	trxExt, err := w.client.ExchangeCreate(ctx, contract, callOpt)

	return trxExt, err
}

// ExchangeInject ...
func (w *Wallet) ExchangeInject(addr string, exchangeID int64, tokenID string, quant int64) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.ExchangeInjectContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.ExchangeId = exchangeID
	contract.TokenId = []byte(tokenID)
	contract.Quant = quant

	trxExt, err := w.client.ExchangeInject(ctx, contract, callOpt)

	return trxExt, err
}

// ExchangeWithdraw ...
func (w *Wallet) ExchangeWithdraw(addr string, exchangeID int64, tokenID string, quant int64) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.ExchangeWithdrawContract{}
	contract.OwnerAddress = utils.Base58DecodeAddr(addr)
	contract.ExchangeId = exchangeID
	contract.TokenId = []byte(tokenID)
	contract.Quant = quant

	trxExt, err := w.client.ExchangeWithdraw(ctx, contract, callOpt)

	return trxExt, err
}

// ExchangeTransaction ...
func (w *Wallet) ExchangeTransaction(addr string, exchangeID int64, tokenID string, quant int64) (*api.TransactionExtention, error) {
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()

	contract := &core.ExchangeTransactionContract{}
	contract.ExchangeId = exchangeID
	contract.TokenId = []byte(tokenID)
	contract.Quant = quant

	trxExt, err := w.client.ExchangeTransaction(ctx, contract, callOpt)

	return trxExt, err
}
