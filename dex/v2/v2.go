package v2

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/Newt6611/apollo"
	"github.com/Newt6611/apollo/serialization/Address"
	apollo_c "github.com/Newt6611/apollo/constants"
	"github.com/Newt6611/apollo/serialization/Metadata"
	"github.com/Newt6611/apollo/serialization/PlutusData"
	"github.com/Newt6611/apollo/serialization/UTxO"
	"github.com/Newt6611/go-minswap/adapter"
	"github.com/Newt6611/go-minswap/constants"
	"github.com/Newt6611/go-minswap/utils"
	"github.com/Salvionied/cbor/v2"
)

type MetadataMessage string

const (
	MetadataMessage_DEPOSIT_ORDER             MetadataMessage = "go-minswap: Deposit Order"
	MetadataMessage_CANCEL_ORDER              MetadataMessage = "go-minswap: Cancel Order"
	MetadataMessage_ZAP_IN_ORDER              MetadataMessage = "go-minswap: Zap Order"
	MetadataMessage_ZAP_OUT_ORDER             MetadataMessage = "go-minswap: Zap Out Order"
	MetadataMessage_SWAP_EXACT_IN_ORDER       MetadataMessage = "go-minswap: Swap Exact In Order"
	MetadataMessage_SWAP_EXACT_IN_LIMIT_ORDER MetadataMessage = "go-minswap: Swap Exact In Limit Order"
	MetadataMessage_SWAP_EXACT_OUT_ORDER      MetadataMessage = "go-minswap: Swap Exact Out Order"
	MetadataMessage_WITHDRAW_ORDER            MetadataMessage = "go-minswap: Withdraw Order"
	MetadataMessage_STOP_ORDER                MetadataMessage = "go-minswap: Stop Order"
	MetadataMessage_OCO_ORDER                 MetadataMessage = "go-minswap: OCO Order"
	MetadataMessage_ROUTING_ORDER             MetadataMessage = "go-minswap: Routing Order"
	MetadataMessage_PARTIAL_SWAP_ORDER        MetadataMessage = "go-minswap: Partial Fill Order"
	MetadataMessage_DONATION_ORDER            MetadataMessage = "go-minswap: Donation Order"
	MetadataMessage_MIXED_ORDERS              MetadataMessage = "go-minswap: Mixed Orders"
	MetadataMessage_CREATE_POOL               MetadataMessage = "go-minswap: Create Pool"
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
		return builder, errors.New("order Utxos are empty")
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
			orderDatum, err = OrderDatumFromPlutusData(datum, constants.NetworkIdTestnet)
			if err != nil {
				return builder ,err
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
			orderDatum, err = OrderDatumFromPlutusData(&p, constants.NetworkIdTestnet)
			if err != nil {
				return builder, err
			}
		} else {
			return builder, errors.New("utxo without Datum Hash or Inline Datum can not be spent")
		}

		if orderDatum.Canceller.Type != AuthorizationMethodType_Signature {
			return builder, errors.New("only support PubKey canceller on this function")
		}

		c := apollo_c.MAINNET
		if networkId != constants.NetworkIdMainnet {
			c = apollo_c.TESTNET
		}
		addr := Address.WalletAddressFromBytes(orderDatum.Canceller.Hash, nil, c)
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
						string(MetadataMessage_CANCEL_ORDER),
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
