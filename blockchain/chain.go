package blockchain

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
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

// A method of BlockChain that accumulates all unspent transaction
// for a given address and returns them as a slice of Transactions.
func (chain *BlockChain) AccumulateUTXN(address string) []Transaction {
	// Define the slice of unspent transactions
	var unspenttxns []Transaction
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

				// Check if the transaction output can be unlocked for the given address
				if output.CanBeUnlocked(address) {
					// Add it to the list of unspent transactions
					unspenttxns = append(unspenttxns, *tx)
				}
			}

			// Check if the transaction is a coinbase transaction
			if !tx.IsCoinbaseTx() {
				// Iterate over the transaction's inputs
				for _, input := range tx.Inputs {
					// Check if the transaction input can unlock for the given address
					if input.CanUnlock(address) {
						// Encode the ID of the transaction input (hash of the reference transaction output)
						inputtxid := hex.EncodeToString(input.ID)
						// Add the output index of the transaction input to spent transactions map
						spenttxos[inputtxid] = append(spenttxos[inputtxid], input.OutIndex)
					}
				}
			}
		}

		// Check if the block is the genesis block and break from the loop
		if len(block.PrevHash) == 0 {
			break
		}
	}

	// Return the accumulated unspent transactions
	return unspenttxns
}

// A method of BlockChain that accumulates all unspent transaction outputs.
// for a given address and returns them as a slice of TxOutputs.
func (chain *BlockChain) AccumulateUTXO(address string) []TxOutput {
	// Declare a slice of transaction outputs
	var unspenttxos []TxOutput
	// Accumulate the unspent transactions of the address
	unspenttxns := chain.AccumulateUTXN(address)

	// Iterate over the unspent transactions
	for _, tx := range unspenttxns {
		// Iterate over transaction's outputs
		for _, output := range tx.Outputs {
			// Check if the transaction ouput can be unlocked by the address
			if output.CanBeUnlocked(address) {
				// Add the unspent transaction output to the list
				unspenttxos = append(unspenttxos, output)
			}
		}
	}

	// Return the accumulated unspent transaction outputs
	return unspenttxos
}

// A method of BlockChain that accumulates unspent transaction outputs for a given
// address until a given amount and returns the accumulated amount and the map of
// transaction output IDs to their output index on the transaction.
func (chain *BlockChain) AccumulateSpendableTXO(address string, amount int) (int, map[string][]int) {
	// Declare a map to collect unspent transaction outputs
	unspenttxos := make(map[string][]int)
	// Accumulate the unspent transactions of the address
	unspenttxns := chain.AccumulateUTXN(address)
	// Declare an integer to accumulate output values
	accumulated := 0

Work:
	// Iterate over the unspent transactions
	for _, tx := range unspenttxns {
		// Encode the transaction hash into a string
		txid := hex.EncodeToString(tx.ID)

		// Iterate over the transaction outputs
		for outindex, output := range tx.Outputs {
			// Check if the output can be unlocked by the address and if the accumulated value is less than the amount
			if output.CanBeUnlocked(address) && accumulated < amount {
				// Add the value of the transaction output to the accumulation
				accumulated += output.Value
				// Add the transaction output id and index to the unspent transaction output map
				unspenttxos[txid] = append(unspenttxos[txid], outindex)

				// Check if the accumulated value has exceeded the amount
				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	// Return the accumulated value and the map of transaction outputs to their indexes
	return accumulated, unspenttxos
}
