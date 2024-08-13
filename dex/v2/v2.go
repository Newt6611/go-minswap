package v2

import "github.com/Newt6611/go-minswap/adapter"

type MetadataMessage string

const (
	DEPOSIT_ORDER             MetadataMessage = "go-minswap: Deposit Order"
	CANCEL_ORDER              MetadataMessage = "go-minswap: Cancel Order"
	ZAP_IN_ORDER              MetadataMessage = "go-minswap: Zap Order"
	ZAP_OUT_ORDER             MetadataMessage = "go-minswap: Zap Out Order"
	SWAP_EXACT_IN_ORDER       MetadataMessage = "go-minswap: Swap Exact In Order"
	SWAP_EXACT_IN_LIMIT_ORDER MetadataMessage = "go-minswap: Swap Exact In Limit Order"
	SWAP_EXACT_OUT_ORDER      MetadataMessage = "go-minswap: Swap Exact Out Order"
	WITHDRAW_ORDER            MetadataMessage = "go-minswap: Withdraw Order"
	STOP_ORDER                MetadataMessage = "go-minswap: Stop Order"
	OCO_ORDER                 MetadataMessage = "go-minswap: OCO Order"
	ROUTING_ORDER             MetadataMessage = "go-minswap: Routing Order"
	PARTIAL_SWAP_ORDER        MetadataMessage = "go-minswap: Partial Fill Order"
	DONATION_ORDER            MetadataMessage = "go-minswap: Donation Order"
	MIXED_ORDERS              MetadataMessage = "go-minswap: Mixed Orders"
	CREATE_POOL               MetadataMessage = "go-minswap: Create Pool"
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
