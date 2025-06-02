package models

import (
)


// MigrateModels contains a list of model objects that should have DB migrations applied
var MigrateModels = []any{
	&Address{},
	&Transaction{},
	&TransactionInput{},
	&TransactionOutput{},
	&Asset{},
	&Datum{},
	&Redeemer{},
	&Witness{},
	&SimpleUTxO{}, // Add SimpleUTxO to the migration list
}
