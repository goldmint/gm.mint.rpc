package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	sumuslib "github.com/void616/gm-sumuslib"
	"github.com/void616/gm-sumuslib/amount"
	"github.com/void616/gm-sumusrpc/conn"
)

// AddTransaction posts new transaction
func AddTransaction(c *conn.Conn, t sumuslib.Transaction, hexdata string) (err error) {
	err = nil

	req := struct {
		TransactionName string `json:"transaction_name,omitempty"`
		TransactionData string `json:"transaction_data,omitempty"`
	}{
		t.String(),
		hexdata,
	}
	res := struct{}{}

	err = RawCall(c, "add-transaction", &req, &res)
	return
}

// WalletState gets state of specific wallet
func WalletState(c *conn.Conn, address string) (state WalletStateResult, err error) {
	state, err = WalletStateResult{
		Balance: WalletBalanceResult{
			Gold: amount.NewInteger(0),
			Mnt:  amount.NewInteger(0),
		},
	}, nil

	req := struct {
		PublicKey string `json:"public_key,omitempty"`
	}{
		address,
	}
	res := struct {
		Balance           *json.RawMessage `json:"balance,omitempty"`
		Exist             int              `json:"exist,string,omitempty"`
		LastTransactionID uint64           `json:"last_transaction_id,string,omitempty"`
		Tags              *json.RawMessage `json:"tags,omitempty"`
	}{}
	type BalanceItem struct {
		AssetCode string `json:"asset_code,omitempty"`
		Amount    string `json:"amount,omitempty"`
	}

	err = RawCall(c, "get-wallet-state", &req, &res)
	if err != nil {
		return
	}

	// balance exists
	if res.Balance == nil {
		err = fmt.Errorf("Balance field not set")
		return
	}

	// balance is array or string - try parse as array
	balance := []BalanceItem{}
	if perr := json.Unmarshal(*res.Balance, &balance); perr == nil {
		for _, v := range balance {
			if parsed := amount.NewFloatString(v.Amount); parsed != nil {
				if token, perr := sumuslib.ParseToken(v.AssetCode); perr == nil {
					switch token {
					case sumuslib.TokenMNT:
						state.Balance.Mnt = parsed
					case sumuslib.TokenGOLD:
						state.Balance.Gold = parsed
					}
				}
			}
		}
	}

	// tags is array or string - try parse as array
	tags := []string{}
	json.Unmarshal(*res.Tags, &tags)

	state.Exists = res.Exist == 1
	state.ApprovedNonce = res.LastTransactionID
	state.Tags = tags
	return
}

