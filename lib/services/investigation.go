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

func (i *InvestigationToolService) GetAllAssociatedTransactionsForAddress(address string, depth int) (*graph.Graph[string, AddressNode], error) {
	g := graph.New(AddressHash)
	// g.AddEdge("dbfskj", "kdfjnsjk")
	node := AddressNode{Address: address}
	err := populateNode(node, &g, depth, i.explorer)
	if err != nil {
		log.Printf("error: %+v\n", err)
		return nil, err
	}

	return &g, nil
	// fmt.Printf("%+v hiiii\n", allEdges)
	// return c.JSON(map[string]any{"Links": allEdges})
}

func AddressHash(c AddressNode) string {
	return c.Address
}

func populateNode(n AddressNode, mainGraph *graph.Graph[string, AddressNode], depth int, exp thirdparty.Explorer) error {
	g := *mainGraph
	addVertexError := g.AddVertex(n)
	if addVertexError != nil {
		if !errors.Is(addVertexError, graph.ErrVertexAlreadyExists) {
			log.Printf("error adding vertex: %+v\n", n)

			return addVertexError
		}
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

		var otherAddress AddressNode
		if from.Address == n.Address {
			otherAddress = to
		} else {
			otherAddress = from
		}

		errFrom := g.AddVertex(from)
		if errFrom != nil {
			if !errors.Is(errFrom, graph.ErrVertexAlreadyExists) {
				log.Printf("error adding vertex: %+v\n", from)

				return errFrom
			}
		}
		errTo := g.AddVertex(to)
		if errTo != nil {
			if !errors.Is(errTo, graph.ErrVertexAlreadyExists) {
				log.Printf("error adding vertex: %+v\n", to)

				return errTo
			}
		}

		weiBigInt, ok := utils.StringToBigInt(tx.Value)
		if !ok {
			return errors.New("failed to parse value to big.int")
		}
		value := utils.WeiToEther(weiBigInt)

		errAddEdge := g.AddEdge(from.Address, to.Address, graph.EdgeAttribute("value", value.String()))
		if errAddEdge != nil {
			if !errors.Is(errAddEdge, graph.ErrEdgeAlreadyExists) {
				log.Printf("error adding Edge from %v to %+v\n", from, to)
				return err
			}
		}
		err = populateNode(otherAddress, &g, depth-1, exp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *InvestigationToolService) GetAssociatedAddressesForAddress(address string, depth int) (*AddressNode, error) {
	addressNode := AddressNode{Address: address, Children: []*AddressNode{}}
	if err := addressNode.PopulateNode(depth, i.explorer); err != nil {
		return nil, err
	}

	return &addressNode, nil
}

type AddressNode struct {
	Address  string         `json:"address"`
	Children []*AddressNode `json:"children,omitempty"`
}

func (n *AddressNode) PopulateNode(depth int, exp thirdparty.Explorer) error {

	addresses, err := getAssociatedAddressesForAddress(n.Address, exp)
	if err != nil {
		return err
	}

	n.Children = addresses

	if depth <= 1 {
		return nil
	}

	for _, child := range n.Children {
		err := child.PopulateNode(depth-1, exp)
		if err != nil {
			return err
		}
	}

	return nil
}

func getAssociatedAddressesForAddress(address string, exp thirdparty.Explorer) ([]*AddressNode, error) {
	txs, err := exp.GetAllTransactionsForAddress(address)
	if err != nil {
		return nil, err
	}
	if len(txs) == 0 {
		return []*AddressNode{}, nil
	}

	addressesSet := map[string]bool{}

	for _, tx := range txs {
		from := tx.From
		to := tx.To

		if from != "" && from != address {
			addressesSet[from] = true
		}
		if to != "" && to != address {
			addressesSet[to] = true

		}
	}

	addresses := []*AddressNode{}
	for k := range addressesSet {
		addresses = append(addresses, &AddressNode{Address: k, Children: []*AddressNode{}})
	}

	return addresses, nil
}
