package indexer

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"

	"github.com/Andamio-Platform/andamio-indexer/config"
	"github.com/Andamio-Platform/andamio-indexer/constants"
	"github.com/Andamio-Platform/andamio-indexer/database"
	plugin "github.com/Andamio-Platform/andamio-indexer/database/plugin"
	"github.com/Andamio-Platform/andamio-indexer/indexer/cache"
	"github.com/Andamio-Platform/andamio-indexer/indexer/filters"
	filter_event "github.com/blinklabs-io/adder/filter/event"
	input_chainsync "github.com/blinklabs-io/adder/input/chainsync"
	output_embedded "github.com/blinklabs-io/adder/output/embedded"
	"github.com/blinklabs-io/adder/pipeline"
	ocommon "github.com/blinklabs-io/gouroboros/protocol/common"
)

func StartIndexer(logger plugin.Logger) error {
	slog.Info("Starting indexer...")

	// Load config
	cfg := config.GetGlobalConfig()
	slog.Info("Configuration loaded.")

	// Create the CursorStore using the BadgerDB blob store
	globalDB := database.GetGlobalDB()

	// Initialize transaction cache with a limit (e.g., 100)
	cache.InitTransactionCache(100)
	slog.Info("Transaction cache initialized.", "limit", 100)

	cursorStore := database.NewCursorStore(globalDB)
	slog.Info("Cursor store created.")

	// Create pipeline
	p := pipeline.New()
	slog.Info("Pipeline created.")

	// Configure pipeline input
	inputOpts := []input_chainsync.ChainSyncOptionFunc{
		input_chainsync.WithBulkMode(true),
		input_chainsync.WithAutoReconnect(true),
		input_chainsync.WithIntersectTip(false),
		input_chainsync.WithLogger(logger),
		input_chainsync.WithStatusUpdateFunc(func(status input_chainsync.ChainSyncStatus) {
			slog.Info("Chain sync status update", "slot", status.SlotNumber, "blockHash", status.BlockHash)
			// Use the cursorStore instance to update the cursor
			blockHashBytes := []byte(status.BlockHash)
			if err := cursorStore.UpdateCursor(status.SlotNumber, blockHashBytes); err != nil {
				slog.Error(
					fmt.Sprintf("failed to update cursor: %s", err),
				)
			}
		}),
		input_chainsync.WithNetworkMagic(cfg.Network.Magic),
		input_chainsync.WithKupoUrl(constants.BLINKLABS_KUPO_ENDPOINT),
		input_chainsync.WithIncludeCbor(true),
		// input_chainsync.WithSocketPath(cfg.Network.SocketPath),
		input_chainsync.WithAddress("preprod-node.play.dev.cardano.org:3001"),
	}

	// Get the last saved cursor state using the cursorStore instance
	slog.Info("Fetching cursor state...")
	cursorState, err := cursorStore.GetCursor() // Use cursorStore.GetCursor
	if err != nil {
		return fmt.Errorf("failed to get cursor state: %w", err)
	}

	// Use the retrieved cursor state for intersection
	// Check if a cursor state was found (SlotNumber will be non-zero if found)
	if cursorState.SlotNumber > 0 || len(cursorState.BlockHash) > 0 {
		slog.Info("Found previous cursor state, intersecting.", "slot", cursorState.SlotNumber, "blockHash", cursorState.BlockHash)
		hashBytes, err := hex.DecodeString(string(cursorState.BlockHash))
		if err != nil {
			return nil
		}
		inputOpts = append(
			inputOpts,
			input_chainsync.WithIntersectPoints(
				[]ocommon.Point{
					{
						Hash: hashBytes,
						Slot: cursorState.SlotNumber,
					},
				},
			))
	} else {

		// If no cursor state is found, intersect at the tip (already handled by WithIntersectTip)
		slog.Info("No previous cursor state found, starting from Andamio genesis.")
	}

	input := input_chainsync.New(
		inputOpts...,
	)
	p.AddInput(input)
	slog.Info("Input configured.")

	// Define type in event filter
	filterEvent := filter_event.New(
		filter_event.WithTypes([]string{"chainsync.transaction"}),
	)
	// Add event filter to pipeline
	p.AddFilter(filterEvent)
	slog.Info("Event filter configured.")

	// Configure pipeline output
	output := output_embedded.New(
		output_embedded.WithCallbackFunc(filters.FilterTxEvent),
	)

	//! For Debug Purpose
	// output := output_embedded.New(
	// 	output_embedded.WithCallbackFunc(func(evt event.Event) error {
	// 		fmt.Println(evt)
	// 		return nil
	// 	},
	// 	),
	// )

	p.AddOutput(output)
	slog.Info("Output configured.")

	// Start pipeline
	slog.Info("Starting pipeline...")
	if err := p.Start(); err != nil {
		slog.Info(fmt.Sprintf("failed to start pipeline: %s\n", err))
	}

	// Start error handler
	for {
		err, ok := <-p.ErrorChan()
		if ok {
			slog.Error("Pipeline error received", "error", err)
			slog.Info(fmt.Sprintf("pipeline failed: %v\n", err)) // Keep existing log for context
			os.Exit(1)                                           // Exit the application on pipeline error
		} else {
			slog.Info("Pipeline error channel closed.")
			break
		}
	}
	return nil

}
