package services

import (
	"github.com/officer47p/addressport/lib/thirdparty"
)

func NewInvestigationToolService(explorer thirdparty.Explorer) InvestigationToolService {
	return InvestigationToolService{explorer: explorer}
}

type InvestigationToolService struct {
	explorer thirdparty.Explorer
}

func (i *InvestigationToolService) GetAssociatedAddressesForAddress(address string, depth int) (*AddressNode, error) {
	addressNode := AddressNode{Address: address, Children: []*AddressNode{}}
	if err := addressNode.PopulateNode(depth, i.explorer); err != nil {
		return nil, err
	}

	return &addressNode, nil
}

type AddressNode struct {
	Address  string         `json:"address"`
	Children []*AddressNode `json:"children,omitempty"`
}

func (n *AddressNode) PopulateNode(depth int, exp thirdparty.Explorer) error {

	addresses, err := getAssociatedAddressesForAddress(n.Address, exp)
	if err != nil {
		return err
	}

	n.Children = addresses

	if depth <= 1 {
		return nil
	}

	for _, child := range n.Children {
		err := child.PopulateNode(depth-1, exp)
		if err != nil {
			return err
		}
	}

	return nil
}

func getAssociatedAddressesForAddress(address string, exp thirdparty.Explorer) ([]*AddressNode, error) {
	txs, err := exp.GetAllTransactionsForAddress(address)
	if err != nil {
		return nil, err
	}
	if len(txs) == 0 {
		return []*AddressNode{}, nil
	}

	addressesSet := map[string]bool{}

	for _, tx := range txs {
		from := tx.From
		to := tx.To

		if from != "" && from != address {
			addressesSet[from] = true
		}
		if to != "" && to != address {
			addressesSet[to] = true

		}
	}

	addresses := []*AddressNode{}
	for k := range addressesSet {
		addresses = append(addresses, &AddressNode{Address: k, Children: []*AddressNode{}})
	}

	return addresses, nil
}
