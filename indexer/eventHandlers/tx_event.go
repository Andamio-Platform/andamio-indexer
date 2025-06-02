package eventHandlers

import (
	"encoding/json"

	"fmt"

	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
	"github.com/Andamio-Platform/andamio-indexer/database/types"
	input_chainsync "github.com/blinklabs-io/adder/input/chainsync"
	lcommon "github.com/blinklabs-io/gouroboros/ledger/common"
	fiberLogger "github.com/gofiber/fiber/v2/log"
)

type Asset struct {
	Name        string `json:"name"`
	NameHex     string `json:"nameHex"`
	PolicyId    string `json:"policyId"`
	Fingerprint string `json:"fingerprint"`
	Amount      uint64 `json:"amount"`
}

func TxEvent(logger *slog.Logger, eventTx input_chainsync.TransactionEvent, eventCtx input_chainsync.TransactionContext, txn *database.Txn) error {
	logger.Info("Processing transaction event",
		"txHash", fmt.Sprintf("%x", eventTx.Transaction.Hash().Bytes()),
		"blockNumber", eventCtx.BlockNumber,
		"slotNumber", eventCtx.SlotNumber,
	)

	//! For Debug Purpose
	// logger.Info("*********************************************************************")
	// logger.Info("*********************************************************************")
	// logger.Info("*********************************************************************")
	// logger.Info("Transaction Event", "transaction", eventTx.Transaction)
	// logger.Info("Transaction Context", "transactionHash", eventCtx.TransactionHash)
	// logger.Info("*********************************************************************")
	// logger.Info("*********************************************************************")
	// logger.Info("*********************************************************************")

	txHash := eventTx.Transaction.Hash().Bytes()

	var inputs []models.TransactionInput
	if len(eventTx.Inputs) != len(eventTx.ResolvedInputs) {
		// This is an unexpected condition, log a warning
		logger.Warn("Mismatch between Inputs and ResolvedInputs length",
			"txHash", fmt.Sprintf("%x", txHash),
			"inputsLength", len(eventTx.Inputs),
			"resolvedInputsLength", len(eventTx.ResolvedInputs),
		)
		return fmt.Errorf("Warning: Mismatch between Inputs and ResolvedInputs length for transaction %x\n", eventTx.Transaction.Hash().Bytes())

	}
	logger.Debug("Processing transaction inputs.", "count", len(eventTx.Inputs))
	for i, input := range eventTx.Inputs {
		resolvedInput := eventTx.ResolvedInputs[i]
		logger.Debug("Processing input", "index", i, "utxoId", fmt.Sprintf("%x", input.Id().Bytes()), "utxoIndex", input.Index())

		inputIdHash := input.Id().Bytes()
		inputIdIndex := input.Index()

		// Convert assets
		var inputAssets []models.Asset
		if resolvedInput.Assets() != nil {
			logger.Debug("Processing input assets.", "inputIndex", i)
			assetData, err := resolvedInput.Assets().MarshalJSON()
			if err != nil {
				logger.Error("failed to marshal input assets to JSON", "error", err)
				return fmt.Errorf("failed to marshal transaction body to JSON: %v", err)
			}
			var assets []Asset
			err = json.Unmarshal(assetData, &assets)
			if err != nil {
				fiberLogger.Error("failed to unmarshal transaction body to JSON: %v", err)

			}
			for _, asset := range assets {
				inputAssets = append(inputAssets, models.Asset{
					UTxOID:      inputIdHash,
					UTxOIDIndex: inputIdIndex,
					Name:        []byte(asset.Name),
					NameHex:     []byte(asset.NameHex),
					PolicyId:    []byte(asset.PolicyId),
					Fingerprint: []byte(asset.Fingerprint),
					Amount:      uint64(asset.Amount),
				})
				logger.Debug("Converted input asset", "inputIndex", i, "fingerprint", asset.Fingerprint, "amount", asset.Amount, "fingerprint_string_value", asset.Fingerprint)
			}
			logger.Debug("Finished processing input assets.", "inputIndex", i, "count", len(inputAssets))
		}

		// Convert datum
		var inputDatum models.Datum
		if resolvedInput.Datum() != nil || resolvedInput.DatumHash() != nil {
			logger.Debug("Processing input datum.", "inputIndex", i, "datumHash", fmt.Sprintf("%x", resolvedInput.DatumHash().Bytes()))
			inputDatum = models.Datum{
				UTxOID:      inputIdHash,
				UTxOIDIndex: inputIdIndex,
				DatumHash:   resolvedInput.DatumHash().Bytes(),
				DatumCbor:   resolvedInput.Datum().Cbor(),
			}
			logger.Debug("Converted input datum.", "inputIndex", i)
		}

		inputs = append(inputs, models.TransactionInput{
			TransactionHash: eventTx.Transaction.Hash().Bytes(),
			UTxOID:          inputIdHash,
			UTxOIDIndex:     inputIdIndex,
			Address:         []byte(resolvedInput.Address().String()),
			Amount:          resolvedInput.Amount(),
			Asset:           inputAssets,
			Datum:           inputDatum,
			Cbor:            resolvedInput.Cbor(),
		})
		logger.Debug("Appended input to list.", "inputIndex", i)
	}
	logger.Debug("Finished processing transaction inputs.", "count", len(inputs))

	// Process and convert eventTx.Outputs to []models.TransactionOutput
	var outputs []models.TransactionOutput
	logger.Debug("Starting output processing loop.")
	logger.Debug("Processing transaction outputs.", "count", len(eventTx.Outputs))
	for i, output := range eventTx.Outputs {
		logger.Debug("Processing output", "index", i)

		outputIdIndex := uint32(i)

		// Convert assets
		var outputAssets []models.Asset
		if output.Assets() != nil {
			logger.Debug("Processing output assets.", "outputIndex", i)
			assetData, err := output.Assets().MarshalJSON()
			if err != nil {
				logger.Error("failed to marshal output assets to JSON", "error", err)
				return fmt.Errorf("failed to unmarshal transaction body to JSON: %v", err)
			}
			var assets []Asset
			err = json.Unmarshal(assetData, &assets)
			if err != nil {
				fiberLogger.Error("failed to unmarshal transaction body to JSON: %v", err)

			}
			for _, asset := range assets {
				outputAssets = append(outputAssets, models.Asset{
					UTxOID:      txHash,
					UTxOIDIndex: outputIdIndex,
					Name:        []byte(asset.Name),
					NameHex:     []byte(asset.NameHex),
					PolicyId:    []byte(asset.PolicyId),
					Fingerprint: []byte(asset.Fingerprint),
					Amount:      uint64(asset.Amount),
				})
				logger.Debug("Converted output asset", "outputIndex", i, "fingerprint", asset.Fingerprint, "amount", asset.Amount, "fingerprint_string_value", asset.Fingerprint)
			}
			logger.Debug("Finished processing output assets.", "outputIndex", i, "count", len(outputAssets))
		}

		// Convert datum
		var outputDatum models.Datum
		if output == nil {
			logger.Error("Transaction output is nil", "outputIndex", i)
		} else if output.DatumHash() != nil && output.Datum() != nil { // Only create Datum if both DatumHash and Datum are not nil
			var datumHashBytes []byte
			if output.DatumHash() != nil {
				datumHashBytes = output.DatumHash().Bytes()
			}
			var datumCborBytes []byte
			if output.Datum() != nil {
				datumCborBytes = output.Datum().Cbor()
			}
			logger.Debug("Processing output datum.", "outputIndex", i, "datumHash", fmt.Sprintf("%x", datumHashBytes))
			outputDatum = models.Datum{
				UTxOID:      txHash,
				UTxOIDIndex: outputIdIndex,
				DatumHash:   datumHashBytes,
				DatumCbor:   datumCborBytes,
			}
			logger.Debug("Converted output datum.", "outputIndex", i)
		}

		outputs = append(outputs, models.TransactionOutput{
			UTxOID:      txHash,
			UTxOIDIndex: outputIdIndex,
			Address:     []byte(output.Address().String()),
			Amount:      output.Amount(),
			Asset:       outputAssets,
			Datum:       outputDatum,
			Cbor:        output.Cbor(),
		})
		logger.Debug("Appended output to list.", "outputIndex", i)
	}
	logger.Debug("Finished processing transaction outputs.", "count", len(outputs))
	// Process and convert eventTx.ReferenceInputs to []models.SimpleUTxO
	var referenceInputs []models.SimpleUTxO
	logger.Debug("Processing transaction reference inputs.", "count", len(eventTx.ReferenceInputs))
	for _, refInput := range eventTx.ReferenceInputs {
		logger.Debug("Processing reference input", "utxoId", fmt.Sprintf("%x", refInput.Id().Bytes()), "utxoIndex", refInput.Index())
		referenceInputs = append(referenceInputs, models.SimpleUTxO{
			TransactionHash: txHash,
			UTxOID:          refInput.Id().Bytes(),
			UTxOIDIndex:     refInput.Index(),
		})
		logger.Debug("Appended reference input to list.")
	}
	logger.Debug("Finished processing transaction reference inputs.", "count", len(referenceInputs))

	// Process and convert eventTx.Witnesses to models.Witness
	var witness models.Witness
	logger.Debug("Processing transaction witnesses.")
	// Convert Redeemers
	var redeemers []models.Redeemer
	redeemersInterface := eventTx.Witnesses.Redeemers()
	redeemerTags := []lcommon.RedeemerTag{
		lcommon.RedeemerTagSpend,
		lcommon.RedeemerTagMint,
		lcommon.RedeemerTagCert,
		lcommon.RedeemerTagReward,
	}
	logger.Debug("Processing redeemers.")
	for _, tag := range redeemerTags {
		indexes := redeemersInterface.Indexes(tag)
		logger.Debug("Processing redeemers for tag", "tag", tag, "count", len(indexes))
		for _, index := range indexes {
			redeemerValue, _ := redeemersInterface.Value(uint(index), tag)
			redeemers = append(redeemers, models.Redeemer{
				Index: uint(index),
				Tag:   uint8(tag), // Convert uint tag to byte slice
				Cbor:  redeemerValue.Cbor(),
			})
			logger.Debug("Converted redeemer", "index", index, "tag", tag)
		}
	}
	logger.Debug("Finished processing redeemers.", "count", len(redeemers))

	var plutusData [][]byte
	logger.Debug("Processing plutus data.", "count", len(eventTx.Witnesses.PlutusData()))
	for _, pd := range eventTx.Witnesses.PlutusData() {
		plutusData = append(plutusData, pd.Cbor())
		logger.Debug("Converted plutus data.")
	}
	logger.Debug("Finished processing plutus data.", "count", len(plutusData))

	witness = models.Witness{
		TransactionHash: eventTx.Transaction.Hash().Bytes(),
		PlutusData:      plutusData,
		PlutusV1Scripts: types.ByteSliceSlice(eventTx.Witnesses.PlutusV1Scripts()),
		PlutusV2Scripts: types.ByteSliceSlice(eventTx.Witnesses.PlutusV2Scripts()),
		PlutusV3Scripts: types.ByteSliceSlice(eventTx.Witnesses.PlutusV3Scripts()),
		Redeemers:       redeemers, // Use the converted redeemers
	}
	logger.Debug("Witness data processed.")

	// Process and convert eventTx.Certificates to [][]byte
	var certificates [][]byte
	logger.Debug("Processing transaction certificates.", "count", len(eventTx.Certificates))
	for _, cert := range eventTx.Certificates {
		certificates = append(certificates, cert.Cbor())
		logger.Debug("Converted certificate.")
	}
	logger.Debug("Finished processing transaction certificates.", "count", len(certificates))

	logger.Info("Saving transaction to database.", "txHash", fmt.Sprintf("%x", txHash))

	err := txn.DB().NewTx(
		[]byte(eventTx.BlockHash),
		eventCtx.BlockNumber,
		eventCtx.SlotNumber,
		eventTx.Transaction.Hash().Bytes(),
		inputs,
		outputs,
		referenceInputs,
		// Check if Metadata is nil before accessing Cbor()
		func() []byte {
			if eventTx.Metadata == nil {
				logger.Debug("eventTx.Metadata is nil")
				return nil
			}
			logger.Debug("eventTx.Metadata is not nil")
			return eventTx.Metadata.Cbor()
		}(),
		eventTx.Fee,
		eventTx.TTL,
		eventTx.Withdrawals,
		witness,
		certificates,
		eventTx.Transaction.Cbor(),
		txn,
	)
	logger.Debug("Transaction CBOR length before saving", "length", len(eventTx.Transaction.Cbor()))
	if err != nil {
		logger.Error("Failed to save transaction to database.", "txHash", fmt.Sprintf("%x", txHash), "error", err)
		return err
	}
	logger.Info("Transaction saved to database successfully.", "txHash", fmt.Sprintf("%x", txHash))

	return nil

}
