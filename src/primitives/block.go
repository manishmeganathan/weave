/*
This module contains the definition and implementation
of the Block structure and its methods
*/
package primitives

import (
	"bytes"
	"log"
	"time"
)

const libversion = "v0.5.0"
const fustianweave = "fustian" // Represents the proof of work weave net

// An interface for all types of consensus headers
type ConsensusHeader interface {
	Mint(*Block) error
	Validate(*Block) bool
}

// A structure that represents the header of a Block
type BlockHeader struct {
	// Represents the consensus parameters of the Block
	Consensus ConsensusHeader

	// Represents the hash of the previous Block
	Priori Hash

	// Represents the timestamp of the Block at the point of creation
	Timestamp int64

	// Represent the merkle root of transactions on the Block
	MerkleRoot Hash

	// Represents the weave network of the Block
	BlockWeave []byte

	// Represents the software version of Block Generator
	Version []byte
}

// A structure that represents a single Block on the Blockchain
type Block struct {
	// Represents the Block Header
	Header BlockHeader

	// Represents the Block Height
	BlockHeight int

	// Represents the no.of Transactions in the Block
	TXCount int

	// Represents the list of Transactions in the Block
	TXList []*Transaction

	// Represents the address of the Block origin (coinbase address)
	BlockOrigin Address

	// Represents the Hash of the Block Header
	BlockHash Hash
}

// A constructor function that generates and returns a BlockHeader
// for a given priori hash, merkle root and weave net address.
func NewBlockHeader(priori Hash, root Hash, weave []byte) *BlockHeader {
	// Create an empty BlockHeader
	header := BlockHeader{}

	// Assign the software version
	header.Version = []byte(libversion)
	// Assign the weave network ID
	header.BlockWeave = weave
	// Assign the block timestamp
	header.Timestamp = time.Now().Unix()
	// Assign hash of the previous block
	header.Priori = priori
	// Assign the merkle root hash
	header.MerkleRoot = root

	// Check the value of the weave network and assign the consensus header
	switch {
	case bytes.Equal(weave, []byte(fustianweave)):
		// TODO: Implement POW module
		// Set the Consensus Header to Proof Of Work
		header.Consensus = NewPOW()

	default:
		log.Fatal("invalid weave at block header creation", string(weave))
	}

	// Return the blockheader
	return &header
}

// A constructor function that generates and returns a new Block
// that has been minted for a given merkle builder, previous block
// hash, block height and a coinbase address.
func NewBlock(merkle *MerkleBuilder, priori Hash, height int, origin Address) *Block {
	// Create and empty Block
	block := Block{}

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
	block.Header = *NewBlockHeader(priori, merkle.MerkleRoot, []byte(fustianweave))

	// Mint the block (sign)
	block.Header.Consensus.Mint(&block)
	// Return the signed block
	return &block
}
