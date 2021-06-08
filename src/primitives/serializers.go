/*
This module contains the functions used to serialize and
deserialize primitive type to binary gob data
*/
package primitives

import (
	"bytes"
	"encoding/gob"

	"github.com/sirupsen/logrus"
)

// A function to serialize a Block into gob of bytes
func BlockSerialize(block *Block) Gob {
	// Create a bytes buffer
	var gobdata bytes.Buffer
	// Create a new Gob encoder with the bytes buffer
	encoder := gob.NewEncoder(&gobdata)
	// Encode the Block into a gob
	err := encoder.Encode(block)
	// Handle any potential errors
	logrus.WithFields(logrus.Fields{
		"struct": "Block",
	}).Fatal("gob encode failed!", err)

	// Return the gob bytes
	return gobdata.Bytes()
}

// A function to deserialize a gob of bytes into a Block
func BlockDeserialize(gobdata Gob) *Block {
	// Declare a Block variable
	var block Block
	// Create a new Gob decoder by reading the gob bytes
	decoder := gob.NewDecoder(bytes.NewReader(gobdata))
	// Decode the gob into a Block
	err := decoder.Decode(&block)
	// Handle any potential errors
	logrus.WithFields(logrus.Fields{
		"struct": "Block",
	}).Fatal("gob encode failed!", err)

	// Return the pointer to the Block
	return &block
}

// A function to serialize a Transaction into a gob of bytes
func TxnSerialize(txn *Transaction) Gob {
	// Create a bytes buffer
	var gobdata bytes.Buffer

	// Create a new Gob encoder with the bytes buffer
	encoder := gob.NewEncoder(&gobdata)
	// Encode the Transaction into a gob
	err := encoder.Encode(txn)
	// Handle any potential errors
	logrus.WithFields(logrus.Fields{
		"struct": "Transaction",
	}).Fatal("gob encode failed!", err)

	// Return the gob bytes
	return gobdata.Bytes()
}

// A function to deserialize a gob of bytes into a Transaction
func TxnDeserialize(gobdata Gob) *Transaction {
	// Declare a Block variable
	var txn Transaction
	// Create a new Gob decoder by reading the gob bytes
	decoder := gob.NewDecoder(bytes.NewReader(gobdata))
	// Decode the gob into a Block
	err := decoder.Decode(&txn)
	// Handle any potential errors
	logrus.WithFields(logrus.Fields{
		"struct": "Transaction",
	}).Fatal("gob decode failed!", err)

	// Return the pointer to the Transaction
	return &txn
}

// A function to serialize a TXOList into gob of bytes
func TXOListSerialize(txos *TXOList) Gob {
	// Create a bytes buffer
	var gobdata bytes.Buffer

	// Create a new Gob encoder with the bytes buffer
	encoder := gob.NewEncoder(&gobdata)
	// Encode the Transaction into a gob
	err := encoder.Encode(txos)
	// Handle any potential errors
	logrus.WithFields(logrus.Fields{
		"struct": "TXOList",
	}).Fatal("gob decode failed!", err)

	// Return the gob bytes
	return gobdata.Bytes()
}

// A function to deserialize a gob of bytes into a TXOList
func TXOListDeserialize(gobdata Gob) *TXOList {
	// Declare a Block variable
	var txos TXOList

	// Create a new Gob decoder by reading the gob bytes
	decoder := gob.NewDecoder(bytes.NewReader(gobdata))
	// Decode the gob into a Block
	err := decoder.Decode(&txos)
	// Handle any potential errors
	logrus.WithFields(logrus.Fields{
		"struct": "TXOList",
	}).Fatal("gob decode failed!", err)

	// Return the pointer to the Transaction
	return &txos
}
