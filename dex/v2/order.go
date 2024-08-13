package v2

import (
	"math/big"

	"github.com/Newt6611/apollo/plutusencoder"
	"github.com/Newt6611/apollo/serialization/Address"
	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/PlutusData"
	"github.com/Newt6611/apollo/serialization/Redeemer"
)

var (
	FIXED_BATCHER_FEE = big.NewInt(2_000_000)
)

type AuthorizationMethodType int

const (
	AuthorizationMethodType_Signature AuthorizationMethodType = iota
	AuthorizationMethodType_Spend_Script
	AuthorizationMethodType_Withdraw_Script
	AuthorizationMethodType_Mint_Script
)

type AuthorizationMethod struct {
	Type AuthorizationMethodType
	Hash []byte
}

func (a AuthorizationMethod) ToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121,
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusBytes,
				Value:          a.Hash,
			},
		},
	}
}

type ExtraDatumType int

const (
	ExtraDatumType_No_Datum = iota
	ExtraDatumType_Datum_Hash
	ExtraDatumType_Inline_Datum
)

type ExtraDatum struct {
	Type ExtraDatumType
	Hash []byte
}

func (e ExtraDatum) ToPlutusData() PlutusData.PlutusData {
	if e.Type == ExtraDatumType_No_Datum {
		return PlutusData.PlutusData{
			TagNr:          121,
			PlutusDataType: PlutusData.PlutusArray,
			Value:          PlutusData.PlutusDefArray{},
		}
	}

	return PlutusData.PlutusData{
		TagNr:          121 + uint64(e.Type),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusDefArray{
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusBytes,
				Value:          e.Hash,
			},
		},
	}
}

type StepType int

const (
	StepType_Swap_Exact_In StepType = iota
	StepType_Stop
	StepType_OCO
	StepType_Swap_Exact_Out
	StepType_Deposit
	StepType_Withdraw
	StepType_Zap_Out
	StepType_Partial_Swap
	StepType_Withdraw_Imbalance
	StepType_Swap_Routing
	StepType_Donation
)

type Direction int

const (
	Direction_B_To_A Direction = iota
	Direction_A_To_B
)

func (d Direction) ToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(d),
		PlutusDataType: PlutusData.PlutusArray,
		Value:          PlutusData.PlutusDefArray{},
	}
}

type AmountType int

const (
	AmountType_Specific_Amount AmountType = iota
	AmountType_All
)

type SwapAmount struct {
	Type   AmountType
	Amount *big.Int
}

func (s SwapAmount) ToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(s.Type),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.Amount,
			},
		},
	}
}

type Killable int

const (
	Killable_Pending_On_Failed Killable = iota
	Killable_Kill_On_Failed
)

func (k Killable) ToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(k),
		PlutusDataType: PlutusData.PlutusArray,
		Value:          PlutusData.PlutusDefArray{},
	}
}

type DepositAmount struct {
	Type           AmountType
	DepositAmountA *big.Int
	DepositAmountB *big.Int
}

func (d DepositAmount) ToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(d.Type),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          d.DepositAmountA,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          d.DepositAmountB,
			},
		},
	}
}

type WithdrawAmount struct {
	Type     AmountType
	LPAmount *big.Int
}

func (w WithdrawAmount) ToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(w.Type),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          w.LPAmount,
			},
		},
	}
}

type Route struct {
	LPAsset   Fingerprint.Fingerprint
	Direction Direction
}

func (r Route) ToPlutusData() PlutusData.PlutusData {
	lpAssetPlutusData, _ := r.LPAsset.ToPlutusData()

	return PlutusData.PlutusData{
		TagNr:          121,
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			lpAssetPlutusData,
			r.Direction.ToPlutusData(),
		},
	}
}

type StepI interface {
	StepToPlutusData() PlutusData.PlutusData
}

type SwapExactIn struct {
	Direction       Direction
	SwapAmount      SwapAmount
	MinimumReceived *big.Int
	Killable        Killable
}

