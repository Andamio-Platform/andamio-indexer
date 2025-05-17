package filters

import (
	"fmt"
	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/indexer/cache"                       // Import the cache package
	eventHandlers "github.com/Andamio-Platform/andamio-indexer/indexer/eventHandlers" // Import the eventHandlers package
	"github.com/blinklabs-io/adder/event"
	input_chainsync "github.com/blinklabs-io/adder/input/chainsync"
)

func FilterTxEvent(evt event.Event) error {
	slog.Debug("Received event from pipeline output", "eventType", evt.Type)
	if evt.Type == "chainsync.transaction" {
		slog.Debug("Processing chainsync.transaction event")

		eventTx := evt.Payload.(input_chainsync.TransactionEvent)
		eventCtx := evt.Context.(input_chainsync.TransactionContext)

		// Get relevant data from cache
		relevantDataCache := cache.GetRelevantDataCache()
		relevantAddresses := relevantDataCache.GetAddresses()
		relevantPolicies := relevantDataCache.GetPolicies()
		slog.Debug("Retrieved relevant data from cache.", "addressesCount", len(relevantAddresses), "policiesCount", len(relevantPolicies))

		shouldProcess := false

		// Check inputs
		slog.Debug("Checking transaction inputs for relevance.", "inputsCount", len(eventTx.ResolvedInputs))
		for _, input := range eventTx.ResolvedInputs {
			// Check address
			inputAddress := input.Address().String()
			slog.Debug("Checking input address", "address", inputAddress)
			for _, addr := range relevantAddresses {
				if inputAddress == addr {
					slog.Debug("Input address is relevant.", "address", inputAddress)
					shouldProcess = true
					break
				}
			}
			if shouldProcess {
				slog.Debug("Transaction marked for processing based on input address.")
				break
			}

			// Check assets policy IDs
			if input.Assets() != nil {
				slog.Debug("Checking input assets policy IDs.")
				for _, policyId := range input.Assets().Policies() { // Iterate over policy IDs
					slog.Debug("Checking input asset policy ID", "policyId", policyId.String())
					for _, relevantPolicy := range relevantPolicies {
						if policyId.String() == relevantPolicy {
							slog.Debug("Input asset policy ID is relevant.", "policyId", policyId.String())
							shouldProcess = true
							break
						}
					}
					if shouldProcess {
						slog.Debug("Transaction marked for processing based on input asset policy ID.")
						break
					}
				}
			}
		}
		slog.Debug("Finished checking transaction inputs.", "shouldProcess", shouldProcess)

		if !shouldProcess {
			// Check outputs if not already marked for processing by inputs
			slog.Debug("Checking transaction outputs for relevance.", "outputsCount", len(eventTx.Outputs))
			for _, output := range eventTx.Outputs {
				// Check address
				outputAddress := output.Address().String()
				slog.Debug("Checking output address", "address", outputAddress)
				for _, addr := range relevantAddresses {
					if outputAddress == addr {
						slog.Debug("Output address is relevant.", "address", outputAddress)
						shouldProcess = true
						break
					}
				}
				if shouldProcess {
					slog.Debug("Transaction marked for processing based on output address.")
					break
				}

				// Check assets policy IDs
				if output.Assets() != nil {
					slog.Debug("Checking output assets policy IDs.")
					for _, policyId := range output.Assets().Policies() { // Iterate over policy IDs
						slog.Debug("Checking output asset policy ID", "policyId", policyId.String())
						for _, relevantPolicy := range relevantPolicies {
							if policyId.String() == relevantPolicy {
								slog.Debug("Output asset policy ID is relevant.", "policyId", policyId.String())
								shouldProcess = true
								break
							}
						}
						if shouldProcess {
							slog.Debug("Transaction marked for processing based on output asset policy ID.")
							break
						}
					}
				}
			}
			slog.Debug("Finished checking transaction outputs.", "shouldProcess", shouldProcess)
		}

		// If the transaction meets filtering criteria and has a certificate, add to batch
		if shouldProcess {
			slog.Info("Transaction meets filtering criteria, adding to batch.", "txHash", fmt.Sprintf("%x", eventTx.Transaction.Hash().Bytes()))
			eventHandlers.AddToTransactionBatch(eventTx, eventCtx)
		} else {
			slog.Debug("Transaction does not meet filtering criteria, skipping.", "txHash", fmt.Sprintf("%x", eventTx.Transaction.Hash().Bytes()))
		}
	} else {
		slog.Debug("Event is not a chainsync.transaction, skipping.", "eventType", evt.Type)
	}
	return nil // Return nil if the event is not a transaction or if filtering passes without error
}
