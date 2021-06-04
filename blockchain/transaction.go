package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/manishmeganathan/animus/wallet"
)

// A structure that represents a transaction on the Animus Blockchain
type Transaction struct {
	ID      []byte     // Represents the hash of the transaction
	Inputs  []TxInput  // Represents the inputs of the transaction
	Outputs []TxOutput // Represents the outputs of the transaction
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
	// Set the ID (hash) for the transaction
	txn.ID = txn.GenerateHash()
	// Return the transaction
	return &txn
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
	inputs := TxInput{ID: []byte{}, OutIndex: -1, Signature: nil, PublicKey: []byte(data)}
	// Create a transaction output with the token reward
	outputs := *NewTxOutput(100, to)

	// Construct a transaction with no ID, and the set of inputs and outputs
	txn := Transaction{ID: nil, Inputs: []TxInput{inputs}, Outputs: []TxOutput{outputs}}
	// Set the ID (hash) for the transaction
	txn.ID = txn.GenerateHash()

	// Return the transaction
	return &txn
}

// A method of Transaction that checks if it is a Coinbase Transaction
func (txn *Transaction) IsCoinbaseTxn() bool {
	return len(txn.Inputs) == 1 && len(txn.Inputs[0].ID) == 0 && txn.Inputs[0].OutIndex == -1
}

// A method of Transaction that generates a hash of the Transaction
func (txn *Transaction) GenerateHash() []byte {
	// Create a copy of the transaction
	txncopy := *txn
	// Remove the ID of the transaction copy
	txncopy.ID = []byte{}

	// Serialize the transaction into a gob and hash it
	hash := sha256.Sum256(TxnSerialize(&txncopy))
	// Return the hash slice
	return hash[:]
}

// A method of Transaction that generates a trimmed version
// of the Transaction that does not include the signature
// and public keys of the transaction inputs
func (txn *Transaction) GenerateTrimCopy() Transaction {
	// Declare a slice of transaction inputs
	var inputs []TxInput

	// Iterate over the transaction inputs
	for _, input := range txn.Inputs {
		// Append the transaction inputs into the slice without the signature and public key
		inputs = append(inputs, TxInput{ID: input.ID, OutIndex: input.OutIndex, Signature: nil, PublicKey: nil})
	}

	// Create a new transaction with the trimmed inputs
	txncopy := Transaction{ID: txn.ID, Inputs: inputs, Outputs: txn.Outputs}
	// Return the trimmed transaction
	return txncopy
}

// A function to serialize a Transaction into gob of bytes
func TxnSerialize(txn *Transaction) []byte {
	// Create a bytes buffer
	var gobdata bytes.Buffer

	// Create a new Gob encoder with the bytes buffer
	encoder := gob.NewEncoder(&gobdata)
	// Encode the Transaction into a gob
	err := encoder.Encode(txn)
	// Handle any potential errors
	Handle(err)

	// Return the gob bytes
	return gobdata.Bytes()
}

// A function to deserialize a gob of bytes into a Transaction
func TxnDeserialize(gobdata []byte) *Transaction {
	// Declare a Block variable
	var txn Transaction
	// Create a new Gob decoder by reading the gob bytes
	decoder := gob.NewDecoder(bytes.NewReader(gobdata))
	// Decode the gob into a Block
	err := decoder.Decode(&txn)
	// Handle any potential errors
	Handle(err)

	// Return the pointer to the Transaction
	return &txn
}

// A structure that represents the inputs in a transaction
// which are really just references to previous outputs
type TxInput struct {
	ID        []byte // Represents the reference transaction of which the output is a part
	OutIndex  int    // Represents the index of output in the reference transaction
	Signature []byte // Represents the signature of the transaction
	PublicKey []byte // Represents the public key of the sending address
}

// A method of TxInput that checks if the input public key is valid for a given public key hash
func (txin *TxInput) CheckKey(pubkeyhash []byte) bool {
	// Generate the hash of the input public key
	lockhash := wallet.GeneratePublicKeyHash(txin.PublicKey)
	// Check if the locking hash is equal to the given hash
	return bytes.Equal(lockhash, pubkeyhash)
}

// A structure that represents the outputs in a transaction
type TxOutput struct {
	Value         int    // Represents the token value of a given transaction output
	PublicKeyHash []byte // Represents the hash of the public key of the recieving address
}

// A constructor function
func NewTxOutput(value int, address string) *TxOutput {
	txo := TxOutput{Value: value, PublicKeyHash: nil}
	txo.Lock([]byte(address))

	return &txo
}

// A method of TxOutput that locks the output for a given address
func (txout *TxOutput) Lock(address []byte) {
	// Decode the address from base58
	publickeyhash := wallet.Base58Decode(address)
	// Isolate public key hash from the checksum and version
	publickeyhash = publickeyhash[1 : len(publickeyhash)-4]
	// Assign the output key hash to public hash of the given address
	txout.PublicKeyHash = publickeyhash
}

// A method of TxOutput that checks if the ouput key hash is valid for a given locking hash
func (txout *TxOutput) CheckLock(lockhash []byte) bool {
	// Check if locking hash is equal to output's key hash
	return bytes.Equal(txout.PublicKeyHash, lockhash)
}
