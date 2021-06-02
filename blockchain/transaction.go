package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

// A structure that represents a transaction on the Animus Blockchain
type Transaction struct {
	ID      []byte     // Represents the hash of the transaction
	Inputs  []TxInput  // Represents the inputs of the transaction
	Outputs []TxOutput // Represents the outputs of the transaction
}

// A structure that represents the inputs in a transaction
// which are really just references to previous outputs
type TxInput struct {
	ID        []byte // Represents the reference transaction of which the output is a part
	OutIndex  int    // Represents the index of output in the reference transaction
	Signature string // Represents the data from which the output key is derived
}

// A structure that represents the outputs in a transaction
type TxOutput struct {
	Value     int    // Represents the token value of a given transaction output
	PublicKey string // Represents the public key of the transaction output required to retrieve the value inside it
}

func NewTransaction() {}

// A method of Transaction that generates assigns the hash of transaction to the ID field.
// The hash is obtained from the Gob data of the transaction.
func (tx *Transaction) GenerateID() {
	// Declare a new bytes buffer and slice of bytes for the hash
	var encoded bytes.Buffer
	var hash [32]byte

	// Create a new gob encoder from the bytes buffer
	encoder := gob.NewEncoder(&encoded)
	// Encode the transaction into a gob of bytes
	err := encoder.Encode(tx)
	// Handle any potential error
	Handle(err)

	// Hash the gob data with SHA256
	hash = sha256.Sum256(encoded.Bytes())
	// Assign the slice of the hash to the ID field of the transcation
	tx.ID = hash[:]
}

// A constructor function that generates and returns a coinbase Transaction.
// A Coinbase transaction refers to a first transaction on a block and does not refer to any
// previous output transactions and contains a token reward for the user who signs the block.
func NewCoinbaseTransaction(to, data string) *Transaction {
	// Check if the data passed is null
	if data == "" {
		// Create a data based on the address of the reciever
		data = fmt.Sprintf("coins to %s", to)
	}

	// Create a transaction input with no reference to a previous output
	inputs := TxInput{ID: []byte{}, OutIndex: -1, Signature: data}
	// Create a transaction output with the token reward
	outputs := TxOutput{Value: 100, PublicKey: to}

	// Construct a transaction with no ID, and the set of inputs and outputs
	tx := Transaction{ID: nil, Inputs: []TxInput{inputs}, Outputs: []TxOutput{outputs}}
	// Generate an ID (hash) for the transaction
	tx.GenerateID()

	// Return the transaction
	return &tx
}

// A method of Transaction that checks if it is a Coinbase Transaction
func (tx *Transaction) IsCoinbaseTx() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].OutIndex == -1
}

// A method of TxInput that checks if the input signature can unlock outputs of an account address
func (txin *TxInput) CanUnlock(address string) bool {
	return txin.Signature == address
}

// A method of TxOutput that checks if account address can be used unlock the outputs
func (txout *TxOutput) CanBeUnlocked(address string) bool {
	return txout.PublicKey == address
}
