package v2

import (
	"errors"
	"fmt"

	c "github.com/Newt6611/apollo/constants"
	"github.com/Newt6611/apollo/plutusencoder"
	"github.com/Newt6611/apollo/serialization/Address"
	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/PlutusData"
	"github.com/Newt6611/apollo/serialization/Redeemer"
	"github.com/Newt6611/go-minswap/utils"
)

const (
	FIXED_BATCHER_FEE = 2_000_000
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

func AuthorizationMethodFromPlutusData(plutusData *PlutusData.PlutusData) (AuthorizationMethod, error) {
	var authorizationMethod AuthorizationMethod
	hash, ok := plutusData.Value.(PlutusData.PlutusIndefArray)[0].Value.([]byte)
	if !ok {
		return authorizationMethod, fmt.Errorf("invalid AuthorizationMethodFromPlutusData")
	}

	authorizationMethod.Type = AuthorizationMethodType(plutusData.TagNr - 121)
	authorizationMethod.Hash = hash
	return authorizationMethod, nil
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

func ExtraDatumFromPlutusData(plutusData *PlutusData.PlutusData) (ExtraDatum, error) {
	var extraDatum ExtraDatum
	data := plutusData.Value.(PlutusData.PlutusDefArray)
	extraDatum.Type = ExtraDatumType(plutusData.TagNr - 121)
	if len(data) == 0 {
		return extraDatum, nil
	}

	hash, ok := data[0].Value.([]byte)
	if !ok {
		return extraDatum, fmt.Errorf("invalid ExtraDatumFromPlutusData")
	}

	extraDatum.Hash = hash
	return extraDatum, nil
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

func DirectionFromPlutusData(plutusData *PlutusData.PlutusData) (Direction, error) {
	direction := Direction(plutusData.TagNr - 121)
	if direction != Direction_B_To_A && direction != Direction_A_To_B {
		return direction, fmt.Errorf("invalid DirectionFromPlutusData")
	}
	return direction, nil
}

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
	Amount uint64
}

func SwapAmountFromPlutusData(plutusData *PlutusData.PlutusData) (SwapAmount, error) {
	var swapAmount SwapAmount
	swapAmount.Type = AmountType(plutusData.TagNr - 121)
	amount, ok := plutusData.Value.(PlutusData.PlutusIndefArray)[0].Value.(uint64)
	if !ok {
		return swapAmount, fmt.Errorf("invalid SwapAmountFromPlutusData")
	}
	swapAmount.Amount = amount
	return swapAmount, nil
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

func KillableFromPlutusData(plutusData *PlutusData.PlutusData) (Killable, error) {
	killable := Killable(plutusData.TagNr - 121)
	if killable != Killable_Pending_On_Failed && killable != Killable_Kill_On_Failed {
		return killable, fmt.Errorf("invalid KillableFromPlutusData")
	}
	return killable, nil
}

func (k Killable) ToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(k),
		PlutusDataType: PlutusData.PlutusArray,
		Value:          PlutusData.PlutusDefArray{},
	}
}

type DepositAmount struct {
	Type           AmountType
	DepositAmountA uint64
	DepositAmountB uint64
}

func DepositAmountFromPlutusData(plutusData *PlutusData.PlutusData) (DepositAmount, error) {
	var depositAmount DepositAmount
	depositAmount.Type = AmountType(plutusData.TagNr - 121)
	data := plutusData.Value.(PlutusData.PlutusIndefArray)

	depositAmountA, ok := data[0].Value.(uint64)
	if !ok {
		return depositAmount, errors.New("invalid DepositAmountFromPlutusData DepositAmountA")
	}
	depositAmount.DepositAmountA = depositAmountA

	depositAmountB, ok := data[1].Value.(uint64)
	if !ok {
		return depositAmount, errors.New("invalid DepositAmountFromPlutusData DepositAmountB")
	}
	depositAmount.DepositAmountB = depositAmountB

	return depositAmount, nil
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

type WithdrawalAmount struct {
	Type     AmountType
	LPAmount uint64
}

func WithdrawAmountFromPlutusData(plutusData *PlutusData.PlutusData) (WithdrawalAmount, error) {
	var withdrawAmount WithdrawalAmount
	withdrawAmount.Type = AmountType(plutusData.TagNr - 121)
	data := plutusData.Value.(PlutusData.PlutusIndefArray)

	lpAmount, ok := data[0].Value.(uint64)
	if !ok {
		return withdrawAmount, errors.New("invalid WithdrawAmountFromPlutusData LPAmount")
	}
	withdrawAmount.LPAmount = lpAmount

	return withdrawAmount, nil
}

func (w WithdrawalAmount) ToPlutusData() PlutusData.PlutusData {
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

func RouteFromPlutusData(plutusData *PlutusData.PlutusData) (Route, error) {
	var err error
	var route Route
	data := plutusData.Value.(PlutusData.PlutusIndefArray)
	route.LPAsset, err = utils.FingerprintFromPlutusData(&data[0])
	if err != nil {
		return route, err
	}

	route.Direction, err = DirectionFromPlutusData(&data[1])
	if err != nil {
		return route, err
	}

	return route, nil
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

func StepFromPlutusData(plutusData *PlutusData.PlutusData) (StepI, error) {
	stepType := StepType(plutusData.TagNr - 121)
	switch stepType {
	case StepType_Swap_Exact_In:
		return SwapExactInFromPlutusData(plutusData)
	case StepType_Stop:
		return StopFromPlutusData(plutusData)
	case StepType_OCO:
		return OCOFromPlutusData(plutusData)
	case StepType_Swap_Exact_Out:
		return SwapExactOutFromPlutusData(plutusData)
	case StepType_Deposit:
		return DepositFromPlutusData(plutusData)
	case StepType_Withdraw:
		return WithdrawFromPlutusData(plutusData)
	case StepType_Zap_Out:
		return ZapOutFromPlutusData(plutusData)
	case StepType_Partial_Swap:
		return PartialSwapFromPlutusData(plutusData)
	case StepType_Withdraw_Imbalance:
		return WithdrawImbalanceFromPlutusData(plutusData)
	case StepType_Swap_Routing:
		return SwapRoutingFromPlutusData(plutusData)
	case StepType_Donation:
		return DonationFromPlutusData(plutusData)
	}

	return nil, fmt.Errorf("invalid StepFromPlutusData")
}

type SwapExactIn struct {
	Type            StepType
	Direction       Direction
	SwapAmount      SwapAmount
	MinimumReceived uint64
	Killable        Killable
}

func SwapExactInFromPlutusData(plutusData *PlutusData.PlutusData) (SwapExactIn, error) {
	var err error
	data := plutusData.Value.(PlutusData.PlutusIndefArray)
	var swapExactIn SwapExactIn

	swapExactIn.Type = StepType_Swap_Exact_In
	swapExactIn.Direction, err = DirectionFromPlutusData(&data[0])
	if err != nil {
		return swapExactIn, err
	}

	swapExactIn.SwapAmount, err = SwapAmountFromPlutusData(&data[1])
	if err != nil {
		return swapExactIn, err
	}

	minimumReceived, ok := data[2].Value.(uint64)
	if !ok {
		return swapExactIn, fmt.Errorf("invalid SwapExactInFromPlutusData")
	}
	swapExactIn.MinimumReceived = minimumReceived

	killable, err := KillableFromPlutusData(&data[3])
	if err != nil {
		return swapExactIn, err
	}
	swapExactIn.Killable = killable

	return swapExactIn, nil
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
	Type         StepType
	Direction    Direction
	SwapAmount   SwapAmount
	StopReceived uint64
}

func StopFromPlutusData(plutusData *PlutusData.PlutusData) (Stop, error) {
	var err error
	data := plutusData.Value.(PlutusData.PlutusIndefArray)
	var stop Stop

	stop.Type = StepType_Stop
	stop.Direction, err = DirectionFromPlutusData(&data[0])
	if err != nil {
		return stop, err
	}

	stop.SwapAmount, err = SwapAmountFromPlutusData(&data[1])
	if err != nil {
		return stop, err
	}

	stopReceived, ok := data[2].Value.(uint64)
	if !ok {
		return stop, fmt.Errorf("invalid StopFromPlutusData")
	}
	stop.StopReceived = stopReceived

	return stop, nil
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
	Type            StepType
	Direction       Direction
	SwapAmount      SwapAmount
	MinimumReceived uint64
	StopReceived    uint64
}

func OCOFromPlutusData(plutusData *PlutusData.PlutusData) (OCO, error) {
	var err error
	data := plutusData.Value.(PlutusData.PlutusIndefArray)
	var oco OCO

	oco.Type = StepType_OCO
	oco.Direction, err = DirectionFromPlutusData(&data[0])
	if err != nil {
		return oco, err
	}

	oco.SwapAmount, err = SwapAmountFromPlutusData(&data[1])
	if err != nil {
		return oco, err
	}

	minimumReceived, ok := data[2].Value.(uint64)
	if !ok {
		return oco, fmt.Errorf("invalid OCOFromPlutusData")
	}
	oco.MinimumReceived = minimumReceived

	stopReceived, ok := data[3].Value.(uint64)
	if !ok {
		return oco, fmt.Errorf("invalid OCOFromPlutusData")
	}
	oco.StopReceived = stopReceived

	return oco, nil
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
	Type              StepType
	Direction         Direction
	MaximumSwapAmount SwapAmount
	ExpectedReceived  uint64
	Killable          Killable
}

func SwapExactOutFromPlutusData(plutusData *PlutusData.PlutusData) (SwapExactOut, error) {
	var err error
	data := plutusData.Value.(PlutusData.PlutusIndefArray)
	var swapExactOut SwapExactOut

	swapExactOut.Type = StepType_Swap_Exact_Out
	swapExactOut.Direction, err = DirectionFromPlutusData(&data[0])
	if err != nil {
		return swapExactOut, err
	}

	swapExactOut.MaximumSwapAmount, err = SwapAmountFromPlutusData(&data[1])
	if err != nil {
		return swapExactOut, err
	}

	expectedReceived, ok := data[2].Value.(uint64)
	if !ok {
		return swapExactOut, fmt.Errorf("invalid SwapExactOutFromPlutusData")
	}
	swapExactOut.ExpectedReceived = expectedReceived

	killable, err := KillableFromPlutusData(&data[3])
	if err != nil {
		return swapExactOut, err
	}
	swapExactOut.Killable = killable

	return swapExactOut, nil
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
	Type          StepType
	DepositAmount DepositAmount
	MinimumLP     uint64
	Killable      Killable
}

func DepositFromPlutusData(plutusData *PlutusData.PlutusData) (Deposit, error) {
	var err error
	data := plutusData.Value.(PlutusData.PlutusIndefArray)
	var deposit Deposit

	deposit.Type = StepType_Deposit
	deposit.DepositAmount, err = DepositAmountFromPlutusData(&data[0])
	if err != nil {
		return deposit, err
	}

	minimumLP, ok := data[1].Value.(uint64)
	if !ok {
		return deposit, fmt.Errorf("invalid DepositFromPlutusData")
	}
	deposit.MinimumLP = minimumLP

	killable, err := KillableFromPlutusData(&data[2])
	if err != nil {
		return deposit, err
	}
	deposit.Killable = killable

	return deposit, nil
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
	Type             StepType
	WithdrawalAmount WithdrawalAmount
	MinimumAssetA    uint64
	MinimumAssetB    uint64
	Killable         Killable
}

func WithdrawFromPlutusData(plutusData *PlutusData.PlutusData) (Withdraw, error) {
	var err error
	data := plutusData.Value.(PlutusData.PlutusIndefArray)
	var withdraw Withdraw

	withdraw.Type = StepType_Withdraw
	withdraw.WithdrawalAmount, err = WithdrawAmountFromPlutusData(&data[0])
	if err != nil {
		return withdraw, err
	}

	minimumAssetA, ok := data[1].Value.(uint64)
	if !ok {
		return withdraw, fmt.Errorf("invalid WithdrawFromPlutusData")
	}
	withdraw.MinimumAssetA = minimumAssetA

	minimumAssetB, ok := data[2].Value.(uint64)
	if !ok {
		return withdraw, fmt.Errorf("invalid WithdrawFromPlutusData")
	}
	withdraw.MinimumAssetB = minimumAssetB

	killable, err := KillableFromPlutusData(&data[3])
	if err != nil {
		return withdraw, err
	}
	withdraw.Killable = killable

	return withdraw, nil
}

func (s Withdraw) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Withdraw),
		PlutusDataType: PlutusData.PlutusArray,
		Value: PlutusData.PlutusIndefArray{
			s.WithdrawalAmount.ToPlutusData(),
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
	Type             StepType
	Direction        Direction
	WithdrawalAmount WithdrawalAmount
	MinimumReceived  uint64
	Killable         Killable
}

func ZapOutFromPlutusData(plutusData *PlutusData.PlutusData) (ZapOut, error) {
	var err error
	data := plutusData.Value.(PlutusData.PlutusIndefArray)
	var zapOut ZapOut

	zapOut.Type = StepType_Zap_Out
	zapOut.Direction, err = DirectionFromPlutusData(&data[0])
	if err != nil {
		return zapOut, err
	}

	zapOut.WithdrawalAmount, err = WithdrawAmountFromPlutusData(&data[1])
	if err != nil {
		return zapOut, err
	}

	minimumReceived, ok := data[2].Value.(uint64)
	if !ok {
		return zapOut, fmt.Errorf("invalid ZapOutFromPlutusData")
	}
	zapOut.MinimumReceived = minimumReceived

	killable, err := KillableFromPlutusData(&data[3])
	if err != nil {
		return zapOut, err
	}
	zapOut.Killable = killable

	return zapOut, nil
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
	Type                      StepType
	Direction                 Direction
	TotalSwapAmount           uint64
	IoRatioNumerator          uint64
	IoRatioDenominator        uint64
	Hops                      uint64
	MinimumSwapAmountRequired uint64
	MaxBatcherFeeEachTime     uint64
}

func PartialSwapFromPlutusData(plutusData *PlutusData.PlutusData) (PartialSwap, error) {
	var err error
	data := plutusData.Value.(PlutusData.PlutusIndefArray)
	var partialSwap PartialSwap

	partialSwap.Type = StepType_Partial_Swap
	partialSwap.Direction, err = DirectionFromPlutusData(&data[0])
	if err != nil {
		return partialSwap, err
	}

	totalSwapAmount, ok := data[1].Value.(uint64)
	if !ok {
		return partialSwap, fmt.Errorf("invalid PartialSwapFromPlutusData")
	}
	partialSwap.TotalSwapAmount = totalSwapAmount

	ioRatioNumerator, ok := data[2].Value.(uint64)
	if !ok {
		return partialSwap, fmt.Errorf("invalid PartialSwapFromPlutusData")
	}
	partialSwap.IoRatioNumerator = ioRatioNumerator

	ioRatioDenominator, ok := data[3].Value.(uint64)
	if !ok {
		return partialSwap, fmt.Errorf("invalid PartialSwapFromPlutusData")
	}
	partialSwap.IoRatioDenominator = ioRatioDenominator

	hops, ok := data[4].Value.(uint64)
	if !ok {
		return partialSwap, fmt.Errorf("invalid PartialSwapFromPlutusData")
	}
	partialSwap.Hops = hops

	minimumSwapAmountRequired, ok := data[5].Value.(uint64)
	if !ok {
		return partialSwap, fmt.Errorf("invalid PartialSwapFromPlutusData")
	}
	partialSwap.MinimumSwapAmountRequired = minimumSwapAmountRequired

	maxBatcherFeeEachTime, ok := data[6].Value.(uint64)
	if !ok {
		return partialSwap, fmt.Errorf("invalid PartialSwapFromPlutusData")
	}
	partialSwap.MaxBatcherFeeEachTime = maxBatcherFeeEachTime

	return partialSwap, nil
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
	Type           StepType
	WithdrawAmount WithdrawalAmount
	RatioAssetA    uint64
	RatioAssetB    uint64
	MinimumAssetA  uint64
	Killable       Killable
}

func WithdrawImbalanceFromPlutusData(plutusData *PlutusData.PlutusData) (WithdrawImbalance, error) {
	var err error
	data := plutusData.Value.(PlutusData.PlutusIndefArray)
	var withdrawImbalance WithdrawImbalance

	withdrawImbalance.Type = StepType_Withdraw_Imbalance
	withdrawImbalance.WithdrawAmount, err = WithdrawAmountFromPlutusData(&data[0])
	if err != nil {
		return withdrawImbalance, err
	}

	ratioAssetA, ok := data[1].Value.(uint64)
	if !ok {
		return withdrawImbalance, fmt.Errorf("invalid WithdrawImbalanceFromPlutusData")
	}
	withdrawImbalance.RatioAssetA = ratioAssetA

	ratioAssetB, ok := data[2].Value.(uint64)
	if !ok {
		return withdrawImbalance, fmt.Errorf("invalid WithdrawImbalanceFromPlutusData")
	}
	withdrawImbalance.RatioAssetB = ratioAssetB

	minimumAssetA, ok := data[3].Value.(uint64)
	if !ok {
		return withdrawImbalance, fmt.Errorf("invalid WithdrawImbalanceFromPlutusData")
	}
	withdrawImbalance.MinimumAssetA = minimumAssetA

	killable, err := KillableFromPlutusData(&data[4])
	if err != nil {
		return withdrawImbalance, err
	}
	withdrawImbalance.Killable = killable

	return withdrawImbalance, nil
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
	Type            StepType
	Routings        []Route
	SwapAmount      SwapAmount
	MinimumReceived uint64
}

func SwapRoutingFromPlutusData(plutusData *PlutusData.PlutusData) (SwapRouting, error) {
	var err error
	data := plutusData.Value.(PlutusData.PlutusIndefArray)
	var swapRouting SwapRouting

	swapRouting.Type = StepType_Swap_Routing

	routingData := data[0].Value.(PlutusData.PlutusIndefArray)
	for i, r := range routingData {
		route, err := RouteFromPlutusData(&r)
		if err != nil {
			return swapRouting, err
		}
		swapRouting.Routings[i] = route
	}

	swapRouting.SwapAmount, err = SwapAmountFromPlutusData(&data[1])
	if err != nil {
		return swapRouting, err
	}

	minimumReceived, ok := data[2].Value.(uint64)
	if !ok {
		return swapRouting, fmt.Errorf("invalid SwapRoutingFromPlutusData")
	}
	swapRouting.MinimumReceived = minimumReceived

	return swapRouting, nil
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

type Donation struct {
	Type StepType
}

func DonationFromPlutusData(plutusData *PlutusData.PlutusData) (Donation, error) {
	var donation Donation
	donation.Type = StepType(plutusData.TagNr - 121)
	return donation, nil
}

func (d Donation) StepToPlutusData() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr:          121 + uint64(StepType_Donation),
		PlutusDataType: PlutusData.PlutusArray,
		Value:          PlutusData.PlutusDefArray{},
	}
}

type ExpirySetting struct {
	ExpiredTime        uint64
	MaxCancellationTip uint64
}

func ExpirySettingFromPlutusData(plutusData *PlutusData.PlutusData) (ExpirySetting, error) {
	var expirySetting ExpirySetting
	data := plutusData.Value.(PlutusData.PlutusDefArray)

	if plutusData.TagNr == 121 {
		expirySetting.ExpiredTime = data[0].Value.(uint64)
		expirySetting.MaxCancellationTip = data[1].Value.(uint64)
	}

	return expirySetting, nil
}

func (e ExpirySetting) ToPlutusData() PlutusData.PlutusData {
	if e.ExpiredTime != 0 || e.MaxCancellationTip != 0 {
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
	MaxBatcherFee        uint64
	ExpiredOptions       ExpirySetting
}

func OrderDatumFromPlutusData(plutusData *PlutusData.PlutusData, networkId c.Network) (OrderDatum, error) {
	var orderDatum OrderDatum
	data := plutusData.Value.(PlutusData.PlutusIndefArray)

	canceller, err := AuthorizationMethodFromPlutusData(&data[0])
	if err != nil {
		return orderDatum, err
	}
	orderDatum.Canceller = canceller

	refundReceiver := plutusencoder.DecodePlutusAddress(data[1], byte(networkId))
	orderDatum.RefundReceiver = refundReceiver

	refundReceiverDatum, err := ExtraDatumFromPlutusData(&data[2])
	if err != nil {
		return orderDatum, err
	}
	orderDatum.RefundReceiverDatum = refundReceiverDatum

	successReceiver := plutusencoder.DecodePlutusAddress(data[3], byte(networkId))
	orderDatum.SuccessReceiver = successReceiver

	successReceiverDatum, err := ExtraDatumFromPlutusData(&data[4])
	if err != nil {
		return orderDatum, err
	}
	orderDatum.SuccessReceiverDatum = successReceiverDatum

	lpAsset, err := utils.FingerprintFromPlutusData(&data[5])
	if err != nil {
		return orderDatum, err
	}
	orderDatum.LpAsset = lpAsset

	orderDatum.Step, err = StepFromPlutusData(&data[6])
	if err != nil {
		return orderDatum, err
	}

	orderDatum.MaxBatcherFee = data[7].Value.(uint64)
	orderDatum.ExpiredOptions, err = ExpirySettingFromPlutusData(&data[8])
	if err != nil {
		return orderDatum, err
	}

	return orderDatum, nil
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
	OrderRedeemer_ApplyOrder = Redeemer.Redeemer{
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
