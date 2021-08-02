package core

import (
	"encoding/gob"
	"time"

	"github.com/manishmeganathan/blockweave/utils"
)

// An interface for all types of consensus headers
type ConsensusHeader interface {
	Mint(utils.GobEncodable) utils.Hash
	Validate(utils.GobEncodable) bool
}

// A structure that represents the header of a Block
type BlockHeader struct {
	// Represents the consensus parameters of the Block
	ConsensusHeader

	// Represents the hash of the previous Block
	Priori utils.Hash

	// Represents the timestamp of the Block at the point of creation
	Timestamp int64

	// Represent the merkle root of transactions on the Block
	MerkleRoot utils.Hash

	// Represents the network version of Block
	Version []byte
}

// A constructor function that generates and returns a BlockHeader
// for a given priori hash, merkle root and weave net address.
func NewBlockHeader(priori, root utils.Hash) *BlockHeader {
	// Generate and return the block header
	return &BlockHeader{
		// Assign the software version
		Version: []byte(utils.SrcVersion),
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

// A method that returns the gob encoded data of the BlockHeader
func (bh *BlockHeader) Serialize() utils.Gob {
	// Register the gob library with the Consensus Header type
	gob.Register(bh.ConsensusHeader)
	// Encode the blockheader as a gob and return it
	return utils.GobEncode(bh)
}

// A method that decodes a gob of bytes into the BlockHeader struct
func (bh *BlockHeader) Deserialize(gobdata utils.Gob) {
	// Decode the gob data into the blockheader
	utils.GobDecode(gobdata, bh)
}
