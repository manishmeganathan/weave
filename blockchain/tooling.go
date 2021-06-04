package blockchain

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"log"
	"os"

	"github.com/dgraph-io/badger"
)

const dbfile = "./tmp/db/blocks/MANIFEST"
const dbpath = "./tmp/db/blocks"

// A function that generates and returns the
// Hex/Bytes representation of an int64
func Hexify(number int64) []byte {
	// Construct a new binary buffer
	buff := new(bytes.Buffer)
	// Write the number as a binary into the buffer in Big Endian order
	err := binary.Write(buff, binary.BigEndian, number)
	// Handle any potential error
	Handle(err)

	// Return the bytes from the binary buffer
	return buff.Bytes()
}

// A function to handle errors.
func Handle(err error) {
	// Check if error is non null
	if err != nil {
		// Log the error and throw a panic
		log.Panic(err)
	}
}

// A function to check if the DB exists
func DBexists() bool {
	// Check if the database MANIFEST file exists
	if _, err := os.Stat(dbfile); errors.Is(err, os.ErrNotExist) {
		// Return false if the file is missing
		return false
	}

	// Return true if the file exists
	return true
}

// A function to open the BadgerDB
func DBopen() *badger.DB {
	// Define the BadgerDB options for the DB path
	opts := badger.DefaultOptions(dbpath)
	// Switch off the Badger Logger
	opts.Logger = nil

	// Open the Badger DB with the defined options
	db, err := badger.Open(opts)
	// Handle any potential error
	Handle(err)

	// Return the DB
	return db
}

// A function to serialize a Block into gob of bytes
func BlockSerialize(block *Block) []byte {
	// Create a bytes buffer
	var gobdata bytes.Buffer
	// Create a new Gob encoder with the bytes buffer
	encoder := gob.NewEncoder(&gobdata)
	// Encode the Block into a gob
	err := encoder.Encode(block)
	// Handle any potential errors
	Handle(err)

	// Return the gob bytes
	return gobdata.Bytes()
}

// A function to deserialize a gob of bytes into a Block
func BlockDeserialize(gobdata []byte) *Block {
	// Declare a Block variable
	var block Block
	// Create a new Gob decoder by reading the gob bytes
	decoder := gob.NewDecoder(bytes.NewReader(gobdata))
	// Decode the gob into a Block
	err := decoder.Decode(&block)
	// Handle any potential errors
	Handle(err)

	// Return the pointer to the Block
	return &block
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

// A function to serialize a UTXO into gob of bytes
func TXOListSerialize(txos *TXOList) []byte {
	// Create a bytes buffer
	var gobdata bytes.Buffer

	// Create a new Gob encoder with the bytes buffer
	encoder := gob.NewEncoder(&gobdata)
	// Encode the Transaction into a gob
	err := encoder.Encode(txos)
	// Handle any potential errors
	Handle(err)

	// Return the gob bytes
	return gobdata.Bytes()
}

// A function to deserialize a gob of bytes into a UTXO
func TXOListDeserialize(gobdata []byte) *TXOList {
	// Declare a Block variable
	var txos TXOList

	// Create a new Gob decoder by reading the gob bytes
	decoder := gob.NewDecoder(bytes.NewReader(gobdata))
	// Decode the gob into a Block
	err := decoder.Decode(&txos)
	// Handle any potential errors
	Handle(err)

	// Return the pointer to the Transaction
	return &txos
}
