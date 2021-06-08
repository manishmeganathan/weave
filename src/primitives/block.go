/*
This module contains the definition and implementation
of the Block structure and its methods
*/
package primitives

var version = []byte("v0.1.0")

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
	Timestamp uint64

	// Represent the merkle root of transactions on the Block
	MerkleRoot Hash

	// Represents the weave network of the Block
	WeaveNet []byte

	// Represents the software version of Block Generator
	Version []byte
}

// A structure that represents a single Block on the Blockchain
type Block struct {
	// Represents the Block Header
	Header BlockHeader

	// Represents the Block Height
	BlockHeight uint64

	// Represents the no.of Transactions in the Block
	TXCount uint

	// Represents the list of Transactions in the Block
	TXList []*Transaction

	// Represents the address of the Block origin (coinbase address)
	BlockOrigin Address

	// Represents the Hash of the Block Header
	BlockHash Hash
}

func NewBlock() {}
