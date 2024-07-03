package types

type Transaction struct {
	From string `json:"from"`
	To   string `json:"to"`
	Hash string `json:"hash"`
}
