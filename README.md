Allows interact with Mint node via JSONRPC. \
By default node listens port TCP 4010.

```sh
go get github.com/void616/gm.mint.rpc
```

See `/examples`

# Messaging

Messages are delimited with zero-terminator (\0).

## Request format

Valid JsonRPC body followed by zero-terminator. \
For instance (via netcat):
```sh
nc localhost 4010
{"method":"get_blockchain_state","params":{},"id":1}
# Press Ctrl+Shift+2 or Ctrl+2 to print zero-terminator, then Enter
```

## Response format

Example response for the request above:
```js
{
	"method": "get_blockchain_state",
	"id": 1,
	"result": {
		"last_block_digest": "umPfEnn7qmARdPtxTbxv4NdXpaP5dZ36GgoRxHMxjji97P2cg",
		"last_block_merkle_root": "Law4XPsP5Q113RMfNKm8xcmn55q2RAB2kdjG76rMTjBBjHHZ8",
		"block_count": "1300",
		"transaction_count": "2984",
		// ...
	}
} // zero-terminator here
```

## Events

Node-side initiated events are delivered on the same connection. For instance, when new blocks are received by the node:
```js
{
	"method": "blocks_synchronized",
	"params": {
		"count": "1",
		"last_id": "1300",
		"last_digest": "w15pNFWu2aAXkh7rnKs2uqBx4yK1M4nH1X3vf2r1ZS7CpRvsx"
	}
}
// zero-terminator here
```

# API

## Blockchain

### `get_blockchain_info`
Node and API version, supported transactions, assets and wallet tags. \
Params:
```js
{}
```

### `get_blockchain_state`
Blockchain height, overall balance, transaction pool state. \
Params:
```js
{}
```

### `get_blockchain_nodes`
Peer nodes list. \
Params:
```js
{}
```

### `get_block`
Binary data of the block defined by ID or digest. \
Params:
```js
{ 
	"id": "1", // numeric block ID, optional
	"digest": "2rY5Kri4t22T8xCsXfJVNV7MUYxT7npGNFYSGQWr5b3UfA1hKS" // block digest, optional
}
```

### `add_transaction`
Add a transaction to the pending pool. \
Params:
```js
{
	"name": "transfer_asset", // predefined transaction name
	"data": "deadbeef" // hex-encoded transaction data
}
```

### `synchronize`
Force the node to synchronize blocks. \
Params:
```js
{
	"mode": "regular", // one of regular|fast|manual
	"nodes": "1,2,3" // node IDs in manual mode, optional
}
``` 

## Address / wallet

### `get_wallet_state`
Balance, nonce and tags of the specific address. \
Params:
```js
{
	"public_key": "111111111111111111111111111111115RyRTN" // address
}
```

### `get_wallet_transactions`
List of incoming/outgoing transactions for the specific address. \
Params:
```js
{
	"public_key": "111111111111111111111111111111115RyRTN", // address
	"count": 1000, // limited output, default is 1000
	"pool_lookup": true, // include pending transactions
	"binary": true, // binary/textual transaction representation
	"incoming": true, // include incoming transactions, optional
	"outgoing": true // include outgoing transactions, optional
}
```

### `dump_wallets`
Dump existing addresses (LZ4 compressed). \
Params:
```js
{
	"with_mnt": true, // include only adresses with non-zero MNT balance, optional
	"with_gold": true // include only adresses with non-zero GOLD balance, optional
}
```
