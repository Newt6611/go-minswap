package utils

import (
	"math/big"
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

func CalculateDepositAmount(amountA uint64, amountB uint64, pool V2PoolState) uint64 {
	amountA_b := big.NewInt(int64(amountA))
	amountB_b := big.NewInt(int64(amountB))
	reserveA_b := big.NewInt(int64(pool.ReserveA))
	reserveB_b := big.NewInt(int64(pool.ReserveB))
	totalLiquidity_b := big.NewInt(int64(pool.TotalLiquidity))

	// ratioA = amountA * totalLiquidity / reserveA
	ratioA := big.NewInt(0).Mul(amountA_b, totalLiquidity_b)
	ratioA = big.NewInt(0).Div(ratioA, reserveA_b)

	ratioB := big.NewInt(0).Mul(amountB_b, totalLiquidity_b)
	ratioB = big.NewInt(0).Div(ratioB, reserveB_b)

	if ratioA.Cmp(ratioB) == 1 {
		swapAmountA := CalculateDepositSwapAmount(amountA, amountB, pool.ReserveA, pool.ReserveB, pool.BaseFeeANumerator)

		// let lp_amount =
		// ( amount_a * swap_amount_a_denominator - swap_amount_a_numerator ) * total_liquidity / (
		//   reserve_a * swap_amount_a_denominator + swap_amount_a_numerator
		// )
		var lpAmount *big.Int
		{
			t1 := new(big.Int).Mul(amountA_b, swapAmountA.Denominator)
			t1.Sub(t1, swapAmountA.Numerator)
			t1.Mul(t1, totalLiquidity_b)

			t2 := new(big.Int).Mul(reserveA_b, swapAmountA.Denominator)
			t2.Add(t2, swapAmountA.Numerator)

			lpAmount = new(big.Int).Div(t1, t2)
		}
		return lpAmount.Uint64()

	} else if ratioA.Cmp(ratioB) == -1 {
		swapAmountB := CalculateDepositSwapAmount(amountB, amountA, pool.ReserveB, pool.ReserveA, pool.BaseFeeBNumerator)

		// let lp_amount =
		// ( amount_a * swap_amount_a_denominator - swap_amount_a_numerator ) * total_liquidity / (
		//   reserve_a * swap_amount_a_denominator + swap_amount_a_numerator
		// )
		var lpAmount *big.Int
		{
			t1 := new(big.Int).Mul(amountB_b, swapAmountB.Denominator)
			t1.Sub(t1, swapAmountB.Numerator)
			t1.Mul(t1, totalLiquidity_b)

			t2 := new(big.Int).Mul(reserveB_b, swapAmountB.Denominator)
			t2.Add(t2, swapAmountB.Numerator)

			lpAmount = new(big.Int).Div(t1, t2)
		}
		return lpAmount.Uint64()
	}
	return ratioA.Uint64()
}

type Fraction struct {
	Numerator   *big.Int
	Denominator *big.Int
}

func CalculateDepositSwapAmount(amountIn, amountOut, reserveIn, reserveOut, tradingFeeNumerator uint64) Fraction {
	amountIn_b := big.NewInt(int64(amountIn))
	amountOut_b := big.NewInt(int64(amountOut))
	reserveIn_b := big.NewInt(int64(reserveIn))
	reserveOut_b := big.NewInt(int64(reserveOut))
	tradingFeeNumerator_b := big.NewInt(int64(tradingFeeNumerator))
	DEFAULT_TRADING_FEE_DENOMINATOR_B := big.NewInt(int64(DEFAULT_TRADING_FEE_DENOMINATOR))

	// x := (amountOut + reserveOut) * reserveIn
	var x *big.Int
	{
		t1 := new(big.Int).Add(amountOut_b, reserveOut_b)
		x = new(big.Int).Mul(t1, reserveIn_b)
	}

	// let y = 4 * ( amount_out + reserve_out ) *
	// (amount_out * Pow(reserve_in) - amount_in * reserve_in * reserve_out)
	var y *big.Int
	{
		t1 := big.NewInt(4)
		t2 := new(big.Int).Add(amountOut_b, reserveOut_b)

		reserveInPow := new(big.Int).Exp(reserveIn_b, big.NewInt(2), nil)
		t3 := new(big.Int).Mul(amountOut_b, reserveInPow)

		t4 := new(big.Int).Mul(amountIn_b, reserveIn_b)
		t4.Mul(t4, reserveOut_b)

		t3.Sub(t3, t4)

		y = t1.Mul(t1, t2)
		y.Mul(y, t3)
	}

	// z := 2 * (amountOut + reserveOut)
	var z *big.Int
	{
		z = big.NewInt(2)
		t1 := new(big.Int).Add(amountOut_b, reserveOut_b)
		z.Mul(z, t1)
	}

	// b := (2 * DEFAULT_TRADING_FEE_DENOMINATOR - tradingFeeNumerator) * x
	var b *big.Int
	{
		b = big.NewInt(2)
		b.Mul(b, DEFAULT_TRADING_FEE_DENOMINATOR_B)
		b.Sub(b, tradingFeeNumerator_b)
		b.Mul(b, x)
	}

	// let a =
	//     calculate_pow(b) -
	//     y * utils.default_fee_denominator *
	//     (utils.default_fee_denominator - trading_fee_numerator)
	var a *big.Int
	{
		t1 := new(big.Int).Exp(b, big.NewInt(2), nil)
		t2 := new(big.Int).Mul(y, DEFAULT_TRADING_FEE_DENOMINATOR_B)
		t3 := new(big.Int).Sub(DEFAULT_TRADING_FEE_DENOMINATOR_B, tradingFeeNumerator_b)

		t := t2.Mul(t2, t3)
		a = t1.Sub(t1, t)
	}


	// numerator := math.Sqrt(a) - b
	var numerator *big.Int
	{
		t1 := new(big.Int).Sqrt(a)
		t1.Sub(t1, b)
		numerator = t1
	}

	// denominator := z * (DEFAULT_TRADING_FEE_DENOMINATOR - tradingFeeNumerator)
	var denominator *big.Int
	{
		denominator = new(big.Int).Sub(DEFAULT_TRADING_FEE_DENOMINATOR_B, tradingFeeNumerator_b)
		denominator.Mul(denominator, z)
	}

	return Fraction{
		Numerator:   numerator,
		Denominator: denominator,
	}
}
