package main

import (
	"fmt"
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/config"
	"github.com/blinklabs-io/adder/event"
	filter_chainsync "github.com/blinklabs-io/adder/filter/chainsync"
	filter_event "github.com/blinklabs-io/adder/filter/event"
	input_chainsync "github.com/blinklabs-io/adder/input/chainsync"
	output_embedded "github.com/blinklabs-io/adder/output/embedded"
	"github.com/blinklabs-io/adder/pipeline"
)

func InitAdder() error {
	// Load config
	cfg := config.GetGlobalConfig()

	// Create pipeline
	p := pipeline.New()

	// Configure pipeline input
	inputOpts := []input_chainsync.ChainSyncOptionFunc{
		input_chainsync.WithBulkMode(true),
		input_chainsync.WithAutoReconnect(true),
		input_chainsync.WithIntersectTip(true),
		input_chainsync.WithStatusUpdateFunc(updateStatus),
		input_chainsync.WithNetworkMagic(cfg.Network.Magic),
		// input_chainsync.WithSocketPath(cfg.SocketPath),
		// Use this if you want to connect to a remote node and not SocketPath
		// IOG cardano node
		input_chainsync.WithAddress("preprod-node.play.dev.cardano.org:3001"),
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
	filterChainsync := filter_chainsync.New(
		filter_chainsync.WithAddresses(
			[]string{
				"addr1q93l79hdpvaeqnnmdkshmr4mpjvxnacqxs967keht465tt2dn0z9uhgereqgjsw33ka6c8tu5um7hqsnf5fd50fge9gq4lu2ql",
			},
		),
	)

	p.AddFilter(filterChainsync)

	// Configure pipeline output
	output := output_embedded.New(
		output_embedded.WithCallbackFunc(handleEvent),
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

func handleEvent(evt event.Event) error {
	slog.Info(fmt.Sprintf("Received event: %v\n", evt))
	return nil
}

func updateStatus(status input_chainsync.ChainSyncStatus) {
	slog.Info(fmt.Sprintf("ChainSync status update: %v\n", status))
}
