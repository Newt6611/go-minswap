package adapter

import (
	"context"

	"github.com/Newt6611/apollo"
	c "github.com/Newt6611/apollo/constants"
	"github.com/Newt6611/apollo/serialization/Fingerprint"
	"github.com/Newt6611/apollo/serialization/UTxO"
	"github.com/Newt6611/apollo/txBuilding/Backend/Base"
	"github.com/Newt6611/go-minswap/utils"
)

type QueryParams struct {
	Count int
	Page  int
	Order string
	From  string
	To    string
}

type Amount struct {
	Unit     string
	Quantity string
}

type Adapter interface {
	NetworkId() c.Network
	ChainContext() Base.ChainContext
	NewBuilder() *apollo.Apollo
	GetV2PoolAll(ctx context.Context) ([]utils.V2PoolState, []error)
	GetV2Pool(ctx context.Context, params QueryParams) ([]utils.V2PoolState, []error)
	GetV2PoolByPair(ctx context.Context, assetA Fingerprint.Fingerprint, assetB Fingerprint.Fingerprint) (utils.V2PoolState, error)
	GetDatumByDatumHash(ctx context.Context, datumHash string) (string, error)
	GetUtxoFromRef(ctx context.Context, txhash string, index int) *UTxO.UTxO
}
