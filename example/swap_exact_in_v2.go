package example

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	c "github.com/Newt6611/apollo/constants"
	"github.com/Newt6611/go-minswap/adapter"
	v2 "github.com/Newt6611/go-minswap/dex/v2"
	"github.com/Newt6611/go-minswap/utils"
	"github.com/blockfrost/blockfrost-go"
)

func SwapExactInExample() {
	ctx := context.Background()

	// Create Adapter
	blockfrostAdapter, err := adapter.NewBlockFrost(blockfrost.APIClientOptions{
		ProjectID: YOUR_BLOCKFROST_API_KEY_HERE,
		Server: blockfrost.CardanoPreProd,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Get A New Builder
	builder := blockfrostAdapter.NewBuilder()
	builder, _ = builder.SetWalletFromMnemonic(YOUR_TEST_SEED_HERE, c.PREPROD)

	dexv2 := v2.NewDexV2(blockfrostAdapter)

	assetA := utils.ADA
	assetB := utils.MIN
	pool, err := blockfrostAdapter.GetV2PoolByPair(ctx, assetA, assetB)
	if err != nil {
		log.Fatal(err)
	}

	swapAmount := uint64(5_000000);
	amountOut := v2.CalculateAmountOut(pool.ReserveA, pool.ReserveB, swapAmount, pool.BaseFeeANumerator)

	// 20%
	slippageTolerance := 20.0 / 100.0
	acceptedAmountOut := utils.ApplySlippage(slippageTolerance, amountOut, utils.SlippageTypeDown)

	swapExactIn := v2.SwapExactIn {
		Direction: v2.Direction_A_To_B,
		SwapAmount: v2.SwapAmount{
			Type:   v2.AmountType_Specific_Amount,
			Amount: swapAmount,
		},
		MinimumReceived: acceptedAmountOut,
		Killable: v2.Killable_Pending_On_Failed,
	}

	builder, err = dexv2.BuildSwapExactInOrder(ctx, builder, swapExactIn, utils.ADA, utils.MIN, 10_000000)
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
