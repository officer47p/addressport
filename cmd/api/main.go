package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/officer47p/addressport/lib/api"
	"github.com/officer47p/addressport/lib/db"
	"github.com/officer47p/addressport/lib/services"
	"github.com/officer47p/addressport/lib/thirdparty"
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

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI).SetTimeout(time.Second*5))
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Database(db.DBNAME).ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		log.Fatalf("warmup connection to database was faild. error: %+v", err)
	}

	// // Dependencies
	// store initialization
	reportsStore := db.NewMongoReportsStore(client, db.DBNAME)
	// thirdparty initialization
	etherscanExplorer := thirdparty.NewEtherscanExplorer(os.Getenv("ETHERSCAN_API_KEY"))
	// service initialization
	reportsService := services.NewReportsService(reportsStore)
	investigationToolService := services.NewInvestigationToolService(etherscanExplorer)

	// handlers initialization
	reportsHandler := api.NewReportsHandler(reportsService)
	investigationToolHandler := api.NewInvestigationToolHandler(investigationToolService)

	app := fiber.New(config)
	app.Use(cors.New())

	apiv1 := app.Group("/api/v1")
	// application := app.Group("/app/investigation/")

	// TODO: add query filters to the handler
	apiv1.Get("/reports", reportsHandler.HandleGetReports)
	apiv1.Post("/reports", reportsHandler.HandlePostReport)
	apiv1.Get("/reports/:address", reportsHandler.HandleGetReportsByAddress)
	apiv1.Delete("/reports/:id", reportsHandler.HandleDeleteReport)
	apiv1.Put("/reports/:id", reportsHandler.HandlePutReportById)

	apiv1.Get("/investigation/tools/address-association/:address", investigationToolHandler.HandleGetAssociatedAddresses)
	apiv1.Get("/investigation/tools/transaction-association/:address", investigationToolHandler.HandleGetAssociatedTransactionsForAddress)

	err = app.Listen(*listenAddr)
	if err != nil {
		panic(err)
	}
}
