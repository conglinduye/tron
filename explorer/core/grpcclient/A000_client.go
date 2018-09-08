package grpcclient

import (
	"google.golang.org/grpc"
)

type _conn struct {
	c          *grpc.ClientConn
	serverAddr string
}

// Connect 尝试建立连接
func (c *_conn) Connect() (err error) {
	c.c, err = grpc.Dial(c.serverAddr, grpc.WithInsecure())
	if nil != err {
		return err
	}

	return
}

// GetState 返回连接状态
func (c *_conn) GetState() string {
	return c.c.GetState().String()
}

// Target 返回连接端
func (c *_conn) Target() string {
	return c.c.Target()
}
