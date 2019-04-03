package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/void616/gm-sumusrpc/conn"
)

// rpcRequest is a request model
type rpcRequest struct {
	// ID is command
	ID string `json:"id,omitempty"`
	// Params are request params
	Params interface{} `json:"params,omitempty"`
}

// rpcResponse is a response model
type rpcResponse struct {
	// ID is command
	ID string `json:"id,omitempty"`
	// Result is error flag (1/0)
	Result string `json:"result,omitempty"`
	// Text is error message
	Text string `json:"text,omitempty"`
	// Params are response params
	Params interface{} `json:"params,omitempty"`
}

// RawCall sends request `req` to the node via connection `c`, waits for response and deserializes it into `res`.
// In case of transport/parsing/node-side/format problems `err` will be non-nil.
func RawCall(c *conn.Conn, command string, req interface{}, res interface{}) error {

	reqModel := rpcRequest{ID: command, Params: req}
	reqBytes, err := json.Marshal(&reqModel)
	if err != nil {
		return fmt.Errorf("Failed to marshal: %v", err.Error())
	}

	// send request and receive response on exact command/ID
	resBytes, err := c.SendThenReceiveMessage(reqBytes, command)
	if err != nil {
		return fmt.Errorf("Failed to make RPC call: %v", err.Error())
	}

	//log.Print("RPC\n", ">>\n", hex.Dump(rpcReqBytes), "<<\n", hex.Dump(rpcResBytes))

	// parse node reponse
	resModel := rpcResponse{Params: res}
	if err := json.Unmarshal(resBytes, &resModel); err != nil {
		return fmt.Errorf("Failed to unmarshal: %v", err.Error())
	}

	// check node result
	if resModel.Result != "0" {
		return fmt.Errorf("Result: %v; Text: %v", resModel.Result, resModel.Text)
	}

	return nil
}
