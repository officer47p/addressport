package services

import (
	"errors"
	"log"

	"github.com/dominikbraun/graph"
	"github.com/officer47p/addressport/lib/thirdparty"
	"github.com/officer47p/addressport/lib/utils"
)

func NewInvestigationToolService(explorer thirdparty.Explorer) InvestigationToolService {
	return InvestigationToolService{explorer: explorer}
}

type InvestigationToolService struct {
	explorer thirdparty.Explorer
}

func (i *InvestigationToolService) GetAllAssociatedTransactionsForAddress(address string, depth int) (*graph.Graph[string, AddressNode], error) {
	g := graph.New(
		func(c AddressNode) string { return c.Address },
	)

	node := AddressNode{Address: address}
	err := populateNode(node, &g, depth, i.explorer)
	if err != nil {
		log.Printf("error while populating the address node: %+v\n", err)
		return nil, err
	}

	return &g, nil
}

func (i *InvestigationToolService) GraphToNodesAndEdges(gRef *graph.Graph[string, AddressNode]) (*[]AddressNode, *[]graph.Edge[AddressNode], error) {
	g := *gRef
	nodesSet := map[string]bool{}
	edges := []graph.Edge[AddressNode]{}

	graphEdges, err := g.Edges()
	if err != nil {
		return nil, nil, err
	}

	for _, e := range graphEdges {
		// first add nodes
		nodesSet[e.Source] = true
		nodesSet[e.Target] = true

		edges = append(edges,
			graph.Edge[AddressNode]{
				Source:     AddressNode{Address: e.Source},
				Target:     AddressNode{Address: e.Target},
				Properties: graph.EdgeProperties{Attributes: map[string]string{"value": e.Properties.Attributes["value"]}},
			})
	}

	nodes := []AddressNode{}
	for k := range nodesSet {
		nodes = append(nodes, AddressNode{Address: k})
	}

	return &nodes, &edges, nil
}

type AddressNode struct {
	Address string `json:"address"`
}

func addNodeEvenIfExists(mainGraph *graph.Graph[string, AddressNode], node AddressNode) error {
	g := *mainGraph
	err := g.AddVertex(node)
	if err != nil {
		if !errors.Is(err, graph.ErrVertexAlreadyExists) {
			log.Printf("error adding vertex: %+v\n", node)

			return err
		}
		return nil
	}
	return nil
}

func addLinkEvenIfExists(mainGraph *graph.Graph[string, AddressNode], from AddressNode, to AddressNode, value string) error {
	g := *mainGraph
	err := g.AddEdge(from.Address, to.Address, graph.EdgeAttribute("value", value))
	if err != nil {
		if !errors.Is(err, graph.ErrEdgeAlreadyExists) {
			log.Printf("error adding Edge from %v to %+v\n", from, to)
			return err
		}
		return nil
	}
	return nil
}

func getOtherAddress(n, from, to AddressNode) AddressNode {
	var otherAddress AddressNode
	if from.Address == n.Address {
		otherAddress = to
	} else {
		otherAddress = from
	}
	return otherAddress
}

func populateNode(n AddressNode, mainGraph *graph.Graph[string, AddressNode], depth int, exp thirdparty.Explorer) error {
	g := *mainGraph
	err := addNodeEvenIfExists(&g, n)
	if err != nil {
		return err
	}

	if depth <= 0 {
		return nil
	}

	txs, err := exp.GetAllTransactionsForAddress(n.Address)
	if err != nil {
		return err
	}

	for _, tx := range txs {
		from := AddressNode{Address: tx.From}
		to := AddressNode{Address: tx.To}

		otherAddress := getOtherAddress(n, from, to)

		if err := addNodeEvenIfExists(&g, from); err != nil {
			return err
		}

		if err := addNodeEvenIfExists(&g, to); err != nil {
			return err
		}

		weiBigInt, ok := utils.StringToBigInt(tx.Value)
		if !ok {
			return errors.New("failed to parse value to big.int")
		}
		value := utils.WeiToEther(weiBigInt)

		if err := addLinkEvenIfExists(&g, from, to, value.String()); err != nil {
			return err
		}
		err := populateNode(otherAddress, &g, depth-1, exp)
		if err != nil {
			return err
		}
	}

	return nil
}
