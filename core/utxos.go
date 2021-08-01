package core

import (
	"bytes"
	"encoding/hex"

	"github.com/dgraph-io/badger"
	"github.com/manishmeganathan/blockweave/utils"
	"github.com/sirupsen/logrus"
)

// A method of BlockChain that accumulates all unspent transactions on
// the chain and returns them as a map transaction ID to TXOList.
func (chain *BlockChain) AccumulateUTX0S() map[string]TXOList {
	// Define a slice of UTXOs
	utxos := make(map[string]TXOList)
	// Define a map to store spent transaction outputs
	spenttxos := make(map[string][]int)

	// Get an iterator for the blockchain and iterate over its block
	iter := NewIterator(chain)
	for {
		// Get a block from the iterator
		block := iter.Next()

		// Iterate over the transactions in the block
		for _, tx := range block.TXList {
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
			if !tx.IsCoinbase() {
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
		if block.BlockHeight == 0 {
			break
		}
	}

	// Return the accumulated unspent transactions list
	return utxos
}

// A method of BlockChain that collects the spendable transaction outputs
// given a public key hash and a target amount upto which to collect.
func (chain *BlockChain) CollectSpendableUTXOS(publickeyhash []byte, amount int) (int, map[string][]int) {
	// Create a map of strings to a slice of ints
	unspenttxos := make(map[string][]int)
	// Declare an accumulation integer
	accumulated := 0

	// Define a View transaction on the database
	_ = chain.DB.Client.View(func(txn *badger.Txn) error {
		// Start a database iterator with the default options
		dbiterator := txn.NewIterator(badger.DefaultIteratorOptions)
		// Defer the closing of the database
		defer dbiterator.Close()

		// Iterate over the database elements that are utxo items
		for dbiterator.Seek(utils.UTXOprefix); dbiterator.ValidForPrefix(utils.UTXOprefix); dbiterator.Next() {
			// Retrieve an item from the database iterator
			item := dbiterator.Item()

			// Retrieve the key of the item
			key := item.Key()
			// Trime the key to not have the utxo prefix
			key = bytes.TrimPrefix(key, utils.UTXOprefix)
			// Encode the key into the transaction ID
			txnid := hex.EncodeToString(key)

			// Declare a transaction output list
			var outputs TXOList
			// Retrieve the value of the item and
			// deserialize it into the output list
			_ = item.Value(func(val []byte) error {
				outputs.Deserialize(val)
				return nil
			})

			// Iterate over the transaction output list
			for outindex, output := range outputs {
				// Checl if the transaction output is locked by the public key
				// and ensure that accumulation has not reached the amount target
				if output.CheckLock(publickeyhash) && accumulated < amount {
					// Add the value of the output into the accumulation
					accumulated += output.Value
					// Add the transaction output's ID and index to the map
					unspenttxos[txnid] = append(unspenttxos[txnid], outindex)
				}
			}
		}
		// Return a nil error
		return nil
	})

	// Return the accumulated amount and the list of unspent transactions
	return accumulated, unspenttxos
}

// A method of BlockChain that fetches the unspent transaction outputs for
// a given public key hash and returns it in a list of transaction outputs
func (chain *BlockChain) FetchUTXOS(publickeyhash []byte) TXOList {
	// Declare a transaction output list
	var utxos TXOList

	// Define a View transaction on the database
	_ = chain.DB.Client.View(func(txn *badger.Txn) error {
		// Start a database iterator with the default options
		dbiterator := txn.NewIterator(badger.DefaultIteratorOptions)
		// Defer the closing of the database
		defer dbiterator.Close()

		// Iterate over the database elements that are utxo items
		for dbiterator.Seek(utils.UTXOprefix); dbiterator.ValidForPrefix(utils.UTXOprefix); dbiterator.Next() {
			// Retrieve the iterator item
			item := dbiterator.Item()
			// Declare a transaction output list
			var txolist TXOList

			// Retrieve the value of the item and deserialize
			// into a transaction output list
			_ = item.Value(func(val []byte) error {
				txolist.Deserialize(val)
				return nil
			})

			// Iterate over the transaction output list
			for _, output := range txolist {
				// Check if the transaction output is locked by the public key
				if output.CheckLock(publickeyhash) {
					// Add the transaction output to the list
					utxos = append(utxos, output)
				}
			}
		}
		// Return a nil error
		return nil
	})

	// Return the list of unspent transaction outputs
	return utxos
}

// A method of BlockChain that counts the number
// of unspent transactions stored on the database
func (chain *BlockChain) CountUTXOS() int {
	// Declare a counter integer
	counter := 0

	// Define a View transaction on the database
	_ = chain.DB.Client.View(func(txn *badger.Txn) error {
		// Start a database iterator with the default options
		dbiterator := txn.NewIterator(badger.DefaultIteratorOptions)
		// Defer the closing of the database
		defer dbiterator.Close()

		// Iterate over the database elements that are utxo items
		for dbiterator.Seek(utils.UTXOprefix); dbiterator.ValidForPrefix(utils.UTXOprefix); dbiterator.Next() {
			// Increment the counter for each item
			counter++
		}
		// Return nil error
		return nil
	})

	// Return the counter value
	return counter
}

// A method of BlockChain that reindexes
// all the utxo layer keys on the database.
func (chain *BlockChain) ReindexUTXOS() {
	// Delete all the UTXOs stored on the database
	chain.DB.DeleteKeyPrefix(utils.UTXOprefix)
	// Accumulate all the UTXOs on the blockchain
	utxos := chain.AccumulateUTX0S()

	// Define an Update transaction on the database
	err := chain.DB.Client.Update(func(txn *badger.Txn) error {
		// Iterate over the UTXOs map
		for txid, txolist := range utxos {
			// Decode the transaction ID
			key, err := hex.DecodeString(txid)
			if err != nil {
				// Return an error if any
				return err
			}

			// Construct the key by adding the UTXO key prefix
			key = append(utils.UTXOprefix, key...)
			// Add the TXOList to the database with the key
			err = txn.Set(key, txolist.Serialize())
			if err != nil {
				// Log a fatal error
				logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to reindex utxos.")
			}
		}

		// Return nil error
		return nil
	})

	// Handle any potential error
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to reindex utxos.")
	}
}

// A method of BlockChain that updates the utxo layer keys
// from the transaction of a Block, given the block.
func (chain *BlockChain) UpdateUTXOS(block *Block) {
	// Define an Update transaction on the database
	err := chain.DB.Client.Update(func(dbtxn *badger.Txn) error {
		// Iterate over the transactions in the block
		for _, txn := range block.TXList {
			// Verify that transaction is not a coinbase
			if !txn.IsCoinbase() {
				// Iterate over the transaction inputs
				for _, input := range txn.Inputs {
					// Create an empty transaction output list
					updatedouts := TXOList{}

					// Create the input ID from the utxo
					// prefix and ID of the transaction input
					inputid := append(utils.UTXOprefix, input.ID...)

					// Retrieve a utxo item from the database
					item, err := dbtxn.Get(inputid)
					if err != nil {
						// Log a fatal error
						logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to update utxos.")
					}

					// Declare a transaction output list
					var outputs TXOList
					// Retrieve the value of the utxo item and
					// deserialize it into the output list
					err = item.Value(func(val []byte) error {
						outputs.Deserialize(val)
						return nil
					})
					// Handle any potential errors
					if err != nil {
						// Log a fatal error
						logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to update utxos.")
					}

					// Iterate over the transaction output list
					for outindex, output := range outputs {
						// Update the transaction outputs for the transaction ID from the input
						if outindex != input.OutIndex {
							updatedouts = append(updatedouts, output)
						}
					}

					// Check if there are any transactions in the updated list
					if len(updatedouts) == 0 {
						// Delete all transaction outputs in the utxo item on the db
						if err := dbtxn.Delete(inputid); err != nil {
							// Log a fatal error
							logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to update utxos.")
						}
					} else {
						// Set the utxo item to the updated list of transaction outputs
						if err := dbtxn.Set(inputid, updatedouts.Serialize()); err != nil {
							// Log a fatal error
							logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to update utxos.")
						}
					}
				}
			}

			// Create a new transaction output list
			newoutputs := TXOList{}
			// Accumulate the transaction outputs to the list
			newoutputs = append(newoutputs, txn.Outputs...)

			// Create the utxo item key from the utxo prefix and transaction ID
			txnid := append(utils.UTXOprefix, txn.ID...)
			// Add the list of transaction outputs to the db
			if err := dbtxn.Set(txnid, newoutputs.Serialize()); err != nil {
				// Log a fatal error
				logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to update utxos.")
			}
		}
		// Return a nil error
		return nil
	})

	// Handle any potential error
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to update utxos.")
	}
}
