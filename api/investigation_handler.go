package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/officer47p/addressport/explorer"
)

type InvestigationHandler struct {
	explorer explorer.Explorer
}

type AddressNode struct {
	Address  string
	Children []*AddressNode
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
func NewInvestigationHandler(explorer explorer.Explorer) *InvestigationHandler {
	return &InvestigationHandler{explorer: explorer}
}

func (h *InvestigationHandler) HandleGetAssociatedAddresses(c *fiber.Ctx) error {
	address := c.Params("address")
	depthString := c.Query("depth", "1")
	depth, err := strconv.Atoi(depthString)
	if err != nil {
		return err
	}

	addressNode := AddressNode{Address: address, Children: []*AddressNode{}}
	addressNode.PopulateNode(depth, h.explorer)

	return c.JSON(addressNode)

	// txs, err := h.explorer.GetAllTransactionsForAddress(address)
	// if err != nil {
	// 	return err
	// }
	// if len(txs) == 0 {
	// 	return c.SendStatus(404)
	// }

	// addresses := map[string]bool{}

	// for _, tx := range txs {
	// 	from := tx.From
	// 	to := tx.To

	// 	if from != address {
	// 		addresses[from] = true
	// 	}
	// 	if to != address {
	// 		addresses[to] = true

	// 	}
	// }

	// return c.JSON(addresses)

	// for txn in txn_data['result']:
	//         from_address = txn['from']
	//         to_address = txn['to']

	//         # Check if 'from' address is flagged
	//         if from_address != address:
	//             from_flagged = check_address_flagged(from_address)
	//             if from_flagged and from_address not in scam_interactions:
	//                 scam_interactions.append(from_address)

	//         # Check if 'to' address is flagged
	//         if to_address != address:
	//             to_flagged = check_address_flagged(to_address)
	//             if to_flagged and to_address not in scam_interactions:
	//                 scam_interactions.append(to_address)
}