func (s SwapExactIn) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Swap_Exact_In),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			s.Direction.ToPlutusData(),
			s.SwapAmount.ToPlutusData(),
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.MinimumReceived,
			},
			s.Killable.ToPlutusData(),
		},
	}
}

type Stop struct {
	Direction    Direction
	SwapAmount   SwapAmount
	StopReceived *big.Int
}

func (s Stop) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Stop),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			s.Direction.ToPlutusData(),
			s.SwapAmount.ToPlutusData(),
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.StopReceived,
			},
		},
	}
}

type OCO struct {
	Direction       Direction
	SwapAmount      SwapAmount
	MinimumReceived *big.Int
	StopReceived    *big.Int
}

func (s OCO) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_OCO),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			s.Direction.ToPlutusData(),
			s.SwapAmount.ToPlutusData(),
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.MinimumReceived,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.StopReceived,
			},
		},
	}
}

type SwapExactOut struct {
	Direction         Direction
	MaximumSwapAmount SwapAmount
	ExpectedReceived  *big.Int
	Killable          Killable
}

func (s SwapExactOut) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Swap_Exact_Out),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			s.Direction.ToPlutusData(),
			s.MaximumSwapAmount.ToPlutusData(),
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.ExpectedReceived,
			},
			s.Killable.ToPlutusData(),
		},
	}
}

type Deposit struct {
	DepositAmount DepositAmount
	MinimumLP     *big.Int
	Killable      Killable
}

func (s Deposit) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Deposit),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			s.DepositAmount.ToPlutusData(),
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.MinimumLP,
			},
			s.Killable.ToPlutusData(),
		},
	}
}

type Withdraw struct {
	WithdrawAmount WithdrawAmount
	MinimumAssetA  *big.Int
	MinimumAssetB  *big.Int
	Killable       Killable
}

func (s Withdraw) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Withdraw),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			s.WithdrawAmount.ToPlutusData(),
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.MinimumAssetA,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.MinimumAssetB,
			},
			s.Killable.ToPlutusData(),
		},
	}
}

type ZapOut struct {
	Direction        Direction
	WithdrawalAmount WithdrawAmount
	MinimumReceived  *big.Int
	Killable         Killable
}

func (z ZapOut) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Zap_Out),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			z.Direction.ToPlutusData(),
			z.WithdrawalAmount.ToPlutusData(),
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          z.MinimumReceived,
			},
			z.Killable.ToPlutusData(),
		},
	}
}

type PartialSwap struct {
	Direction                 Direction
	TotalSwapAmount           *big.Int
	IoRatioNumerator          *big.Int
	IoRatioDenominator        *big.Int
	Hops                      *big.Int
	MinimumSwapAmountRequired *big.Int
	MaxBatcherFeeEachTime     *big.Int
}

func (p PartialSwap) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Partial_Swap),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			p.Direction.ToPlutusData(),
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          p.TotalSwapAmount,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          p.IoRatioNumerator,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          p.IoRatioDenominator,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          p.Hops,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          p.MinimumSwapAmountRequired,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          p.MaxBatcherFeeEachTime,
			},
		},
	}
}

type WithdrawImbalance struct {
	WithdrawAmount WithdrawAmount
	RatioAssetA    *big.Int
	RatioAssetB    *big.Int
	MinimumAssetA  *big.Int
	Killable       Killable
}

func (w WithdrawImbalance) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Withdraw_Imbalance),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			w.WithdrawAmount.ToPlutusData(),
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          w.RatioAssetA,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          w.RatioAssetB,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          w.MinimumAssetA,
			},
			w.Killable.ToPlutusData(),
		},
	}
}

type SwapRouting struct {
	Routings        []Route
	SwapAmount      SwapAmount
	MinimumReceived *big.Int
}

