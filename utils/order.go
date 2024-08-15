package utils

import (
	"math"
)

type SlippageType int

const (
	SlippageTypeDown SlippageType = 0
	SlippageTypeUp   SlippageType = 1
)

func ApplySlippage(slippage float64, amount uint64, slippageType SlippageType) uint64 {
	if slippageType == SlippageTypeUp {
		amountF := float64(amount)
		slippageAdjustedAmount := amountF * (1 + slippage)
		return uint64(slippageAdjustedAmount)

	} else {
		amountF := float64(amount)
		slippageAdjustedAmount := 1 / (1 + slippage)
		slippageAdjustedAmount = slippageAdjustedAmount * amountF
		return uint64(slippageAdjustedAmount)
	}
}

// TODO: Haven't tested this yet
func CalculateDepositAmount(amountA uint64, amountB uint64, pool V2PoolState) uint64 {
	ratioA := (amountA * pool.TotalLiquidity) / pool.ReserveA
	ratioB := (amountB * pool.TotalLiquidity) / pool.ReserveB

	if ratioA > ratioB {
		swapAmountA := CalculateDepositSwapAmount(amountA, amountB, pool.ReserveA, pool.ReserveB, pool.BaseFeeANumerator)
		lpAmount := ((amountA * swapAmountA.Denominator - swapAmountA.Numerator) * pool.TotalLiquidity) /
        (pool.ReserveA * swapAmountA.Denominator + swapAmountA.Numerator);
		return lpAmount

	} else if ratioA < ratioB {
		swapAmountB := CalculateDepositSwapAmount(amountB, amountA, pool.ReserveB, pool.ReserveA, pool.BaseFeeBNumerator)
		lpAmount := ((amountB * swapAmountB.Denominator - swapAmountB.Numerator) * pool.TotalLiquidity) /
			(pool.ReserveB * swapAmountB.Denominator + swapAmountB.Numerator);
		return lpAmount
	}

	return ratioA
}

type Fraction struct {
	Numerator   uint64
	Denominator uint64
}

// TODO: Haven't tested this yet
func CalculateDepositSwapAmount(amountIn, amountOut, reserveIn, reserveOut, tradingFeeNumerator uint64) Fraction {
	x := (amountOut + reserveOut) * reserveIn
	y := 4 * (amountOut + reserveOut) *
		(amountOut * reserveIn * reserveIn - amountIn * reserveIn * reserveOut)
	z := 2 * (amountOut + reserveOut)
	a := uint64Pow(x) * uint64Pow(uint64(2) *
		DEFAULT_TRADING_FEE_DENOMINATOR - tradingFeeNumerator) - y *
		DEFAULT_TRADING_FEE_DENOMINATOR * (DEFAULT_TRADING_FEE_DENOMINATOR - tradingFeeNumerator)

	b := (uint64(2) * DEFAULT_TRADING_FEE_DENOMINATOR - tradingFeeNumerator) * x
	numerator := uint64(math.Sqrt(float64(a)) - float64(b))
	denominator := z * (DEFAULT_TRADING_FEE_DENOMINATOR - tradingFeeNumerator)

	return Fraction{
		Numerator:   numerator,
		Denominator: denominator,
	}
}

func int64Pow(i int64) int64 {
	return i * i
}

func uint64Pow(i uint64) uint64 {
	return i * i
}
