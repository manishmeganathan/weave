package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/manishmeganathan/weave/utils"
	"github.com/manishmeganathan/weave/wallet"
	"github.com/sirupsen/logrus"
)

// A structure that represents a transaction on the Animus Blockchain
type Transaction struct {
	// Represents the ID of the transaction obtained from its hash
	ID utils.Hash

	// Represents the list of transaction inputs
	Inputs TXIList

	// Represents the list of transaction outputs
	Outputs TXOList
}

// A constructor function that generates and returns a Transaction
// given the to and from addresses and the amount to transact.
func NewTransaction(from, to wallet.Address, amount int, chain *BlockChain) *Transaction {
	// Declare slices of transaction outputs and inputs
	var txinputs TXIList
	var txoutputs TXOList

	// Create the wallet store
	wallets := wallet.NewJBOK()
	// Fetch the wallet from the wallet store for the given address
	w := wallets.FetchWallet(from.String)

	// // Generate the public key hash for the wallet's public key
	// publickeyhash := wallet.GeneratePublicKeyHash(w.PublicKey)

	// Collect the spendable transaction outputs of the account up to the given amount with the public key hash
	accumulated, validoutputs := chain.CollectSpendableUTXOS(from.PublicKeyHash, amount)

	// Check if the account has enough funds
	if accumulated < amount {
		log.Panic("error: insufficient funds!")
	}

	// Iterate over the spendable transaction output IDs
	for txnid, outputs := range validoutputs {
		// Decode the transaction ID
		txid, _ := hex.DecodeString(txnid)

		// Iterate over the the output indexes
		for _, output := range outputs {
			// Create a transaction input with the transaction ID, output index and from address signature
			input := TXI{ID: txid, OutIndex: output, Signature: nil, PublicKey: w.PublicKey}
			// Add the transaction input into the slice
			txinputs = append(txinputs, input)
		}
	}

	// Add a transaction output with the amount to the address
	txoutputs = append(txoutputs, *NewTXO(amount, to))

	// Check if there is a balance in the accumulated amounted
	if accumulated > amount {
		// Add a transaction output with the balance amount back to the original sender
		txoutputs = append(txoutputs, *NewTXO(accumulated-amount, from))
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
func NewCoinbaseTransaction(to wallet.Address) *Transaction {
	// Create a slice a bytes
	randdata := make([]byte, 24)
	// Add random data to the slice of bytes
	_, err := rand.Read(randdata)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to generate random bytes.")
	}

	// Collect the data from the hexadecimal interpretation of the random bytes
	data := fmt.Sprintf("%x", randdata)

	// Create a transaction input with no reference to a previous output
	inputs := TXI{ID: []byte{}, OutIndex: -1, Signature: nil, PublicKey: []byte(data)}
	// Create a transaction output with the token reward
	outputs := *NewTXO(25, to)

	// Construct a transaction with no ID, and the set of inputs and outputs
	txn := Transaction{
		ID:      nil,
		Inputs:  TXIList{inputs},
		Outputs: TXOList{outputs},
	}

	// Set the ID (hash) for the transaction
	txn.ID = txn.GenerateHash()

	// Return the transaction
	return &txn
}

// A method of Transaction that checks if it is a Coinbase Transaction
func (txn *Transaction) IsCoinbase() bool {
	return len(txn.Inputs) == 1 && len(txn.Inputs[0].ID) == 0 && txn.Inputs[0].OutIndex == -1
}

// A method of Transaction that generates a hash of the Transaction
func (txn *Transaction) GenerateHash() utils.Hash {
	// Create a copy of the transaction
	txncopy := *txn
	// Remove the ID of the transaction copy
	txncopy.ID = utils.Hash{}

	// Serialize the transaction into a gob and hash it
	hash := utils.Hash256(txncopy.Serialize())
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

// A method that returns the gob encoded data of the Transaction
func (txn *Transaction) Serialize() utils.Gob {
	// Encode the blockheader as a gob and return it
	return utils.GobEncode(txn)
}

// A method that decodes a gob of bytes into the Transaction struct
func (txn *Transaction) Deserialize(gobdata utils.Gob) {
	// Decode the gob data into the blockheader
	utils.GobDecode(gobdata, txn)
}
