package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/void616/gm-sumusrpc/conn"
	"github.com/void616/gm-sumusrpc/rpc"
)

func main() {
	c, err := conn.New(os.Getenv("MINTRPC"), conn.DefaultOptions.WithStdLogger())
	if err != nil {
		log.Fatalln("failed to connect:", err)
	}
	defer c.Close()
	go c.Serve()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rctx, rcancel := c.Receive(ctx)
	defer rcancel()

	msg, err := c.Request(rctx, "get_blockchain_state", nil)
	if err != nil {
		log.Println(err)
		return
	}

	switch m := msg.(type) {
	case *rpc.Result:
		log.Println("OK")
	case *rpc.Error:
		log.Println(m.Err())
	}
}
