package v2_test

import (
	"testing"

	c "github.com/Newt6611/apollo/constants"
	"github.com/Newt6611/apollo/serialization/Address"
	"github.com/Newt6611/apollo/serialization/AssetName"
	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/Policy"
	v2 "github.com/Newt6611/go-minswap/dex/v2"
	"github.com/Newt6611/go-minswap/utils"
)

func TestBuildOrderAddress(t *testing.T) {
	correctOrderAddr := "addr_test1zrdf2f2x8pq3wwk3yv936ksmt59rz94mm66yzge8zj9pk75dvyatxqgacw6azcshwwywv0nkxkdp2l2uq5qrn628mw0s7r0f20"
	userAddr, _ := Address.DecodeAddress("addr_test1qqpfnkrx0wv3zty7c3uw9j6tvvj2tvuh7spgw80v8c353yydvyatxqgacw6azcshwwywv0nkxkdp2l2uq5qrn628mw0ssws0re")
	orderAddr := v2.BuildOrderAddress(userAddr, c.TESTNET)
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
		25434557769208,
		413300185707175,
		1000000000,
		30)

	anwser := uint64(16200168971)
	if out != anwser {
		t.Errorf("Expected %d, but get %d\n", anwser, out)
	}
}
