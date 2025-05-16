package filters

import (
	"github.com/Andamio-Platform/andamio-indexer/indexer/cache" // Import the cache package
	eventHandlers "github.com/Andamio-Platform/andamio-indexer/indexer/eventHandlers" // Import the eventHandlers package
	"github.com/blinklabs-io/adder/event"
	input_chainsync "github.com/blinklabs-io/adder/input/chainsync"
)

func FilterTxEvent(evt event.Event) error {
	if evt.Type == "chainsync.transaction" {

		eventTx := evt.Payload.(input_chainsync.TransactionEvent)
		eventCtx := evt.Context.(input_chainsync.TransactionContext)

		// Get relevant data from cache
		relevantDataCache := cache.GetRelevantDataCache()
		relevantAddresses := relevantDataCache.GetAddresses()
		relevantPolicies := relevantDataCache.GetPolicies()

		shouldProcess := false

		// Check inputs
		for _, input := range eventTx.ResolvedInputs {
			// Check address
			inputAddress := input.Address().String()
			for _, addr := range relevantAddresses {
				if inputAddress == addr {
					shouldProcess = true
					break
				}
			}
			if shouldProcess {
				break
			}

			// Check assets policy IDs
			if input.Assets() != nil {
				for _, policyId := range input.Assets().Policies() { // Iterate over policy IDs
					for _, relevantPolicy := range relevantPolicies {
						if policyId.String() == relevantPolicy {
							shouldProcess = true
							break
						}
					}
					if shouldProcess {
						break
					}
				}
			}
		}

		if !shouldProcess {
			// Check outputs if not already marked for processing by inputs
			for _, output := range eventTx.Outputs {
				// Check address
				outputAddress := output.Address().String()
				for _, addr := range relevantAddresses {
					if outputAddress == addr {
						shouldProcess = true
						break
					}
				}
				if shouldProcess {
					break
				}

				// Check assets policy IDs
				if output.Assets() != nil {
					for _, policyId := range output.Assets().Policies() { // Iterate over policy IDs
						for _, relevantPolicy := range relevantPolicies {
							if policyId.String() == relevantPolicy {
								shouldProcess = true
								break
							}
						}
						if shouldProcess {
							break
						}
					}
				}
			}
		}

		// If the transaction meets filtering criteria and has a certificate, add to batch
		if shouldProcess {
			eventHandlers.AddToTransactionBatch(eventTx, eventCtx)
		}
	}
	return nil // Return nil if the event is not a transaction or if filtering passes without error
}
