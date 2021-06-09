package core

import (
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/manishmeganathan/blockweave/src/primitives"
	"github.com/manishmeganathan/blockweave/src/utils"
	"github.com/sirupsen/logrus"
)

type BlockChain struct {
	DB            *Database
	LastHash      primitives.Hash
	CurrentHeight int
}

func SeedBlockChain(address primitives.Address) (*BlockChain, error) {
	// Declare a slice a bytes to collect the hash value
	var lasthash primitives.Hash

	// Check if a blockchain already exists by checking if the DB exists
	if DBExists() {
		return &BlockChain{}, fmt.Errorf("blockchain already exists exist")
	}

	// Create a null blockchain
	blockchain := BlockChain{}
	// Set up the database client
	blockchain.DB = NewDatabase()

	// Generate a coinbase transaction for the genesis block
	coinbase := NewCoinbaseTransaction(address)

	merkle := NewMerkleBuilder()
	go merkle.Build()

	merkle.BuildQueue <- coinbase
	close(merkle.BuildQueue)

	// Generate a Genesis Block for the chain with a coinbase transaction
	genesisblock := NewBlock(merkle, []byte{}, 0, address, []byte(utils.WeavePOW))

	logrus.WithFields(logrus.Fields{
		"address": address.String,
		"reward":  coinbase.Outputs[0].Value,
	}).Info("genesis block has been minted!")

	// Retrieve the hash of the Genesis Block
	lasthash = genesisblock.BlockHash

	// Define an Update Transaction on the BadgerDB
	err := blockchain.DB.Client.Update(func(txn *badger.Txn) error {

		// Add the Block to the DB with its hash as the key and its gob data as the value
		err := txn.Set(genesisblock.BlockHash, primitives.BlockSerialize(genesisblock))
		// Handle any potential error
		logrus.Fatal("genesis block could not be stored!", err)

		// Set the last hash of the chain in the DB to the Genesis Block's Hash
		err = txn.Set([]byte("lh"), lasthash)
		return err
	})

	// Log the error
	logrus.Fatal("chain failed to be seeded!", err)

	// Assign the last hash of the chain
	blockchain.LastHash = lasthash
	// Assign the current height of the chain
	blockchain.CurrentHeight = 0

	// Return the blockchain
	return &blockchain, nil
}
