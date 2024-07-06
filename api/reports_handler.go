package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/officer47p/addressport/db"
	"github.com/officer47p/addressport/types"
)

type ReportsHandler struct {
	reportsStore db.ReportsStore
}

func NewReportsHandler(reportsStore db.ReportsStore) *ReportsHandler {
	return &ReportsHandler{reportsStore: reportsStore}
}

func (h *ReportsHandler) HandlePostReport(c *fiber.Ctx) error {
	var params types.CreateReportParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errs := params.Validate(); len(errs) > 0 {
		return errors.Join(errs...)
	}

	report, err := types.NewReportFromParams(params)
	if err != nil {
		return err
	}

	insertedReport, err := h.reportsStore.InsertReport(c.Context(), report)
	if err != nil {
		return err
	}

	return c.JSON(insertedReport)
}

func (h *ReportsHandler) HandleGetReports(c *fiber.Ctx) error {

	reports, err := h.reportsStore.GetReports(c.Context())
	if err != nil {
		return err
	}
	// fmt.Printf("%+v\n", reports)
	return c.JSON(&reports)
}

func (h *ReportsHandler) HandleGetReportsByAddress(c *fiber.Ctx) error {
	address := c.Params("address")

	reports, err := h.reportsStore.GetReportsByAddress(c.Context(), address)
	if err != nil {
		return err
	}

	return c.JSON(reports)
}

func (h *ReportsHandler) HandleDeleteReport(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.reportsStore.DeleteReport(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": id})
}

func (h *ReportsHandler) HandlePutReportById(c *fiber.Ctx) error {
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
	if err := h.reportsStore.UpdateReport(c.Context(), id, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": id})
}
