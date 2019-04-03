package rpc

import (
	"encoding/json"
	"testing"

	"github.com/void616/gm-sumusrpc/conn"
)

var node = ":4010"

func TestRawCall(t *testing.T) {

	c, _ := conn.New(node, conn.Options{
		Logger: func(s string) {
			t.Log(s)
		},
	})
	defer c.Close()

	req := struct{}{}
	var res *json.RawMessage

	if err := RawCall(c, "get-blockchain-state", &req, &res); err != nil {
		t.Fatal(err)
	}
	t.Log(string([]byte(*res)))
}
