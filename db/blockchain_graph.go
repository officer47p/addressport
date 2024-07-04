package db

import (
	"errors"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/officer47p/addressport/types"
	"github.com/officer47p/addressport/utils"
)

type BlockchainGraph interface {
	Dropper
	GetAddresses(label string) (*[]dbtype.Node, error)
	AddTransaction(types.Transaction) error
}

type Neo4jBlockchainGraph struct {
	Database Neo4jGraphDatabase
}

func NewNeo4jBlockchainGraph(db Neo4jGraphDatabase) Neo4jBlockchainGraph {
	return Neo4jBlockchainGraph{Database: db}
}

func (graph *Neo4jBlockchainGraph) GetAddresses() (*[]dbtype.Node, error) {
	result, err := neo4j.ExecuteQuery(*graph.Database.Context, *graph.Database.Driver,
		"MATCH (n:Address) RETURN n",
		map[string]any{},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		return nil, err
	}

	nodes := []dbtype.Node{}

	for _, record := range result.Records {
		n := record.Values[0].(dbtype.Node)
		nodes = append(nodes, n)
	}
	return &nodes, nil

}

func (graph *Neo4jBlockchainGraph) AddTransaction(tx types.Transaction) error {
	from, err := graph.addAddress(tx.From)
	if err != nil {
		return err
	}

	to, err := graph.addAddress(tx.To)
	if err != nil {
		return err
	}

	v, ok := utils.StringToBigInt(tx.Value)
	if !ok {
		return errors.New("error when converting string to big.int")
	}
	value, _ := utils.WeiToEther(v).Float64()
	return graph.addRelation(*from, *to, value, tx.TxHash)

}

func (graph *Neo4jBlockchainGraph) addAddress(address string) (*dbtype.Node, error) {
	result, err := neo4j.ExecuteQuery(*graph.Database.Context, *graph.Database.Driver,
		`
		MERGE (n:Address {address: $address})
		RETURN n
		`,
		map[string]any{
			"address": address,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	if err != nil {
		return nil, err
	}
	node := result.Records[0].Values[0].(dbtype.Node)
	return &node, nil
}

func (graph *Neo4jBlockchainGraph) addRelation(from dbtype.Node, to dbtype.Node, value float64, hash string) error {
	result, err := neo4j.ExecuteQuery(*graph.Database.Context, *graph.Database.Driver,
		`
		MATCH (from:Address {address: $from}), (to:Address {address: $to})
		MERGE (from)-[:SEND {value: $value, hash: $hash}]->(to)
		`,
		map[string]any{
			"from":  from.Props["address"],
			"to":    to.Props["address"],
			"value": value,
			"hash":  hash,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", result.Summary)
	return nil
}
