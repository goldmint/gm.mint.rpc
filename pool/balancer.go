package pool

import (
	"fmt"
)

// DefaultBalancer is default balancer
type DefaultBalancer struct {
	Balancer
}

// Switch picks a node from the pool of nodes
func (b *DefaultBalancer) get(nodes map[string]*nodePool) (*nodePool, error) {
	at := ""
	weight := int64((^uint64(0)) >> 1)
	for i, p := range nodes {
		if p.Available() {
			w := int64(p.ConsumedConnections()) + int64(p.PendingConsumers())
			if weight > w {
				weight = w
				at = i
			}
		}
	}

	if at == "" {
		return nil, fmt.Errorf("Unfortunately, there isn't a free node")
	}

	return nodes[at], nil
}
