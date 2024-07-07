package types

type Block struct {
	Network           string
	BlockNumber       int64
	BlockHash         string
	PreviousBlockHash string
	Transactions      []Transaction
}

type Transaction struct {
	BlockNumber int64
	BlockHash   string
	Network     string
	Currency    string
	TxHash      string
	Value       string
	From        string
	To          string
}

type Network struct {
	Name                string `json:"name"`
	Currency            string `json:"currency"`
	ChainID             int64  `json:"chainID"`
	Decimals            int64  `json:"decimals"`
	StartingBlockNumber int64  `json:"startingBlockNumber"`
}
