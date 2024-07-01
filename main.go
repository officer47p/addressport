package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/officer47p/addressport/api"
	"github.com/officer47p/addressport/db"
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

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI).SetTimeout(time.Second*5))
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Database(db.DBNAME).ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatalf("warmup connection to database was faild. error: %+v", err)
	}

	addressStore := db.NewMongoAddressStore(client, db.DBNAME)
	// handlers initialization
	addressHandler := api.NewAddressHandler(addressStore)

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/address", addressHandler.HandleGetAddresses)
	apiv1.Get("/address/:address", addressHandler.HandleGetAddress)
	apiv1.Post("/address", addressHandler.HandlePostAddress)

	// not needed now
	// apiv1.Delete("/address/:id", addressHandler.HandleDeleteAddress)
	// apiv1.Put("/address/:id", addressHandler.HandlePutAddress)

	err = app.Listen(*listenAddr)
	if err != nil {
		panic(err)
	}
}
