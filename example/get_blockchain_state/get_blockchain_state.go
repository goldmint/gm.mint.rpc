package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/void616/gm.mint.rpc/conn"
	"github.com/void616/gm.mint.rpc/request"
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

	res, rerr, err := request.GetBlockchainState(ctx, c)
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
