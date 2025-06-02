package router

import (
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/config"
	"github.com/Andamio-Platform/andamio-indexer/database"
	address_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/address_handlers"
	asset_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/asset_handlers"
	metrics_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/metrics_handlers"
	redeemer_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/redeemer_handlers"
	transaction_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/transaction_handlers"
	"github.com/Andamio-Platform/andamio-indexer/internal/logutils" // Add this import

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	fiberMiddlewareLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

// RouterInit initializes the Fiber router with middleware and routes.
// It configures the router, sets up middleware, defines API versioned routes,
// and starts the server listening on the configured host.
func RouterInit(router *fiber.App, db *database.Database, logger *slog.Logger) {
	// middlewares
	api := router.Group("/api")
	api.Use(helmet.New())
	api.Use(cors.New())
	api.Use(etag.New())
	api.Use(idempotency.New())
	api.Use(fiberMiddlewareLogger.New(fiberMiddlewareLogger.Config{
		Format: "[${pid}] [${ip}]:${port} ${status} - ${method} ${path} ${latency} ${bytesReceived} ${bytesSent} ${reqHeaders} ${resHeaders} ${body} ${error}\n",
		Output: logutils.NewSlogWriter(logger),
	}))
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
	logger.Info("Setting up Swagger handler", "path", "/docs/*")
	indexer.Get("/docs/*", swagger.New(swagger.Config{URL: config.GetGlobalConfig().Indexer.SwaggerURL}))

	// Addresses handlers
	addresses := indexer.Group("/addresses")
	addresses.Post("/", address_handlers.AddAddressHandler(globalDB, logger))
	addresses.Delete("/remove-address", address_handlers.RemoveAddressHandler(globalDB, logger))
	addresses.Get("/:address/transactions", address_handlers.GetTransactionsByAddressHandler(globalDB))
	addresses.Get("/:address/assets", address_handlers.GetAssetsByAddressHandler(globalDB, logger))

	// Transaction handlers
	transactions := indexer.Group("/transactions")
	transactions.Get("/by-slot-range", transaction_handlers.GetTransactionsBySlotRangeHandler(globalDB, logger))
	transactions.Get("/by-block-number/:block_number", transaction_handlers.GetTransactionsByBlockNumberHandler(globalDB))
	transactions.Get("/:tx_hash", transaction_handlers.GetTransactionByTxHashHandler(globalDB))
	transactions.Get("/:tx_hash/utxos", transaction_handlers.GetUTxOsByTransactionHandler(globalDB))
	transactions.Get("/:tx_hash/utxos/inputs", transaction_handlers.GetUTxOsInputsByTransactionHandler(globalDB))
	transactions.Get("/:tx_hash/utxos/outputs", transaction_handlers.GetUTxOsOutputsByTransactionHandler(globalDB))

	// Asset handlers
	asset := indexer.Group("/assets")
	asset.Get("/policy/:policyId/transactions", asset_handlers.GetTransactionsByPolicyIdHandler(globalDB))
	asset.Get("/token/:tokenname/transactions", asset_handlers.GetTransactionsByTokenNameHandler(globalDB))
	asset.Get("/fingerprint/:asset_fingerprint/transactions", asset_handlers.GetTransactionsByAssetFingerprintHandler(globalDB))
	asset.Get("/policy/:policyId/token/:tokenname/transactions", asset_handlers.GetTransactionsByPolicyIdAndTokenNameHandler(globalDB))
	asset.Get("/fingerprint/:asset_fingerprint/addresses", asset_handlers.GetAddressesByAssetFingerprintHandler(globalDB, logger))
	asset.Get("/fingerprint/:asset_fingerprint/utxos", asset_handlers.GetUTxOsByAssetFingerprintHandler(globalDB, logger))

	// Metrics handlers
	metrics := indexer.Group("/metrics")
	metrics.Get("/addresses/count", metrics_handlers.GetAddressesCountHandler(globalDB, logger))
	metrics.Get("/assets/count", metrics_handlers.GetAssetsCountHandler(globalDB, logger))
	metrics.Get("/latest-block", metrics_handlers.GetLatestBlockHandler(globalDB, logger))
	metrics.Get("/transactions/count", metrics_handlers.GetTransactionsCountHandler(globalDB, logger))

	// Redeemer handlers
	redeemers := indexer.Group("/redeemers")
	redeemers.Get("/:tx_hash", redeemer_handlers.GetRedeemersByTxHashHandler(globalDB, logger))
}
