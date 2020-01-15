package main

import (
	"log"
	"os"
	"time"

	"github.com/void616/gm.mint.rpc/conn"
)

func main() {
	c, err := conn.New(os.Getenv("MINTRPC"), conn.DefaultOptions.WithStdLogger())
	if err != nil {
		log.Fatalln("failed to connect:", err)
	}
	defer c.Close()
	go c.Serve()

	if err := c.Heartbeat(time.Second * 5); err != nil {
		log.Println(err)
	} else {
		log.Println("OK")
	}
}
