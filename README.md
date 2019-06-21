# GM Sumus RPC library
Allows interact with Sumus node via RPC (kinda JSONRPC). \
Sumus node listens at port 4010.

Import:
```sh
go get github.com/void616/gm-sumusrpc
```

## conn package

Conn serves RPC connection. \
Useful for "new-block-mined" event interception on the network. For other cases see pool/rpc below.\
\
Usage:
```go
import "github.com/void616/gm-sumusrpc/conn"

// connect
c, _ := conn.New("127.0.0.1:4010", conn.Options{})
defer c.Close()

// get events channel
cch := c.Subscribe()
defer c.Unsubscribe()

// read channel
event := <- cch
if event.Error != nil {
  // error during recv/send
}
bytes := event.Message
// ...
```

## /pool package

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

// new pool
nodeMaxConns := 32
p := pool.New(&pool.DefaultBalancer{})
defer p.Close()

// add node
p.AddNode("127.0.0.1:4010", nodeMaxConns, conn.Options{})

// get free pool's connection holder
pconn, _ := p.Get(time.Second * 10)
// release back to the pool
defer pconn.Close()

// use the connection
c := pconn.Conn()
// ...
```

## /rpc package

RPC package defines RPC methods and utilizes `*conn.Conn` instance. \
This is a general approach to interact with Sumus node via RPC.

Usage:
```go
import (
  "github.com/void616/gm-sumusrpc/conn"
  "github.com/void616/gm-sumusrpc/rpc"
)

// get a connection
c, _ := conn.New("127.0.0.1:4010", conn.Options{})
defer c.Close()

// call RPC endpoint to get blockchain state
bcstate, rpcCode, err := rpc.BlockchainState(c)
if err != nil {
  // transport error
}
if rpcCode != rpc.ECSuccess {
  // RPC error code, see below
}
// ...
```

## RPC error codes
| Name                                 | Code | Description |
| ---                                  | ---: | --- |
| `ECSuccess`                          | `0`  | Success |
| `ECUnclassified`                     | `1`  | Unclassified error |
| `ECJSONBadRequest`                   | `10` | JSON request error: bad request that cannot be parsed |
| `ECJSONRequestIDNotFound`            | `11` | JSON request error: request with not specified field 'id' |
| `ECJSONUnknownDebugRequest`          | `12` | JSON request error: debug request with unknown ID |
| `ECJSONUnknownJSONRequest`           | `13` | JSON request error: JSON request with unknown ID |
| `ECJSONBadRequestFormat`             | `14` | JSON request error: bad JSON request format e.g. at least one mandatory field not found |
| `ECGetBlockFailure`                  | `20` | Blockchain manager error: get block failure |
| `ECBlockNotFound`                    | `21` | Blockchain manager error: block not found in DB |
| `ECConsensusIsAlreadyStarted`        | `31` | Synchronization error: consensus is already started |
| `ECSynchronizationIsAlreadyStarted`  | `32` | Synchronization error: synchronization is already started |
| `ECBadNodeCount`                     | `33` | Synchronization error: bad node count: number of nodes is greater than number of nodes for voting |
| `ECUnknownNode`                      | `34` | Synchronization error: specified node for manual synchronization is not found in blockchain |
| `ECTransactionNotSigned`             | `40` | Transaction pool error: transaction is not signed |
| `ECBadTransactionSignature`          | `41` | Transaction pool error: bad transaction signature |
| `ECVotingPoolOverflow`               | `42` | Transaction pool error: voting pool overflow |
| `ECPendingPoolOverflow`              | `43` | Transaction pool error: pending pool overflow |
| `ECTransactionWalletNotFound`        | `44` | Transaction pool error: existing transaction wallet not found |
| `ECBadTransactionID`                 | `45` | Transaction pool error: transaction ID is leser or equal to last transaction ID in last approved block |
| `ECBadTransactionDeltaID`            | `46` | Transaction pool error: delta ID is exceeded doubled max size of the block |
| `ECBadTransaction`                   | `47` | Transaction pool error: bad transaction that cannot be applied to its wallet |
| `ECTransactionIDExistsInVotingPool`  | `48` | Transaction pool error: transaction with specified ID already exists in voting pool |
| `ECTransactionIDExistsInPendingPool` | `49` | Transaction pool error: transaction with specified ID already exists in pending pool |
| `ECMaxSizeOfTransactionPackExceeded` | `50` | Transaction pool error: max size of transaction pack exceeded |


## TODO
- Pool tests
- Panic on pool.Close() instead of hanging (leaked/unreleased connections)
