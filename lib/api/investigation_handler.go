package api

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
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

func (h *InvestigationHandler) HandleGetAssociatedTransactionsForAddress(c *fiber.Ctx) error {
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

	result, err := h.investigationService.GetAllAssociatedTransactionsForAddress(address, depth)
	if err != nil {
		return err
	}
	fmt.Println(result)

	nodes, links, err := h.investigationService.GraphToNodesAndEdges(result)
	if err != nil {
		return err
	}

	graphNodes := []opts.GraphNode{}
	for _, n := range *nodes {
		graphNodes = append(graphNodes, opts.GraphNode{Name: n.Address})
	}

	graphLinks := []opts.GraphLink{}
	for _, l := range *links {
		graphLinks = append(graphLinks, opts.GraphLink{Source: l.Source.Address, Target: l.Target.Address})
	}

	page := components.NewPage()
	graph := charts.NewGraph()
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Addressport",
		}))

	graph.AddSeries("graph", graphNodes, graphLinks).
		SetSeriesOptions(
			charts.WithGraphChartOpts(opts.GraphChart{
				Layout: "circular",
				// Force:              &opts.GraphForce{Repulsion: 10},
				Roam:               opts.Bool(true),
				FocusNodeAdjacency: opts.Bool(true),
			}),
			charts.WithEmphasisOpts(opts.Emphasis{
				Label: &opts.Label{
					Show:     opts.Bool(true),
					Color:    "black",
					Position: "left",
				},
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Curveness: 0.3,
			}),
		)

	page.AddCharts(graph)
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Send((page.RenderContent()))

	// return c.JSON(map[string]any{"nodes": nodes, "links": links})

}
