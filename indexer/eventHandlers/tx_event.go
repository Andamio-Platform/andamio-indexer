package eventHandlers

import (
	"encoding/json"

	"fmt"

	"log/slog"

	"github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
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

func TxEvent(eventTx input_chainsync.TransactionEvent, eventCtx input_chainsync.TransactionContext) error {
	slog.Info("Processing transaction event",
		"txHash", fmt.Sprintf("%x", eventTx.Transaction.Hash().Bytes()),
		"blockNumber", eventCtx.BlockNumber,
		"slotNumber", eventCtx.SlotNumber,
	)

	//! For Debug Purpose
	slog.Info("*********************************************************************")
	slog.Info("*********************************************************************")
	slog.Info("*********************************************************************")
	slog.Info("Transaction Event", "transaction", eventTx.Transaction)
	slog.Info("Transaction Context", "transactionHash", eventCtx.TransactionHash)
	slog.Info("*********************************************************************")
	slog.Info("*********************************************************************")
	slog.Info("*********************************************************************")

	globalDB := database.GetGlobalDB()
	txn := globalDB.Transaction(true)
	slog.Debug("Database transaction started for transaction event.")

	txHash := eventTx.Transaction.Hash().Bytes()

	var inputs []models.TransactionInput
	if len(eventTx.Inputs) != len(eventTx.ResolvedInputs) {
		// This is an unexpected condition, log a warning
		slog.Warn("Mismatch between Inputs and ResolvedInputs length",
			"txHash", fmt.Sprintf("%x", txHash),
			"inputsLength", len(eventTx.Inputs),
			"resolvedInputsLength", len(eventTx.ResolvedInputs),
		)
		return fmt.Errorf("Warning: Mismatch between Inputs and ResolvedInputs length for transaction %x\n", eventTx.Transaction.Hash().Bytes())

	}
	slog.Debug("Processing transaction inputs.", "count", len(eventTx.Inputs))
	for i, input := range eventTx.Inputs {
		resolvedInput := eventTx.ResolvedInputs[i]
		slog.Debug("Processing input", "index", i, "utxoId", fmt.Sprintf("%x", input.Id().Bytes()), "utxoIndex", input.Index())

		inputIdHash := input.Id().Bytes()
		inputIdIndex := input.Index()

		// Convert assets
		var inputAssets []models.Asset
		if resolvedInput.Assets() != nil {
			slog.Debug("Processing input assets.", "inputIndex", i)
			assetData, err := resolvedInput.Assets().MarshalJSON()
			if err != nil {
				slog.Error("failed to marshal input assets to JSON", "error", err)
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
				slog.Debug("Converted input asset", "inputIndex", i, "fingerprint", asset.Fingerprint, "amount", asset.Amount)
			}
			slog.Debug("Finished processing input assets.", "inputIndex", i, "count", len(inputAssets))
		}

		// Convert datum
		var inputDatum models.Datum
		if resolvedInput.Datum() != nil || resolvedInput.DatumHash() != nil {
			slog.Debug("Processing input datum.", "inputIndex", i, "datumHash", fmt.Sprintf("%x", resolvedInput.DatumHash().Bytes()))
			inputDatum = models.Datum{
				UTxOID:      inputIdHash,
				UTxOIDIndex: inputIdIndex,
				DatumHash:   resolvedInput.DatumHash().Bytes(),
				DatumCbor:   resolvedInput.Datum().Cbor(),
			}
			slog.Debug("Converted input datum.", "inputIndex", i)
		}

		inputs = append(inputs, models.TransactionInput{
			TransactionHash: eventTx.Transaction.Hash().Bytes(),
			UTxOID:          inputIdHash,
			UTxOIDIndex:     inputIdIndex,
			Address:         resolvedInput.Address().Bytes(),
			Amount:          resolvedInput.Amount(),
			Asset:           inputAssets,
			Datum:           inputDatum,
			Cbor:            resolvedInput.Cbor(),
		})
		slog.Debug("Appended input to list.", "inputIndex", i)
	}
	slog.Debug("Finished processing transaction inputs.", "count", len(inputs))

	// Process and convert eventTx.Outputs to []models.TransactionOutput
	var outputs []models.TransactionOutput
	slog.Debug("Processing transaction outputs.", "count", len(eventTx.Outputs))
	for i, output := range eventTx.Outputs {
		slog.Debug("Processing output", "index", i)

		outputIdIndex := uint32(i)

		// Convert assets
		var outputAssets []models.Asset
		if output.Assets() != nil {
			slog.Debug("Processing output assets.", "outputIndex", i)
			assetData, err := output.Assets().MarshalJSON()
			if err != nil {
				slog.Error("failed to marshal output assets to JSON", "error", err)
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
				slog.Debug("Converted output asset", "outputIndex", i, "fingerprint", asset.Fingerprint, "amount", asset.Amount)
			}
			slog.Debug("Finished processing output assets.", "outputIndex", i, "count", len(outputAssets))
		}

		// Convert datum
		var outputDatum models.Datum
		if output.Datum() != nil || output.DatumHash() != nil {
			slog.Debug("Processing output datum.", "outputIndex", i, "datumHash", fmt.Sprintf("%x", output.DatumHash().Bytes()))
			outputDatum = models.Datum{
				UTxOID:      txHash,
				UTxOIDIndex: outputIdIndex,
				DatumHash:   output.DatumHash().Bytes(),
				DatumCbor:   output.Datum().Cbor(),
			}
			slog.Debug("Converted output datum.", "outputIndex", i)
		}

		outputs = append(outputs, models.TransactionOutput{
			UTxOID:      txHash,
			UTxOIDIndex: outputIdIndex,
			Address:     output.Address().Bytes(),
			Amount:      output.Amount(),
			Asset:       outputAssets,
			Datum:       outputDatum,
			Cbor:        output.Cbor(),
		})
		slog.Debug("Appended output to list.", "outputIndex", i)
	}
	slog.Debug("Finished processing transaction outputs.", "count", len(outputs))

	// Process and convert eventTx.ReferenceInputs to []models.SimpleUTxO
	var referenceInputs []models.SimpleUTxO
	slog.Debug("Processing transaction reference inputs.", "count", len(eventTx.ReferenceInputs))
	for _, refInput := range eventTx.ReferenceInputs {
		slog.Debug("Processing reference input", "utxoId", fmt.Sprintf("%x", refInput.Id().Bytes()), "utxoIndex", refInput.Index())
		referenceInputs = append(referenceInputs, models.SimpleUTxO{
			TransactionHash: txHash,
			UTxOID:          refInput.Id().Bytes(),
			UTxOIDIndex:     refInput.Index(),
		})
		slog.Debug("Appended reference input to list.")
	}
	slog.Debug("Finished processing transaction reference inputs.", "count", len(referenceInputs))

	// Process and convert eventTx.Witnesses to models.Witness
	var witness models.Witness
	slog.Debug("Processing transaction witnesses.")
	// Convert Redeemers
	var redeemers []models.Redeemer
	redeemersInterface := eventTx.Witnesses.Redeemers()
	redeemerTags := []lcommon.RedeemerTag{
		lcommon.RedeemerTagSpend,
		lcommon.RedeemerTagMint,
		lcommon.RedeemerTagCert,
		lcommon.RedeemerTagReward,
	}
	slog.Debug("Processing redeemers.")
	for _, tag := range redeemerTags {
		indexes := redeemersInterface.Indexes(tag)
		slog.Debug("Processing redeemers for tag", "tag", tag, "count", len(indexes))
		for _, index := range indexes {
			redeemerValue, _ := redeemersInterface.Value(uint(index), tag)
			redeemers = append(redeemers, models.Redeemer{
				Index: uint(index),
				Tag:   uint8(tag), // Convert uint tag to byte slice
				Cbor:  redeemerValue.Cbor(),
			})
			slog.Debug("Converted redeemer", "index", index, "tag", tag)
		}
	}
	slog.Debug("Finished processing redeemers.", "count", len(redeemers))

	var plutusData [][]byte
	slog.Debug("Processing plutus data.", "count", len(eventTx.Witnesses.PlutusData()))
	for _, pd := range eventTx.Witnesses.PlutusData() {
		plutusData = append(plutusData, pd.Cbor())
		slog.Debug("Converted plutus data.")
	}
	slog.Debug("Finished processing plutus data.", "count", len(plutusData))

	witness = models.Witness{
		TransactionHash: eventTx.Transaction.Hash().Bytes(),
		PlutusData:      plutusData,
		PlutusV1Scripts: eventTx.Witnesses.PlutusV1Scripts(),
		PlutusV2Scripts: eventTx.Witnesses.PlutusV2Scripts(),
		PlutusV3Scripts: eventTx.Witnesses.PlutusV3Scripts(),
		Redeemers:       redeemers, // Use the converted redeemers
	}
	slog.Debug("Witness data processed.")

	// Process and convert eventTx.Certificates to [][]byte
	var certificates [][]byte
	slog.Debug("Processing transaction certificates.", "count", len(eventTx.Certificates))
	for _, cert := range eventTx.Certificates {
		certificates = append(certificates, cert.Cbor())
		slog.Debug("Converted certificate.")
	}
	slog.Debug("Finished processing transaction certificates.", "count", len(certificates))

	slog.Info("Saving transaction to database.", "txHash", fmt.Sprintf("%x", txHash))
	err := globalDB.NewTx(
		[]byte(eventTx.BlockHash),
		eventCtx.BlockNumber,
		eventCtx.SlotNumber,
		eventTx.Transaction.Hash().Bytes(),
		inputs,
		outputs,
		referenceInputs,
		eventTx.Metadata.Cbor(),
		eventTx.Fee,
		eventTx.TTL,
		eventTx.Withdrawals,
		witness,
		certificates,
		eventTx.Transaction.Cbor(),
		txn,
	)
	if err != nil {
		slog.Error("Failed to save transaction to database.", "txHash", fmt.Sprintf("%x", txHash), "error", err)
		return err
	}
	slog.Info("Transaction saved to database successfully.", "txHash", fmt.Sprintf("%x", txHash))

	return nil

}
