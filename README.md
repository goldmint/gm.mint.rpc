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
