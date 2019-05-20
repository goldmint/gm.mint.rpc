package rpc

import (
	"encoding/json"
	"testing"

	"github.com/void616/gm-sumusrpc/conn"
)

var node = ":4010"

func TestRawCall(t *testing.T) {

	c, err := conn.New(node, conn.Options{
		Logger: func(s string) {
			t.Log(s)
		},
	})
	if err != nil {
		t.Fatal("Failed to connect")
	}
	defer c.Close()

	req := struct {
		Pub     string `json:"public_key,omitempty"`
		Count   string `json:"count,omitempty"`
		RawData string `json:"raw_data,omitempty"`
	}{
		"6z2L3uqqcUtSKA1AXFaWmW4A5Rs8fBuB5F7zeb7MhSFUV6Zv6",
		"100", "yes",
	}
	var res *json.RawMessage

	if err := RawCall(c, "get-wallet-transactions", &req, &res); err != nil {
		t.Fatal(err)
	}
	t.Log(string([]byte(*res)))
}
