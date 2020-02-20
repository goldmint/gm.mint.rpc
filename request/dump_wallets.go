package request

import (
	"context"
	"io"

	"github.com/pierrec/lz4"
	mint "github.com/void616/gm.mint"
	"github.com/void616/gm.mint.rpc/conn"
	"github.com/void616/gm.mint.rpc/rpc"
	"github.com/void616/gm.mint/serializer"
)

// WalletsDump model
type WalletsDump struct {
	Count        uint64    `json:"count"`
	Uncompressed uint64    `json:"uncompressed"`
	Dump         ByteArray `json:"dump"`
}

// DumpWallets method
func DumpWallets(ctx context.Context, c *conn.Conn, withMNT, withGOLD bool) (res WalletsDump, rerr *rpc.Error, err error) {
	res, rerr, err = WalletsDump{}, nil, nil

	rctx, rcancel := c.Receive(ctx)
	defer rcancel()

	params := struct {
		WithMNT  bool `json:"with_mnt,omitempty"`
		WithGOLD bool `json:"with_gold,omitempty"`
	}{
		withMNT,
		withGOLD,
	}

	msg, err := c.Request(rctx, "dump_wallets", params)
	if err != nil {
		return
	}

	switch m := msg.(type) {
	case *rpc.Error:
		rerr = m
		return
	case *rpc.Result:
		err = m.Parse(&res)
		return
	}
	return
}

// Wallets parses dump
func (wd WalletsDump) Wallets() (map[mint.PublicKey]WalletState, error) {
	un := make([]byte, wd.Uncompressed)
	ulen, err := lz4.UncompressBlock(wd.Dump[:], un)
	if err != nil {
		return nil, err
	}
	des := serializer.NewDeserializer(un[:ulen])

	states := make(map[mint.PublicKey]WalletState)
	for {

		pub := des.GetPublicKey()
		mnt := des.GetAmount()
		gold := des.GetAmount()
		tagsCount := des.GetUint16()
		tags := make([]string, tagsCount)
		for i := uint16(0); i < tagsCount; i++ {
			tags[i] = mint.WalletTag(des.GetByte()).String()
		}

		if err := des.Error(); err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		states[pub] = WalletState{
			Exist: true,
			Balance: Balance{
				Mnt:  mnt,
				Gold: gold,
			},
			Tags: tags,
		}
	}

	return states, nil
}
