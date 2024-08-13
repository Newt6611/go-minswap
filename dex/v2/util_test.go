package v2_test

import (
	"math/big"
	"testing"

	"github.com/Newt6611/apollo/serialization/Address"
	"github.com/Newt6611/apollo/serialization/AssetName"
	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/Policy"
	"github.com/Newt6611/go-minswap/constants"
	v2 "github.com/Newt6611/go-minswap/dex/v2"
	"github.com/Newt6611/go-minswap/utils"
)

func TestBuildOrderAddress(t *testing.T) {
	correctOrderAddr := "addr_test1zrdf2f2x8pq3wwk3yv936ksmt59rz94mm66yzge8zj9pk75dvyatxqgacw6azcshwwywv0nkxkdp2l2uq5qrn628mw0s7r0f20"
	userAddr, _ := Address.DecodeAddress("addr_test1qqpfnkrx0wv3zty7c3uw9j6tvvj2tvuh7spgw80v8c353yydvyatxqgacw6azcshwwywv0nkxkdp2l2uq5qrn628mw0ssws0re")
	orderAddr := v2.BuildOrderAddress(userAddr, constants.NetworkIdTestnet)
	if orderAddr.String() != correctOrderAddr {
		t.Errorf("BuildOrderAddress expect %s but get %s\n", correctOrderAddr, orderAddr.String())
	}
}

func TestComputeLPAsset(t *testing.T) {
	correctAssetName := "6c3ea488e6ff940bb6fb1b18fd605b5931d9fefde6440117015ba484cf321200"

	assetA := utils.ADA
	// testnet
	minPolicyId, _ := Policy.New("e16c2dc8ae937e8d3790c7fd7168d7b994621ba14ca11415f39fed72")
	minAssetName := AssetName.NewAssetNameFromHexString("4d494e")
	minFingerprint := Fingerprint.New(*minPolicyId, *minAssetName)

	result, err := v2.ComputeLPAsset(assetA.PolicyId.String(), assetA.AssetName.HexString(),
		minFingerprint.PolicyId.String(), minFingerprint.AssetName.HexString())
	if err != nil {
		t.Error(err)
	}
	if result.HexString() != correctAssetName {
		t.Errorf("ComputeLPAsset expect %s, but get %s\n", correctAssetName, result)
	}
}

func TestCalculateAmountOut(t *testing.T) {
	out := v2.CalculateAmountOut(
		big.NewInt(25434557769208),
		big.NewInt(413300185707175),
		big.NewInt(1000000000),
		big.NewInt(30))

	anwser := big.NewInt(16200168971)
	if anwser.Cmp(out) != 0 {
		t.Errorf("Expected %s, but get %s\n", anwser.String(), out.String())
	}
}
