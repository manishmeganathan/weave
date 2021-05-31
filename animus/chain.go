package animus

import (
	"log"

	"github.com/dgraph-io/badger"
)

const dbpath = "./db/blocks"

type BlockChain struct {
	Database *badger.DB
	LastHash []byte
}

// A constructor function that generates a BlockChain from the database.
// If the blockchain does not exist on the database, a new is created
// and initialized with a Genesis Block.
func NewBlockChain() *BlockChain {
	// Declare a slice a bytes to collect the hash value
	var lasthash []byte

	// Open the Badger DB with the default option and db path
	db, err := badger.Open(badger.DefaultOptions(dbpath))
	// Handle any potential error
	Handle(err)

	// Define an Update Transaction on the BadgerDB
	err = db.Update(func(txn *badger.Txn) error {

		// Check if the last hash key 'lh' has been set in the DB
		if item, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			// Print that no blockchain currently exists
			log.Println("No Existing BlockChain!")

			// Generate a Genesis Block for the chain
			genesisblock := NewBlock("Genesis", []byte{})
			log.Println("Genesis Block Signed!")

			// Add the Block to the DB with its hash as the key and its gob data as the value
			err = txn.Set(genesisblock.Hash, BlockSerialize(genesisblock))
			// Handle any potential error
			Handle(err)

			// Retrieve the hash of the Genesis Block
			lasthash = genesisblock.Hash

			// Set the last hash of the chain in the DB to the Genesis Block's Hash
			err = txn.Set([]byte("lh"), lasthash)
			return err

		} else {

			// Retrieve the value of the last hash item
			err = item.Value(func(val []byte) error {
				lasthash = val
				return nil
			})
			return err
		}

	})
	// Handle any potential error
	Handle(err)

	// Construct a blockchain with the BadgerDB and the last hash of the chain
	blockchain := BlockChain{Database: db, LastHash: lasthash}
	// Return the blockchain
	return &blockchain
}

// A method of BlockChain that adds a new Block to the chain
func (chain *BlockChain) AddBlock(blockdata string) {
	// Declare a slice a bytes to collect the hash value
	var lasthash []byte

	// Define a View Transaction on the BadgerDB
	err := chain.Database.View(func(txn *badger.Txn) error {

		// Get the value of the last hash key in the database
		item, err := txn.Get([]byte("lh"))
		// Handle any potential error
		Handle(err)

		// Retrieve the value of the last hash item
		err = item.Value(func(val []byte) error {
			lasthash = val
			return nil
		})
		return err
	})

	// Handle any potential error
	Handle(err)
	// Generate a new Block from the given block data and the hash of the previous block
	block := NewBlock(blockdata, lasthash)

	// Define an Update Transaction on the BadgerDB
	err = chain.Database.Update(func(txn *badger.Txn) error {

		// Add the Block to the DB with its hash as the key and its gob data as the value
		err := txn.Set(block.Hash, BlockSerialize(block))
		// Handle any potential error
		Handle(err)

		// Assign the hash of the block as the last hash of the chain
		chain.LastHash = block.Hash
		// Set the last hash key of the database to the hash of the block
		err = txn.Set([]byte("lh"), block.Hash)
		return err
	})

	// Handle any potential error
	Handle(err)
}
