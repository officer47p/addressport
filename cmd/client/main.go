package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-resty/resty/v2"
)

type Node struct {
	Address string `json:"address"`
}

type Link struct {
	Source struct {
		Address string `json:"address"`
	}
	Target struct {
		Address string `json:"address"`
	}
}

type AddressPortInvestigationResponse struct {
	Links []Link `json:"links"`
	Nodes []Node `json:"nodes"`
}

func graphNpmDep() *charts.Graph {
	graph := charts.NewGraph()
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Addressport",
		}))

	// Create a Resty Client
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"depth": "3",
		}).
		SetHeader("Accept", "application/json").
		// SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
		Get(fmt.Sprintf("http://localhost:3001/api/v1/investigation/tools/transaction-association/%s", "0xC1D8E8f14b6AA1cf2F2321348Cbb51d94dc73152"))

	if err != nil {
		log.Printf("explorer: error while fetching transactions of an address. err: %e\n", err)
		log.Fatalln(err)
	}

	var result AddressPortInvestigationResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		log.Fatalln(err)
	}

	nodes := []opts.GraphNode{}
	for _, n := range result.Nodes {
		nodes = append(nodes, opts.GraphNode{Name: n.Address})
	}

	links := []opts.GraphLink{}
	for _, l := range result.Links {
		links = append(links, opts.GraphLink{Source: l.Source.Address, Target: l.Target.Address})
	}

	fmt.Printf("nodes: %+v\n", nodes)
	fmt.Printf("links: %+v\n", links)
	// fmt.Printf("number of links: %d\n", len(links))

	graph.AddSeries("graph", nodes, links).
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
	return graph
}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	page := components.NewPage()
	page.AddCharts(
		graphNpmDep(),
	)

	page.Render(w)

}

func main() {
	http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8081", nil)
}

type AddressNode struct {
	Address  string         `json:"address"`
	Children []*AddressNode `json:"children,omitempty"`
}
