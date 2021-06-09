/*
This module contains the definition and implementation
of the Transaction struct and its methods
*/
package primitives

import (
	"fmt"
	"strings"

	"github.com/manishmeganathan/blockweave/src/utils"
)

// A structure that represents a transaction on the Animus Blockchain
type Transaction struct {
	// Represents the ID of the transaction obtained from its hash
	ID Hash

	// Represents the list of transaction inputs
	Inputs TXIList

	// Represents the list of transaction outputs
	Outputs TXOList
}

// A method of Transaction that checks if it is a Coinbase Transaction
func (txn *Transaction) IsCoinbaseTxn() bool {
	return len(txn.Inputs) == 1 && len(txn.Inputs[0].ID) == 0 && txn.Inputs[0].OutIndex == -1
}

// A method of Transaction that generates a hash of the Transaction
func (txn *Transaction) GenerateHash() Hash {
	// Create a copy of the transaction
	txncopy := *txn
	// Remove the ID of the transaction copy
	txncopy.ID = Hash{}

	// Serialize the transaction into a gob and hash it
	hash := utils.Hash256(TxnSerialize(&txncopy))
	// Return the hash slice
	return hash[:]
}

// A method of Transaction that generates a safe version
// of the Transaction that does not include the signature
// and public keys of its transaction inputs
func (txn *Transaction) GenerateSafeCopy() Transaction {
	// Declare a slice of transaction inputs
	var inputs TXIList

	// Iterate over the transaction inputs
	for _, input := range txn.Inputs {
		// Append the transaction inputs into the slice without the signature and public key
		inputs = append(inputs, TXI{ID: input.ID, OutIndex: input.OutIndex, Signature: nil, PublicKey: nil})
	}

	// Create a new transaction with the trimmed inputs
	txncopy := Transaction{ID: txn.ID, Inputs: inputs, Outputs: txn.Outputs}
	// Return the trimmed transaction
	return txncopy
}

// A method of Transaction that generates the string representation
// of a transaction and all its inputs and outputs.
// TODO: NEEDS REWORK FOR FORMATTING
func (txn *Transaction) GenerateString() string {
	lines := []string{fmt.Sprintf("--- Transaction %x:", txn.ID)}

	for i, input := range txn.Inputs {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TxnID:     %x", input.ID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.OutIndex))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PublicKey))
	}

	for i, output := range txn.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PublicKeyHash))
	}

	lines = append(lines, "---")

	return strings.Join(lines, "\n")
}
