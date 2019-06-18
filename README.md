# GM Sumus RPC library

Allows interact with Sumus node via RPC.

Import:
```sh
go get github.com/void616/gm-sumusrpc
```

## /conn

RPC enpoint connection.

Usage:
```go
import "github.com/void616/gm-sumusrpc/conn"

c, _ := conn.New("127.0.0.1:4010", conn.Options{})
defer c.Close()
```

## /pool

Pool holds connections to each added node.
```
                    |--[Connection]
         |--[Node]--|--[Connection]
[Pool]---|
         |--[Node]--|--[Connection]
                    |--[Connection]
```

Usage:
```go
import (
  "github.com/void616/gm-sumusrpc/conn"
  "github.com/void616/gm-sumusrpc/pool"
)

nodeMaxConns := 32
p := pool.New(&pool.DefaultBalancer{})
defer p.Close()

// add node
p.AddNode("127.0.0.1:4010", nodeMaxConns, conn.Options{})

// get free connection with timeout
c, _ := p.Get(time.Second * 10)
// release back to the pool
defer c.Close()

// use connection
c.Conn()
```

## /rpc

Contains methods to communicate to Sumus node via RPC.

Usage:
```go
import (
  "github.com/void616/gm-sumusrpc/conn"
  "github.com/void616/gm-sumusrpc/rpc"
)

c, _ := conn.New("127.0.0.1:4010", conn.Options{})
defer c.Close()

bcstate, code, err := rpc.BlockchainState(c)
if err != nil {
  // transport error
}
if code != rpc.ECSuccess {
  // request error
}
...
```

Error codes:

|ECSuccess                          |Success|
|ECUnclassified                     |Unclassified error|
|ECJSONBadRequest                   |JSON request error: bad request that cannot be parsed|
|ECJSONRequestIDNotFound            |JSON request error: request with not specified field 'id'|
|ECJSONUnknownDebugRequest          |JSON request error: debug request with unknown ID|
|ECJSONUnknownJSONRequest           |JSON request error: JSON request with unknown ID|
|ECJSONBadRequestFormat             |JSON request error: bad JSON request format e.g. at least one mandatory field not found|
|ECGetBlockFailure                  |Blockchain manager error: get block failure|
|ECBlockNotFound                    |Blockchain manager error: block not found in DB|
|ECConsensusIsAlreadyStarted        |Synchronization error: consensus is already started|
|ECSynchronizationIsAlreadyStarted  |Synchronization error: synchronization is already started|
|ECBadNodeCount                     |Synchronization error: bad node count: number of nodes is greater than number of nodes for voting|
|ECUnknownNode                      |Synchronization error: specified node for manual synchronization is not found in blockchain|
|ECTransactionNotSigned             |Transaction pool error: transaction is not signed|
|ECBadTransactionSignature          |Transaction pool error: bad transaction signature|
|ECVotingPoolOverflow               |Transaction pool error: voting pool overflow|
|ECPendingPoolOverflow              |Transaction pool error: pending pool overflow|
|ECTransactionWalletNotFound        |Transaction pool error: existing transaction wallet not found|
|ECBadTransactionID                 |Transaction pool error: transaction ID is leser or equal to last transaction ID in last approved block|
|ECBadTransactionDeltaID            |Transaction pool error: delta ID is exceeded doubled max size of the block|
|ECBadTransaction                   |Transaction pool error: bad transaction that cannot be applied to its wallet|
|ECTransactionIDExistsInVotingPool  |Transaction pool error: transaction with specified ID already exists in voting pool|
|ECTransactionIDExistsInPendingPool |Transaction pool error: transaction with specified ID already exists in pending pool|
|ECMaxSizeOfTransactionPackExceeded |Transaction pool error: max size of transaction pack exceeded|


## TODO
- Logging;
- Panic on pool.Close() instead of hanging (leaked connections);