# Golang Minswap SDK

Still in heavy development, please don't use it in production.

Get MIN/ADA pool example:
```go
ctx := context.Background()

blockfrostAdapter, err := adapter.NewBlockFrost(blockfrost.APIClientOptions{
	ProjectID: YOUR_BLOCKFROST_API_KEY_HERE,
	Server: blockfrost.CardanoPreProd,
})
if err != nil {
	log.Fatal(err)
}

pools, errs := blockfrostAdapter.GetV2Pool(ctx, adapter.QueryParams{ })
if len(errs) != 0 {
	log.Fatal(errs[0])
}

if len(pools) <= 0 {
	log.Fatal("can't find pools")
	return
}

// find MIN/ADA pool
var MinAdaPool utils.V2PoolState
for _, pool := range pools {
	poolAssetA := pool.AssetA.PolicyId.Value + pool.AssetA.AssetName.Value
	poolAssetB := pool.AssetB.PolicyId.Value + pool.AssetB.AssetName.Value

	if poolAssetA == utils.ADA.PolicyId.Value + utils.ADA.AssetName.Value &&
		poolAssetB == utils.MIN.PolicyId.Value + utils.MIN.AssetName.Value {
		MinAdaPool = pool
		break
	}

	if poolAssetB == utils.ADA.PolicyId.Value + utils.ADA.AssetName.Value &&
		poolAssetA == utils.MIN.PolicyId.Value + utils.MIN.AssetName.Value {
		MinAdaPool = pool
		break
	}
}

if len(MinAdaPool.Datum) == 0 {
	log.Fatal("can't find MIN/ADA pool")
	return
}
reserveA, _ := MinAdaPool.ReserveA.Float64()
reserveB, _ := MinAdaPool.ReserveB.Float64()
fmt.Println(reserveA / reserveB)
```


### TODO:
- [ ] V1
	- [ ] Order Creation
	- [ ] Order Cancellation
- [ ] V2
	- [x] Get Pool Data
	- [X] Order Creation
	- [X] Order Cancellation
	- [ ] Pool Creation
- [ ] StableSwap
	- [ ] Order Creation
	- [ ] Order Cancellation
