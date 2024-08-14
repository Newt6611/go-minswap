package v2_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/Newt6611/apollo/serialization/Address"
	"github.com/Newt6611/apollo/serialization/AssetName"
	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/PlutusData"
	"github.com/Newt6611/apollo/serialization/Policy"
	v2 "github.com/Newt6611/go-minswap/dex/v2"
	"github.com/Salvionied/cbor/v2"
)

var (
	testSenderAddress, _          = Address.DecodeAddress("addr_test1qpssc0r090a9u0pyvdr9y76sm2xzx04n6d4j0y5hukcx6rxz4dtgkhfdynadkea0qezv99wljdl076xkg2krm96nn8jszmh3w7")
	testReceiverDatumHash         = "b8b912cdbcc998f3f0c18e951928ca179de85735c4fc2d82e8d10777"
	testReceiverDatumHashBytes, _ = hex.DecodeString(testReceiverDatumHash)
)

func buildOrderDatum(step v2.StepI) []v2.OrderDatum {
	policy, _ := Policy.New("d6aae2059baee188f74917493cf7637e679cd219bdfbbf4dcbeb1d0b")
	assetName := AssetName.NewAssetNameFromHexString("e08460587b08cca542bd2856b8d5e1d23bf3f63f9916fb81f6d95fda0910bf69")
	fingerprint := Fingerprint.New(*policy, *assetName)
	return []v2.OrderDatum{
		{
			Canceller: v2.AuthorizationMethod{
				Type: v2.AuthorizationMethodType_Signature,
				Hash: testSenderAddress.PaymentPart,
			},
			RefundReceiver: testSenderAddress,
			RefundReceiverDatum: v2.ExtraDatum{
				Type: v2.ExtraDatumType_No_Datum,
			},
			SuccessReceiver: testSenderAddress,
			SuccessReceiverDatum: v2.ExtraDatum{
				Type: v2.ExtraDatumType_No_Datum,
			},
			Step:           step,
			LpAsset:        *fingerprint,
			MaxBatcherFee:  v2.FIXED_BATCHER_FEE,
			ExpiredOptions: v2.ExpirySetting{},
		},
		{
			Canceller: v2.AuthorizationMethod{
				Type: v2.AuthorizationMethodType_Signature,
				Hash: testSenderAddress.PaymentPart,
			},
			RefundReceiver: testSenderAddress,
			RefundReceiverDatum: v2.ExtraDatum{
				Type: v2.ExtraDatumType_Datum_Hash,
				Hash: testReceiverDatumHashBytes,
			},
			SuccessReceiver: testSenderAddress,
			SuccessReceiverDatum: v2.ExtraDatum{
				Type: v2.ExtraDatumType_Datum_Hash,
				Hash: testReceiverDatumHashBytes,
			},
			Step:           step,
			LpAsset:        *fingerprint,
			MaxBatcherFee:  v2.FIXED_BATCHER_FEE,
			ExpiredOptions: v2.ExpirySetting{},
		},
		{
			Canceller: v2.AuthorizationMethod{
				Type: v2.AuthorizationMethodType_Signature,
				Hash: testSenderAddress.PaymentPart,
			},
			RefundReceiver: testSenderAddress,
			RefundReceiverDatum: v2.ExtraDatum{
				Type: v2.ExtraDatumType_Inline_Datum,
				Hash: testReceiverDatumHashBytes,
			},
			SuccessReceiver: testSenderAddress,
			SuccessReceiverDatum: v2.ExtraDatum{
				Type: v2.ExtraDatumType_Inline_Datum,
				Hash: testReceiverDatumHashBytes,
			},
			Step:           step,
			LpAsset:        *fingerprint,
			MaxBatcherFee:  v2.FIXED_BATCHER_FEE,
			ExpiredOptions: v2.ExpirySetting{},
		},
		{
			Canceller: v2.AuthorizationMethod{
				Type: v2.AuthorizationMethodType_Signature,
				Hash: testSenderAddress.PaymentPart,
			},
			RefundReceiver: testSenderAddress,
			RefundReceiverDatum: v2.ExtraDatum{
				Type: v2.ExtraDatumType_No_Datum,
			},
			SuccessReceiver: testSenderAddress,
			SuccessReceiverDatum: v2.ExtraDatum{
				Type: v2.ExtraDatumType_No_Datum,
			},
			Step:          step,
			LpAsset:       *fingerprint,
			MaxBatcherFee: v2.FIXED_BATCHER_FEE,
			ExpiredOptions: v2.ExpirySetting{
				ExpiredTime:        big.NewInt(1721010208050),
				MaxCancellationTip: big.NewInt(300_000),
			},
		},
	}
}