// BlockchainState gets blockchain state info
func BlockchainState(c *conn.Conn) (state BlockchainStateResult, err error) {
	state, err = BlockchainStateResult{
		BlockCount:              new(big.Int),
		LastBlockDigest:         "",
		LastBlockMerkleRoot:     "",
		TransactionCount:        new(big.Int),
		NodeCount:               new(big.Int),
		NonEmptyWalletCount:     new(big.Int),
		VotingTransactionCount:  new(big.Int),
		PendingTransactionCount: new(big.Int),
		BlockchainState:         "",
		ConsensusRound:          "",
		VotingNodes:             "",
	}, nil

	req := struct{}{}
	res := struct {
		BlockCount              string `json:"block_count,omitempty"`
		LastBlockDigest         string `json:"last_block_digest,omitempty"`
		LastBlockMerkleRoot     string `json:"last_block_merkle_root,omitempty"`
		TransactionCount        string `json:"transaction_count,omitempty"`
		NodeCount               string `json:"node_count,omitempty"`
		NonEmptyWalletCount     string `json:"non_empty_wallet_count,omitempty"`
		VotingTransactionCount  string `json:"voing_transaction_count,omitempty"`
		PendingTransactionCount string `json:"pending_transaction_count,omitempty"`
		BlockchainState         string `json:"blockchain_state,omitempty"`
		ConsensusRound          string `json:"consensus_round,omitempty"`
		VotingNodes             string `json:"voting_nodes,omitempty"`
	}{}

	err = RawCall(c, "get-blockchain-state", &req, &res)
	if err != nil {
		return
	}

	// blocks count
	if i, ok := new(big.Int).SetString(res.BlockCount, 10); ok {
		state.BlockCount = i
	} else {
		err = errors.New("Failed to parse blocks count")
		return
	}

	// tx count
	if i, ok := new(big.Int).SetString(res.TransactionCount, 10); ok {
		state.TransactionCount = i
	} else {
		err = errors.New("Failed to parse transactions count")
		return
	}

	// node count
	if i, ok := new(big.Int).SetString(res.NodeCount, 10); ok {
		state.NodeCount = i
	} else {
		err = errors.New("Faield to parse nodes count")
		return
	}

	// non-empty wallets
	if i, ok := new(big.Int).SetString(res.NonEmptyWalletCount, 10); ok {
		state.NonEmptyWalletCount = i
	} else {
		err = errors.New("Faield to parse non-empty wallets count")
		return
	}

	// voting transactions
	if i, ok := new(big.Int).SetString(res.VotingTransactionCount, 10); ok {
		state.VotingTransactionCount = i
	} else {
		err = errors.New("Faield to parse voting transactions count")
		return
	}

	// pending transactions
	if i, ok := new(big.Int).SetString(res.PendingTransactionCount, 10); ok {
		state.PendingTransactionCount = i
	} else {
		err = errors.New("Faield to parse pending transactions count")
		return
	}

	state.LastBlockDigest = res.LastBlockDigest
	state.LastBlockMerkleRoot = res.LastBlockMerkleRoot
	state.BlockchainState = res.BlockchainState
	state.ConsensusRound = res.ConsensusRound
	state.VotingNodes = res.VotingNodes

	return
}

// BlockData gets raw block data by ID
func BlockData(c *conn.Conn, id *big.Int) (data string, err error) {
	data, err = "", nil

	req := struct {
		BlockID string `json:"block_id,omitempty"`
	}{
		id.String(),
	}
	res := struct {
		BlockData string `json:"block_data,omitempty"`
	}{}

	err = RawCall(c, "get-block", &req, &res)
	if err != nil {
		return
	}
	data = res.BlockData
	return
}

// Nodes gets blockchain nodes list
func Nodes(c *conn.Conn) (nodes []NodeResult, err error) {

	nodes, err = make([]NodeResult, 0), nil

	req := struct{}{}
	type Item struct {
		Index   string `json:"index,omitempty"`
		Wallet  string `json:"wallet_public_key,omitempty"`
		Address string `json:"address,omitempty"`
	}
	res := struct {
		NodeList []Item `json:"node_list,omitempty"`
	}{}

	err = RawCall(c, "get-nodes", &req, &res)
	if err != nil {
		return
	}

	// list exists
	if res.NodeList == nil || len(res.NodeList) == 0 {
		err = fmt.Errorf("Node list is empty")
		return
	}

	for _, v := range res.NodeList {
		nodes = append(nodes, NodeResult{
			Index:   v.Index,
			Address: v.Wallet,
			IP:      v.Address,
		})
	}
	return
}

// WalletTransactions returns wallet incoming/outgoing transaction list
func WalletTransactions(c *conn.Conn, count uint16, address string) (list []WalletTransactionsResult, err error) {

	list, err = make([]WalletTransactionsResult, 0), nil

	req := struct {
		PublicKey string `json:"public_key,omitempty"`
		Count     string `json:"count,omitempty"`
	}{
		address, fmt.Sprintf("%v", count),
	}
	type Item struct {
		Digest string `json:"digest,omitempty"`
		Type   string `json:"type,omitempty"`
	}
	res := struct {
		TxList []Item `json:"transaction_list,omitempty"`
	}{}

	err = RawCall(c, "get-wallet-transactions", &req, &res)
	if err != nil {
		return
	}

	for _, v := range res.TxList {
		list = append(list, WalletTransactionsResult{
			Digest: v.Digest,
			Status: v.Type,
		})
	}
	return
}
