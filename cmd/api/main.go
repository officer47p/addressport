package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"github.com/officer47p/addressport/pkg/api"
	"github.com/officer47p/addressport/pkg/services"
	"github.com/officer47p/addressport/pkg/thirdparty"
)

func main() {

	listenAddr := flag.String("listenAddr", ":3000", "The listen address")
	flag.Parse()

	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file. err: %s", err.Error())
		}
	}

	// thirdparty initialization
	etherscanExplorer := thirdparty.NewEtherscanExplorer(os.Getenv("ETHERSCAN_API_KEY"))
	// service initialization
	investigationToolService := services.NewInvestigationToolService(etherscanExplorer)
	// handlers initialization
	investigationToolHandler := api.NewInvestigationToolHandler(investigationToolService)

	fiberConfig, err := getFiberConfig()
	if err != nil {
		log.Fatalf("error when getting fiber configurations. err: %s", err.Error())
	}

	app := fiber.New(*fiberConfig)
	app.Use(cors.New())

	dashboardRouter := app.Group("/")
	dashboardRouter.Get("/", investigationToolHandler.HandleAddressInfoForm)

	apiv1Router := app.Group("/api/v1")
	apiv1Router.Get("/investigation/tools/transaction-association/:address", investigationToolHandler.HandleGetAssociatedTransactionsForAddress)

	err = app.Listen(*listenAddr)
	if err != nil {
		log.Fatalf("error when exiting server. err: %s", err.Error())
	}
}

func getFiberConfig() (*fiber.Config, error) {
	engine := html.New("./pkg/views", ".html")
	if err := engine.Load(); err != nil {
		return nil, fmt.Errorf("error when loding templates. err: %s", err.Error())
	}
	return &fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.JSON(map[string]string{"error": err.Error()})
		},
	}, nil
}