func (s SwapRouting) StepToPlutusData() PlutusData.PlutusData {
	routingsPlutusDataArray := PlutusData.PlutusIndefArray{}
	for _, r := range s.Routings {
		routingsPlutusDataArray = append(routingsPlutusDataArray, r.ToPlutusData())
	}
	routingsPlutusData := PlutusData.PlutusData{
		TagNr:          0,
		PlutusDataType: PlutusData.PlutusArray,
		Value:          routingsPlutusDataArray,
	}

	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Swap_Routing),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			routingsPlutusData,
			s.SwapAmount.ToPlutusData(),
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.MinimumReceived,
			},
		},
	}
}

type Donation struct{}

func (d Donation) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Donation),
		PlutusDataType: PlutusData.PlutusArray,
		Value:          PlutusData.PlutusDefArray{},
	}
}

type ExpirySetting struct {
	ExpiredTime        *big.Int
	MaxCancellationTip *big.Int
}

func (e ExpirySetting) ToPlutusData() PlutusData.PlutusData {
	if e.ExpiredTime != nil || e.MaxCancellationTip != nil {
		return PlutusData.PlutusData{
			TagNr:          121,
			PlutusDataType: PlutusData.PlutusArray,
			Value: PlutusData.PlutusDefArray{
				PlutusData.PlutusData{
					TagNr:          0,
					PlutusDataType: PlutusData.PlutusInt,
					Value:          e.ExpiredTime,
				},
				PlutusData.PlutusData{
					TagNr:          0,
					PlutusDataType: PlutusData.PlutusInt,
					Value:          e.MaxCancellationTip,
				},
			},
		}
	} else {
		return PlutusData.PlutusData{
			TagNr:          121 + 1,
			PlutusDataType: PlutusData.PlutusArray,
			Value:          PlutusData.PlutusDefArray{},
		}
	}
}

type OrderDatum struct {
	Canceller            AuthorizationMethod
	RefundReceiver       Address.Address
	RefundReceiverDatum  ExtraDatum
	SuccessReceiver      Address.Address
	SuccessReceiverDatum ExtraDatum
	LpAsset              Fingerprint.Fingerprint
	Step                 StepI
	MaxBatcherFee        *big.Int
	ExpiredOptions       ExpirySetting
}

func (o OrderDatum) ToPlutusData() PlutusData.PlutusData {
	refundReceiverPlutusData, _ := plutusencoder.GetAddressPlutusData(o.RefundReceiver)
	successReceiverPlutusData, _ := plutusencoder.GetAddressPlutusData(o.SuccessReceiver)
	lpAssetPlutusData, _ := o.LpAsset.ToPlutusData()

	return PlutusData.PlutusData{
		TagNr:          121,
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			o.Canceller.ToPlutusData(),
			*refundReceiverPlutusData,
			o.RefundReceiverDatum.ToPlutusData(),
			*successReceiverPlutusData,
			o.SuccessReceiverDatum.ToPlutusData(),
			lpAssetPlutusData,
			o.Step.StepToPlutusData(),
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          o.MaxBatcherFee,
			},
			o.ExpiredOptions.ToPlutusData(),
		},
	}
}

var (
	OrderRedeemer_ApplyOrder = Redeemer.Redeemer {
		Tag:   Redeemer.SPEND,
		Index: 0,
		Data: PlutusData.PlutusData{
			TagNr:          121,
			PlutusDataType: PlutusData.PlutusArray,
			Value:          PlutusData.PlutusIndefArray{},
		},
	}

	OrderRedeemer_CancelOrderByOwner = Redeemer.Redeemer{
		Tag:   Redeemer.SPEND,
		Index: 0,
		Data: PlutusData.PlutusData{
			TagNr:          121 + 1,
			PlutusDataType: PlutusData.PlutusArray,
			Value:          PlutusData.PlutusIndefArray{},
		},
	}

	OrderRedeemer_CancelExpiredOrderByAnyone = Redeemer.Redeemer{
		Tag:   Redeemer.SPEND,
		Index: 0,
		Data: PlutusData.PlutusData{
			TagNr:          121 + 2,
			PlutusDataType: PlutusData.PlutusArray,
			Value:          PlutusData.PlutusIndefArray{},
		},
	}
)
