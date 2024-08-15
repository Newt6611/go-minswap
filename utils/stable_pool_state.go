package utils

import (
	"github.com/Newt6611/apollo/serialization/PlutusData"
)

type StablePoolState struct {
	Balances       []uint64
	TotalLiquidity uint64
	AMP            uint64
	OrderHash      []byte
}

func ConvertToStablePoolState(plutusData PlutusData.PlutusData) (StablePoolState, error) {
	var stablePoolState StablePoolState
	_ = stablePoolState
	data := plutusData.Value.(PlutusData.PlutusIndefArray)

	balanceArr := data[0].Value.(PlutusData.PlutusIndefArray)
	for _, balance := range balanceArr {
		i := balance.Value.(uint64)
		stablePoolState.Balances = append(stablePoolState.Balances, i)
	}

	stablePoolState.TotalLiquidity = data[1].Value.(uint64)
	stablePoolState.AMP = data[2].Value.(uint64)
	stablePoolState.OrderHash = data[3].Value.([]byte)

	return stablePoolState, nil
}
