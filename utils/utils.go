package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/Newt6611/apollo/serialization/Address"
	"github.com/Newt6611/apollo/serialization/AssetName"
	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/PlutusData"
	"github.com/Newt6611/apollo/serialization/Policy"
	"golang.org/x/crypto/sha3"
)

var (
	ADA = Fingerprint.Fingerprint{
		PolicyId:  Policy.PolicyId{Value: ""},
		AssetName: AssetName.AssetName{Value: ""},
	}
	MIN = Fingerprint.Fingerprint{
		PolicyId:  Policy.PolicyId{Value: "e16c2dc8ae937e8d3790c7fd7168d7b994621ba14ca11415f39fed72"},
		AssetName: AssetName.AssetName{Value: "4d494e"},
	}
)

type CredentialType int

const (
	CredentialTypeKey    CredentialType = 0
	CredentialTypeScript CredentialType = 1
)

type Credential struct {
	Type CredentialType
	Hash []byte
}

func CredentialFromPlutusData(plutusData *PlutusData.PlutusData) (Credential, error) {
	var credential Credential
	data := plutusData.Value.(PlutusData.PlutusIndefArray)[0]
	switch data.TagNr {
	case 121:
		credential.Type = CredentialTypeKey
		credential.Hash = data.Value.(PlutusData.PlutusIndefArray)[0].Value.([]byte)
	case 121 + 1:
		credential.Type = CredentialTypeScript
		credential.Hash = data.Value.(PlutusData.PlutusIndefArray)[0].Value.([]byte)
	default:
		return credential, fmt.Errorf("invalid CredentialFromPlutusData")
	}
	return credential, nil
}

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

func IsScriptAddress(address Address.Address) bool {
	if address.AddressType != Address.SCRIPT_KEY &&
		address.AddressType != Address.SCRIPT_SCRIPT &&
		address.AddressType != Address.SCRIPT_POINTER &&
		address.AddressType != Address.SCRIPT_NONE {
		return false
	}
	return true
}

func FingerprintFromPlutusData(plutusData *PlutusData.PlutusData) (Fingerprint.Fingerprint, error) {
	var fingerprint Fingerprint.Fingerprint
	policyHash, ok := plutusData.Value.(PlutusData.PlutusIndefArray)[0].Value.([]byte)
	if !ok {
		return fingerprint, fmt.Errorf("invalid FingerprintFromPlutusData Policy\n")
	}

	assetNameHash, ok := plutusData.Value.(PlutusData.PlutusIndefArray)[1].Value.([]byte)
	if !ok {
		return fingerprint, fmt.Errorf("invalid FingerprintFromPlutusData AssetName\n")
	}

	fingerprint.PolicyId.Value = hex.EncodeToString(policyHash)
	fingerprint.AssetName.Value = hex.EncodeToString(assetNameHash)
	return fingerprint, nil
}
