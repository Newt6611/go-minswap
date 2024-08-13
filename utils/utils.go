package utils

import (
	"encoding/hex"
	"math/big"

	"github.com/Newt6611/apollo/serialization/AssetName"
	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/Policy"
	"golang.org/x/crypto/sha3"
)

var (
	ADA = Fingerprint.Fingerprint {
		PolicyId: Policy.PolicyId{ Value: "" },
		AssetName: AssetName.AssetName{ Value: "" },
	}
	MIN = Fingerprint.Fingerprint {
		PolicyId: Policy.PolicyId{ Value: "e16c2dc8ae937e8d3790c7fd7168d7b994621ba14ca11415f39fed72" },
		AssetName: AssetName.AssetName{ Value: "4d494e" },
	}
)

func Sha3(hexString string) (string, error) {
	data, err := hex.DecodeString(hexString)
	if err != nil {
		return "", err
	}

	hash := sha3.New256()
	hash.Write(data)
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString, nil
}

type SlippageType int
const (
	SlippageTypeDown SlippageType = 0
	SlippageTypeUp   SlippageType = 1
)

func ApplySlippage(slippage float64, amount *big.Int, slippageType SlippageType) *big.Int {
	if slippageType == SlippageTypeUp {
		amountF, _ := amount.Float64()
		slippageAdjustedAmount := amountF * (1 + slippage)
		return big.NewInt(int64(slippageAdjustedAmount))

	} else {
		amountF, _ := amount.Float64()
		slippageAdjustedAmount := 1 / (1 + slippage)
		slippageAdjustedAmount = slippageAdjustedAmount * amountF
		return big.NewInt(int64(slippageAdjustedAmount))
	}
}