func TestSwapExactInStepConverter(t *testing.T) {
	step1 := v2.SwapExactIn{
		SwapAmount: v2.SwapAmount{
			Type:   v2.AmountType_Specific_Amount,
			Amount: big.NewInt(10000),
		},
		Direction:       v2.Direction_A_To_B,
		MinimumReceived: big.NewInt(1),
		Killable:        v2.Killable_Pending_On_Failed,
	}
	step2 := v2.SwapExactIn{
		SwapAmount: v2.SwapAmount{
			Type:   v2.AmountType_All,
			Amount: big.NewInt(10000),
		},
		Direction:       v2.Direction_B_To_A,
		MinimumReceived: big.NewInt(1),
		Killable:        v2.Killable_Kill_On_Failed,
	}
	datums := []v2.OrderDatum{}
	datums = append(datums, buildOrderDatum(step1)...)
	datums = append(datums, buildOrderDatum(step2)...)
	for _, datum := range datums {
		plutusData := datum.ToPlutusData()
		b, err := plutusData.MarshalCBOR()
		if err != nil {
			t.Errorf("SwapExactIn ToPlutusData error: %s\n", err)
		}

		var convertedPlutusData PlutusData.PlutusData
		err = cbor.Unmarshal(b, &convertedPlutusData)
		if err != nil {
			t.Errorf("SwapExactIn ToPlutusData error: %s\n", err)
		}

		if !convertedPlutusData.Equal(plutusData) {
			t.Error("TestSwapExactInStepConverter error")
		}
	}
}

func TestStopStepConverter(t *testing.T) {
	step1 := v2.Stop{
		SwapAmount: v2.SwapAmount{
			Type:   v2.AmountType_Specific_Amount,
			Amount: big.NewInt(10000),
		},
		Direction:    v2.Direction_A_To_B,
		StopReceived: big.NewInt(1),
	}
	step2 := v2.Stop{
		SwapAmount: v2.SwapAmount{
			Type:   v2.AmountType_Specific_Amount,
			Amount: big.NewInt(10000),
		},
		Direction:    v2.Direction_B_To_A,
		StopReceived: big.NewInt(1),
	}
	datums := []v2.OrderDatum{}
	datums = append(datums, buildOrderDatum(step1)...)
	datums = append(datums, buildOrderDatum(step2)...)
	for _, datum := range datums {
		plutusData := datum.ToPlutusData()
		b, err := plutusData.MarshalCBOR()
		if err != nil {
			t.Errorf("Stop ToPlutusData error: %s\n", err)
		}

		var convertedPlutusData PlutusData.PlutusData
		err = cbor.Unmarshal(b, &convertedPlutusData)
		if err != nil {
			t.Errorf("Stop ToPlutusData error: %s\n", err)
		}

		if !convertedPlutusData.Equal(plutusData) {
			t.Error("TestStopStepConverter error")
		}
	}
}

func TestOCOStepConverter(t *testing.T) {
	step1 := v2.OCO{
		SwapAmount: v2.SwapAmount{
			Type:   v2.AmountType_Specific_Amount,
			Amount: big.NewInt(10000),
		},
		Direction:       v2.Direction_A_To_B,
		StopReceived:    big.NewInt(1),
		MinimumReceived: big.NewInt(1),
	}
	step2 := v2.OCO{
		SwapAmount: v2.SwapAmount{
			Type:   v2.AmountType_All,
			Amount: big.NewInt(10000),
		},
		Direction:       v2.Direction_B_To_A,
		StopReceived:    big.NewInt(1),
		MinimumReceived: big.NewInt(1),
	}
	datums := []v2.OrderDatum{}
	datums = append(datums, buildOrderDatum(step1)...)
	datums = append(datums, buildOrderDatum(step2)...)
	for _, datum := range datums {
		plutusData := datum.ToPlutusData()
		b, err := plutusData.MarshalCBOR()
		if err != nil {
			t.Errorf("OCO ToPlutusData error: %s\n", err)
		}

		var convertedPlutusData PlutusData.PlutusData
		err = cbor.Unmarshal(b, &convertedPlutusData)
		if err != nil {
			t.Errorf("OCO ToPlutusData error: %s\n", err)
		}

		if !convertedPlutusData.Equal(plutusData) {
			t.Error("TestOCOStepConverter error")
		}
	}
}

