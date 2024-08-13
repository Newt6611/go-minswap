package adapter

import (
	"context"

	"github.com/Newt6611/apollo/serialization/Fingerprint"
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
	GetV2PoolAll (ctx context.Context) ([]utils.V2PoolState, []error)
	GetV2Pool (ctx context.Context, params QueryParams) ([]utils.V2PoolState, []error)
	GetV2PoolByPair(ctx context.Context, assetA Fingerprint.Fingerprint, assetB Fingerprint.Fingerprint) (utils.V2PoolState, error)
}