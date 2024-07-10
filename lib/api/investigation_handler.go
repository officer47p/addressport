package api

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/officer47p/addressport/lib/thirdparty"
)

type InvestigationHandler struct {
	explorer thirdparty.Explorer
}

func NewInvestigationToolHandler(explorer thirdparty.Explorer) *InvestigationHandler {
	return &InvestigationHandler{explorer: explorer}
}

func (h *InvestigationHandler) HandleGetAssociatedAddresses(c *fiber.Ctx) error {
	tempId := rand.Intn(100_000_000_000)
	reqId := strconv.Itoa(tempId)
	log.Printf("%s %s request(%s)\n", c.Method(), c.OriginalURL(), string(reqId))
	start := time.Now()
	defer func() {
		log.Printf("request(%s) took %d ms\n", reqId, time.Since(start).Milliseconds())
	}()

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
