package adapter

import "github.com/Newt6611/apollo/serialization/Fingerprint"

func normalizeAssets(assetA Fingerprint.Fingerprint, assetB Fingerprint.Fingerprint) (Fingerprint.Fingerprint, Fingerprint.Fingerprint){
	if assetA.String() == "lovelace" {
		return assetA, assetB
	}
	if assetB.String() == "lovelace" {
		return assetB, assetA
	}

	assetALen := len(assetA.PolicyId.Value) + len(assetA.AssetName.Value)
	assetBLen := len(assetB.PolicyId.Value) + len(assetB.AssetName.Value)

	if assetALen < assetBLen {
		return assetA, assetB
	} else {
		return assetB, assetA
	}
}
