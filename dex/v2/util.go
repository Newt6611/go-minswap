package v2

import (
	"encoding/hex"
	"errors"
	"math/big"

	c "github.com/Newt6611/apollo/constants"
	"github.com/Newt6611/apollo/serialization/Address"
	"github.com/Newt6611/apollo/serialization/AssetName"
	"github.com/Newt6611/go-minswap/constants"
	"github.com/Newt6611/go-minswap/utils"
)

func BuildOrderAddress(senderAddr Address.Address, network c.Network) Address.Address {
	senderStakePart := senderAddr.StakingPart
	orderAddrHex := constants.V2Config[network].OrderEnterpriseAddress
	orderAddr, _ := Address.DecodeAddress(orderAddrHex)

	apolloNetwork := c.MAINNET
	networkByte := 0b0001
	if network != c.MAINNET {
		apolloNetwork = c.TESTNET
		networkByte = 0b0000
	}

	result := Address.WalletAddressFromBytes(orderAddr.PaymentPart, senderStakePart, apolloNetwork)
	result.HeaderByte = Address.SCRIPT_KEY<<4 | byte(networkByte)
	return *result
}

func GetOrderScriptHash(networkId c.Network) (string, error) {
	orderAddress := constants.V2Config[networkId].OrderEnterpriseAddress
	orderAddr, err := Address.DecodeAddress(orderAddress)
	if err != nil {
		return "", err
	}
	if !utils.IsScriptAddress(orderAddr) {
		return "", errors.New("Order address should be a script address")
	}

	return hex.EncodeToString(orderAddr.PaymentPart), nil
}

func ComputeLPAsset(assetAPolicyId, assetAName, assetBPolicyId, assetBName string) (AssetName.AssetName, error) {
	k1, err := utils.Sha3(assetAPolicyId + assetAName)
	if err != nil {
		return AssetName.AssetName{}, err
	}
	k2, err := utils.Sha3(assetBPolicyId + assetBName)
	if err != nil {
		return AssetName.AssetName{}, err
	}

	hex, err := utils.Sha3(k1 + k2)
	if err != nil {
		return AssetName.AssetName{}, err
	}

	result := AssetName.NewAssetNameFromHexString(hex)

	return *result, nil
}

/*
Functions using for DexV2 properties calculation
pub fn calculate_amount_out(

	reserve_in: Int,
	reserve_out: Int,
	amount_in: Int,
	trading_fee_numerator: Int,
	) -> Int {
	  let diff = utils.default_fee_denominator - trading_fee_numerator
	  let in_with_fee = diff * amount_in
	  let numerator = in_with_fee * reserve_out
	  let denominator = utils.default_fee_denominator * reserve_in + in_with_fee
	  numerator / denominator
	}
*/
func CalculateAmountOut(reserveIn, reserveOut, amountIn, tradingFeeNumerator uint64) uint64 {
	diff := big.NewInt(0).
		Sub(big.NewInt(int64(utils.DEFAULT_TRADING_FEE_DENOMINATOR)),
			big.NewInt(int64(tradingFeeNumerator)))

	inWithFee := big.NewInt(0).Mul(diff, big.NewInt(int64(amountIn)))

	numerator := big.NewInt(0).Mul(inWithFee, big.NewInt(int64(reserveOut)))

	denominator := big.NewInt(0).Mul(big.NewInt(int64(utils.DEFAULT_TRADING_FEE_DENOMINATOR)), big.NewInt(int64(reserveIn)))
	denominator = big.NewInt(0).Add(denominator, inWithFee)

	return numerator.Div(numerator, denominator).Uint64()
}

func getTagNr(tag uint64) uint64 {
	if tag < 7 {
		return 121 + tag
	} else {
		return 1280 + tag - 7
	}
}
