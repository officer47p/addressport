package api

import (
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

	txs, err := h.explorer.GetAllTransactionsForAddress(address)
	if err != nil {
		return err
	}
	if len(txs) == 0 {
		return c.SendStatus(404)
	}

	addresses := map[string]bool{}

	for _, tx := range txs {
		from := tx.From
		to := tx.To

		if from != address {
			addresses[from] = true
		}
		if to != address {
			addresses[to] = true

		}
	}

	return c.JSON(addresses)

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
