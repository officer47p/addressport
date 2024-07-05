package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/officer47p/addressport/api"
	"github.com/officer47p/addressport/db"
	"github.com/officer47p/addressport/explorer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	listenAddr := flag.String("listenAddr", ":3000", "The listen address")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// neo4jDB, err := db.NewNeo4jGraphDatabase(
	// 	os.Getenv("NEO4J_URI"),
	// 	os.Getenv("NEO4J_USER"),
	// 	os.Getenv("NEO4J_PASSWORD"),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// graphDB := db.NewNeo4jBlockchainGraph(*neo4jDB)

	// nodes, err := graphDB.GetAddresses()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(nodes)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI).SetTimeout(time.Second*5))
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Database(db.DBNAME).ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatalf("warmup connection to database was faild. error: %+v", err)
	}

	// stores initialization
	addressStore := db.NewMongoAddressStore(client, db.DBNAME)
	// explorer initialization
	etherscanExplorer := explorer.NewEtherscanExplorer(os.Getenv("ETHERSCAN_API_KEY"))
	// handlers initialization
	addressHandler := api.NewAddressHandler(addressStore, etherscanExplorer)

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/address", addressHandler.HandleGetAddresses)
	apiv1.Post("/address", addressHandler.HandlePostAddress)
	apiv1.Get("/address/:address", addressHandler.HandleGetAddressByAddress)
	apiv1.Get("/address/:address/associated", addressHandler.HandleGetAssociatedAddresses)

	// not needed now
	// apiv1.Delete("/address/:id", addressHandler.HandleDeleteAddress)
	// apiv1.Put("/address/:id", addressHandler.HandlePutAddress)

	err = app.Listen(*listenAddr)
	if err != nil {
		panic(err)
	}
}