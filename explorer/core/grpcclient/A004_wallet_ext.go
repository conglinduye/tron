package grpcclient

import (
	"fmt"

	"github.com/tronprotocol/grpc-gateway/api"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
)

// WalletExt grpc wallet client wrapper
type WalletExt struct {
	_conn
	client api.WalletExtensionClient
}

// NewWalletExt create new wallet grpc client
func NewWalletExt(serverAddr string) *WalletExt {
	ret := &WalletExt{}
	ret.serverAddr = serverAddr
	return ret
}

// GetRandomWalletExt ...
func GetRandomWalletExt() *WalletExt {
	serverAddr := fmt.Sprintf("%v:%v", utils.GetRandSolidityNodeAddr(), utils.DefaultGrpPort)
	ret := &WalletExt{}
	ret.serverAddr = serverAddr
	ret.Connect()
	return ret
}

// Connect estable connect to server
func (g *WalletExt) Connect() (err error) {
	err = g._conn.Connect()
	if nil != err {
		return err
	}

	g.client = api.NewWalletExtensionClient(g.c)

	if nil == g.client {
		return utils.ErrorCreateGrpClient
	}

	return nil
}

/*
	GetTransactionsFromThis(ctx context.Context, in *AccountPaginated, opts ...grpc.CallOption) (*TransactionList, error)
	// Use this function instead of GetTransactionsFromThis.
	GetTransactionsFromThis2(ctx context.Context, in *AccountPaginated, opts ...grpc.CallOption) (*TransactionListExtention, error)
	// Please use GetTransactionsToThis2 instead of this function.
	GetTransactionsToThis(ctx context.Context, in *AccountPaginated, opts ...grpc.CallOption) (*TransactionList, error)
	// Use this function instead of GetTransactionsToThis.
	GetTransactionsToThis2(ctx context.Context, in *AccountPaginated, opts ...grpc.CallOption) (*TransactionListExtention, error)
}
*/

// GetTransactionsFromThis ...
func (g *WalletExt) GetTransactionsFromThis(addr string, offset, limit int64) ([]*core.Transaction, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	msg := &api.AccountPaginated{}
	account := &core.Account{}
	account.Address = utils.Base58DecodeAddr(addr)
	msg.Account = account
	msg.Offset = offset
	msg.Limit = limit

	tranList, err := g.client.GetTransactionsFromThis(ctx, msg, callOpt)

	if nil == tranList {
		return nil, err
	}

	return tranList.Transaction, err

}

// GetTransactionsFromThis2 ...
func (g *WalletExt) GetTransactionsFromThis2(addr string, offset, limit int64) ([]*api.TransactionExtention, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	msg := &api.AccountPaginated{}
	account := &core.Account{}
	account.Address = utils.Base58DecodeAddr(addr)
	msg.Account = account
	msg.Offset = offset
	msg.Limit = limit

	tranList, err := g.client.GetTransactionsFromThis2(ctx, msg, callOpt)

	if nil == tranList {
		return nil, err
	}

	return tranList.Transaction, err

}

// GetTransactionsToThis ...
func (g *WalletExt) GetTransactionsToThis(addr string, offset, limit int64) ([]*core.Transaction, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	msg := &api.AccountPaginated{}
	account := &core.Account{}
	account.Address = utils.Base58DecodeAddr(addr)
	msg.Account = account
	msg.Offset = offset
	msg.Limit = limit

	tranList, err := g.client.GetTransactionsToThis(ctx, msg, callOpt)

	if nil == tranList {
		return nil, err
	}

	return tranList.Transaction, err

}

// GetTransactionsToThisi2 ...
// 默认就是倒叙
func (g *WalletExt) GetTransactionsToThisi2(addr string, offset, limit int64) ([]*api.TransactionExtention, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	msg := &api.AccountPaginated{}
	account := &core.Account{}
	account.Address = utils.Base58DecodeAddr(addr)
	msg.Account = account
	msg.Offset = offset
	msg.Limit = limit

	tranList, err := g.client.GetTransactionsToThis2(ctx, msg, callOpt)

	if nil == tranList {
		return nil, err
	}

	return tranList.Transaction, err

}
