package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/void616/gm-sumuslib/amount"
	"github.com/void616/gm-sumuslib/signer"
	"github.com/void616/gm-sumuslib/transaction"

	sumuslib "github.com/void616/gm-sumuslib"
	"github.com/void616/gm-sumusrpc/conn"
	"github.com/void616/gm-sumusrpc/request"
)

func main() {
	c, err := conn.New(os.Getenv("MINTRPC"), conn.DefaultOptions)
	if err != nil {
		log.Fatalln("failed to connect:", err)
	}
	defer c.Close()
	go c.Serve()

	// sender
	pk, err := sumuslib.ParsePrivateKey(os.Getenv("MINTSENDERPK"))
	if err != nil {
		log.Fatalln("failed to parse private key:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// receiver
	to, _ := sumuslib.ParsePublicKey("2gqgdwYHECniC5KNZq7Gi2876qJNUVDSKSLSGSuczdJDdNrWi8")

	// tx
	tx := transaction.TransferAsset{
		Address: to,
		Amount:  amount.MustFromString("0.1337"),
		Token:   sumuslib.TokenGOLD,
	}

	// tx nonce
	nonce := uint64(1)

	// signed tx
	stx, err := tx.Sign(signer.FromPrivateKey(pk), nonce)
	if err != nil {
		log.Fatalln("failed to sign transaction:", err)
	}

	// add
	res, rerr, err := request.AddTransaction(ctx, c, tx.Code(), stx.Data)
	if err != nil {
		log.Println("error:", err)
		return
	}
	if rerr != nil {
		log.Println("rpc error:", rerr.Err())
		return
	}

	log.Printf("Result: %+v", res)
}
