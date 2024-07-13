package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/gofiber/fiber/v2"
	"github.com/officer47p/addressport/lib/services"
	"github.com/officer47p/addressport/lib/utils"
)

type InvestigationHandler struct {
	investigationService services.InvestigationToolService
}

func NewInvestigationToolHandler(investigationService services.InvestigationToolService) *InvestigationHandler {
	return &InvestigationHandler{investigationService: investigationService}
}

func (h *InvestigationHandler) HandleGetAssociatedTransactionsForAddress(c *fiber.Ctx) error {
	start, endFunc := utils.LogReuqest(c)
	defer endFunc(start)

	address := c.Params("address")
	address = strings.ToLower(address)

	depthString := c.Query("depth", "1")
	format := c.Query("format", "html")
	depth, err := strconv.Atoi(depthString)
	if err != nil {
		return err
	}

	if depth > 3 {
		return errors.New("depth greater than 3 is not currently supported")
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

	if format == "nodesandlinks" {
		return c.JSON(map[string]any{"nodes": nodes, "links": links})
	}

	subtitle := fmt.Sprintf("Transaction analysis for address: %s with depth of %d", address, depth)
	renderedHtmlPage := createGraphChart(address, nodes, links, subtitle)
	// renderedHtmlPage := graphNpmDep()
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Send(renderedHtmlPage)
}

func createGraphChart(address string, nodes *[]services.AddressNode, links *[]graph.Edge[services.AddressNode], subtitle string) (html []byte) {
	graphNodes := []opts.GraphNode{}
	for _, n := range *nodes {
		nodeColor := "#000000"
		// it's the main address
		if n.Address == address {
			nodeColor = "#20aa20"
		}
		graphNodes = append(graphNodes, opts.GraphNode{Name: n.Address, ItemStyle: &opts.ItemStyle{Color: nodeColor}})
	}

	graphLinks := []opts.GraphLink{}
	for _, l := range *links {
		graphLinks = append(graphLinks, opts.GraphLink{Source: l.Source.Address, Target: l.Target.Address})
	}

	page := components.NewPage()
	page.SetLayout(components.PageCenterLayout)

	graph := charts.NewGraph()
	graph.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			// BackgroundColor: "#ffffff",
			Width:  "100vw",
			Height: "100vh",
		}),
		charts.WithColorsOpts(
			opts.Colors{"0x000000"},
		),
		charts.WithTitleOpts(opts.Title{
			Title:    "Addressport",
			Subtitle: subtitle,
		}))

	graph.AddSeries("graph", graphNodes, graphLinks).
		SetSeriesOptions(
			charts.WithGraphChartOpts(opts.GraphChart{
				Force: &opts.GraphForce{Repulsion: 8000},
				Roam:  opts.Bool(true),
				// FocusNodeAdjacency: opts.Bool(true),
				Draggable: opts.Bool(true),
				// Layout: "circular",
				// // Force:              &opts.GraphForce{Repulsion: 10},
				// Roam:               opts.Bool(true),
				// FocusNodeAdjacency: opts.Bool(true),
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

	return page.RenderContent()
}
