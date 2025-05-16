package indexer

import (
	"fmt"
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/config"
	"github.com/Andamio-Platform/andamio-indexer/constants"
	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/indexer/cache" // Import cache
	"github.com/Andamio-Platform/andamio-indexer/indexer/filters"
	filter_event "github.com/blinklabs-io/adder/filter/event"
	input_chainsync "github.com/blinklabs-io/adder/input/chainsync"
	output_embedded "github.com/blinklabs-io/adder/output/embedded"
	"github.com/blinklabs-io/adder/pipeline"
	ocommon "github.com/blinklabs-io/gouroboros/protocol/common"
)

func StartIndexer() error {
	// Load config
	cfg := config.GetGlobalConfig()

	// Create the CursorStore using the BadgerDB blob store
	globalDB := database.GetGlobalDB()

	// Initialize transaction cache with a limit (e.g., 100)
	cache.InitTransactionCache(100)

	cursorStore := database.NewCursorStore(globalDB)

	// Create pipeline
	p := pipeline.New()

	// Configure pipeline input
	inputOpts := []input_chainsync.ChainSyncOptionFunc{
		input_chainsync.WithBulkMode(true),
		input_chainsync.WithAutoReconnect(true),
		input_chainsync.WithIntersectTip(true),
		input_chainsync.WithStatusUpdateFunc(func(status input_chainsync.ChainSyncStatus) {
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
		// input_chainsync.WithSocketPath(cfg.SocketPath),
		input_chainsync.WithAddress("preprod-node.play.dev.cardano.org:3001"),
	}

	// Get the last saved cursor state using the cursorStore instance
	cursorState, err := cursorStore.GetCursor() // Use cursorStore.GetCursor
	if err != nil {
		return fmt.Errorf("failed to get cursor state: %w", err)
	}

	// Use the retrieved cursor state for intersection
	// Check if a cursor state was found (SlotNumber will be non-zero if found)
	if cursorState.SlotNumber > 0 || len(cursorState.BlockHash) > 0 {
		inputOpts = append(
			inputOpts,
			input_chainsync.WithIntersectPoints(
				[]ocommon.Point{
					{
						Hash: cursorState.BlockHash,
						Slot: cursorState.SlotNumber,
					},
				},
			))
	} else {
		// If no cursor state is found, intersect at the tip (already handled by WithIntersectTip)
		slog.Info("no previous cursor state found, starting from tip")
	}

	input := input_chainsync.New(
		inputOpts...,
	)
	p.AddInput(input)

	// Define type in event filter
	filterEvent := filter_event.New(
		filter_event.WithTypes([]string{"chainsync.transaction"}),
	)
	// Add event filter to pipeline
	p.AddFilter(filterEvent)

	// Configure pipeline output
	output := output_embedded.New(
		output_embedded.WithCallbackFunc(filters.FilterTxEvent),
	)

	p.AddOutput(output)

	// Start pipeline
	if err := p.Start(); err != nil {
		slog.Info(fmt.Sprintf("failed to start pipeline: %s\n", err))
	}

	// Start error handler
	for {
		err, ok := <-p.ErrorChan()
		if ok {
			slog.Info(fmt.Sprintf("pipeline failed: %v\n", err))
		} else {
			break
		}
	}
	return nil

}
