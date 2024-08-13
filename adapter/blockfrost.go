package adapter

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/go-minswap/constants"
	"github.com/Newt6611/go-minswap/utils"
	"github.com/blinklabs-io/gouroboros/cbor"
	"github.com/blockfrost/blockfrost-go"
)

type BlockFrost struct {
	client    blockfrost.APIClient
	network   constants.NetworkId
	options   blockfrost.APIClientOptions
}

func NewBlockFrost(options blockfrost.APIClientOptions) *BlockFrost {
	client := blockfrost.NewAPIClient(options)
	network := constants.NetworkIdMainnet

	if options.Server == blockfrost.CardanoMainNet || options.Server == "" {
		network = constants.NetworkIdMainnet
		options.Server = blockfrost.CardanoMainNet
	} else {
		network = constants.NetworkIdTestnet
	}

	return &BlockFrost{
		client:    client,
		network:   network,
		options:   options,
	}
}

func (b *BlockFrost) GetV2PoolAll (ctx context.Context) ([]utils.V2PoolState, []error) {
	address := constants.V2Config[b.network].PoolScriptHashBech32
	asset := constants.V2Config[b.network].PoolAuthenAsset

	poolStates := []utils.V2PoolState{}
	errs := []error{}

	resultChan := b.client.AddressUTXOsAssetAll(ctx, address, asset)
	for {
		result, keep := <-resultChan
		if result.Err != nil {
			errs = append(errs, result.Err)
		}

		utxos := result.Res
		poolState, es := convertUtxosToPoolState(utxos, errs)
		poolStates = append(poolStates, poolState...)
		errs = append(errs, es...)

		if !keep {
			break
		}
	}
	return poolStates, errs
}

func (b *BlockFrost) GetV2Pool(ctx context.Context, params QueryParams) ([]utils.V2PoolState, []error) {
	errs := []error{}

	address := constants.V2Config[b.network].PoolScriptHashBech32
	asset := constants.V2Config[b.network].PoolAuthenAsset
	p := blockfrost.APIQueryParams{
		Count: params.Count,
		Page:  params.Page,
		Order: params.Order,
		From:  params.From,
		To:    params.To,
	}
	utxos, err := b.client.AddressUTXOsAsset(ctx, address, asset, p)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	poolStates, errs := convertUtxosToPoolState(utxos, errs)

	return poolStates, errs
}

func (b *BlockFrost) GetV2PoolByPair(ctx context.Context, assetA Fingerprint.Fingerprint, assetB Fingerprint.Fingerprint) (utils.V2PoolState, error) {
	normalizedAssetA, normalizedAssetB := normalizeAssets(assetA, assetB)
	pools, errs := b.GetV2PoolAll(ctx)
	if len(errs) != 0 {
		return utils.V2PoolState{}, errs[0]
	}

	for _, pool := range pools {
		if pool.AssetA.String() == normalizedAssetA.String() && 
			pool.AssetB.String() == normalizedAssetB.String() {
			return pool, nil
		}
	}

	return utils.V2PoolState{}, errors.New("pool not found")
}

func (b *BlockFrost) GetDatumByDatumHash(ctx context.Context, datumHash string) (string, error) {
	url := fmt.Sprintf("%s/scripts/datum/%s/cbor", b.options.Server, datumHash)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("project_id", b.options.ProjectID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	return data["cbor"], nil
}

func convertUtxosToPoolState(utxos []blockfrost.AddressUTXO, errs []error) ([]utils.V2PoolState, []error) {
	poolStates := []utils.V2PoolState{}
	for _, utxo := range utxos {
		if utxo.InlineDatum == nil {
			continue
		}

		decodedHex, _ := hex.DecodeString(*utxo.InlineDatum)
		var cborData cbor.Constructor
		_, err := cbor.Decode(decodedHex, &cborData)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		pool, err := utils.ConvertToV2PoolState(cborData)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		pool.Datum = *utxo.InlineDatum
		poolStates = append(poolStates, pool)
	}

	return poolStates, errs
}
