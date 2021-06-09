package primitives

import (
	"sync"

	"github.com/manishmeganathan/blockweave/src/utils"
	"github.com/sirupsen/logrus"
)

// A structure that represents a Node on the Merkle Tree
type MerkleNode struct {
	Left  Hash // Represents the hash data of the left child
	Right Hash // Represents the hash data of the right child
	Data  Hash // Represents the hash data of the merkle node
}

// A structure that represents a Merkle Tree Builder
type MerkleBuilder struct {
	// Represents the channel that accepts transactions to add to the Merkle Tree
	BuildQueue chan *Transaction

	// Represents the wait group for the tree builder tasks
	BuildGroup sync.WaitGroup

	// Represents the root hash of the Merkle Tree
	MerkleRoot Hash

	// Represents the Transaction inside the Merkle Tree
	Transactions []*Transaction

	// Represents the number of Transaction inside the Merkle Tree
	Count int
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

// A constructor function that generates and returns
// a MerkleBuilder with its initialized fields
func NewMerkleBuilder() *MerkleBuilder {
	return &MerkleBuilder{
		BuildQueue: make(chan *Transaction),
		BuildGroup: sync.WaitGroup{},
		MerkleRoot: nil,
	}
}

// A method of MerkelBuilder that begins the merkle tree construction
// for a given a slice of Transactions. Internally handles trigger for
// the build runtime and closing the buildqueue.
func (mb *MerkleBuilder) BuildWithTransactions(txns []*Transaction) {
	// Start the build runtime
	go mb.Build()

	// Iterate over the transactions
	for _, txn := range txns {
		// Feed the transactions in the BuildQueue
		mb.BuildQueue <- txn
	}

	// Close the BuildQueue
	close(mb.BuildQueue)
}

// A method of MerkleBuilder that begins the construction of the merkle tree
// based on the Transaction recieved on its build queue. The transaction are
// accumulated into the structure and the resulting merkle root is stored.
// Wait on the BuildGroup field to confirm the build completion
func (mb *MerkleBuilder) Build() {
	// Declare a slice of MerkleNodes
	var nodes []MerkleNode
	// Set the BuildGroup counter for two process layers
	mb.BuildGroup.Add(2)

	// Iterate over the BuildQueue
	for lefttxn := range mb.BuildQueue {
		// Add the left transaction to the Merkle builder's transactions
		mb.Transactions = append(mb.Transactions, lefttxn)

		// Collect another transaction from the BuildQueue
		righttxn, ok := <-mb.BuildQueue
		// Check if the channel has closed and value is nil
		if !ok {
			// Copy the right transaction from the left transaction.
			// This ensures an even number of merkle leaves
			righttxn = lefttxn
		} else {
			// Add the right transaction to Merkle builder's transaction. This
			// is only done if the right transaction is not copied from the left
			mb.Transactions = append(mb.Transactions, righttxn)
		}

		// Generate a MerkleNode for the transaction pair (as a base node)
		merklenode := NewMerkleNode(
			TxnSerialize(lefttxn),
			TxnSerialize(righttxn),
			true)

		// Add the base merkle node to slice
		nodes = append(nodes, *merklenode)
	}

	// Assign the transaction count
	mb.Count = len(mb.Transactions)
	// Decrement the BuildGroup counter
	mb.BuildGroup.Done()

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
		logrus.Fatal("merkle builder logic failed!")
	}

	// Set the merkle builder's root
	mb.MerkleRoot = nodes[0].Data
	/// Decrement the BuildGroup counter (completes build)
	mb.BuildGroup.Done()
}
