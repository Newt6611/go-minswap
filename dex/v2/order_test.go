package v2_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/Newt6611/apollo/serialization/Address"
	"github.com/Newt6611/apollo/serialization/AssetName"
	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/Policy"
	v2 "github.com/Newt6611/go-minswap/dex/v2"
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
			Step:           step,
			LpAsset:        *fingerprint,
			MaxBatcherFee:  v2.FIXED_BATCHER_FEE,
			ExpiredOptions: v2.ExpirySetting{
				ExpiredTime: big.NewInt(1721010208050),
				MaxCancellationTip: big.NewInt(300_000),
			},
		},
	}
}

// TODO:
func TestSwapExactInStepToPlutusData(t *testing.T) {
	step1 := v2.SwapExactIn {
		SwapAmount: v2.SwapAmount{
			Type: v2.AmountType_Specific_Amount,
			Amount: big.NewInt(10000),
		},
		Direction: v2.Direction_A_To_B,
		MinimumReceived: big.NewInt(1),
		Killable: v2.Killable_Pending_On_Failed,
	}
	step2 := v2.SwapExactIn {
		SwapAmount: v2.SwapAmount{
			Type: v2.AmountType_All,
			Amount: big.NewInt(10000),
		},
		Direction: v2.Direction_B_To_A,
		MinimumReceived: big.NewInt(1),
		Killable: v2.Killable_Kill_On_Failed,
	}
	datums := []v2.OrderDatum{}
	datums = append(datums, buildOrderDatum(step1)...)
	datums = append(datums, buildOrderDatum(step2)...)
	for _, datum := range datums {
		_ = datum
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
