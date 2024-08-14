package utils

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/PlutusData"
)

type FeeSharingOpt struct {
	Numerator *big.Int
	Enable    bool
}

type V2PoolState struct {
	Datum string
	// pool_batching_stake_credential: StakeCredential,
	PoolBatchingStakeCredential Credential
	// The Pool's Asset A
	AssetA Fingerprint.Fingerprint
	// The Pool's Asset B
	AssetB Fingerprint.Fingerprint
	// Total Share of Liquidity Providers
	TotalLiquidity *big.Int
	// Asset A's balance of Liquidity Providers
	ReserveA *big.Int
	// Asset B's balance of Liquidity Providers
	ReserveB *big.Int
	// Numerator of Trading Fee on Asset A side
	BaseFeeANumerator *big.Int
	// Numerator of Trading Fee on Asset B side
	BaseFeeBNumerator *big.Int
	// (Optional) Numerator of Fee Sharing percentage.
	// This is the percentage of Trading Fee. (eg, Trading Fee is 3%, Profit Sharing is 1/6 -> Profit Sharing = 1/6 * 3%)
	FeeSharingNumeratorOpt FeeSharingOpt
	// Allow Batcher can decide volatility fee for each batch transaction
	AllowDynamicFee bool
}

func ConvertToV2PoolState(plutusData PlutusData.PlutusData) (V2PoolState, error) {
	poolState := V2PoolState{}
	var err error
	data := plutusData.Value.(PlutusData.PlutusIndefArray)
	if plutusData.TagNr != 121 {
		errStr := fmt.Sprintf("index of pool datum must be 0, actual: %d", plutusData.TagNr)
		return V2PoolState{}, errors.New(errStr)
	}

	poolState.PoolBatchingStakeCredential, err = CredentialFromPlutusData(&data[0])
	if err != nil {
		return V2PoolState{}, err
	}

	poolState.AssetA, err = FingerprintFromPlutusData(&data[1])
	if err != nil {
		return V2PoolState{}, err
	}
	poolState.AssetB, err = FingerprintFromPlutusData(&data[2])
	if err != nil {
		return V2PoolState{}, err
	}

	poolState.TotalLiquidity = big.NewInt(int64(data[3].Value.(uint64)))
	poolState.ReserveA = big.NewInt(int64(data[4].Value.(uint64)))
	poolState.ReserveB = big.NewInt(int64(data[5].Value.(uint64)))
	poolState.BaseFeeANumerator = big.NewInt(int64(data[6].Value.(uint64)))
	poolState.BaseFeeBNumerator = big.NewInt(int64(data[7].Value.(uint64)))

	if data[8].TagNr == 121 {
		poolState.FeeSharingNumeratorOpt.Enable = true
		v := data[8].Value.(PlutusData.PlutusIndefArray)[0].Value.(uint64)
		poolState.FeeSharingNumeratorOpt.Numerator = big.NewInt(int64(v))
	} else {
		poolState.FeeSharingNumeratorOpt.Enable = false
	}

	poolState.AllowDynamicFee = data[9].TagNr == 122
	return poolState, nil
}
