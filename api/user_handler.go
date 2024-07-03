package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/officer47p/addressport/db"
	"github.com/officer47p/addressport/explorer"
	"github.com/officer47p/addressport/types"
)

type AddressHandler struct {
	addressStore db.AddressStore
	explorer     explorer.Explorer
}

func NewAddressHandler(addressStore db.AddressStore, explorer explorer.Explorer) *AddressHandler {
	return &AddressHandler{addressStore: addressStore, explorer: explorer}
}

func (h *AddressHandler) HandlePostAddress(c *fiber.Ctx) error {
	var params types.CreateAddressParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errs := params.Validate(); len(errs) > 0 {
		return errors.Join(errs...)
	}

	address, err := types.NewAddressFromParams(params)
	if err != nil {
		return err
	}

	insertedAddress, err := h.addressStore.InsertAddress(c.Context(), address)
	if err != nil {
		return err
	}

	return c.JSON(insertedAddress)
}

func (h *AddressHandler) HandleGetAddresses(c *fiber.Ctx) error {

	users, err := h.addressStore.GetAddresses(c.Context())
	if err != nil {
		return err
	}
	// fmt.Printf("%+v\n", users)
	return c.JSON(&users)
}

func (h *AddressHandler) HandleGetAddressByAddress(c *fiber.Ctx) error {
	address := c.Params("address")

	addresses, err := h.addressStore.GetAddressesByAddress(c.Context(), address)
	if err != nil {
		return err
	}

	return c.JSON(addresses)
}

func (h *AddressHandler) HandleGetAssociatedAddresses(c *fiber.Ctx) error {
	address := c.Params("address")

	txs, err := h.explorer.GetAllTransactionsForAddress(address)
	if err != nil {
		return err
	}
	if len(txs) == 0 {
		return c.SendStatus(404)
	}

	addresses := []string{}

	for _, tx := range txs {
		from := tx.From
		to := tx.To

		if from != address {
			addresses = append(addresses, from)
		}
		if to != address {
			addresses = append(addresses, to)
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

// func (h *AddressHandler) HandleDeleteAddress(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	err := h.addressStore.DeleteAddress(c.Context(), id)
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(map[string]string{"deleted": id})
// }

// func (h *AddressHandler) HandlePutAddress(c *fiber.Ctx) error {
// 	var (
// 		id     = c.Params("id")
// 		params types.UpdateAddressParams
// 	)
// 	oid, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	if err := c.BodyParser(&params); err != nil {
// 		return err
// 	}

// 	if errs := params.Validate(); len(errs) > 0 {
// 		return errors.Join(errs...)
// 	}

// 	filter := bson.M{"_id": oid}
// 	if err := h.addressStore.UpdateAddress(c.Context(), filter, params); err != nil {
// 		return err
// 	}
// 	return c.JSON(map[string]string{"updated": id})
// }
