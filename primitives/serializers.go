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
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to encode Block as Gob.")
	}

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
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to decode Gob as Block.")
	}

	// Return the pointer to the Block
	return &block
}

// A function to serialize a BlockHeader into gob of bytes
func BlockHeaderSerialize(header *BlockHeader) Gob {
	// Create a bytes buffer
	var gobdata bytes.Buffer
	// Create a new Gob encoder with the bytes buffer
	encoder := gob.NewEncoder(&gobdata)
	// Encode the Block into a gob
	err := encoder.Encode(header)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to encode BlockHeader as Gob.")
	}

	// Return the gob bytes
	return gobdata.Bytes()
}

// A function to deserialize a gob of bytes into a BlockHeader
func BlockHeaderDeserialize(gobdata Gob) *BlockHeader {
	// Declare a Block variable
	var header BlockHeader
	// Create a new Gob decoder by reading the gob bytes
	decoder := gob.NewDecoder(bytes.NewReader(gobdata))
	// Decode the gob into a Block
	err := decoder.Decode(&header)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to decode Gob as BlockHeader.")
	}

	// Return the pointer to the BlockHeader
	return &header
}

// A function to serialize a Transaction into a gob of bytes
func TxnSerialize(txn *Transaction) Gob {
	// Create a bytes buffer
	var gobdata bytes.Buffer

	// Create a new Gob encoder with the bytes buffer
	encoder := gob.NewEncoder(&gobdata)
	// Encode the Transaction into a gob
	err := encoder.Encode(txn)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to encode Transaction as Gob.")
	}

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
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to decode Gob as Transaction.")
	}

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
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to encode TXOList as Gob.")
	}

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
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to decode Gob as TXOList.")
	}

	// Return the pointer to the Transaction
	return &txos
}
