package viewmodel

import (
	"encoding/hex"

	"github.com/Andamio-Platform/andamio-indexer/database/plugin/metadata/sqlite/models"
)

// Helper function to convert a slice of models.TransactionInput to a slice of viewmodel.TransactionInput
func ConvertTransactionInputsToViewModels(inputs []models.TransactionInput) []TransactionInput {
	inputViewModels := []TransactionInput{}
	for _, input := range inputs {
		inputViewModels = append(inputViewModels, TransactionInput{
			TransactionHash: hex.EncodeToString(input.TransactionHash),
			UTxOID:          hex.EncodeToString(input.UTxOID),
			UTxOIDIndex:     input.UTxOIDIndex,
			Address:         string(input.Address),
			Amount:          input.Amount,
			Cbor:            hex.EncodeToString(input.Cbor),
			Asset: ConvertAssetModelsToViewModels(input.Asset),
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
			UTxOID:      hex.EncodeToString(output.UTxOID),
			UTxOIDIndex: output.UTxOIDIndex,
			Address:     string(output.Address),
			Amount:      output.Amount,
			Cbor:        hex.EncodeToString(output.Cbor),
			Asset: ConvertAssetModelsToViewModels(output.Asset),
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
			UTxOID:      hex.EncodeToString(asset.UTxOID),
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
		UTxOID:      hex.EncodeToString(datum.UTxOID),
		UTxOIDIndex: datum.UTxOIDIndex,
		DatumHash:   hex.EncodeToString(datum.DatumHash),
		DatumCbor:   hex.EncodeToString(datum.DatumCbor), // Assuming DatumCbor should be a string representation of CBOR
	}
}

// Helper function to convert a slice of models.SimpleUTxO to a slice of viewmodel.SimpleUTxO
func ConvertSimpleUTxOModelsToViewModels(utxos []models.SimpleUTxO) []SimpleUTxO {
	utxoViewModels := []SimpleUTxO{}
	for _, utxo := range utxos {
		utxoViewModels = append(utxoViewModels, SimpleUTxO{
			TransactionHash: hex.EncodeToString(utxo.TransactionHash),
			UTxOID:          hex.EncodeToString(utxo.UTxOID),
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
			TransactionHash: hex.EncodeToString(redeemer.TransactionHash),
			Index:           redeemer.Index,
			Tag:             redeemer.Tag,
			Cbor:            hex.EncodeToString(redeemer.Cbor), // CBOR string representation
		})
	}
	return redeemerViewModels
}

// Helper function to convert a models.Witness to a viewmodel.Witness
func ConvertWitnessModelToViewModel(witness models.Witness) Witness {
	return Witness{
		TransactionHash: hex.EncodeToString(witness.TransactionHash),
		PlutusData:      ConvertByteSliceSliceToStringSlice(witness.PlutusData),
		PlutusV1Scripts: ConvertByteSliceSliceToStringSlice(witness.PlutusV1Scripts),
		PlutusV2Scripts: ConvertByteSliceSliceToStringSlice(witness.PlutusV2Scripts),
		PlutusV3Scripts: ConvertByteSliceSliceToStringSlice(witness.PlutusV3Scripts),
		Redeemers:       ConvertRedeemersToViewModels(witness.Redeemers),
	}
}

// Helper function to convert a slice of byte slices to a slice of strings
func ConvertByteSliceSliceToStringSlice(byteSlices [][]byte) []string {
	stringSlice := []string{}
	for _, byteSlice := range byteSlices {
		stringSlice = append(stringSlice, hex.EncodeToString(byteSlice))
	}
	return stringSlice
}
