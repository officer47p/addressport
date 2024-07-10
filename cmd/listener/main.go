package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/officer47p/addressport/lib/db"
	"github.com/officer47p/addressport/lib/thirdparty"
	"github.com/officer47p/addressport/lib/types"
)

func main() {
	// Shortest path query for future:
	// PROFILE
	// MATCH
	// (addr1:Address {address: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48'}),
	// (addr2:Address {address: '0x422F5acCC812C396600010f224b320a743695f85'}),
	// p = SHORTESTPATH((addr1)-[:SEND*]-(addr2))
	// // WHERE all(r IN relationships(p) WHERE r.role IS NOT NULL)
	// RETURN p

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	neo4jDB, err := db.NewNeo4jGraphDatabase(
		os.Getenv("NEO4J_URI"),
		os.Getenv("NEO4J_USER"),
		os.Getenv("NEO4J_PASSWORD"),
	)
	if err != nil {
		log.Fatal(err)
	}

	graphDB := db.NewNeo4jBlockchainGraph(*neo4jDB)

	provider, err := thirdparty.NewEvmProvider(os.Getenv("ETHREUM_PROVIDER_URI"), types.Network{Name: "ethereum", Currency: "ETH", ChainID: 1, Decimals: 18, StartingBlockNumber: 0})
	if err != nil {
		log.Fatal(err)
	}

	for i := 20229300; i < 20229303; i++ {
		block, err := provider.GetBlockByNumber(int64(i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("block number: ", block.BlockNumber)
		if err != nil {
			log.Fatal(err)
		}
		for _, tx := range block.Transactions {
			fmt.Println("transaction: ", tx.TxHash)
			if err := graphDB.AddTransaction(tx); err != nil {
				log.Fatal(err)
			}

		}
		fmt.Println(block)
	}
	// nodes, err := graphDB.GetAddresses()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(nodes)
}
