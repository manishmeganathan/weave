package core

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/manishmeganathan/blockweave/consensus"
	"github.com/manishmeganathan/blockweave/primitives"
	"github.com/manishmeganathan/blockweave/utils"
	"github.com/manishmeganathan/blockweave/wallet"
)

// A constructor function that generates and returns a new
// transaction output given a token value and address
func NewTxOutput(value int, address primitives.Address) *primitives.TXO {
	txo := primitives.TXO{Value: value, PublicKeyHash: nil}
	txo.Lock(address)

	return &txo
}

// A constructor function that generates and returns a Transaction
// given the to and from addresses and the amount to transact.
func NewTransaction(from, to primitives.Address, amount int, chain *BlockChain) *primitives.Transaction {
	// Declare slices of transaction outputs and inputs
	var txinputs primitives.TXIList
	var txoutputs primitives.TXOList

	// Create the wallet store
	walletstore := wallet.NewWalletStore()
	// Fetch the wallet from the wallet store for the given address
	w := walletstore.FetchWallet(from.String)

	// // Generate the public key hash for the wallet's public key
	// publickeyhash := wallet.GeneratePublicKeyHash(w.PublicKey)

	// Collect the spendable transaction outputs of the account up to the given amount with the public key hash
	accumulated, validoutputs := chain.CollectSpendableUTXOS(from.PublicKeyHash, amount)

	// Check if the account has enough funds
	if accumulated < amount {
		log.Panic("error: insufficient funds!")
	}

	// Iterate over the spendable transaction output IDs
	for txnid, outputs := range validoutputs {
		// Decode the transaction ID
		txid, _ := hex.DecodeString(txnid)

		// Iterate over the the output indexes
		for _, output := range outputs {
			// Create a transaction input with the transaction ID, output index and from address signature
			input := primitives.TXI{ID: txid, OutIndex: output, Signature: nil, PublicKey: w.PublicKey}
			// Add the transaction input into the slice
			txinputs = append(txinputs, input)
		}
	}

	// Add a transaction output with the amount to the address
	txoutputs = append(txoutputs, *NewTxOutput(amount, to))

	// Check if there is a balance in the accumulated amounted
	if accumulated > amount {
		// Add a transaction output with the balance amount back to the original sender
		txoutputs = append(txoutputs, *NewTxOutput(accumulated-amount, from))
	}

	// Create a Transaction with the list of input and outputs
	txn := primitives.Transaction{ID: nil, Inputs: txinputs, Outputs: txoutputs}
	// Set the ID (hash) for the transaction
	txn.ID = txn.GenerateHash()
	// Sign the transaction using the wallet's private key
	chain.SignTransaction(&txn, w.PrivateKey)

	// Return the transaction
	return &txn
}

// A constructor function that generates and returns a coinbase Transaction.
// A Coinbase transaction refers to a first transaction on a block and does not refer to any
// previous output transactions and contains a token reward for the user who signs the block.
func NewCoinbaseTransaction(to primitives.Address) *primitives.Transaction {
	// Create a slice a bytes
	randdata := make([]byte, 24)
	// Add random data to the slice of bytes
	_, err := rand.Read(randdata)
	// Handle any potential errors
	utils.HandleErrorLog(err, "coinbase transaction generation failed!")

	// Collect the data from the hexadecimal interpretation of the random bytes
	data := fmt.Sprintf("%x", randdata)

	// Create a transaction input with no reference to a previous output
	inputs := primitives.TXI{ID: []byte{}, OutIndex: -1, Signature: nil, PublicKey: []byte(data)}
	// Create a transaction output with the token reward
	outputs := *NewTxOutput(25, to)

	// Construct a transaction with no ID, and the set of inputs and outputs
	txn := primitives.Transaction{
		ID:      nil,
		Inputs:  primitives.TXIList{inputs},
		Outputs: primitives.TXOList{outputs},
	}

	// Set the ID (hash) for the transaction
	txn.ID = txn.GenerateHash()

	// Return the transaction
	return &txn
}

// A constructor function that generates and returns a BlockHeader
// for a given priori hash, merkle root and weave net address.
func NewBlockHeader(priori primitives.Hash, root primitives.Hash, weave []byte) *primitives.BlockHeader {
	// Generate and return the block header
	return &primitives.BlockHeader{
		// Assign the software version
		Version: []byte(utils.WeaverVersion),
		// Assign the weave network ID
		BlockWeave: weave,
		// Assign the block timestamp
		Timestamp: time.Now().Unix(),
		// Assign hash of the previous block
		Priori: priori,
		// Assign the merkle root hash
		MerkleRoot: root,
		// Assign a nil consensus header
		ConsensusHeader: nil,
	}
}

// A constructor function that generates and returns a new Block
// that has been minted for a given merkle builder, previous block
// hash, block height and a coinbase address.
func NewBlock(
	merkle *MerkleBuilder, priori primitives.Hash, height int,
	origin primitives.Address, weave []byte) *primitives.Block {

	// Create and empty Block
	block := primitives.Block{}

	// Set the BlockHash to nil
	block.BlockHash = nil
	// Set the block height
	block.BlockHeight = height
	// Set the block origin address
	block.BlockOrigin = origin

	// Wait fot the merkle builder to finish building
	merkle.BuildGroup.Wait()

	// Extract transactions from the merkle builder
	block.TXList = merkle.Transactions
	// Extract transaction count from the merkle builder
	block.TXCount = merkle.Count

	// Create and assign the block header
	block.BlockHeader = *NewBlockHeader(priori, merkle.MerkleRoot, weave)

	// Check the value of the weave network and assign the consensus header
	switch {
	case bytes.Equal(weave, []byte(utils.WeavePOW)):
		// TODO: Implement POW module
		// Set the Consensus Header to Proof Of Work
		block.BlockHeader.ConsensusHeader = consensus.NewPOW()

	default:
		log.Fatal("invalid weave at block header creation", string(weave))
	}

	// Mint the block (sign)
	block.Mint(&block)
	// Return the signed block
	return &block
}