func TestSwapExactOutStepConverter(t *testing.T) {
	step1 := v2.SwapExactOut{
		MaximumSwapAmount: v2.SwapAmount{
			Type:   v2.AmountType_Specific_Amount,
			Amount: big.NewInt(10000),
		},
		Direction:        v2.Direction_A_To_B,
		ExpectedReceived: big.NewInt(1),
		Killable:         v2.Killable_Pending_On_Failed,
	}
	step2 := v2.SwapExactOut{
		MaximumSwapAmount: v2.SwapAmount{
			Type:   v2.AmountType_All,
			Amount: big.NewInt(10000),
		},
		Direction:        v2.Direction_B_To_A,
		ExpectedReceived: big.NewInt(1),
		Killable:         v2.Killable_Kill_On_Failed,
	}
	datums := []v2.OrderDatum{}
	datums = append(datums, buildOrderDatum(step1)...)
	datums = append(datums, buildOrderDatum(step2)...)
	for _, datum := range datums {
		plutusData := datum.ToPlutusData()
		b, err := plutusData.MarshalCBOR()
		if err != nil {
			t.Errorf("SwapExactOut ToPlutusData error: %s\n", err)
		}

		var convertedPlutusData PlutusData.PlutusData
		err = cbor.Unmarshal(b, &convertedPlutusData)
		if err != nil {
			t.Errorf("SwapExactOut ToPlutusData error: %s\n", err)
		}

		if !convertedPlutusData.Equal(plutusData) {
			t.Error("TestSwapExactOutStepConverter error")
		}
	}
}

func TestDepositStepConverter(t *testing.T) {
	step1 := v2.Deposit{
		DepositAmount: v2.DepositAmount{
			Type:           v2.AmountType_Specific_Amount,
			DepositAmountA: big.NewInt(10000),
			DepositAmountB: big.NewInt(10000),
		},
		MinimumLP: big.NewInt(1),
		Killable:  v2.Killable_Pending_On_Failed,
	}
	step2 := v2.Deposit{
		DepositAmount: v2.DepositAmount{
			Type:           v2.AmountType_All,
			DepositAmountA: big.NewInt(10000),
			DepositAmountB: big.NewInt(10000),
		},
		MinimumLP: big.NewInt(1),
		Killable:  v2.Killable_Kill_On_Failed,
	}
	datums := []v2.OrderDatum{}
	datums = append(datums, buildOrderDatum(step1)...)
	datums = append(datums, buildOrderDatum(step2)...)
	for _, datum := range datums {
		plutusData := datum.ToPlutusData()
		b, err := plutusData.MarshalCBOR()
		if err != nil {
			t.Errorf("Deposit ToPlutusData error: %s\n", err)
		}

		var convertedPlutusData PlutusData.PlutusData
		err = cbor.Unmarshal(b, &convertedPlutusData)
		if err != nil {
			t.Errorf("Deposit ToPlutusData error: %s\n", err)
		}

		if !convertedPlutusData.Equal(plutusData) {
			t.Error("TestDepositStepConverter error")
		}
	}
}

func TestWithdrawStepConverter(t *testing.T) {
	step1 := v2.Withdraw{
		WithdrawalAmount: v2.WithdrawalAmount{
			Type:     v2.AmountType_Specific_Amount,
			LPAmount: big.NewInt(10000),
		},
		MinimumAssetA: big.NewInt(1),
		MinimumAssetB: big.NewInt(1),
		Killable:      v2.Killable_Pending_On_Failed,
	}
	step2 := v2.Withdraw{
		WithdrawalAmount: v2.WithdrawalAmount{
			Type:     v2.AmountType_All,
			LPAmount: big.NewInt(10000),
		},
		MinimumAssetA: big.NewInt(1),
		MinimumAssetB: big.NewInt(1),
		Killable:      v2.Killable_Kill_On_Failed,
	}
	datums := []v2.OrderDatum{}
	datums = append(datums, buildOrderDatum(step1)...)
	datums = append(datums, buildOrderDatum(step2)...)
	for _, datum := range datums {
		plutusData := datum.ToPlutusData()
		b, err := plutusData.MarshalCBOR()
		if err != nil {
			t.Errorf("Withdraw ToPlutusData error: %s\n", err)
		}

		var convertedPlutusData PlutusData.PlutusData
		err = cbor.Unmarshal(b, &convertedPlutusData)
		if err != nil {
			t.Errorf("Withdraw ToPlutusData error: %s\n", err)
		}

		if !convertedPlutusData.Equal(plutusData) {
			t.Error("TestWithdrawStepConverter error")
		}
	}
}

