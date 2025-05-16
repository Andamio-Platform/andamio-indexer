package eventHandlers

import (
	"encoding/json"
	"fmt"

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

	globalDB := database.GetGlobalDB()
	txn := globalDB.Transaction(true)

	txHash := eventTx.Transaction.Hash().Bytes()

	var inputs []models.TransactionInput
	if len(eventTx.Inputs) != len(eventTx.ResolvedInputs) {
		// This is an unexpected condition, log a warning
		return fmt.Errorf("Warning: Mismatch between Inputs and ResolvedInputs length for transaction %x\n", eventTx.Transaction.Hash().Bytes())

	}

	for i, input := range eventTx.Inputs {
		resolvedInput := eventTx.ResolvedInputs[i]

		inputIdHash := input.Id().Bytes()
		inputIdIndex := input.Index()

		// Convert assets
		var inputAssets []models.Asset
		if resolvedInput.Assets() != nil {
			assetData, err := resolvedInput.Assets().MarshalJSON()
			if err != nil {
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
			}
		}

		// Convert datum
		var inputDatum models.Datum
		if resolvedInput.Datum() != nil || resolvedInput.DatumHash() != nil {
			inputDatum = models.Datum{
				UTxOID:      inputIdHash,
				UTxOIDIndex: inputIdIndex,
				DatumHash:   resolvedInput.DatumHash().Bytes(),
				DatumCbor:   resolvedInput.Datum().Cbor(),
			}
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
	}

	// Process and convert eventTx.Outputs to []models.TransactionOutput
	var outputs []models.TransactionOutput
	for i, output := range eventTx.Outputs {

		outputIdIndex := uint32(i)

		// Convert assets
		var outputAssets []models.Asset
		if output.Assets() != nil {
			assetData, err := output.Assets().MarshalJSON()
			if err != nil {
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
			}
		}

		// Convert datum
		var outputDatum models.Datum
		if output.Datum() != nil || output.DatumHash() != nil {
			outputDatum = models.Datum{
				UTxOID:      txHash,
				UTxOIDIndex: outputIdIndex,
				DatumHash:   output.DatumHash().Bytes(),
				DatumCbor:   output.Datum().Cbor(),
			}
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
	}

	// Process and convert eventTx.ReferenceInputs to []models.SimpleUTxO
	var referenceInputs []models.SimpleUTxO
	for _, refInput := range eventTx.ReferenceInputs {
		referenceInputs = append(referenceInputs, models.SimpleUTxO{
			TransactionHash: txHash,
			UTxOID:          refInput.Id().Bytes(),
			UTxOIDIndex:     refInput.Index(),
		})
	}

	// Process and convert eventTx.Witnesses to models.Witness
	var witness models.Witness
	// Convert Redeemers
	var redeemers []models.Redeemer
	redeemersInterface := eventTx.Witnesses.Redeemers()
	redeemerTags := []lcommon.RedeemerTag{
		lcommon.RedeemerTagSpend,
		lcommon.RedeemerTagMint,
		lcommon.RedeemerTagCert,
		lcommon.RedeemerTagReward,
	}
	for _, tag := range redeemerTags {
		indexes := redeemersInterface.Indexes(tag)
		for _, index := range indexes {
			redeemerValue, _ := redeemersInterface.Value(uint(index), tag)
			redeemers = append(redeemers, models.Redeemer{
				Index: uint(index),
				Tag:   uint8(tag), // Convert uint tag to byte slice
				Cbor:  redeemerValue.Cbor(),
			})
		}
	}

	var plutusData [][]byte
	for _, pd := range eventTx.Witnesses.PlutusData() {
		plutusData = append(plutusData, pd.Cbor())
	}

	witness = models.Witness{
		TransactionHash: eventTx.Transaction.Hash().Bytes(),
		PlutusData:      plutusData,
		PlutusV1Scripts: eventTx.Witnesses.PlutusV1Scripts(),
		PlutusV2Scripts: eventTx.Witnesses.PlutusV2Scripts(),
		PlutusV3Scripts: eventTx.Witnesses.PlutusV3Scripts(),
		Redeemers:       redeemers, // Use the converted redeemers
	}

	// Process and convert eventTx.Certificates to [][]byte
	var certificates [][]byte
	for _, cert := range eventTx.Certificates {
		certificates = append(certificates, cert.Cbor())
	}

	err := globalDB.NewTx(
		eventTx.BlockHash,
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
		return err
	}

	return nil

}
