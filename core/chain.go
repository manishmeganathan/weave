package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/dgraph-io/badger"
	"github.com/manishmeganathan/blockweave/merkle"
	"github.com/manishmeganathan/blockweave/persistence"
	"github.com/manishmeganathan/blockweave/utils"
	"github.com/manishmeganathan/blockweave/wallet"
	"github.com/sirupsen/logrus"
)

// A structure that represents the blockchain
type BlockChain struct {
	// Represents the database where the chain is stored
	DB *persistence.DatabaseClient

	// Represents the hash of the latest block
	ChainHead utils.Hash

	// Represents the number of block on the chain (last block height+1)
	ChainHeight int
}

// A constructor function that seeds a new blockchain i.e creates one.
// Returns an error if an Animus Blockchain already exists.
func SeedBlockChain(address wallet.Address) (*BlockChain, error) {
	// Check if a blockchain already exists by checking if the DB exists
	if persistence.DBExists() {
		return &BlockChain{}, fmt.Errorf("blockchain already exists exist")
	}

	// Create a null blockchain
	blockchain := BlockChain{}
	// Set up the database client
	blockchain.DB = persistence.NewDatabaseClient()

	// Generate a coinbase transaction for the genesis block
	coinbase := NewCoinbaseTransaction(address)

	// Create a merkle builder
	merkletree := merkle.NewMerkleTree()
	// Start the merkle builder
	go merkletree.Build()
	// Send the coinbase transaction to the merkle build queue
	merkletree.BuildQueue <- coinbase
	// Close the build queue
	close(merkletree.BuildQueue)

	// Generate a Genesis Block for the chain with a coinbase transaction
	genesisblock := NewBlock(merkletree, []byte{}, 0, address)
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
		if err := txn.Set(genesisblock.BlockHash, genesisblock.Serialize()); err != nil {
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
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to seed chain.")
	}

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
	if !persistence.DBExists() {
		return &BlockChain{}, fmt.Errorf("blockchain does not exist")
	}

	// Create a null blockchain
	blockchain := BlockChain{}
	// Set up the database client
	blockchain.DB = persistence.NewDatabaseClient()

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
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to animate chain.")
	}

	// Assign the current chain head
	blockchain.ChainHead = chainhead
	// Assign the current chain height
	blockchain.ChainHeight = utils.HexDecode(chainheight)

	// Return the blockchain
	return &blockchain, nil
}

// A method of BlockChain that adds a new Block to the chain and returns it
func (chain *BlockChain) AddBlock(blocktxns []*Transaction, addr wallet.Address) *Block {
	// Create a merkle builder
	merkletree := merkle.NewMerkleTree()
	// Start the merkle builder
	go merkletree.Build()

	// Iterate over the block transactions
	for _, txn := range blocktxns {
		// Send each transaction to the merkle build queue
		merkletree.BuildQueue <- txn
	}
	// Close the build queue
	close(merkletree.BuildQueue)

	// Generate a new Block
	block := NewBlock(merkletree, chain.ChainHead, chain.ChainHeight, addr)

	// Assign the hash of the block as the chain head
	chain.ChainHead = block.BlockHash
	// Increment the chain height
	chain.ChainHeight++

	// Define an Update Transaction on the BadgerDB
	err := chain.DB.Client.Update(func(txn *badger.Txn) error {

		// Add the Block to the DB with its hash as the key and its gob data as the value
		if err := txn.Set(block.BlockHash, block.Serialize()); err != nil {
			// Return any potential error
			return fmt.Errorf("block data could not be stored! error - %v", err)
		}

		// Set the last hash key of the database to the hash of the block
		if err := txn.Set([]byte("chainhead"), chain.ChainHead); err != nil {
			// Return any potential error
			return fmt.Errorf("chain head could not be updated! error - %v", err)
		}

		// Set the last hash key of the database to the hash of the block
		if err := txn.Set([]byte("chainheight"), utils.HexEncode(chain.ChainHeight)); err != nil {
			// Return any potential error
			return fmt.Errorf("chain height could not be updated! error - %v", err)
		}

		return nil
	})

	// Handle any potential errors
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to add block to chain.")
	}

	// Return the block
	return block
}

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
