package blockchain

import "crypto/sha256"

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}
	var hash [32]byte

	if !(left == nil && right == nil) {
		data = append(left.Data, right.Data...)
	}

	hash = sha256.Sum256(data)
	node.Data = hash[:]

	node.Left = left
	node.Right = right

	return &node
}

func NewMerkleTree(merkledata [][]byte) *MerkleTree {

	count := len(merkledata)
	var nodes []MerkleNode

	if count%2 != 0 {
		merkledata = append(merkledata, merkledata[count-1])
	}

	for _, data := range merkledata {
		node := NewMerkleNode(nil, nil, data)
		nodes = append(nodes, *node)
	}

	for i := 0; i < count/2; i++ {
		var level []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			level = append(level, *node)
		}

		nodes = level
	}

	tree := MerkleTree{RootNode: &nodes[0]}
	return &tree
}
