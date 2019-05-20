package rpc

import (
	"math/big"

	sumuslib "github.com/void616/gm-sumuslib"
	"github.com/void616/gm-sumuslib/amount"
)

// BlockchainStateResult model
type BlockchainStateResult struct {
	BlockCount              *big.Int
	LastBlockDigest         string
	LastBlockMerkleRoot     string
	TransactionCount        *big.Int
	NodeCount               *big.Int
	NonEmptyWalletCount     *big.Int
	VotingTransactionCount  *big.Int
	PendingTransactionCount *big.Int
	BlockchainState         string
	ConsensusRound          string
	VotingNodes             string
}

// TransactionResult model
type TransactionResult struct {
	Name   string
	Hash   string
	Nonce  uint64
	From   string
	To     string
	Amount *amount.Amount
	Token  sumuslib.Token
	Digest string
	Status string
}

// WalletStateResult model
type WalletStateResult struct {
	Balance       WalletBalanceResult
	Exists        bool
	ApprovedNonce uint64
	Tags          []string
}

// WalletBalanceResult model
type WalletBalanceResult struct {
	Gold *amount.Amount
	Mnt  *amount.Amount
}

// NodeResult model
type NodeResult struct {
	Index   string
	Address string
	IP      string
}

// WalletTransactionsResult model
type WalletTransactionsResult struct {
	Digest string
	Status string
}
