package api

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/officer47p/addressport/explorer"
)

type InvestigationHandler struct {
	explorer explorer.Explorer
}

func NewInvestigationHandler(explorer explorer.Explorer) *InvestigationHandler {
	return &InvestigationHandler{explorer: explorer}
}

func (h *InvestigationHandler) HandleGetAssociatedAddresses(c *fiber.Ctx) error {
	address := c.Params("address")
	address = strings.ToLower(address)

	depthString := c.Query("depth", "1")
	depth, err := strconv.Atoi(depthString)
	if err != nil {
		return err
	}

	addressNode := AddressNode{Address: address, Children: []*AddressNode{}}
	if err = addressNode.PopulateNode(depth, h.explorer); err != nil {
		return err
	}

	return c.JSON(addressNode)

}

type AddressNode struct {
	Address  string         `json:"address"`
	Children []*AddressNode `json:"children,omitempty"`
}

func (n *AddressNode) PopulateNode(depth int, exp explorer.Explorer) error {

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

func getAssociatedAddressesForAddress(address string, exp explorer.Explorer) ([]*AddressNode, error) {
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

		if from != address {
			addressesSet[from] = true
		}
		if to != address {
			addressesSet[to] = true

		}
	}

	addresses := []*AddressNode{}
	for k, _ := range addressesSet {
		addresses = append(addresses, &AddressNode{Address: k, Children: []*AddressNode{}})
	}

	return addresses, nil
}