func TestZapOutStepConverter(t *testing.T) {
	step1 := v2.ZapOut{
		WithdrawalAmount: v2.WithdrawalAmount{
			Type:     v2.AmountType_Specific_Amount,
			LPAmount: big.NewInt(10000),
		},
		Direction:       v2.Direction_A_To_B,
		MinimumReceived: big.NewInt(1),
		Killable:        v2.Killable_Pending_On_Failed,
	}
	step2 := v2.ZapOut{
		WithdrawalAmount: v2.WithdrawalAmount{
			Type:     v2.AmountType_All,
			LPAmount: big.NewInt(10000),
		},
		Direction:       v2.Direction_B_To_A,
		MinimumReceived: big.NewInt(1),
		Killable:        v2.Killable_Kill_On_Failed,
	}
	datums := []v2.OrderDatum{}
	datums = append(datums, buildOrderDatum(step1)...)
	datums = append(datums, buildOrderDatum(step2)...)
	for _, datum := range datums {
		plutusData := datum.ToPlutusData()
		b, err := plutusData.MarshalCBOR()
		if err != nil {
			t.Errorf("ZapOut ToPlutusData error: %s\n", err)
		}

		var convertedPlutusData PlutusData.PlutusData
		err = cbor.Unmarshal(b, &convertedPlutusData)
		if err != nil {
			t.Errorf("ZapOut ToPlutusData error: %s\n", err)
		}

		if !convertedPlutusData.Equal(plutusData) {
			t.Error("TestZapOutStepConverter error")
		}
	}
}

func TestPartialSwapStepConverter(t *testing.T) {
	step1 := v2.PartialSwap{
		TotalSwapAmount: big.NewInt(10000),
		IoRatioNumerator: big.NewInt(1),
		IoRatioDenominator: big.NewInt(1),
		Hops: big.NewInt(3),
		Direction: v2.Direction_A_To_B,
		MaxBatcherFeeEachTime: big.NewInt(0).Mul(v2.FIXED_BATCHER_FEE, big.NewInt(3)),
		MinimumSwapAmountRequired: big.NewInt(1000),
	}
	step2 := v2.PartialSwap{
		TotalSwapAmount: big.NewInt(10000),
		IoRatioNumerator: big.NewInt(1),
		IoRatioDenominator: big.NewInt(1),
		Hops: big.NewInt(3),
		Direction: v2.Direction_B_To_A,
		MaxBatcherFeeEachTime: big.NewInt(0).Mul(v2.FIXED_BATCHER_FEE, big.NewInt(3)),
		MinimumSwapAmountRequired: big.NewInt(1000),
	}
	datums := []v2.OrderDatum{}
	datums = append(datums, buildOrderDatum(step1)...)
	datums = append(datums, buildOrderDatum(step2)...)
	for _, datum := range datums {
		plutusData := datum.ToPlutusData()
		b, err := plutusData.MarshalCBOR()
		if err != nil {
			t.Errorf("PartialSwap ToPlutusData error: %s\n", err)
		}

		var convertedPlutusData PlutusData.PlutusData
		err = cbor.Unmarshal(b, &convertedPlutusData)
		if err != nil {
			t.Errorf("PartialSwap ToPlutusData error: %s\n", err)
		}

		if !convertedPlutusData.Equal(plutusData) {
			t.Error("TestPartialSwapStepConverter error")
		}
	}
}

func TestWithdrawImbalanceStepConverter(t *testing.T) {
	step1 := v2.WithdrawImbalance{
		WithdrawAmount: v2.WithdrawalAmount{
			Type:     v2.AmountType_Specific_Amount,
			LPAmount: big.NewInt(10000),
		},
		Killable: v2.Killable_Pending_On_Failed,
		RatioAssetA: big.NewInt(1),
		RatioAssetB: big.NewInt(1),
		MinimumAssetA: big.NewInt(1000),
	}
	step2 := v2.WithdrawImbalance{
		WithdrawAmount: v2.WithdrawalAmount{
			Type:     v2.AmountType_All,
			LPAmount: big.NewInt(10000),
		},
		Killable: v2.Killable_Kill_On_Failed,
		RatioAssetA: big.NewInt(1),
		RatioAssetB: big.NewInt(1),
		MinimumAssetA: big.NewInt(1000),
	}
	datums := []v2.OrderDatum{}
	datums = append(datums, buildOrderDatum(step1)...)
	datums = append(datums, buildOrderDatum(step2)...)
	for _, datum := range datums {
		plutusData := datum.ToPlutusData()
		b, err := plutusData.MarshalCBOR()
		if err != nil {
			t.Errorf("WithdrawImbalance ToPlutusData error: %s\n", err)
		}

		var convertedPlutusData PlutusData.PlutusData
		err = cbor.Unmarshal(b, &convertedPlutusData)
		if err != nil {
			t.Errorf("WithdrawImbalance ToPlutusData error: %s\n", err)
		}

		if !convertedPlutusData.Equal(plutusData) {
			t.Error("TestWithdrawImbalanceStepConverter error")
		}
	}
}

