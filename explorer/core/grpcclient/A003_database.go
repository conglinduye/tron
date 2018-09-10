package grpcclient

import (
	"fmt"

	"github.com/tronprotocol/grpc-gateway/api"
	"github.com/tronprotocol/grpc-gateway/core"
	"github.com/wlcy/tron/explorer/core/utils"
)

// Database grpc wallet client wrapper
//	目前了解到 Database 是 FullNode 实现
type Database struct {
	_conn
	client api.DatabaseClient
}

// NewDatabase create new wallet grpc client
func NewDatabase(serverAddr string) *Database {
	ret := &Database{}
	ret.serverAddr = serverAddr
	return ret
}

// GetRandomDatabase ...
func GetRandomDatabase() *Database {
	serverAddr := fmt.Sprintf("%v:%v", utils.GetRandFullNodeAddr(), utils.DefaultGrpPort)
	ret := &Database{}
	ret.serverAddr = serverAddr
	ret.Connect()
	return ret
}

// Connect estable connect to server
func (d *Database) Connect() (err error) {
	err = d._conn.Connect()
	if nil != err {
		return err
	}

	d.client = api.NewDatabaseClient(d.c)

	if nil == d.client {
		return utils.ErrorCreateGrpClient
	}

	return nil
}

// GetBlockReference ...
func (d *Database) GetBlockReference() (*api.BlockReference, error) {
	//(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*BlockReference, error)

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	msg := &api.EmptyMessage{}

	blockRef, err := d.client.GetBlockReference(ctx, msg, callOpt)

	return blockRef, err

}

// GetDynamicProperties ...
func (d *Database) GetDynamicProperties() (*core.DynamicProperties, error) {
	// (ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*core.DynamicProperties, error)

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	msg := &api.EmptyMessage{}

	prop, err := d.client.GetDynamicProperties(ctx, msg, callOpt)

	return prop, err

}

// GetNowBlock ...
func (d *Database) GetNowBlock() (*core.Block, error) {
	//  (ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*core.Block, error)
	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	emptyMsg := &api.EmptyMessage{}

	block, err := d.client.GetNowBlock(ctx, emptyMsg, callOpt)

	return block, err
}

// GetBlockByNum ...
func (d *Database) GetBlockByNum(num int64) (*core.Block, error) {

	ctx, cancel := getTimeoutContext(defaultTimeout)
	defer cancel()
	callOpt := getDefaultCallOptions()
	numMsg := &api.NumberMessage{Num: num}

	block, err := d.client.GetBlockByNum(ctx, numMsg, callOpt)
	return block, err
}
