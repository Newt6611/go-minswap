package stable

import (
	"github.com/Newt6611/apollo/plutusencoder"
	"github.com/Newt6611/apollo/serialization/Address"
	"github.com/Newt6611/apollo/serialization/PlutusData"
)

type StepType int

const (
	StepType_Swap = iota
	StepType_Deposit
	StepType_Withdraw
	StepType_WithdrawImbalance
	StepType_Zapout
)

type SwapStep struct {
    Type StepType
    AssetInIndex uint64
    AssetOutIndex uint64
    MinimumAssetOut uint64
}

func (s SwapStep) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Swap),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.AssetInIndex,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.AssetOutIndex,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.MinimumAssetOut,
			},
		},
	}
}

type DepositStep struct {
    Type StepType
    MinimumLP uint64
}

func (s DepositStep) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Deposit),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.MinimumLP,
			},
		},
	}
}

type WithdrawStep struct {
    Type StepType
    MinimumAmounts []uint64
}

func (s WithdrawStep) StepToPlutusData() PlutusData.PlutusData {
	minimumAmountsPlutusDataArray := PlutusData.PlutusIndefArray{}
	for _, amount := range s.MinimumAmounts {
		minimumAmountsPlutusDataArray = append(minimumAmountsPlutusDataArray, PlutusData.PlutusData{
			TagNr:          0,
			PlutusDataType: PlutusData.PlutusInt,
			Value:          amount,
		})
	}

	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Withdraw),
		PlutusDataType: PlutusData.PlutusArray,
		Value:  minimumAmountsPlutusDataArray,
	}
}

type WithdrawImbalanceStep struct {
	Type StepType
	WithdrawAmounts []uint64
}

func (s WithdrawImbalanceStep) StepToPlutusData() PlutusData.PlutusData {
	withdrawAmountsPlutusDataArray := PlutusData.PlutusIndefArray{}
	for _, amount := range s.WithdrawAmounts {
		withdrawAmountsPlutusDataArray = append(withdrawAmountsPlutusDataArray, PlutusData.PlutusData{
			TagNr:          0,
			PlutusDataType: PlutusData.PlutusInt,
			Value:          amount,
		})
	}

	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_WithdrawImbalance),
		PlutusDataType: PlutusData.PlutusArray,
		Value:  withdrawAmountsPlutusDataArray,
	}
}

type ZapOutStep struct {
    Type StepType
    AssetOutIndex uint64
    MinimumAssetOut uint64
}

func (s ZapOutStep) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Zapout),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.AssetOutIndex,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          s.MinimumAssetOut,
			},
		},
    }
}

type IStep interface {
	StepToPlutusData() PlutusData.PlutusData
}

type OrderDatum struct {
	Sender            Address.Address
	Receiver          Address.Address
	ReceiverDatumHash []byte
	Step              IStep
	BatcherFee        uint64
	OutputAda         uint64
}

func (o OrderDatum) ToPlutusData() PlutusData.PlutusData {
	senderPlutusData, _ := plutusencoder.GetAddressPlutusData(o.Sender)
	receiverPlutusData, _ := plutusencoder.GetAddressPlutusData(o.Receiver)

	var receiverDatumHashPlutusData PlutusData.PlutusData
	if o.ReceiverDatumHash != nil && len(o.ReceiverDatumHash) != 0 {
		receiverDatumHashPlutusData = PlutusData.PlutusData{
			TagNr:          121,
			PlutusDataType: PlutusData.PlutusArray,
			Value: PlutusData.PlutusDefArray{
				PlutusData.PlutusData{
					TagNr:          0,
					PlutusDataType: PlutusData.PlutusBytes,
					Value:          o.ReceiverDatumHash,
				},
			},
		}
	} else {
		receiverDatumHashPlutusData = PlutusData.PlutusData{
			TagNr:          121 + 1,
			PlutusDataType: PlutusData.PlutusArray,
			Value:          PlutusData.PlutusDefArray{},
		}
	}

	return PlutusData.PlutusData{
		TagNr:          121,
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			*senderPlutusData,
			*receiverPlutusData,
			receiverDatumHashPlutusData,
			o.Step.StepToPlutusData(),
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          o.BatcherFee,
			},
			PlutusData.PlutusData{
				TagNr:          0,
				PlutusDataType: PlutusData.PlutusInt,
				Value:          o.OutputAda,
			},
		},
	}
}
