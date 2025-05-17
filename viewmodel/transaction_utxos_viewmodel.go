package viewmodel

// TransactionUTxOs represents the view model for transaction UTxOs (inputs and outputs).
type TransactionUTxOs struct {
	Inputs  []TransactionInput  `json:"inputs"`
	Outputs []TransactionOutput `json:"outputs"`
}