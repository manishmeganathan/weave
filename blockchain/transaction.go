package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

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

	// Create the wallet store
	walletstore := wallet.NewWalletStore()
	// Fetch the wallet from the wallet store for the given address
	w := walletstore.FetchWallet(from)

	// Generate the public key hash for the wallet's public key
	publickeyhash := wallet.GeneratePublicKeyHash(w.PublicKey)
	// Collect the spendable transaction outputs of the account up to the given amount with the public key hash
	accumulated, validoutputs := chain.CollectSpendableTXOS(publickeyhash, amount)

	// Check if the account has enough funds
	if accumulated < amount {
		log.Panic("error: insufficient funds!")
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
			input := TxInput{ID: txid, OutIndex: output, Signature: nil, PublicKey: w.PublicKey}
			// Add the transaction input into the slice
			txinputs = append(txinputs, input)
		}
	}

	// Add a transaction output with the amount to the address
	txoutputs = append(txoutputs, *NewTxOutput(amount, to))

	// Check if there is a balance in the accumulated amounted
	if accumulated > amount {
		// Add a transaction output with the balance amount back to the original sender
		txoutputs = append(txoutputs, *NewTxOutput(accumulated-amount, from))
	}

	// Create a Transaction with the list of input and outputs
	txn := Transaction{ID: nil, Inputs: txinputs, Outputs: txoutputs}
	// Set the ID (hash) for the transaction
	txn.ID = txn.GenerateHash()
	// Sign the transaction using the wallet's private key
	chain.SignTransaction(&txn, w.PrivateKey)

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

// A method of Transaction that generates a string representation
// of the transaction and all its inputs and outputs.
func (txn *Transaction) String() string {
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

// A method of Transaction that signs the transaction given the private key
// of the wallet and a map of previous Transaction IDs to their Transactions.
func (txn *Transaction) Sign(privatekey ecdsa.PrivateKey, prevtxns map[string]Transaction) {
	// Check if the transaction is a coinbase (cannot sign coinbase txns)
	if txn.IsCoinbaseTxn() {
		return
	}

	// Iterate over the transaction inputs
	for _, input := range txn.Inputs {
		// Check if the previous transactions value is consistent
		if prevtxns[hex.EncodeToString(input.ID)].ID == nil {
			Handle(fmt.Errorf("error: previous transaction is not correct"))
		}
	}

	// Generate a trim copy of the transaction
	txncopy := txn.GenerateTrimCopy()

	// Iterate over the inputs of the trimmed transaction
	for inpindex, input := range txncopy.Inputs {
		// Retrive the corresponding the previous transaction
		prevtxn := prevtxns[hex.EncodeToString(input.ID)]

		// Set the input signature to nil
		txncopy.Inputs[inpindex].Signature = nil
		// Set the input public key with the public key hash from the prev txn
		txncopy.Inputs[inpindex].PublicKey = prevtxn.Outputs[input.OutIndex].PublicKeyHash
		// Generate the hash of the trimmed transaction and assign it to its ID
		txncopy.ID = txncopy.GenerateHash()
		// Set the input public key to nil
		txncopy.Inputs[inpindex].PublicKey = nil

		// Sign the transaction with the ECDSA method using the private key and ID of the transaction trim
		r, s, err := ecdsa.Sign(rand.Reader, &privatekey, txncopy.ID)
		// Handle any potential error
		Handle(err)

		// Append method outputs to form the signature
		signature := append(r.Bytes(), s.Bytes()...)
		// Assign the signature of the Transaction
		txn.Inputs[inpindex].Signature = signature
	}

}

// A method of Transaction that verifies if the transaction signature is
// valid for the given of map Transaction IDs to their Transactions.
func (txn *Transaction) Verify(prevtxns map[string]Transaction) bool {
	// Check if the transaction is a coinbase (cannot verify coinbase txns)
	if txn.IsCoinbaseTxn() {
		return true
	}

	// Iterate over the transaction inputs
	for _, input := range txn.Inputs {
		// Check if the previous transactions value is consistent
		if prevtxns[hex.EncodeToString(input.ID)].ID == nil {
			Handle(fmt.Errorf("error: previous transaction is not correct"))
		}
	}

	// Generate a trim copy of the transaction
	txncopy := txn.GenerateTrimCopy()

	// Iterate over the inputs of the trimmed transaction
	for inpindex, input := range txncopy.Inputs {
		// Retrive the corresponding the previous transaction
		prevtxn := prevtxns[hex.EncodeToString(input.ID)]

		// Set the input signature to nil
		txncopy.Inputs[inpindex].Signature = nil
		// Set the input public key with the public key hash from the prev txn
		txncopy.Inputs[inpindex].PublicKey = prevtxn.Outputs[input.OutIndex].PublicKeyHash
		// Generate the hash of the trimmed transaction and assign it to its ID
		txncopy.ID = txncopy.GenerateHash()
		// Set the input public key to nil
		txncopy.Inputs[inpindex].PublicKey = nil

		// Declare r and s as big Ints
		r := big.Int{}
		s := big.Int{}
		// Retrieve the size of the signature
		signaturesize := len(input.Signature)
		// Split the signature into r and s values
		r.SetBytes(input.Signature[:(signaturesize / 2)])
		s.SetBytes(input.Signature[(signaturesize / 2):])

		// Declare the x and y as big Ints
		x := big.Int{}
		y := big.Int{}
		// Retrieve the size of the public key
		keysize := len(input.PublicKey)
		// Split the public key into its x and y coordinates
		x.SetBytes(input.PublicKey[:(keysize / 2)])
		y.SetBytes(input.PublicKey[(keysize / 2):])

		// Create an ECDSA public key from sepc256r1 curve and the x, y coordinates
		rawpublickey := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}

		// Check if the transaction has been signed with the public key's private pair
		if !ecdsa.Verify(&rawpublickey, txncopy.ID, &r, &s) {
			return false
		}
	}

	// Return true if all transactions are verified
	return true
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
func (txin *TxInput) CheckKey(publickeyhash []byte) bool {
	// Generate the hash of the input public key
	lockhash := wallet.GeneratePublicKeyHash(txin.PublicKey)
	// Check if the locking hash is equal to the given hash
	return bytes.Equal(lockhash, publickeyhash)
}

// A structure that represents the outputs in a transaction
type TxOutput struct {
	Value         int    // Represents the token value of a given transaction output
	PublicKeyHash []byte // Represents the hash of the public key of the recieving address
}

type TXOList []TxOutput

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
