package config

import (
	"github.com/Salvionied/apollo/txBuilding/Backend/MaestroChainContext"
	"github.com/Andamio-Platform/andamio-indexer/constants"
)

var (
	// OGMIOS    OgmiosChainContext.OgmiosChainContext
	// CHAIN_CTX BlockFrostChainContext.BlockFrostChainContext
	CHAIN_CTX MaestroChainContext.MaestroChainContext
)

func ChainCTXSetup() error {
	// OGMIOS = OgmiosChainContext.NewOgmiosChainContext(*ogmigo.New(ogmigo.WithEndpoint(constants.OGMIGO_ENDPOINT)), *kugo.New(kugo.WithEndpoint(constants.KUGO_ENDPOINT)))

	// BFC, err := BlockFrostChainContext.NewBlockfrostChainContext(
	// 	constants.BFC_API_URL,
	// 	constants.BFC_NETWORK_ID,
	// 	constants.BFC_API_KEY,
	// )

	MC, err := MaestroChainContext.NewMaestroChainContext(
		constants.MAESTRO_NETWORK_ID,
		constants.MAESTRO_API_KEY,
	)

	if err != nil {
		return err
	} else {
		CHAIN_CTX = MC
		// CHAIN_CTX = BFC
	}
	return nil
}