func TestSwapRoutingStepConverter(t *testing.T) {
	step := v2.SwapRouting{
		SwapAmount: v2.SwapAmount{
			Type:   v2.AmountType_Specific_Amount,
			Amount: big.NewInt(10000),
		},
		MinimumReceived: big.NewInt(1),
		Routings: []v2.Route {
			{
				LPAsset: Fingerprint.Fingerprint{
					PolicyId: Policy.PolicyId{Value: "f5808c2c990d86da54bfc97d89cee6efa20cd8461616359478d96b4c"},
					AssetName: AssetName.AssetName{Value: "ef4530398e53eea75ee3d02a982e87a5c680776904b5d610e63bf6970c528a12"},
				},
				Direction: v2.Direction_A_To_B,
			},
			{
				LPAsset: Fingerprint.Fingerprint{
					PolicyId: Policy.PolicyId{Value: "f5808c2c990d86da54bfc97d89cee6efa20cd8461616359478d96b4c"},
					AssetName: AssetName.AssetName{Value: "eebaae50fe9a09938558096cfebe0aec7dd2728dedadb3d96f02f19e756ca9b8"},
				},
				Direction: v2.Direction_B_To_A,
			},
		},
	}
	datums := []v2.OrderDatum{}
	datums = append(datums, buildOrderDatum(step)...)
	for _, datum := range datums {
		plutusData := datum.ToPlutusData()
		b, err := plutusData.MarshalCBOR()
		if err != nil {
			t.Errorf("SwapRouting ToPlutusData error: %s\n", err)
		}

		var convertedPlutusData PlutusData.PlutusData
		err = cbor.Unmarshal(b, &convertedPlutusData)
		if err != nil {
			t.Errorf("SwapRouting ToPlutusData error: %s\n", err)
		}

		if !convertedPlutusData.Equal(plutusData) {
			t.Error("TestSwapRoutingStepConverter error")
		}
	}
}

func TestDonationStepConverter(t *testing.T) {
	step := v2.Donation{ }
	datums := []v2.OrderDatum{}
	datums = append(datums, buildOrderDatum(step)...)
	for _, datum := range datums {
		plutusData := datum.ToPlutusData()
		b, err := plutusData.MarshalCBOR()
		if err != nil {
			t.Errorf("Donation ToPlutusData error: %s\n", err)
		}

		var convertedPlutusData PlutusData.PlutusData
		err = cbor.Unmarshal(b, &convertedPlutusData)
		if err != nil {
			t.Errorf("Donation ToPlutusData error: %s\n", err)
		}

		if !convertedPlutusData.Equal(plutusData) {
			t.Error("TestDonationStepConverter error")
		}
	}
}

func TestOrderRedeemer(t *testing.T) {
	b, _ := v2.OrderRedeemer_ApplyOrder.Data.MarshalCBOR()
	if hex.EncodeToString(b) != "d8799fff" {
		t.Errorf("OrderRedeemer_ApplyOrder excepted d8799fff but get %s\n", hex.EncodeToString(b))
	}

	b, _ = v2.OrderRedeemer_CancelOrderByOwner.Data.MarshalCBOR()
	if hex.EncodeToString(b) != "d87a9fff" {
		t.Errorf("OrderRedeemer_ApplyOrder excepted d87a9fff but get %s\n", hex.EncodeToString(b))
	}

	b, _ = v2.OrderRedeemer_CancelExpiredOrderByAnyone.Data.MarshalCBOR()
	if hex.EncodeToString(b) != "d87b9fff" {
		t.Errorf("OrderRedeemer_ApplyOrder excepted d87b9fff but get %s\n", hex.EncodeToString(b))
	}
}
