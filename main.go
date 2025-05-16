package main

import (
	"flag"
	"fmt"
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

	err := config.Load(cmdlineFlags.configFile)
	if err != nil {
		fmt.Printf("Failed to load config: %s\n", err)
		os.Exit(1)
	}

	logger := slog.Default()

	db, err := database.New(logger, config.GlobalConfig.Database.DatabaseDIR)
	if err != nil {
		fmt.Printf("Failed to init database: %s\n", err)
		os.Exit(1)
	}

	database.SetGlobalDB(db)

	if err := indexer.StartIndexer(); err != nil {
		fmt.Printf("Failed to init adder: %s\n", err)
		os.Exit(1)
	}

	defer db.Close()

	router.RouterInit(db)

}
