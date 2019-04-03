package pool

/*
import (
	"github.com/void616/gm-sumusrpc/conn"
)

// Node implements NodePool
type Node struct {
	ConnectionOpener
	pool  *nodePool
	addr  string
	copts conn.Options
}

// newNode creates new Node instance
func newNode(addr string, concurrency uint16, copts conn.Options) *Node {
	ret := &Node{
		pool:  nil,
		addr:  addr,
		copts: copts,
	}
	ret.pool = newNodePool(ret, nodePoolOptions{
		MaxConnections: concurrency,
	})
	return ret
}

// nodePool gets underlying nodePool
func (n *Node) nodePool() *nodePool {
	return n.pool
}

// OpenConnection creates actual Sumus RPC connection
func (n *Node) OpenConnection() (ConnectionController, error) {
	return newConnHolder(n.addr, n.copts)
}
*/
