package api

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/officer47p/addressport/lib/services"
)

type InvestigationHandler struct {
	investigationService services.InvestigationToolService
}

func NewInvestigationToolHandler(investigationService services.InvestigationToolService) *InvestigationHandler {
	return &InvestigationHandler{investigationService: investigationService}
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

	result, err := h.investigationService.GetAssociatedAddressesForAddress(address, depth)
	if err != nil {
		return err
	}

	return c.JSON(result)

}
