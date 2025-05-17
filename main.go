package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/Andamio-Platform/andamio-indexer/config"
	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/router"
	"github.com/Andamio-Platform/andamio-indexer/indexer"
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

//	@host		0.0.0.0:42069
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	db, err := database.New(logger, config.GlobalConfig.Database.DatabaseDIR)
	if err != nil {
		slog.Error("Failed to init database", "error", err)
		os.Exit(1)
	}

	database.SetGlobalDB(db)

	go func() {
		if err := indexer.StartIndexer(logger); err != nil {
			slog.Error("Failed to init adder", "error", err)
			os.Exit(1)
		}
	}()

	defer db.Close()

	router.RouterInit(db)

}
