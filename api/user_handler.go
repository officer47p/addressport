package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/officer47p/addressport/db"
	"github.com/officer47p/addressport/types"
)

type AddressHandler struct {
	addressStore db.AddressStore
}

func NewAddressHandler(addressStore db.AddressStore) *AddressHandler {
	return &AddressHandler{addressStore: addressStore}
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

func (h *AddressHandler) HandleGetAddress(c *fiber.Ctx) error {
	address := c.Params("address")

	addresses, err := h.addressStore.GetAddressByAddress(c.Context(), address)
	if err != nil {
		return err
	}

	return c.JSON(addresses)
}

func (h *AddressHandler) HandleGetAddresses(c *fiber.Ctx) error {

	users, err := h.addressStore.GetAddresses(c.Context())
	if err != nil {
		return err
	}
	// fmt.Printf("%+v\n", users)
	return c.JSON(&users)
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
