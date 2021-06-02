package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
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

// A constructor function that generates and returns a Transaction
// given the to and from addresses and the amount to transact.
func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	// Declare slices of transaction outputs and inputs
	var txinputs []TxInput
	var txoutputs []TxOutput
	// Accumulate the spendable transaction outputs of the account up to the given amount
	accumulated, validoutputs := chain.AccumulateSpendableTXO(from, amount)

	// Check if the account has enough funds
	if accumulated < amount {
		log.Panic("Error: Insufficient Funds!")
	}

	// Iterate over the spendable transaction output IDs
	for txnid, outputs := range validoutputs {
		// Decode the transaction ID
		txid, err := hex.DecodeString(txnid)
		// Handle any potential error
		Handle(err)

		// Iterate over the the output indexes
		for _, output := range outputs {
			// Create a transaction input with the transaction ID, output index and from address signature
			input := TxInput{ID: txid, OutIndex: output, Signature: from}
			// Add the transaction input into the slice
			txinputs = append(txinputs, input)
		}
	}

	// Add a transaction output with the amount to the address
	txoutputs = append(txoutputs, TxOutput{Value: amount, PublicKey: to})

	// Check if there is a balance in the accumulated amounted
	if accumulated > amount {
		// Add a transaction output with the balance amount back to the original sender
		txoutputs = append(txoutputs, TxOutput{Value: accumulated - amount, PublicKey: from})
	}

	// Create a Transaction with the list of input and outputs
	txn := Transaction{ID: nil, Inputs: txinputs, Outputs: txoutputs}
	// Generate an ID for the transaction
	txn.GenerateID()
	// Return the transaction
	return &txn
}

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
