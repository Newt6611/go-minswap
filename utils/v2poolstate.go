package utils

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Newt6611/apollo/serialization/AssetName"
	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/Policy"
	"github.com/blinklabs-io/gouroboros/cbor"
)

type FeeSharingOpt struct {
	Numerator uint64
	Enable    bool
}

type V2PoolState struct {
	Datum string
	// pool_batching_stake_credential: StakeCredential,
	// The Pool's Asset A
	AssetA *Fingerprint.Fingerprint
	// The Pool's Asset B
	AssetB *Fingerprint.Fingerprint
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
	// AllowDynamicFee bool
}


func ConvertToV2PoolState(data cbor.Constructor) (V2PoolState, error) {
	poolState := V2PoolState{}

	if data.Constructor() != 0 {
		errStr := fmt.Sprintf("index of pool datum must be 0, actual: %d", data.Constructor())
		return V2PoolState{}, errors.New(errStr)
	}

	fields := data.Fields()
	poolState.AssetA = convertToAsset(fields[1].(cbor.Constructor))
	poolState.AssetB = convertToAsset(fields[2].(cbor.Constructor))
	poolState.TotalLiquidity = big.NewInt(int64(fields[3].(uint64)))
	poolState.ReserveA = big.NewInt(int64(fields[4].(uint64)))
	poolState.ReserveB = big.NewInt(int64(fields[5].(uint64)))
	poolState.BaseFeeANumerator = big.NewInt(int64(fields[6].(uint64)))
	poolState.BaseFeeBNumerator = big.NewInt(int64(fields[7].(uint64)))
	poolState.FeeSharingNumeratorOpt = convertToFeeSharingNumerator(fields[8].(cbor.Constructor))

	return poolState, nil
}

func convertToAsset(data cbor.Constructor) *Fingerprint.Fingerprint {
	fields := data.Fields()
	b_policyId := fields[0].(cbor.ByteString)
	b_assetName := fields[1].(cbor.ByteString)
	policyId := &Policy.PolicyId { Value: b_policyId.String() }
	assetName := AssetName.NewAssetNameFromHexString(b_assetName.String())

	return Fingerprint.New(*policyId, *assetName)
}

func convertToFeeSharingNumerator(data cbor.Constructor) FeeSharingOpt {
	if data.Constructor() != 0 {
		return FeeSharingOpt{
			Enable: false,
		}
	}

	numerator := data.Fields()[0].(uint64)
	return FeeSharingOpt{
		Enable:    true,
		Numerator: numerator,
	}
}
