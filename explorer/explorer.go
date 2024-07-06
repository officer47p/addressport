package explorer

import (
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/officer47p/addressport/types"
)

type Explorer interface {
	GetAllTransactionsForAddress(string) ([]types.Transaction, error)
}

type etherscanTransaction struct {
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	ContractAddress string `json:"contractAddress"`
}

type etherscanGetAllTransactionsReponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Result  []etherscanTransaction `json:"result"`
}

type EtherscanExplorer struct {
	apiKey string
}

func NewEtherscanExplorer(apiKey string) *EtherscanExplorer {
	return &EtherscanExplorer{apiKey: apiKey}
}

func (ex *EtherscanExplorer) GetAllTransactionsForAddress(address string) ([]types.Transaction, error) {
	// Create a Resty Client
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"module":     "account",
			"action":     "txlist",
			"address":    address,
			"startblock": "0",
			"endblock":   "99999999",
			"page":       "1",
			"offset":     "10",
			"sort":       "asc",
			"apikey":     ex.apiKey,
		}).
		SetHeader("Accept", "application/json").
		// SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
		Get("https://api.etherscan.io/api")

	if err != nil {
		log.Printf("explorer: error while fetching transactions of an address. err: %e\n", err)
		return []types.Transaction{}, err
	}

	var result etherscanGetAllTransactionsReponse
	err = json.Unmarshal(resp.Body(), &result)

	if err != nil {
		log.Printf("explorer: error while fetching transactions of an address. err: %e\n response: %s\n", err, resp.Body())
		return []types.Transaction{}, err
	}

	transactions := []types.Transaction{}
	for _, tx := range result.Result {
		transactions = append(transactions, types.Transaction{From: tx.From, To: tx.To, TxHash: tx.Hash})
	}

	return transactions, nil
}
