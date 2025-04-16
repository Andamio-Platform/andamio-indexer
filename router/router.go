package router

import (
	"log"

	"github.com/Andamio-Platform/andamio-indexer/config"
	"github.com/Andamio-Platform/andamio-indexer/database"
	address_handlers "github.com/Andamio-Platform/andamio-indexer/handlers/v1/address_handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gorilla/mux"
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
	api := router.Group("/api", logger.New())
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

	version := api.Group("/v1")

	indexer := version.Group("/indexer")

	address := indexer.Group("/address")

	r := mux.NewRouter()
	address_handlers.AddAddressRoute(r, db, nil)

	address.Delete("/remove-address", address_handlers.RemoveAddressHandler)
	address.Get("/{address}/utxos", func(c *fiber.Ctx) error {
		return address_handlers.GetUTxOsByAddressHandler(c, db)
	})

	// event := indexer.Group("/event")

	// Running server
	log.Fatal(router.Listen(config.HOST))
}
