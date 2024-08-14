package utils

import (
	"math"
	"math/big"
)

type SlippageType int

const (
	SlippageTypeDown SlippageType = 0
	SlippageTypeUp   SlippageType = 1
)

func ApplySlippage(slippage float64, amount *big.Int, slippageType SlippageType) *big.Int {
	if slippageType == SlippageTypeUp {
		amountF, _ := amount.Float64()
		slippageAdjustedAmount := amountF * (1 + slippage)
		return big.NewInt(int64(slippageAdjustedAmount))

	} else {
		amountF, _ := amount.Float64()
		slippageAdjustedAmount := 1 / (1 + slippage)
		slippageAdjustedAmount = slippageAdjustedAmount * amountF
		return big.NewInt(int64(slippageAdjustedAmount))
	}
}

func CalculateDepositAmount(amountA *big.Int, amountB *big.Int, pool V2PoolState) *big.Int {
	ratioA := (amountA.Int64() * pool.TotalLiquidity.Int64()) / pool.ReserveA.Int64()
	ratioB := (amountB.Int64() * pool.TotalLiquidity.Int64()) / pool.ReserveB.Int64()

	if ratioA > ratioB {
		swapAmountA := CalculateDepositSwapAmount(amountA, amountB, pool.ReserveA, pool.ReserveB, pool.BaseFeeANumerator)
		lpAmount := ((amountA.Int64() * swapAmountA.Denominator - swapAmountA.Numerator) * pool.TotalLiquidity.Int64()) /
        (pool.ReserveA.Int64() * swapAmountA.Denominator + swapAmountA.Numerator);
		return big.NewInt(lpAmount)

	} else if ratioA < ratioB {
		swapAmountB := CalculateDepositSwapAmount(amountB, amountA, pool.ReserveB, pool.ReserveA, pool.BaseFeeBNumerator)
		lpAmount := ((amountB.Int64() * swapAmountB.Denominator - swapAmountB.Numerator) * pool.TotalLiquidity.Int64()) /
			(pool.ReserveB.Int64() * swapAmountB.Denominator + swapAmountB.Numerator);
		return big.NewInt(lpAmount)
	}

	return big.NewInt(ratioA)
}

type Fraction struct {
	Numerator   int64
	Denominator int64
}

func CalculateDepositSwapAmount(amountIn, amountOut, reserveIn, reserveOut, tradingFeeNumerator *big.Int) Fraction {
	amountIni := amountIn.Int64()
	amountOuti := amountOut.Int64()
	reserveIni := reserveIn.Int64()
	reserveOuti := reserveOut.Int64()
	tradingFeeNumeratori := tradingFeeNumerator.Int64()

	x := (amountOuti + reserveOuti) * reserveIni
	y := 4 * (amountOuti + reserveOuti) *
		(amountOuti*reserveIni*reserveIni - amountIni*reserveIni*reserveOuti)
	z := 2 * (amountOuti + reserveOuti)
	a := int64Pow(x)*int64Pow(int64(2)*
		DEFAULT_TRADING_FEE_DENOMINATOR-tradingFeeNumeratori) - y*
		DEFAULT_TRADING_FEE_DENOMINATOR*(DEFAULT_TRADING_FEE_DENOMINATOR-tradingFeeNumeratori)

	b := (int64(2)*DEFAULT_TRADING_FEE_DENOMINATOR - tradingFeeNumeratori) * x
	numerator := int64(math.Sqrt(float64(a)) - float64(b))
	denominator := z * (DEFAULT_TRADING_FEE_DENOMINATOR - tradingFeeNumeratori)
	return Fraction{
		Numerator:   numerator,
		Denominator: denominator,
	}
}

func int64Pow(i int64) int64 {
	return i * i
}
