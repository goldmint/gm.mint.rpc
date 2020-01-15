package main

import (
	"bytes"
	"context"
	"log"
	"math/big"
	"os"
	"time"

	mint "github.com/void616/gm.mint"
	"github.com/void616/gm.mint/block"
	"github.com/void616/gm.mint/serializer"
	"github.com/void616/gm.mint/transaction"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// get by id
	blockBytes, rerr, err := request.GetBlockByID(ctx, c, big.NewInt(0))
	if err != nil {
		log.Println("error:", err)
		return
	}
	if rerr != nil {
		log.Println("rpc error:", rerr.Err())
		return
	}
	log.Printf("Block length: %v", len(blockBytes))

	// ---

	var digest mint.Digest

	headerCbk := func(h *block.Header) error {
		digest = h.Digest
		log.Println("ID:", h.BlockID.String())
		log.Println("Version:", h.Version)
		log.Println("Prev block:", h.PrevBlockDigest.String())
		log.Println("Merkle root:", h.MerkleRoot.String())
		log.Println("Signers:", h.SignersCount)
		log.Println("Consensus round:", h.ConsensusRound)
		log.Println("Transactions:", h.TransactionsCount)
		log.Println("Time:", mint.StampToTime(h.Timestamp).String())
		return nil
	}

	txCbk := func(t transaction.Code, d *serializer.Deserializer, h *block.Header) error {
		txer, err := transaction.CodeToTransaction(t)
		if err != nil {
			return err
		}
		_, err = txer.Parse(d.Source())
		if err != nil {
			return err
		}
		log.Println("Transaction: ", txer.Code().String())
		return nil
	}

	// ---

	// parse block
	if err := block.Parse(bytes.NewBuffer(blockBytes), headerCbk, txCbk); err != nil {
		log.Println("block parser error:", err)
		return
	}

	// ---

	log.Println("")
	log.Println("Block by digest: ", digest.String())

	// same block by digest
	blockBytes, rerr, err = request.GetBlockByDigest(ctx, c, digest)
	if err != nil {
		log.Println("error:", err)
		return
	}
	if rerr != nil {
		code, desc, ok := rerr.GetReason()
		if ok {
			log.Println("rpc error reason code:", int(code), "/", code.String())
			log.Println("rpc error reason desc:", desc)
			return
		}
		log.Println("rpc error:", rerr.Err())
		return
	}

	// parse block
	if err := block.Parse(bytes.NewBuffer(blockBytes), headerCbk, txCbk); err != nil {
		log.Println("block parser error:", err)
		return
	}
}
