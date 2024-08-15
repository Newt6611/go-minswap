package example

import (
	"context"
	"fmt"
	"log"

	"github.com/Newt6611/apollo/serialization/AssetName"
	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/Policy"
	"github.com/Newt6611/go-minswap/adapter"
	"github.com/blockfrost/blockfrost-go"
)

func GetStablePoolPrice() {
	ctx := context.Background()

	blockfrostAdapter, err := adapter.NewBlockFrost(blockfrost.APIClientOptions{
		ProjectID: YOUR_BLOCKFROST_API_KEY_HERE,
		Server:    blockfrost.CardanoPreProd,
	})
	if err != nil {
		log.Fatal(err)
	}

	nftHex := "06fe1ba957728130154154d5e5b25a7b533ebe6c4516356c0aa69355646a65642d697573642d76312e342d6c70"
	nft := Fingerprint.Fingerprint{
		PolicyId:  Policy.PolicyId{Value: nftHex[:56]},
		AssetName: AssetName.AssetName{Value: nftHex[56:]},
	}

	pool, err := blockfrostAdapter.GetStableByNFT(ctx, nft)
	if err != nil {
		log.Fatal(err)
	}
	b1 := float64(pool.Balances[0])
	b2 := float64(pool.Balances[1])
	fmt.Println(b1 / b2)
}
