package request

import (
	"context"

	"github.com/void616/gm-sumusrpc/conn"
	"github.com/void616/gm-sumusrpc/rpc"
)

// BlockchainNode model
type BlockchainNode struct {
	Index     uint32 `json:"index"`
	PublicKey string `json:"public_key"`
	IP        string `json:"ip"`
}

// GetBlockchainNodes method
func GetBlockchainNodes(ctx context.Context, c *conn.Conn) (res []BlockchainNode, rerr *rpc.Error, err error) {
	res, rerr, err = []BlockchainNode{}, nil, nil

	rctx, rcancel := c.Receive(ctx)
	defer rcancel()

	msg, err := c.Request(rctx, "get_blockchain_nodes", nil)
	if err != nil {
		return
	}

	switch m := msg.(type) {
	case *rpc.Error:
		rerr = m
		return
	case *rpc.Result:
		err = m.Parse(&res)
		return
	}
	return
}