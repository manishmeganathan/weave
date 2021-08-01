package merkle

import (
	"fmt"
	"sync"

	"github.com/manishmeganathan/blockweave/utils"
	"github.com/sirupsen/logrus"
)

// A structure that represents a Merkle Tree
type MerkleTree struct {
	// Represents the root hash of the Merkle Tree
	MerkleRoot utils.Hash

	// Represents the items inside the Merkle Tree
	// Items must implement the utils.GobEncodable interface
	Items []utils.GobEncodable

	// Represents the number of Items inside the Merkle Tree
	Count int

	// Represents the channel that accepts items to add to the Merkle Tree
	BuildQueue chan utils.GobEncodable

	// Represents the wait group for the tree builder tasks
	BuildGroup *sync.WaitGroup
}

// A constructor function that generates and returns a null MerkleTree
func NewMerkleTree() *MerkleTree {
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(1)

	return &MerkleTree{
		BuildQueue: make(chan utils.GobEncodable),
		BuildGroup: waitgroup,
		MerkleRoot: nil,
	}
}

// A method of MerkleTree that builds a full tree from a slice of encodable items.
// Internally builds the tree item by item and closes the build queue.
func (mt *MerkleTree) BuildFull(items []utils.GobEncodable) {
	// Start the build runtime
	go mt.Build()

	// Iterate over the items
	for _, item := range items {
		// Feed the item into the BuildQueue
		mt.BuildQueue <- item
	}

	// Close the BuildQueue
	close(mt.BuildQueue)
}

// A method of MerkleTree that begins the construction of the merkle tree
// based on the Items received on its build queue. The items are accumulated
// into the tree and the resulting merkle root is stored into the object.
// Wait on the BuildGroup field to confirm the build completion.
func (mt *MerkleTree) Build() {
	// Declare a slice of MerkleNodes
	var nodes []MerkleNode

	// Iterate over the BuildQueue
	for leftitem := range mt.BuildQueue {
		// Add the left item to the Merkle builder's items
		mt.Items = append(mt.Items, leftitem)

		// Collect another item from the BuildQueue
		rightitem, ok := <-mt.BuildQueue
		// Check if the channel has closed and value is nil
		if !ok {
			// Copy the right item from the left item.
			// This ensures an even number of merkle leaves.
			rightitem = leftitem
		} else {
			// Add the right ietm to Merkle builder's items.
			// This is only done if the right item is not copied from the left
			mt.Items = append(mt.Items, rightitem)
		}

		// Generate a MerkleNode for the item pair (as a base node)
		merklenode := NewMerkleNode(
			leftitem.Serialize(),
			rightitem.Serialize(),
			true,
		)

		// Add the base merkle node to the slice of nodes
		nodes = append(nodes, *merklenode)
	}

	// Assign the item count
	mt.Count = len(mt.Items)

	// Collect the intial size of the base nodes collection
	count := len(nodes)
	// Run loop for half the number of leaf nodes
	for i := 0; i < count/2; i++ {
		// Re/Set the level slice of Merkle Node
		var level []MerkleNode

		// Iterate over every two leaf nodes
		for j := 0; j < len(nodes); j += 2 {
			// Generate a MerkleNode from the left and right MerkleNodes
			node := NewMerkleNode(nodes[j].Data, nodes[j+1].Data, false)
			// Append the merklenode to the tree level
			level = append(level, *node)
		}

		// Set the full node list to the level nodes
		nodes = level
	}

	// Check if the final node list has just one node
	if len(nodes) != 1 {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": fmt.Errorf("build logic failure")}).Panicln("failed to build merkle tree")
	}

	// Set the merkle builder's root
	mt.MerkleRoot = nodes[0].Data
	/// Decrement the BuildGroup counter (completes build)
	mt.BuildGroup.Done()
}
