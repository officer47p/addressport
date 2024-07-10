package api

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/officer47p/addressport/lib/services"
	"github.com/officer47p/addressport/lib/types"
)

type ReportsHandler struct {
	reportsService services.ReportsService
}

func NewReportsHandler(reportsService services.ReportsService) *ReportsHandler {
	return &ReportsHandler{reportsService: reportsService}
}

func (h *ReportsHandler) HandlePostReport(c *fiber.Ctx) error {
	tempId := rand.Intn(100_000_000_000)
	reqId := strconv.Itoa(tempId)
	log.Printf("%s %s request(%s) body: %s\n", c.Method(), c.OriginalURL(), string(reqId), string(c.Body()))
	start := time.Now()
	defer func() {
		log.Printf("request(%s) took %d ms\n", reqId, time.Since(start).Milliseconds())
	}()

	var params types.CreateReportParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errs := params.Validate(); len(errs) > 0 {
		return errors.Join(errs...)
	}

	report, err := h.reportsService.CreateReport(c.Context(), params)
	if err != nil {
		return err
	}

	return c.JSON(report)
}

func (h *ReportsHandler) HandleGetReports(c *fiber.Ctx) error {
	tempId := rand.Intn(100_000_000_000)
	reqId := strconv.Itoa(tempId)
	log.Printf("%s %s request(%s)\n", c.Method(), c.OriginalURL(), string(reqId))
	start := time.Now()
	defer func() {
		log.Printf("request(%s) took %d ms\n", reqId, time.Since(start).Milliseconds())
	}()

	reports, err := h.reportsService.GetAllReports(c.Context())
	if err != nil {
		return err
	}
	// fmt.Printf("%+v\n", reports)
	return c.JSON(&reports)
}

func (h *ReportsHandler) HandleGetReportsByAddress(c *fiber.Ctx) error {
	tempId := rand.Intn(100_000_000_000)
	reqId := strconv.Itoa(tempId)
	log.Printf("%s %s request(%s)\n", c.Method(), c.OriginalURL(), string(reqId))
	start := time.Now()
	defer func() {
		log.Printf("request(%s) took %d ms\n", reqId, time.Since(start).Milliseconds())
	}()

	address := c.Params("address")

	reports, err := h.reportsService.GetReportsForAddress(c.Context(), address)
	if err != nil {
		return err
	}

	return c.JSON(reports)
}

func (h *ReportsHandler) HandleDeleteReport(c *fiber.Ctx) error {
	tempId := rand.Intn(100_000_000_000)
	reqId := strconv.Itoa(tempId)
	log.Printf("%s %s request(%s)\n", c.Method(), c.OriginalURL(), string(reqId))
	start := time.Now()
	defer func() {
		log.Printf("request(%s) took %d ms\n", reqId, time.Since(start).Milliseconds())
	}()

	id := c.Params("id")
	err := h.reportsService.DeleteReportById(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": id})
}

func (h *ReportsHandler) HandlePutReportById(c *fiber.Ctx) error {
	tempId := rand.Intn(100_000_000_000)
	reqId := strconv.Itoa(tempId)
	log.Printf("%s %s request(%s)\n", c.Method(), c.OriginalURL(), string(reqId))
	start := time.Now()
	defer func() {
		log.Printf("request(%s) took %d ms\n", reqId, time.Since(start).Milliseconds())
	}()

	var (
		id     = c.Params("id")
		params types.UpdateReportParams
	)

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errs := params.Validate(); len(errs) > 0 {
		return errors.Join(errs...)
	}

	// filter := bson.M{"_id": oid}
	if err := h.reportsService.UpdateReportById(c.Context(), id, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": id})
}
