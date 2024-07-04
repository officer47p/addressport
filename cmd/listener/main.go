package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/officer47p/addressport/db"
	"github.com/officer47p/addressport/providers"
	"github.com/officer47p/addressport/types"
)

func main() {

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

	provider, err := providers.NewEvmProvider(os.Getenv("ETHREUM_PROVIDER_URI"), types.Network{Name: "ethereum", Currency: "ETH", ChainID: 1, Decimals: 18, StartingBlockNumber: 0})
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
