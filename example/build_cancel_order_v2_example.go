package example

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	c "github.com/Newt6611/apollo/constants"
	"github.com/Newt6611/go-minswap/adapter"
	"github.com/Newt6611/go-minswap/constants"
	v2 "github.com/Newt6611/go-minswap/dex/v2"
	"github.com/blockfrost/blockfrost-go"
)

func BuildCancelOrderV2Example() {
	ctx := context.Background()
	blockfrostAdapter, err := adapter.NewBlockFrost(blockfrost.APIClientOptions{
		ProjectID: YOUR_BLOCKFROST_API_KEY_HERE,
		Server: blockfrost.CardanoPreProd,
	})
	if err != nil {
		log.Fatal(err)
	}

	builder := blockfrostAdapter.NewBuilder()
	builder, _ = builder.SetWalletFromMnemonic(YOUR_TEST_SEED_HERE, c.PREPROD)

	dexv2 := v2.NewDexV2(blockfrostAdapter)
	builder, err = dexv2.BuildCancelOrder(ctx, builder, []constants.OutRef {
		{ "f1d7aaf16f1a8385580d7a3045ddc9b1010eab5949d9ccbd23a2760c11ab3b23", 0},
		{ "1d936e663e9e8015afddf53df18d93a8b265c01373369e1c4e09c48cc62a374b", 0},
	})

	if err != nil {
		log.Fatal(err)
	}

	builder = builder.Sign()
	id, err := builder.Submit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hex.EncodeToString(id.Payload))
}
