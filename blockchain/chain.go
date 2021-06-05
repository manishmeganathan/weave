package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
)

var (
	utxoprefix   = []byte("utxo-")
	prefixlength = len(utxoprefix)
)

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

// A constructor function that animates an existing blockchain i.e brings it to life.
// Returns an error if no Animus Blockchain exists.
func AnimateBlockChain() (*BlockChain, error) {
	// Declare a slice a bytes to collect the hash value
	var lasthash []byte

	// Check if a blockchain already exists by checking if the DB exists
	if !DBexists() {
		return &BlockChain{}, fmt.Errorf("blockchain does not exist")
	}

	// Open the BadgerDB
	db := DBopen()

	// Define an Update Transaction on the BadgerDB
	err := db.Update(func(txn *badger.Txn) error {
		// Get the last hash item from the DB
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

	// Handle any potential errors
	Handle(err)

	// Construct a blockchain with the BadgerDB and the last hash of the chain
	blockchain := BlockChain{Database: db, LastHash: lasthash}
	// Return the blockchain
	return &blockchain, nil
}

// A constructor function that seeds a new blockchain i.e creates one.
// Returns an error if an Animus Blockchain already exists.
func SeedBlockChain(address string) (*BlockChain, error) {
	// Declare a slice a bytes to collect the hash value
	var lasthash []byte

	// Check if a blockchain already exists by checking if the DB exists
	if DBexists() {
		return &BlockChain{}, fmt.Errorf("blockchain already exists exist")
	}

	// Open the BadgerDB
	db := DBopen()

	// Define an Update Transaction on the BadgerDB
	err := db.Update(func(txn *badger.Txn) error {
		// Generate a coinbase transaction for the genesis block
		coinbase := NewCoinbaseTransaction(address, "First Transaction from Genesis")
		// Generate a Genesis Block for the chain with a coinbase transaction
		genesisblock := NewBlock([]*Transaction{coinbase}, []byte{})
		log.Println("Genesis Block Signed!")

		// Add the Block to the DB with its hash as the key and its gob data as the value
		err := txn.Set(genesisblock.Hash, BlockSerialize(genesisblock))
		// Handle any potential error
		Handle(err)

		// Retrieve the hash of the Genesis Block
		lasthash = genesisblock.Hash
		// Set the last hash of the chain in the DB to the Genesis Block's Hash
		err = txn.Set([]byte("lh"), lasthash)
		return err
	})

	// Handle any potential errors
	Handle(err)

	// Construct a blockchain with the BadgerDB and the last hash of the chain
	blockchain := BlockChain{Database: db, LastHash: lasthash}
	// Return the blockchain
	return &blockchain, nil
}

// A method of BlockChain that adds a new Block to the chain
func (chain *BlockChain) AddBlock(blocktxns []*Transaction) {
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
	// Generate a new Block from the given block transactions and the hash of the previous block
	block := NewBlock(blocktxns, lasthash)

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

// A method of BlockChain that finds a transaction
// from the chain given a valid Transaction ID
func (chain *BlockChain) FindTransaction(txnid []byte) (Transaction, error) {

	// Get an iterator for the blockchain and iterate over its block
	iter := NewIterator(chain)
	for {
		// Get a block from the iterator
		block := iter.Next()

		// Iterate over the transactions of the block
		for _, txn := range block.Transactions {
			// Check if the transaction ID matches
			if bytes.Equal(txn.ID, txnid) {
				// Return the transaction with a nil error
				return *txn, nil
			}
		}

		// Check if the block is last block on the chain
		if len(block.PrevHash) == 0 {
			break
		}
	}

	// Return a nil Transaction with an error
	return Transaction{}, fmt.Errorf("transaction does not exist")
}

// A method of BlockChain that signs a transaction given a private key
func (chain *BlockChain) SignTransaction(txn *Transaction, privatekey ecdsa.PrivateKey) {
	// Create a map of transaction IDs to Transactions
	prevtxns := make(map[string]Transaction)

	// Iterate over the inputs of the transaction
	for _, input := range txn.Inputs {
		// Find the Transaction with ID on the input from the blockchain
		prevtxn, err := chain.FindTransaction(input.ID)
		// Handle any potential errors
		Handle(err)

		// Add the transaction to the map
		prevtxns[hex.EncodeToString(prevtxn.ID)] = prevtxn
	}

	// Sign the transaction with the map of previous
	// transactions and the ECDSA private key
	txn.Sign(privatekey, prevtxns)
}

// A method of BlockChain that verifies the signature of a transaction given a private key
func (chain *BlockChain) VerifyTransaction(txn *Transaction, privatekey ecdsa.PrivateKey) bool {
	// Create a map of transaction IDs to Transactions
	prevtxns := make(map[string]Transaction)

	// Iterate over the inputs of the transaction
	for _, input := range txn.Inputs {
		// Find the Transaction with ID on the input from the blockchain
		prevtxn, err := chain.FindTransaction(input.ID)
		// Handle any potential errors
		Handle(err)

		// Add the transaction to the map
		prevtxns[hex.EncodeToString(prevtxn.ID)] = prevtxn
	}

	// Verify the transaction signature and return the result
	return txn.Verify(prevtxns)
}

// A method of BlockChain that accumulates all unspent transactions on
// the chain and returns them as a map transaction ID to TXOList.
func (chain *BlockChain) AccumulateUTX0() map[string]TXOList {
	// Define the slice of unspent transactions
	utxos := make(map[string]TXOList)
	// Define a map to store spent transaction outputs
	spenttxos := make(map[string][]int)

	// Get an iterator for the blockchain and iterate over its block
	iter := NewIterator(chain)
	for {
		// Get a block from the iterator
		block := iter.Next()

		// Iterate over the transactions in the block
		for _, tx := range block.Transactions {
			// Encode the transaction hash into a string
			txid := hex.EncodeToString(tx.ID)

		Outputs:
			// Iterate over the transaction's outputs
			for outindex, output := range tx.Outputs {
				// Check if the transaction outputs have been spent
				if spenttxos[txid] != nil {
					// Iterate over the index of the spent transaction outputs
					for _, spentout := range spenttxos[txid] {
						// Check if the spent transaction output is the current transaction output
						if spentout == outindex {
							// Break from loop and check next transaction output
							continue Outputs
						}
					}
				}

				// Retrieve the transaction output list from the map for the transaction ID
				txolist := utxos[txid]
				// Add the output to the list
				txolist = append(txolist, output)
				// Update the list on the map
				utxos[txid] = txolist
			}

			// Check if the transaction is a coinbase transaction
			if !tx.IsCoinbaseTxn() {
				// Iterate over the transaction's inputs
				for _, input := range tx.Inputs {
					// Encode the ID of the transaction input (hash of the reference transaction output)
					inputtxid := hex.EncodeToString(input.ID)
					// Add the output index of the transaction input to spent transactions map
					spenttxos[inputtxid] = append(spenttxos[inputtxid], input.OutIndex)
				}
			}
		}

		// Check if the block is the genesis block and break from the loop
		if len(block.PrevHash) == 0 {
			break
		}
	}

	// Return the accumulated unspent transactions list
	return utxos
}

// A method of BlockChain that deletes all entries
// with a given prefix from the Badger DB.
func (chain *BlockChain) DeleteKeyPrefix(prefix []byte) {

	// Define a function that accepts a 2D slice of byte keys to delete
	DeleteKeys := func(keystodelete [][]byte) error {

		// Define an Update transaction on the database
		err := chain.Database.Update(func(txn *badger.Txn) error {
			// Iterate over the keys to delete
			for _, key := range keystodelete {
				// Delete the key
				if err := txn.Delete(key); err != nil {
					// Return any potential error
					return err
				}
			}
			// Return nil error
			return nil
		})

		// Check if key deletion transaction has comlpeted without error
		if err != nil {
			// Return the error
			return err
		}
		// Return nil error
		return nil
	}

	// Define the size limit of key accumulation. This value
	// is based on BadgerDBs optimal object handling limit.
	collectlimit := 100000

	// Define a View transaction on the database
	chain.Database.View(func(txn *badger.Txn) error {

		// Set up the default iteration options for the database
		opts := badger.DefaultIteratorOptions
		// Set value pre-fetching to off (We only need check the key's prefix to accumulate it)
		opts.PrefetchValues = false
		// Create an iterator with the options
		dbiterator := txn.NewIterator(opts)
		// Defer the closing of the iterator
		defer dbiterator.Close()

		// Declare collection counter
		keyscollected := 0
		// Create a 2D slice of bytes to collect keys
		keystodelete := make([][]byte, 0)

		// Start the iterator and seek the keys with the provided prefix (validate the keys for the prefix as well)
		for dbiterator.Seek(prefix); dbiterator.ValidForPrefix(prefix); dbiterator.Next() {

			// Make a copy of the key value
			key := dbiterator.Item().KeyCopy(nil)
			// Add the key to slice
			keystodelete = append(keystodelete, key)
			// Increment the counter
			keyscollected++

			// Check if counter has reached the collection limit
			if keyscollected == collectlimit {
				// Delete all keys accumulated so far
				if err := DeleteKeys(keystodelete); err != nil {
					// Handle any potential error
					Handle(err)
				}

				// Reset the key accumulation
				keystodelete = make([][]byte, 0)
				// Reset the key accumulation counter
				keyscollected = 0
			}
		}

		// Check if there any keys that have been collected (but less than the collection limit)
		if keyscollected > 0 {
			// Delete all the accumulated keys
			if err := DeleteKeys(keystodelete); err != nil {
				// Handle any potential errors
				Handle(err)
			}
		}

		// Return a nil error for the transaction
		return nil
	})
}
