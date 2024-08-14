package v2

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/Newt6611/apollo"
	"github.com/Newt6611/apollo/serialization/Address"
	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/Metadata"
	"github.com/Newt6611/apollo/serialization/PlutusData"
	"github.com/Newt6611/apollo/serialization/Policy"
	"github.com/Newt6611/apollo/serialization/UTxO"
	"github.com/Newt6611/go-minswap/adapter"
	"github.com/Newt6611/go-minswap/constants"
	"github.com/Newt6611/go-minswap/utils"
	"github.com/Salvionied/cbor/v2"
)


type DexV2 struct {
	adapter adapter.Adapter
}

func NewDexV2(adapter adapter.Adapter) *DexV2 {
	return &DexV2{
		adapter: adapter,
	}
}
func (d *DexV2) CreateBulkOrdersTx() {

}

func (d *DexV2) BuildSwapExactInOrder(ctx context.Context,
	builder *apollo.Apollo,
	swapExactIn SwapExactIn,
	assetA Fingerprint.Fingerprint,
	assetB Fingerprint.Fingerprint,
	lovelace int,
	units ...apollo.Unit) (*apollo.Apollo, error) {

	networkId := d.adapter.NetworkId()
	builderAddr := builder.GetWallet().GetAddress()
	orderAddr := BuildOrderAddress(*builderAddr, networkId)

	assetName, err := ComputeLPAsset(assetA.PolicyId.Value, assetA.AssetName.Value,
		assetB.PolicyId.Value, assetB.AssetName.Value)
	if err != nil {
		return builder, err
	}
	lpPolicy, err := Policy.New(constants.V2Config[networkId].LpPolicyId)
	if err != nil {
		return builder, err
	}

	orderDatum := OrderDatum{
		Canceller: AuthorizationMethod{
			Type: AuthorizationMethodType_Signature,
			Hash: builderAddr.PaymentPart,
		},
		RefundReceiver: *builderAddr,
		RefundReceiverDatum: ExtraDatum{
			Type: ExtraDatumType_No_Datum,
		},
		SuccessReceiver: *builderAddr,
		LpAsset:         *Fingerprint.New(*lpPolicy, assetName),
		Step:            swapExactIn,
		MaxBatcherFee:   FIXED_BATCHER_FEE, //TODO: caculate batcher fee
		ExpiredOptions:  ExpirySetting{},
	}

	orderDatumPlutusData := orderDatum.ToPlutusData()
	builder, err = builder.
		SetWalletAsChangeAddress().
		PayToContract(orderAddr, &orderDatumPlutusData, lovelace, true, units...).
		SetShelleyMetadata(Metadata.ShelleyMaryMetadata{
			Metadata: Metadata.Metadata{
				674: struct {
					Msg []string `json:"msg"`
				}{
					Msg: []string{
						string(utils.MetadataMessage_SWAP_EXACT_IN_ORDER),
					},
				},
			},
		}).Complete()

	if err != nil {
		return builder, err
	}
	return builder, nil
}

func (d *DexV2) BuildCancelOrder(ctx context.Context, builder *apollo.Apollo, outRefs []constants.OutRef) (*apollo.Apollo, error) {
	networkId := d.adapter.NetworkId()
	v2OrderScriptHash, err := GetOrderScriptHash(networkId)
	if err != nil {
		return builder, err
	}

	var orderUtxos []*UTxO.UTxO
	for _, outRef := range outRefs {
		orderUtxo := d.adapter.GetUtxoFromRef(ctx, outRef.TxHash, outRef.Index)
		if orderUtxo != nil {
			orderUtxos = append(orderUtxos, orderUtxo)
		}
	}
	if len(orderUtxos) == 0 {
		return builder, errors.New("order utxos are empty")
	}

	depolyedOrderScript := constants.V2DeployedScripts[networkId].Order
	orderRef := d.adapter.GetUtxoFromRef(ctx, depolyedOrderScript.TxHash, depolyedOrderScript.Index)
	if orderRef == nil {
		return builder, errors.New("cannot find deployed script for V2 Order")
	}

	redeemer := OrderRedeemer_CancelOrderByOwner
	builder = builder.AddReferenceInput(depolyedOrderScript.TxHash, depolyedOrderScript.Index)

	for _, orderUtxo := range orderUtxos {
		orderUtxoAddr := orderUtxo.Output.GetAddress()
		if !utils.IsScriptAddress(orderUtxoAddr) {
			return builder, errors.New("utxo is not belonged Minswap's order address, utxo: " + orderUtxo.GetKey())
		}
		if hex.EncodeToString(orderUtxoAddr.PaymentPart) != v2OrderScriptHash {
			return builder, errors.New("utxo is not belonged Minswap's order address, utxo: " + orderUtxo.GetKey())
		}

		var orderDatum OrderDatum
		if datum := orderUtxo.Output.GetDatum(); datum != nil {
			orderDatum, err = OrderDatumFromPlutusData(datum, networkId)
			if err != nil {
				return builder, err
			}
		} else if dataHash := orderUtxo.Output.GetDatumHash(); dataHash != nil {
			rawdatum, err := d.adapter.GetDatumByDatumHash(ctx, hex.EncodeToString(dataHash.Payload))
			if err != nil {
				return builder, err
			}
			b, _ := hex.DecodeString(rawdatum)
			var p PlutusData.PlutusData
			err = cbor.Unmarshal(b, &p)
			if err != nil {
				return builder, err
			}
			orderDatum, err = OrderDatumFromPlutusData(&p, networkId)
			if err != nil {
				return builder, err
			}
		} else {
			return builder, errors.New("utxo without Datum Hash or Inline Datum can not be spent")
		}

		if orderDatum.Canceller.Type != AuthorizationMethodType_Signature {
			return builder, errors.New("only support PubKey canceller on this function")
		}

		addr := Address.WalletAddressFromBytes(orderDatum.Canceller.Hash, nil, networkId)
		builder = builder.CollectFrom(*orderUtxo, redeemer).
			AddRequiredSignerFromAddress(*addr, true, false)
	}

	builder = builder.SetWalletAsChangeAddress().
		SetShelleyMetadata(Metadata.ShelleyMaryMetadata{
			Metadata: Metadata.Metadata{
				674: struct {
					Msg []string `json:"msg"`
				}{
					Msg: []string{
						string(utils.MetadataMessage_CANCEL_ORDER),
					},
				},
			},
		})

	builder, err = builder.Complete()
	if err != nil {
		return builder, err
	}

	return builder, nil
}
