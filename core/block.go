package core

import (
	"encoding/gob"

	"github.com/manishmeganathan/blockweave/consensus"
	"github.com/manishmeganathan/blockweave/merkle"
	"github.com/manishmeganathan/blockweave/utils"
	"github.com/manishmeganathan/blockweave/wallet"
)

// A structure that represents a single Block on the Blockchain
type Block struct {
	// Represents the Block Header
	BlockHeader

	// Represents the Hash of the Block Header
	BlockHash utils.Hash

	// Represents the Block Height
	BlockHeight int

	// Represents the address of the Block origin (miner address)
	BlockOrigin wallet.Address

	// Represents the no.of Transactions in the Block
	TXCount int

	// Represents the list of Transactions in the Block
	TXList []*Transaction
}

// A constructor function that generates and returns a new Block
// that has been minted for a given merkle builder, previous block
// hash, block height and a coinbase address.
func NewBlock(merkletree *merkle.MerkleTree, priori utils.Hash, height int, origin wallet.Address) *Block {

	// Create and empty Block
	block := Block{}

	// Set the BlockHash to nil
	block.BlockHash = nil
	// Set the block height
	block.BlockHeight = height
	// Set the block origin address
	block.BlockOrigin = origin

	// Wait fot the merkle builder to finish building
	merkletree.BuildGroup.Wait()

	txns := make([]*Transaction, merkletree.Count)
	for i, item := range merkletree.Items {
		txns[i] = item.(*Transaction)
	}

	// Extract transactions from the merkle builder
	block.TXList = txns
	// Extract transaction count from the merkle builder
	block.TXCount = merkletree.Count

	// Create and assign the block header
	block.BlockHeader = *NewBlockHeader(priori, merkletree.MerkleRoot)
	// Set the Consensus Header to Proof Of Work
	block.BlockHeader.ConsensusHeader = consensus.NewPOW()
	// Mint the block (sign)
	block.BlockHash = block.Mint(&block.BlockHeader)

	// Return the signed block
	return &block
}

// A construcor function that generates and returns a null Block
// This a temporary constructor until a deeper integration with the consensus package is implemented.
func NullBlock() *Block {
	// Create an empty block object
	block := &Block{}
	// Set the consensus header to null pow block
	block.BlockHeader.ConsensusHeader = consensus.NewPOW()
	// Return the block
	return block
}

// A method that returns the gob encoded data of the Block
func (block *Block) Serialize() utils.Gob {
	// Register the gob library with the Consensus Header type
	gob.Register(block.BlockHeader.ConsensusHeader)
	// Encode the block as a gob and return it
	return utils.GobEncode(block)
}

// A method that decodes a gob of bytes into the Block struct
func (block *Block) Deserialize(gobdata utils.Gob) {
	// Register the gob library with the Consensus Header type
	gob.Register(block.BlockHeader.ConsensusHeader)
	// Decode the gob data into the block
	utils.GobDecode(gobdata, block)
}
