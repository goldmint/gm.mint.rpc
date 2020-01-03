package main

import (
	"context"
	"log"
	"os"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var pub sumuslib.PublicKey

	res, rerr, err := request.GetWalletTransactionsBinary(ctx, c, pub, 10, true, true, true)
	if err != nil {
		log.Println("error:", err)
		return
	}
	if rerr != nil {
		log.Println("rpc error:", rerr.Err())
		return
	}

	log.Printf("Result: %v transactions", len(res))

	res2, rerr, err := request.GetWalletTransactionsTextual(ctx, c, pub, 10, true, true, true)
	if err != nil {
		log.Println("error:", err)
		return
	}
	if rerr != nil {
		log.Println("rpc error:", rerr.Err())
		return
	}

	log.Printf("Result: %+v", res2)
}
