package eventHandlers

import (
	"github.com/blinklabs-io/adder/event"
)

func TxEvent(evt event.Event) error {
	// if evt.Type == "chainsyc.transaction" {
	// 	fmt.Println(evt)
	// }

	// if evt.Type == "chainsync.transaction" {

	// 	eventTx := evt.Payload.(input_chainsync.TransactionEvent)
	// 	eventCtx := evt.Context.(input_chainsync.TransactionContext)

	// 	globalDB := database.GetGlobalDB()
	// 	// txn := globalDB.Transaction(true) // Keep txn as it might be used later for database operations

	// 	// Keep slot, txHash, txId if they are intended for later use, otherwise remove.
	// 	// Assuming they might be used for database operations related to the transaction.
	// 	// slot := eventCtx.SlotNumber
	// 	// txHash := eventCtx.TransactionHash
	// 	// txId := eventCtx.TransactionIdx

	// 	// Remove unused variables
	// 	// cert := eventTx.Certificates[0].Cbor()
	// 	// native := eventTx.Witnesses.NativeScripts()[0].Item()
	// 	// pl := eventTx.Witnesses.PlutusData()[0].Cbor()

	// 	redeemersInterface := eventTx.Witnesses.Redeemers()

	// 	redeemerTags := []common.RedeemerTag{
	// 		common.RedeemerTagSpend,
	// 		common.RedeemerTagMint,
	// 		common.RedeemerTagCert,
	// 		common.RedeemerTagReward,
	// 	}

	// 	for _, tag := range redeemerTags {
	// 		indexes := redeemersInterface.Indexes(tag)
	// 		for _, index := range indexes {
	// 			// Get the redeemer value and execution units for the specific index and tag
	// 			redeemerValue, exUnits := redeemersInterface.Value(index, tag)

	// 			// Get the raw CBOR bytes
	// 			cborBytes := redeemerValue.Cbor()

	// 			// Now you have the cborBytes ([]byte) and exUnits (ExUnits)
	// 			// You can store cborBytes in your database.
	// 			// You might also want to store the tag and index for context.

	// 			// Example: Print the details (replace with your database storage logic)
	// 			fmt.Printf("Redeemer Tag: %s, Index: %d, CBOR Bytes Length: %d, ExUnits: %+v\n", tag, index, len(cborBytes), exUnits)

	// 			// Store cborBytes in your database here
	// 			// db.SaveRedeemer(eventTx.Hash(), tag, index, cborBytes, exUnits)
	// 		}
	// 	}

	// 	// datum := eventTx.ResolvedInputs[0].Cbor()
	// 	// Removed redeclared datum variable and type mismatch
	// 	// datum := eventTx.Inputs[0].String()

	// 	// err := globalDB.NewTx(tx.Hash().Bytes(), 0, slot, pkh, skh, tx.Cbor(), txn)

	// }

	return nil

}
