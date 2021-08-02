package merkle

import "github.com/manishmeganathan/weave/utils"

// A structure that represents a Node on the Merkle Tree
type MerkleNode struct {
	// Represents the hash data of the left child
	Left utils.Hash

	// Represents the hash data of the right child
	Right utils.Hash

	// Represents the hash data of the merkle node
	Data utils.Hash
}

// A constructor function that generates and returns a MerkleNode
// for a given pair of bytes payloads and flag that indicates if
// the generated MerkleNode is base node (no children/leaves)
func NewMerkleNode(leftdata, rightdata []byte, isbase bool) *MerkleNode {
	// Concatenate the left and right data
	data := append(leftdata, rightdata...)
	// Hash256 the accumulated data
	hash := utils.Hash256(data)

	// Declare a new MerkleNode
	var merklenode MerkleNode

	// Check the base generation flag
	if isbase {
		// Construct a base MerkleNodes that contains no children/leaves
		merklenode = MerkleNode{Left: nil, Right: nil, Data: hash}
	} else {
		// Construct a MerkleNode with the left, right and self hashes
		merklenode = MerkleNode{Left: leftdata, Right: rightdata, Data: hash}
	}

	// Return the merkle node
	return &merklenode
}
