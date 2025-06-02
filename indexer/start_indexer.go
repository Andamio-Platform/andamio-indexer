package indexer

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"

	"github.com/Andamio-Platform/andamio-indexer/config"
	"github.com/Andamio-Platform/andamio-indexer/database"
	plugin "github.com/Andamio-Platform/andamio-indexer/database/plugin"
	"github.com/Andamio-Platform/andamio-indexer/indexer/cache"
	"github.com/Andamio-Platform/andamio-indexer/indexer/filters"
	"github.com/blinklabs-io/adder/event" // Import the event package
	filter_event "github.com/blinklabs-io/adder/filter/event"
	input_chainsync "github.com/blinklabs-io/adder/input/chainsync"
	output_embedded "github.com/blinklabs-io/adder/output/embedded"
	"github.com/blinklabs-io/adder/pipeline"
	ocommon "github.com/blinklabs-io/gouroboros/protocol/common"
)

func StartIndexer(db *database.Database, logger plugin.Logger) error {
	slog.Info("Starting indexer...")

	// Load config
	cfg := config.GetGlobalConfig()
	slog.Info("Configuration loaded.")

	// Create the CursorStore using the BadgerDB blob store

	// Initialize transaction cache with a limit (e.g., 100)
	cache.InitTransactionCache(cfg.Indexer.TrancactionCacheLimit)
	slog.Info("Transaction cache initialized.", "limit", cfg.Indexer.TrancactionCacheLimit)

	cursorStore := database.NewCursorStore(db)
	slog.Info("Cursor store created.")

	// Create pipeline
	p := pipeline.New()
	slog.Info("Pipeline created.")

	// Configure pipeline input
	inputOpts := []input_chainsync.ChainSyncOptionFunc{
		input_chainsync.WithBulkMode(false),
		input_chainsync.WithAutoReconnect(true),
		input_chainsync.WithIncludeCbor(true),
		input_chainsync.WithLogger(logger),
		input_chainsync.WithStatusUpdateFunc(func(status input_chainsync.ChainSyncStatus) {
			slog.Info("Chain sync status update", "slot", status.SlotNumber, "blockHash", status.BlockHash)
			blockHashBytes := []byte(status.BlockHash)
			if err := cursorStore.UpdateCursor(status.SlotNumber, blockHashBytes); err != nil {
				slog.Error(
					fmt.Sprintf("failed to update cursor: %s", err),
				)
			}
		}),
		input_chainsync.WithNetworkMagic(cfg.Network.Magic),
		// input_chainsync.WithIntersectTip(true),
		// input_chainsync.WithKupoUrl(cfg.Network.BlinklabKupoEndpoint),
		input_chainsync.WithKupoUrl(cfg.Network.LocalKupoEndpoint),
		// input_chainsync.WithSocketPath(cfg.Network.LocalCardanoNodeSocketPath), // we cant use this becsue the code in the addre for this opton is not complete
		input_chainsync.WithAddress(cfg.Network.LocalCardanoNodeEndpoint),
		// input_chainsync.WithAddress(cfg.Network.CFCardanoNodeEndpoint),
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
		slog.Info("Found previous cursor state, intersecting.", "slot", cursorState.SlotNumber, "blockHash", string(cursorState.BlockHash))
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
		hashBytes, err := hex.DecodeString(string(cfg.Indexer.IntercerptHash))
		if err != nil {
			return nil
		}
		inputOpts = append(
			inputOpts,
			input_chainsync.WithIntersectPoints(
				[]ocommon.Point{
					{
						Hash: hashBytes,
						Slot: cfg.Indexer.InterceptSlot,
					},
				},
			))
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
	// Configure pipeline output
	// Create an adapter function to pass the database instance to FilterTxEvent
	filterTxEventAdapter := func(evt event.Event) error {
		return filters.FilterTxEvent(db, evt)
	}

	output := output_embedded.New(
		output_embedded.WithCallbackFunc(filterTxEventAdapter),
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
