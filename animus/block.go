package animus

import (
	"bytes"
	"crypto/sha256"
)

// A structure that represents a single
// Block on the Animus BlockChain
type Block struct {
	// Represents the hash of the Block
	Hash []byte

	// Represents the data of the Block
	Data []byte

	// Represents the hash of the previous Block
	PrevHash []byte
}

// A constructor function that generates a new block
// from hash of the previous block and a block data.
func NewBlock(data string, prevHash []byte) *Block {
	// Construct a new Block and asign the block data and hash of the previous block
	block := Block{Data: []byte(data), PrevHash: prevHash}
	// Generate the hash of the block
	block.GenerateHash()
	// Return the block
	return &block
}

// A method of Block that generates the hash for the block
func (b *Block) GenerateHash() {
	// Combine the block data and hash of the prev block
	hashinfo := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	// Generate the SHA256 checksum of the hashinfo
	hash := sha256.Sum256(hashinfo)
	// Assign the hash to the Block
	b.Hash = hash[:]
}
