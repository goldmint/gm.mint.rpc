package rpc

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	sumuslib "github.com/void616/gm-sumuslib"
	"github.com/void616/gm-sumuslib/signer"
	"github.com/void616/gm-sumuslib/transaction"
	"github.com/void616/gm-sumusrpc/conn"
)

var testNode = "127.0.0.1:4010"
var testPK = "PVT"
var testRecipient = "PUB"

func connect(t *testing.T) *conn.Conn {
	var node = testNode
	c, err := conn.New(node, conn.Options{
		Logger: func(s string) {
			t.Log(s)
		},
	})
	if err != nil {
		t.Fatal("Failed to connect:", err)
	}
	return c
}

func jsony(t *testing.T, v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal("Failed to marshal:", err)
	}
	return string(b)
}

func TestBlockchainState(t *testing.T) {
	v, code, err := BlockchainState(connect(t))
	if err != nil {
		t.Fatal("Fail:", err)
	}
	t.Logf("Code: %d %v, Result:\n%v", code, code, jsony(t, v))
}

func TestNodes(t *testing.T) {
	v, code, err := Nodes(connect(t))
	if err != nil {
		t.Fatal("Fail:", err)
	}
	t.Logf("Code: %d %v, Result:\n%v", code, code, jsony(t, v))
}

func TestAddTransaction(t *testing.T) {

	pk := testPK
	nonce := uint64(0)

	pkbytes, _ := sumuslib.Unpack58(pk)
	sig, err := signer.FromBytes(pkbytes)
	if err != nil {
		t.Fatal("Fail:", err)
	}

	// txtype := sumuslib.TransactionTransferAssets
	// dst, _ := sumuslib.UnpackAddress58(testRecipient)
	// tx := transaction.TransferAsset{
	// 	Address: dst,
	// 	Token:   sumuslib.TokenGOLD,
	// 	Amount:  amount.NewInteger(1),
	// }

	// txtype := sumuslib.TransactionRegisterNode
	// tx := transaction.RegisterNode{
	// 	NodeAddress: "127.0.0.1",
	// }

	txtype := sumuslib.TransactionUnregisterNode
	tx := transaction.UnregisterNode{}

	// ---

	_ = tx

	stx, err := tx.Construct(sig, nonce)
	if err != nil {
		t.Fatal("Fail:", err)
	}
	v, code, err := AddTransaction(connect(t), txtype, hex.EncodeToString(stx.Data))
	if err != nil {
		t.Fatal("Fail:", err)
	}
	atec := AddTransactionErrorCode(code)
	t.Log(
		"\nCode:", uint16(code), code.String(),
		"\nNonce behind:", atec.NonceBehind(),
		"\nNonce ahead:", atec.NonceAhead(),
		"\nAdded already:", atec.AddedAlready(),
		"\nWallet problem:", atec.WalletInconsistency(),
		"\nResult:\n", jsony(t, v),
	)
}

func TestWalletTransactions(t *testing.T) {
	v, code, err := WalletTransactions(connect(t), 1000, testRecipient)
	if err != nil {
		t.Fatal("Fail:", err)
	}
	t.Logf("Code: %d %v, Result:\n%v", code, code, jsony(t, v))
}
