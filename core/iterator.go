package core

import (
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/manishmeganathan/blockweave/persistence"
	"github.com/sirupsen/logrus"
)

// A structure that represents an Iterator for the BlockChain
type BlockChainIterator struct {
	Cursor   []byte                      // Represents the hash of the block that the iterator is currently on
	Database *persistence.DatabaseClient // Represents the reference to the chain database
}

// A constructor function that generates and returns an iterator for the BlockChain
func NewIterator(chain *BlockChain) *BlockChainIterator {
	// Assign the values of the BlockChainIterator and return it
	return &BlockChainIterator{Cursor: chain.ChainHead, Database: chain.DB}
}

// A method of BlockChainIterator that iterates over chain and returns the
// next block on the chain (backwards) from the chain DB and returns it
func (iter *BlockChainIterator) Next() *Block {
	// Declare a Block variable
	var block Block

	// Define a View Transaction on the BadgerDB
	err := iter.Database.Client.View(func(txn *badger.Txn) error {

		// Get the block item for the current hash of the iterator
		item, err := txn.Get(iter.Cursor)
		// Return any potential error
		if err != nil {
			return fmt.Errorf("block item retrieval failed! error - %v", err)
		}

		// Declare a slice of bytes for the gob of block data
		var blockgob []byte
		// Retrieve the value of the gob data
		if err = item.Value(func(val []byte) error {
			blockgob = val
			return nil

		}); err != nil {
			// Return any potential error
			return fmt.Errorf("block item value retrival failed! error - %v", err)
		}

		// Convert the block gob data into a Block object
		block = Block{}
		block.Deserialize(blockgob)

		return err
	})

	// Handle any potential error
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to iterate over chain.")
	}

	// Update the iterator's cursor to the hash of block before the current block
	iter.Cursor = block.BlockHeader.Priori
	// Return the block
	return &block
}
