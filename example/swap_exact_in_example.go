package example

import (
	"context"
	"math/big"

	"github.com/Newt6611/apollo"
	"github.com/Newt6611/apollo/serialization/Address"
	"github.com/Newt6611/go-minswap/adapter"
	v2 "github.com/Newt6611/go-minswap/dex/v2"
	"github.com/Newt6611/go-minswap/utils"
)

// TODO:
func SwapExactInExample(ctx context.Context, apollo *apollo.Apollo, adapter adapter.Adapter, address Address.Address) {
	assetA := utils.ADA
	assetB := utils.MIN
	pool, err := adapter.GetV2PoolByPair(ctx, assetA, assetB)
	if err != nil {
		panic(err)
	}

	swapAmount := big.NewInt(5_000_000);
	amountOut := v2.CalculateAmountOut(pool.ReserveA, pool.ReserveB, swapAmount, pool.BaseFeeANumerator)
	_ = amountOut

	// 20%
	// slippageTolerance = new BigNumber(20).div(100);
	// const acceptedAmountOut = Slippage.apply({
	// 	slippage: slippageTolerance,
	// 	amount: amountOut,
	// 	type: "down",
	// 	});
}

