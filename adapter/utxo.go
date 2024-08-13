package adapter

type UTxO struct {
	Address             string
	TxHash              string
	OutputIndex         int
	Amount              []Amount
	Block               string
	DataHash            *string
	InlineDatum         *string
	ReferenceScriptHash *string
}
