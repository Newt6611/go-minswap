package utils

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

const (
	DEFAULT_TRADING_FEE_DENOMINATOR int64 = 10000
)
