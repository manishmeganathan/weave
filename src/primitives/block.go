/*
This module contains the definition and implementation
of the Block structure and its methods
*/
package primitives

// A structure that represents a single Block on the Animus BlockChain
type Block struct {
	Hash         []byte         // Represents the hash of the Block
	PrevHash     []byte         // Represents the hash of the previous Block
	Transactions []*Transaction // Represents the transaction on the Block
	Nonce        int            // Represents the nonce number that signed the block
	Difficulty   int            // Represent the difficulty value to sign the block
}
