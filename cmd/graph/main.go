package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {

	n := Node{Name: "binance1", Children: []*Node{}}

	n.PopulateNode(4)

	j, err := json.Marshal(n)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(j))
	// fmt.Printf("%+v\n",)

}

type Node struct {
	Name     string
	Children []*Node
}

func (n *Node) PopulateNode(depth int) {
	addresses := getTransactionsForAddress(n.Name)
	n.Children = addresses

	if depth <= 1 {
		return
	}

	for _, child := range n.Children {
		child.PopulateNode(depth - 1)
	}
}

func getTransactionsForAddress(address string) []*Node {
	if address == "binance1" {
		return []*Node{
			{
				Name:     "binance1.1",
				Children: []*Node{},
			},
			{
				Name:     "binance1.2",
				Children: []*Node{},
			},
			{
				Name:     "binance1.3",
				Children: []*Node{},
			},
		}
	}
	if address == "binance1.2" {
		return []*Node{
			{
				Name:     "binance1.2.1",
				Children: []*Node{},
			},
			{
				Name:     "binance1.2.2",
				Children: []*Node{},
			},
			{
				Name:     "binance1.2.3",
				Children: []*Node{},
			},
		}
	}
	if address == "binance1.2.3" {
		return []*Node{
			{
				Name:     "binance1.2.3.1",
				Children: []*Node{},
			},
			{
				Name:     "binance1.2.3.2",
				Children: []*Node{},
			},
			{
				Name:     "binance1.2.3.3",
				Children: []*Node{},
			},
		}
	}

	return []*Node{}
}

// n1 := Node{
// 	Name: "binance1",
// children: []Node{
// 	{
// 		Name:     "binance1.1",
// 		children: []Node{},
// 	},
// 	{
// 		Name: "binance1.2",
// 		children: []Node{
// {
// 	Name:     "binance1.2.1",
// 	children: []Node{},
// },
// {
// 	Name:     "binance1.2.2",
// 	children: []Node{},
// },
// {
// 	Name:     "binance1.2.3",
// 	children: []Node{
// {
// 	Name:     "binance1.2.3.1",
// 	children: []Node{},
// },
// {
// 	Name:     "binance1.2.3.2",
// 	children: []Node{},
// },
// {
// 	Name:     "binance1.2.3.3",
// 	children: []Node{},
// },
// 	},
// },
// 		},
// 	},
// 	{
// 		Name:     "binance1.3",
// 		children: []Node{},
// 	},
// },
// }
