package thirdparty

import (
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
)

type Explorer interface {
	GetAllTransactionsForAddress(string) ([]Transaction, error)
}

type EtherscanExplorer struct {
	apiKey string
}

func NewEtherscanExplorer(apiKey string) *EtherscanExplorer {
	return &EtherscanExplorer{apiKey: apiKey}
}

func (ex *EtherscanExplorer) GetAllTransactionsForAddress(address string) ([]Transaction, error) {
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"module":     "account",
			"action":     "txlist",
			"address":    address,
			"startblock": "0",
			"endblock":   "99999999",
			"page":       "1",
			"offset":     "20", // can be maximum 10_000
			"sort":       "asc",
			"apikey":     ex.apiKey,
		}).
		SetHeader("Accept", "application/json").
		Get("https://api.etherscan.io/api")

	if err != nil {
		log.Printf("explorer: error while fetching transactions of an address. err: %e\n", err)
		return []Transaction{}, err
	}

	var result etherscanGetAllTransactionsReponse
	err = json.Unmarshal(resp.Body(), &result)

	if err != nil {
		log.Printf("explorer: error while parsing transactions of an address. err: %e\n response: %s\n", err, resp.Body())
		return []Transaction{}, err
	}

	transactions := []Transaction{}
	for _, tx := range result.Result {
		transactions = append(transactions, Transaction{From: tx.From, To: tx.To, TxHash: tx.Hash, Value: tx.Value})
	}

	return transactions, nil
}

type etherscanTransaction struct {
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	ContractAddress string `json:"contractAddress"`
}

type etherscanGetAllTransactionsReponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Result  []etherscanTransaction `json:"result"`
}
