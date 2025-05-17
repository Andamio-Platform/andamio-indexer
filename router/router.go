package router

import (
	"log"

	"github.com/Andamio-Platform/andamio-indexer/config"
	"github.com/Andamio-Platform/andamio-indexer/database"
	address_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/address_handlers"
	asset_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/asset_handlers"
	datum_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/datum_handlers"
	metrics_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/metrics_handlers"
	redeemer_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/redeemer_handlers"
	transaction_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/transaction_handlers"
	witness_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/witness_handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// RouterInit initializes the Fiber router with middleware and routes.
// It configures the router, sets up middleware, defines API versioned routes,
// and starts the server listening on the configured host.
func RouterInit(db *database.Database) {
	router := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Andamio-Indexer",
		AppName:       "Andamio Indexer 1.0.0",
	})

	// middlewares
	api := router.Group("/api", logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency} ${ua}\n",
	}))
	api.Use(helmet.New())
	api.Use(cors.New())
	api.Use(etag.New())
	api.Use(idempotency.New())
	api.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// api.Use(csrf.New(csrf.Config{
	// 	KeyLookup:      "header:X-Csrf-Token",
	// 	CookieName:     "csrf_",
	// 	CookieSameSite: "Lax",
	// 	Expiration:     1 * time.Hour,
	// 	KeyGenerator:   utils.UUIDv4,
	// }))

	// api.Use(compress.New(compress.Config{
	// 	Level: compress.LevelBestCompression, // 1
	// }))

	globalDB := database.GetGlobalDB()

	version := api.Group("/v1")
	indexer := version.Group("/indexer")

	// Serve swagger.json file
	indexer.Static("/docs/swagger.json", "./docs/swagger.json")

	// Setting up Swagger handler
	log.Printf("Setting up Swagger handler for path: /docs/*")
	indexer.Get("/docs/*", swagger.New(swagger.Config{URL: "swagger.json"}))

	// Addresses handlers
	addresses := indexer.Group("/addresses")
	addresses.Post("/", address_handlers.AddAddressHandler(globalDB, nil))
	addresses.Delete("/remove-address", address_handlers.RemoveAddressHandler(globalDB, nil))
	addresses.Get("/:address/utxos", address_handlers.GetUTxOsByAddressHandler(globalDB))
	addresses.Get("/:address/transactions", address_handlers.GetTransactionsByAddressHandler(globalDB))
	addresses.Get("/:address/utxos/inputs", address_handlers.GetUTxOsInputsByAddressHandler(globalDB))
	addresses.Get("/:address/utxos/outputs", address_handlers.GetUTxOsOutputsByAddressHandler(globalDB))
	addresses.Get("/:address/assets", address_handlers.GetAssetsByAddressHandler(globalDB, nil))

	// Transaction handlers
	transactions := indexer.Group("/transactions")
	transactions.Get("/:tx_hash", transaction_handlers.GetTransactionByHashHandler(globalDB))
	transactions.Get("/:tx_hash/utxos", transaction_handlers.GetUTxOsByTransactionHandler(globalDB))
	transactions.Get("/:block_number", transaction_handlers.GetTransactionsByBlockNumberHandler(globalDB))
	transactions.Get("/:tx_hash/utxos/inputs", transaction_handlers.GetUTxOsInputsByTransactionHandler(globalDB))
	transactions.Get("/:tx_hash/utxos/outputs", transaction_handlers.GetUTxOsOutputsByTransactionHandler(globalDB))
	transactions.Get("/:tx_hash/inputs/:input_index/datum", transaction_handlers.GetInputDatumHandler(globalDB, nil))
	transactions.Get("/:tx_hash/outputs/:output_index/datum", transaction_handlers.GetOutputDatumHandler(globalDB, nil))
	transactions.Get("/by-slot-range", transaction_handlers.GetTransactionsBySlotRangeHandler(globalDB, nil))

	// Asset handlers
	asset := indexer.Group("/assets")
	asset.Get("/policy/:policyId/transactions", asset_handlers.GetTransactionsByPolicyIdHandler(globalDB))
	asset.Get("/token/:tokenname/transactions", asset_handlers.GetTransactionsByTokenNameHandler(globalDB))
	asset.Get("/fingerprint/:asset_fingerprint/transactions", asset_handlers.GetTransactionsByAssetFingerprintHandler(globalDB))
	asset.Get("/policy/:policyId/token/:tokenname/transactions", asset_handlers.GetTransactionsByPolicyIdAndTokenNameHandler(globalDB))
	asset.Get("/fingerprint/:asset_fingerprint/addresses", asset_handlers.GetAddressesByAssetFingerprintHandler(globalDB, nil))
	asset.Get("/fingerprint/:asset_fingerprint/utxos", asset_handlers.GetUTxOsByAssetFingerprintHandler(globalDB, nil))

	// Datum handlers
	datum := indexer.Group("/datums")
	datum.Get("/:datum_hash", datum_handlers.GetDatumByHashHandler(globalDB, nil))
	datum.Get("/:datum_hash/transactions", datum_handlers.GetTransactionsByDatumHashHandler(globalDB, nil))

	// Metrics handlers
	metrics := indexer.Group("/metrics")
	metrics.Get("/addresses/count", metrics_handlers.GetAddressesCountHandler(globalDB, nil))
	metrics.Get("/assets/count", metrics_handlers.GetAssetsCountHandler(globalDB, nil))
	metrics.Get("/latest-block", metrics_handlers.GetLatestBlockHandler(globalDB, nil))
	metrics.Get("/transactions/count", metrics_handlers.GetTransactionsCountHandler(globalDB, nil))

	// Redeemer handlers
	redeemers := indexer.Group("/redeemers")
	redeemers.Get("/:tx_hash", redeemer_handlers.GetRedeemersByTxHashHandler(globalDB, nil))
	redeemers.Get("/:tx_hash/inputs/:input_index/redeemer", transaction_handlers.GetInputRedeemerHandler(globalDB, nil))

	// Witness handlers
	witnesses := indexer.Group("/witnesses")
	witnesses.Get("/:witness_id/transactions", witness_handlers.GetTransactionsByWitnessIdHandler(globalDB, nil))
	witnesses.Get("/:witness_id", witness_handlers.GetWitnessByIdHandler(globalDB, nil))

	// Running server
	log.Printf("Server listening on: %s", config.HOST)
	log.Fatal(router.Listen(config.HOST))
}
