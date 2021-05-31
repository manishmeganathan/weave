package blockchain

import (
	"log"

	"github.com/dgraph-io/badger"
)

const dbpath = "./db/blocks"

// A structure that represents the Animus BlockChain
type BlockChain struct {
	Database *badger.DB // Represents the reference to the chain database
	LastHash []byte     // Represents the hash of the last block on the chain
}

// A structure that represents an Iterator for the Animus BlockChain
type BlockChainIterator struct {
	CursorHash []byte     // Represents the hash of the block that the iterator is currently on
	Database   *badger.DB // Represents the reference to the chain database
}

// A constructor function that generates a BlockChain from the database.
// If the blockchain does not exist on the database, a new is created
// and initialized with a Genesis Block.
func NewBlockChain() *BlockChain {
	// Declare a slice a bytes to collect the hash value
	var lasthash []byte

	// Define the BadgerDB options for the DB path
	opts := badger.DefaultOptions(dbpath)
	// Switch off the Badger Logger
	opts.Logger = nil

	// Open the Badger DB with the defined options
	db, err := badger.Open(opts)
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

// A constructor function that generates an iterator for the BlockChain
func NewIterator(chain *BlockChain) *BlockChainIterator {
	// Assign the values of the BlockChainIterator and return it
	return &BlockChainIterator{CursorHash: chain.LastHash, Database: chain.Database}
}

// A method of BlockChainIterator that iterates over chain and returns the
// next block on the chain (backwards) from the chain DB and returns it
func (iter *BlockChainIterator) Next() *Block {
	// Declare the Block variable
	var block Block

	// Define a View Transaction on the BadgerDB
	err := iter.Database.View(func(txn *badger.Txn) error {

		// Get the block item for the current hash of the iterator
		item, err := txn.Get(iter.CursorHash)
		// Handle any potential errors
		Handle(err)

		// Declare a slice of bytes for the gob of block data
		var blockgob []byte
		// Retrieve the value of the gob data
		err = item.Value(func(val []byte) error {
			blockgob = val
			return nil
		})

		// Convert the block gob data into a Block object
		block = *BlockDeserialize(blockgob)
		return err
	})

	// Handle any potential error
	Handle(err)
	// Update the iterator's cursor to the hash of the previous block
	iter.CursorHash = block.PrevHash
	// Return the block
	return &block
}
