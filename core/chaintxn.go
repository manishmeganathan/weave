package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/manishmeganathan/weave/utils"
	"github.com/sirupsen/logrus"
)

// A method of BlockChain that finds a transaction
// from the chain given a valid Transaction ID
func (chain *BlockChain) FindTransaction(txnid []byte) (Transaction, error) {

	// Get an iterator for the blockchain and iterate over its block
	iter := NewIterator(chain)

	for {
		// Get a block from the iterator
		block := iter.Next()

		// Iterate over the transactions of the block
		for _, txn := range block.TXList {
			// Check if the transaction ID matches
			if bytes.Equal(txn.ID, txnid) {
				// Return the transaction with a nil error
				return *txn, nil
			}
		}

		// Check if the block is genesis block
		if block.BlockHeight == 0 {
			break
		}
	}

	// Return a nil Transaction with an error
	return Transaction{}, fmt.Errorf("transaction does not exist")
}

// A method of BlockChain that signs a transaction given a private key
func (chain *BlockChain) SignTransaction(txn *Transaction, privatekey ecdsa.PrivateKey) {
	// Check if the transaction is a coinbase (cannot sign coinbase txns)
	if txn.IsCoinbase() {
		return
	}

	// Create a map of transaction IDs to Transactions
	prevtxns := make(map[string]Transaction)

	// Iterate over the inputs of the transaction
	for _, input := range txn.Inputs {
		// Find the Transaction with ID on the input from the blockchain
		prevtxn, err := chain.FindTransaction(input.ID)
		if err != nil {
			// Log a fatal error
			logrus.WithFields(logrus.Fields{"error": err}).Fatalf("failed to sign transaction. could not find transaction associated with txn input - %s\n", input.ID)
		}

		// Add the transaction to the map
		prevtxns[hex.EncodeToString(prevtxn.ID)] = prevtxn
	}

	// Generate a safe copy of the transaction
	txncopy := txn.GenerateSafeCopy()

	// Iterate over the inputs of the trimmed transaction
	for inpindex, input := range txncopy.Inputs {
		// Retrive the corresponding the previous transaction
		prevtxn := prevtxns[hex.EncodeToString(input.ID)]
		// Retrieve the public key hash of the corresponding transaction output
		pubkeyhash := prevtxn.Outputs[input.OutIndex].PublicKeyHash

		// Set the input signature to nil
		txncopy.Inputs[inpindex].Signature = nil

		// Set the input public key with the public key hash from the prev txn
		txncopy.Inputs[inpindex].PublicKey = utils.PublicKey(pubkeyhash)
		// Generate the hash of the trimmed transaction and assign it to its ID
		txncopy.ID = txncopy.GenerateHash()
		// Set the input public key to nil
		txncopy.Inputs[inpindex].PublicKey = nil

		// Sign the transaction with the ECDSA method using the private key and ID of the transaction trim
		r, s, err := ecdsa.Sign(rand.Reader, &privatekey, txncopy.ID)
		if err != nil {
			// Log a fatal error
			logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to sign transaction.")
		}

		// Append method outputs to form the signature
		signature := append(r.Bytes(), s.Bytes()...)
		// Assign the signature of the Transaction
		txn.Inputs[inpindex].Signature = signature
	}
}

// A method of BlockChain that verifies the signature of a transaction given a private key
func (chain *BlockChain) VerifyTransaction(txn *Transaction, privatekey ecdsa.PrivateKey) bool {
	// Check if transaction is a coinbase
	if txn.IsCoinbase() {
		return true
	}

	// Create a map of transaction IDs to Transactions
	prevtxns := make(map[string]Transaction)

	// Iterate over the inputs of the transaction
	for _, input := range txn.Inputs {
		// Find the Transaction with ID on the input from the blockchain
		prevtxn, err := chain.FindTransaction(input.ID)
		if err != nil {
			// Log a fatal error
			logrus.WithFields(logrus.Fields{"error": err}).Fatalf("failed to verify transaction. could not find transaction associated with txn input - %s\n", input.ID)
		}

		// Add the transaction to the map
		prevtxns[hex.EncodeToString(prevtxn.ID)] = prevtxn
	}

	// Generate a safe copy of the transaction
	txncopy := txn.GenerateSafeCopy()

	// Iterate over the inputs of the trimmed transaction
	for inpindex, input := range txncopy.Inputs {
		// Retrive the corresponding the previous transaction
		prevtxn := prevtxns[hex.EncodeToString(input.ID)]
		// Retrieve the public key hash of the corresponding transaction output
		pubkeyhash := prevtxn.Outputs[input.OutIndex].PublicKeyHash

		// Set the input signature to nil
		txncopy.Inputs[inpindex].Signature = nil

		// Set the input public key with the public key hash from the prev txn
		txncopy.Inputs[inpindex].PublicKey = utils.PublicKey(pubkeyhash)
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

	// Return true if transactions is verified
	return true
}
