package filters

import (
	"fmt"
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/config"
	"github.com/Andamio-Platform/andamio-indexer/constants"
	"github.com/Andamio-Platform/andamio-indexer/indexer/eventHandlers"
	filter_chainsync "github.com/blinklabs-io/adder/filter/chainsync"
	filter_event "github.com/blinklabs-io/adder/filter/event"
	input_chainsync "github.com/blinklabs-io/adder/input/chainsync"
	output_embedded "github.com/blinklabs-io/adder/output/embedded"
	"github.com/blinklabs-io/adder/pipeline"
)

func FilterByAndamioAddresses() error {
	// Load config
	cfg := config.GetGlobalConfig()

	// Create pipeline
	p := pipeline.New()

	// Configure pipeline input
	inputOpts := []input_chainsync.ChainSyncOptionFunc{
		input_chainsync.WithBulkMode(true),
		input_chainsync.WithAutoReconnect(true),
		input_chainsync.WithIntersectTip(true),
		// input_chainsync.WithStatusUpdateFunc(UpdateSlotTimestamp),
		input_chainsync.WithNetworkMagic(cfg.Network.Magic),
		input_chainsync.WithKupoUrl(constants.BLINKLABS_KUPO_ENDPOINT),
		input_chainsync.WithIncludeCbor(true),
		// input_chainsync.WithSocketPath(cfg.SocketPath),
		// Use this if you want to connect to a remote node and not SocketPath
		// IOG cardano node
		input_chainsync.WithAddress("preprod-node.play.dev.cardano.org:3001"),
		// input_chainsync.WithIntersectPoints()
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

	// Define address in chainsync filter
	addresses := cfg.Andamio.GetAllAndamioAddresses()

	filterChainsync := filter_chainsync.New(
		filter_chainsync.WithAddresses(addresses),
		filter_chainsync.WithPolicies(cfg.Andamio.GetAllAndamioPolicies()),
		filter_chainsync.WithAssetFingerprints(cfg.Andamio.GetAllAndamioAssetFingerprints()),
	)

	p.AddFilter(filterChainsync)

	// Configure pipeline output
	output := output_embedded.New(
		output_embedded.WithCallbackFunc(eventHandlers.TxEvent),
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
