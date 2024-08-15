package example

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/Newt6611/go-minswap/adapter"
	"github.com/blockfrost/blockfrost-go"
)

func GetAllStablePools() {
	ctx := context.Background()
	blockfrostAdapter, err := adapter.NewBlockFrost(blockfrost.APIClientOptions{
		ProjectID: YOUR_BLOCKFROST_API_KEY_HERE,
		Server:    blockfrost.CardanoPreProd,
	})
	if err != nil {
		log.Fatal(err)
	}

	pools, errs := blockfrostAdapter.GetAllStablePools(ctx)
	if len(errs) > 0 {
		log.Fatal(errs)
	}

	fmt.Println("------------")
	for _, pool := range pools {
		fmt.Println("balances ", pool.Balances)
		fmt.Println("total liquidity ", pool.TotalLiquidity)
		fmt.Println("AMP ", pool.AMP)
		fmt.Println("order hash ", hex.EncodeToString(pool.OrderHash))
		fmt.Println("------------")
	}
}
