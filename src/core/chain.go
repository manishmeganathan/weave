package core

import (
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/manishmeganathan/blockweave/src/primitives"
	"github.com/manishmeganathan/blockweave/src/utils"
	"github.com/sirupsen/logrus"
)

type BlockChain struct {
	// Represents the database where the chain is stored
	DB *DatabaseClient

	// Represents the hash of the latest block
	ChainHead primitives.Hash

	// Represents the number of block on the chain (last block height+1)
	ChainHeight int
}

// A constructor function that seeds a new blockchain i.e creates one.
// Returns an error if an Animus Blockchain already exists.
func SeedBlockChain(address primitives.Address) (*BlockChain, error) {
	// Check if a blockchain already exists by checking if the DB exists
	if DBExists() {
		return &BlockChain{}, fmt.Errorf("blockchain already exists exist")
	}

	// Create a null blockchain
	blockchain := BlockChain{}
	// Set up the database client
	blockchain.DB = NewDatabaseClient()

	// Generate a coinbase transaction for the genesis block
	coinbase := NewCoinbaseTransaction(address)

	merkle := NewMerkleBuilder()
	go merkle.Build()

	merkle.BuildQueue <- coinbase
	close(merkle.BuildQueue)

	// Generate a Genesis Block for the chain with a coinbase transaction
	genesisblock := NewBlock(merkle, []byte{}, 0, address, []byte(utils.WeavePOW))
	// Log the minting of the genesis block
	logrus.WithFields(logrus.Fields{
		"address": address.String,
		"reward":  coinbase.Outputs[0].Value,
	}).Info("genesis block has been minted!")

	// Retrieve the hash of the Genesis Block
	chainhead := genesisblock.BlockHash

	// Define an Update Transaction on the BadgerDB
	err := blockchain.DB.Client.Update(func(txn *badger.Txn) error {

		// Add the Block to the DB with its hash as the key and its gob data as the value
		if err := txn.Set(genesisblock.BlockHash, primitives.BlockSerialize(genesisblock)); err != nil {
			// Return any potential error
			return fmt.Errorf("genesis block could not be stored! error - %v", err)
		}

		// Set the chain head of the chain in the DB to the hash of the genesis block
		if err := txn.Set([]byte("chainhead"), chainhead); err != nil {
			// Return any potential error
			return fmt.Errorf("chain head could not be stored! error - %v", err)
		}

		// Set the height of the chain in the DB as 1 (genesis block height + 1)
		if err := txn.Set([]byte("chainheight"), utils.HexEncode(1)); err != nil {
			// Return any potential error
			return fmt.Errorf("chain height could not be stored! error - %v", err)
		}

		return nil
	})
	// Handle any potential errors
	utils.HandleErrorLog(err, "chain seed failed!")

	// Assign the current chain head
	blockchain.ChainHead = chainhead
	// Assign the current chain height
	blockchain.ChainHeight = 1

	// Return the blockchain
	return &blockchain, nil
}

// A constructor function that animates an existing blockchain i.e brings it to life.
// Returns an error if no Animus Blockchain exists.
func AnimateBlockChain() (*BlockChain, error) {
	// Declare a slice a bytes to collect the chain head and height
	var chainhead []byte
	var chainheight []byte

	// Check if a blockchain already exists by checking if the DB exists
	if !DBExists() {
		return &BlockChain{}, fmt.Errorf("blockchain does not exist")
	}

	// Create a null blockchain
	blockchain := BlockChain{}
	// Set up the database client
	blockchain.DB = NewDatabaseClient()

	// Define an Update Transaction on the BadgerDB
	err := blockchain.DB.Client.Update(func(txn *badger.Txn) error {

		// Get the chain head item from the DB
		chainheaditem, err := txn.Get([]byte("chainhead"))
		// Return any potential error
		if err != nil {
			return fmt.Errorf("chain head could not be retrived! error - %v", err)
		}

		// Retrieve the value of the chain head hash item
		if err = chainheaditem.Value(func(val []byte) error {
			chainhead = val
			return nil

		}); err != nil {
			// Return any potential error
			return fmt.Errorf("chain head could not be set! error - %v", err)
		}

		// Get the chain height item from the DB
		chainheightitem, err := txn.Get([]byte("chainheight"))
		// Return any potential error
		if err != nil {
			return fmt.Errorf("chain height could not be retrived! error - %v", err)
		}

		// Retrieve the value of the chain head hash item
		if err = chainheightitem.Value(func(val []byte) error {
			chainhead = val
			return nil

		}); err != nil {
			// Return any potential error
			return fmt.Errorf("chain height could not be set! error - %v", err)
		}

		return nil
	})
	// Handle any potential errors
	utils.HandleErrorLog(err, "chain animate failed!")

	// Assign the current chain head
	blockchain.ChainHead = chainhead
	// Assign the current chain height
	blockchain.ChainHeight = utils.HexDecode(chainheight)

	// Return the blockchain
	return &blockchain, nil
}
