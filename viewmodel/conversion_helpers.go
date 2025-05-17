package viewmodel

import (
	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
)

// Helper function to convert a slice of models.TransactionInput to a slice of viewmodel.TransactionInput
func ConvertTransactionInputsToViewModels(inputs []models.TransactionInput) []TransactionInput {
	inputViewModels := []TransactionInput{}
	for _, input := range inputs {
		inputViewModels = append(inputViewModels, TransactionInput{
			TransactionHash: string(input.TransactionHash),
			UTxOID:          string(input.UTxOID),
			UTxOIDIndex:     input.UTxOIDIndex,
			Address:         string(input.Address),
			Amount:          input.Amount,
			Cbor:            string(input.Cbor), // CBOR string representation
			// Convert nested Asset models to view models
			Asset: ConvertAssetModelsToViewModels(input.Asset),
			// Convert nested Datum model to view model
			Datum: ConvertDatumModelToViewModel(input.Datum),
		})
	}
	return inputViewModels
}

// Helper function to convert a slice of models.TransactionOutput to a slice of viewmodel.TransactionOutput
func ConvertTransactionOutputsToViewModels(outputs []models.TransactionOutput) []TransactionOutput {
	outputViewModels := []TransactionOutput{}
	for _, output := range outputs {
		outputViewModels = append(outputViewModels, TransactionOutput{
			UTxOID:      string(output.UTxOID),
			UTxOIDIndex: output.UTxOIDIndex,
			Address:     string(output.Address),
			Amount:      output.Amount,
			Cbor:        string(output.Cbor), // CBOR string representation
			// Convert nested Asset models to view models
			Asset: ConvertAssetModelsToViewModels(output.Asset),
			// Convert nested Datum model to view model
			Datum: ConvertDatumModelToViewModel(output.Datum),
		})
	}
	return outputViewModels
}

// Helper function to convert a slice of models.Asset to a slice of viewmodel.Asset
func ConvertAssetModelsToViewModels(assets []models.Asset) []Asset {
	assetViewModels := []Asset{}
	for _, asset := range assets {
		assetViewModels = append(assetViewModels, Asset{
			UTxOID:      string(asset.UTxOID),
			UTxOIDIndex: asset.UTxOIDIndex,
			Name:        string(asset.Name),
			NameHex:     string(asset.NameHex),
			PolicyId:    string(asset.PolicyId),
			Fingerprint: string(asset.Fingerprint),
			Amount:      asset.Amount,
		})
	}
	return assetViewModels
}

// Helper function to convert a models.Datum to a viewmodel.Datum
func ConvertDatumModelToViewModel(datum models.Datum) Datum {
	return Datum{
		UTxOID:      string(datum.UTxOID),
		UTxOIDIndex: datum.UTxOIDIndex,
		DatumHash:   string(datum.DatumHash),
		DatumCbor:   string(datum.DatumCbor), // Assuming DatumCbor should be a string representation of CBOR
	}
}

// Helper function to convert a slice of models.SimpleUTxO to a slice of viewmodel.SimpleUTxO
func ConvertSimpleUTxOModelsToViewModels(utxos []models.SimpleUTxO) []SimpleUTxO {
	utxoViewModels := []SimpleUTxO{}
	for _, utxo := range utxos {
		utxoViewModels = append(utxoViewModels, SimpleUTxO{
			TransactionHash: string(utxo.TransactionHash),
			UTxOID:          string(utxo.UTxOID),
			UTxOIDIndex:     utxo.UTxOIDIndex,
		})
	}
	return utxoViewModels
}

// Helper function to convert a slice of models.Redeemer to a slice of viewmodel.Redeemer
func ConvertRedeemersToViewModels(redeemers []models.Redeemer) []Redeemer {
	redeemerViewModels := []Redeemer{}
	for _, redeemer := range redeemers {
		redeemerViewModels = append(redeemerViewModels, Redeemer{
			TransactionHash: string(redeemer.TransactionHash),
			Index:           redeemer.Index,
			Tag:             redeemer.Tag,
			Cbor:            string(redeemer.Cbor), // CBOR string representation
		})
	}
	return redeemerViewModels
}

// Helper function to convert a models.Witness to a viewmodel.Witness
func ConvertWitnessModelToViewModel(witness models.Witness) Witness {
	return Witness{
		TransactionHash: string(witness.TransactionHash),
		PlutusData:      ConvertByteSliceSliceToStringSlice(witness.PlutusData),      // Convert [][]byte to []string
		PlutusV1Scripts: ConvertByteSliceSliceToStringSlice(witness.PlutusV1Scripts), // Convert [][]byte to []string
		PlutusV2Scripts: ConvertByteSliceSliceToStringSlice(witness.PlutusV2Scripts), // Convert [][]byte to []string
		PlutusV3Scripts: ConvertByteSliceSliceToStringSlice(witness.PlutusV3Scripts), // Convert [][]byte to []string
		Redeemers:       ConvertRedeemersToViewModels(witness.Redeemers),
	}
}

// Helper function to convert a slice of byte slices to a slice of strings
func ConvertByteSliceSliceToStringSlice(byteSlices [][]byte) []string {
	stringSlice := []string{}
	for _, byteSlice := range byteSlices {
		stringSlice = append(stringSlice, string(byteSlice))
	}
	return stringSlice
}
