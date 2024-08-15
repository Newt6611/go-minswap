package utils_test

import (
	"math/big"
	"testing"

	"github.com/Newt6611/go-minswap/utils"
)

func TestCalculateDepositAmount(t *testing.T) {
	//swap amount is(29932544594688 / 19979880000) ~ 1498
	pool1 := utils.V2PoolState{
		ReserveA: 1_000_000,
		ReserveB: 1_000_000,
		TotalLiquidity: 1_000_000,
		BaseFeeANumerator: 30,
		BaseFeeBNumerator: 100,
		FeeSharingNumeratorOpt: utils.FeeSharingOpt{
			Enable: true,
			Numerator: 5000,
		},
	}
	result1 := utils.CalculateDepositAmount(5000, 2000, pool1)
	if result1 != 3496 {
		t.Errorf("CalculateDepositAmount() got = %v, want %v", result1, 3496)
	}

	result2 := utils.CalculateDepositAmount(2000, 5000, pool1)
	if result2 != 3491 {
		t.Errorf("CalculateDepositAmount() got = %v, want %v", result2, 3491)
	}
}

func TestCalculateDepositSwapAmount(t *testing.T) {
	result := utils.CalculateDepositSwapAmount(5000, 2000, 1_000_000, 1_000_000, 30)

	anwserNumerator := big.NewInt(29932544594688)
	anwserDenominator := big.NewInt(19979880000)

	if anwserNumerator.Cmp(result.Numerator) != 0 || anwserDenominator.Cmp(result.Denominator) != 0 {
		t.Errorf("CalculateDepositSwapAmount() got = %v, want %v", result, 29932544594688)
	}
}
