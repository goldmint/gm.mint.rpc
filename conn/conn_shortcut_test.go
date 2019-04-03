package conn

import (
	"testing"
)

const node = "127.0.0.1:4010"

func TestRPC_Heartbeat(t *testing.T) {
	r, err := New(node, Options{
		ConnTimeout: 5,
		Logger: func(s string) {
			t.Log(s)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	err = r.Heartbeat()
	if err != nil {
		t.Fatal(err)
	}
}
