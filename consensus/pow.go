package consensus

import (
	"fmt"
	"math"
	"math/big"

	"github.com/manishmeganathan/blockweave/src/primitives"
	"github.com/manishmeganathan/blockweave/src/utils"
)

// A value that represents the difficulty value with a max value of 255
// Currently fixed but eventually will change based on the network hashrate.
var WorkDifficulty uint8 = 20

// A structure that represents the Proof Of Work consensus
// header that implements the ConsensusHeader interface
type POW struct {
	// Represents the target value of POW algorithm
	Target *big.Int

	// Represents the nonce of the block upon being minted
	Nonce int
}

// A constructor function that generates and return a POW
// with its target value for the algorithm set.
func NewPOW() *POW {
	// Create a new POW
	pow := POW{Nonce: 0, Target: big.NewInt(1)}
	// Generate the difficulty target
	pow.GenerateTarget()
	//Return the POW
	return &pow
}

// A method of POW that generates the target value
// for the POW algorithm to mint the block.
func (pow *POW) GenerateTarget() {
	// Generate new big integer with value 1
	target := big.NewInt(1)
	// Left Shift the big integer by the difference between the max hash
	// size and the block's work difficulty. target = 2^(256-difficulty)
	target.Lsh(target, 256-uint(WorkDifficulty))

	// Assign the target to the POW
	pow.Target = target
}

// A method of POW that runs the Proof Of Work Algorithm
// to generate the hash of the block and mint it.
func (pow *POW) Mint(block *primitives.Block) error {
	// Declare a big Int version of the hash
	var inthash big.Int
	// Declare an slice of bytes for the hash
	var hash []byte
	// Reset the Nonce
	pow.Nonce = 0

	// Iterate until nonce reaches the maximum int value
	for pow.Nonce < math.MaxInt64 {
		// Compose the block data
		data := primitives.BlockHeaderSerialize(&block.BlockHeader)
		// Generate the hash256 for the composed data
		hash = utils.Hash256(data)

		// Print the hash (with a line reset)
		fmt.Printf("\r%x", hash)

		// Set the inthash with the hash
		inthash.SetBytes(hash)

		// Check if the inthash is lesser than the proof target
		if inthash.Cmp(pow.Target) == -1 {
			// Block Minted! Break from the loop
			break
		} else {
			// Increment the nonce and retry
			pow.Nonce++
		}
	}

	// Print an empty line for spacing
	fmt.Println()
	// Return a nil error
	return nil
}

// A method of POW that validates the block data for the target
func (pow *POW) Validate(block *primitives.Block) bool {
	// Declare a big Int version of the hash
	var inthash big.Int
	// Compose the block data
	data := primitives.BlockHeaderSerialize(&block.BlockHeader)
	// Generate the hash256 for the composed data
	hash := utils.Hash256(data)
	// Set the inthash with the hash
	inthash.SetBytes(hash)

	// Check if the inthash is lesser than the proof target
	// If the hash of the block data with the given nonce is
	// less than the proof target, the block signature is valid.
	return inthash.Cmp(pow.Target) == -1
}
