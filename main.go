package main

import (
	"flag"
	"log/slog"
	"os"
	"time" // Add this import for timeouts

	"github.com/Andamio-Platform/andamio-indexer/internal/logutils"
	"github.com/gofiber/fiber/v2" // Add this import
	fiberLogger "github.com/gofiber/fiber/v2/log"
	"github.com/lmittmann/tint"

	"github.com/Andamio-Platform/andamio-indexer/config"
	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/indexer"
	"github.com/Andamio-Platform/andamio-indexer/router"
)

var cmdlineFlags struct {
	configFile string
}

//	@title			Andamio Indexer 1.0.0
//	@version		1.0.0
//	@description	Indexer APIs for Andamio dapp-system.

//	@contact.name	Andamio Support
//	@contact.url	https://www.andamio.com/support
//	@contact.email	dev@andamio.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		142.132.201.159:42069
//	@BasePath	/api/v1/indexer

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	// Load config
	flag.StringVar(
		&cmdlineFlags.configFile,
		"config",
		"",
		"path to config file to load",
	)
	flag.Parse()

	if cmdlineFlags.configFile == "" {
		slog.Error("Config file not specified. Use the -config flag.")
		os.Exit(1)
	}

	err := config.Load(cmdlineFlags.configFile)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelInfo,
		// Level: slog.LevelDebug,
		AddSource: true, // adds file and line number
	}))
	slog.SetDefault(logger)

	// Redirect Fiber's internal logger output to slog
	fiberLogger.SetOutput(logutils.NewSlogWriter(logger))
	fiberLogger.SetLevel(fiberLogger.LevelDebug)
	// fiberLogger.SetLevel(fiberLogger.LevelInfo)

	var db *database.Database
	if !fiber.IsChild() {
		var err error
		db, err = database.New(logger, config.GlobalConfig.Database.DatabaseDIR)
		if err != nil {
			slog.Error("Failed to init database", "error", err)
			os.Exit(1)
		}
		database.SetGlobalDB(db)

		go func() {
			if err := indexer.StartIndexer(db, logger); err != nil {
				slog.Error("Failed to init adder", "error", err)
				os.Exit(1)
			}
		}()

		defer db.Close()
	}

	app := fiber.New(fiber.Config{
		Prefork:           false, // Set to false to avoid database lock issues with file-based DB
		CaseSensitive:     true,
		StrictRouting:     true,
		ServerHeader:      "Andamio-Indexer",
		AppName:           "Andamio Indexer 1.0.0",
		ReadTimeout:       10 * time.Second, // Timeout for reading new connections
		WriteTimeout:      10 * time.Second, // Timeout for writing data to the client
		IdleTimeout:       60 * time.Second, // Timeout for idle connections
		EnablePrintRoutes: true,             // Set to true for debugging, false for production
		ReduceMemoryUsage: false,            // Set to true to reduce memory usage, might impact performance
		BodyLimit:         10 * 1024 * 1024, // 10 MB limit for request bodies
	})

	router.RouterInit(app, db, logger)

	host := config.GetGlobalConfig().Indexer.Host
	logger.Info("Attempting to listen", "host", host)
	logger.Error(app.Listen(host).Error())
}
