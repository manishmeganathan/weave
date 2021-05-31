package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/big"
)

const WorkDifficulty = 18

// A structure that represents a single Block on the Animus BlockChain
type Block struct {
	Hash       []byte // Represents the hash of the Block
	Data       []byte // Represents the data of the Block
	PrevHash   []byte // Represents the hash of the previous Block
	Nonce      int    // Represents the nonce number that signed the block
	Difficulty int    // Represent the difficulty value to sign the block
}

// A constructor function that generates and returns a new
// block from th hash of the previous block and a block data.
// Adds the data to a block and signs it using the PoW algorithm.
func NewBlock(data string, prevHash []byte) *Block {
	// Construct a new Block and assign the block data, the
	// hash of the previous block and the block difficulty
	block := Block{
		Data:       []byte(data),
		PrevHash:   prevHash,
		Difficulty: WorkDifficulty,
	}

	// Run the PoW algorithm to sign the block
	nonce, hash := block.Sign()

	// Set the Hash of the Block
	block.Hash = hash[:]
	// Set the Nonce of the Block
	block.Nonce = nonce

	// Return the block
	return &block
}

// A method of Block that composes and returns the block
// data as slice of bytes for a given nonce number.
//
// Considers the block data, the hash of the previous block,
// the block work difficulty and the given nonce number.
func (block *Block) Compose(nonce int) []byte {
	// Combine the block data, the previous block hash, the
	// given nonce number and the block's work difficulty
	data := bytes.Join(
		[][]byte{
			block.PrevHash,
			block.Data,
			Hexify(int64(nonce)),
			Hexify(int64(block.Difficulty)),
		},
		[]byte{},
	)

	// Return the composed data
	return data
}

// A method of Block that generates the max value of
// the hash to sign the block. Returns a big integer
func (block *Block) GenerateProofTarget() *big.Int {
	// Generate new big integer with value 1
	targethash := big.NewInt(1)
	// Left Shift the big integer by the difference between the max hash
	// size and the block's work difficulty. target = 2^(256-difficulty)
	targethash.Lsh(targethash, uint(256-block.Difficulty))

	// Return the hash target
	return targethash
}

// A method of Block that runs the Proof of Work algorithm
// to generate the hash of the block and to sign it.
// Returns the nonce number that signed the block and the hash of the block
func (block *Block) Sign() (int, []byte) {
	// Declare a big Int version of the hash
	var inthash big.Int
	// Declare an array of bytes for the hash
	var hash [32]byte
	// Initialize the nonce
	nonce := 0

	// Iterate until nonce reaches the maximum int64 value
	for nonce < math.MaxInt64 {
		// Compose the block data
		data := block.Compose(nonce)
		// Generate the hash for the composed data
		hash = sha256.Sum256(data)

		// Print the hash (with a line reset)
		fmt.Printf("\r%x", hash)

		// Set the inthash with the hash slice
		inthash.SetBytes(hash[:])
		// Check if the inthash is lesser than the proof target
		if inthash.Cmp(block.GenerateProofTarget()) == -1 {
			// Block Signed! Break from the loop
			break
		} else {
			// Increment the nonce and retry
			nonce++
		}
	}
	// Print an empty line for spacing
	fmt.Println()
	// Return the block nonce and hash
	return nonce, hash[:]
}

// A method of Block that validates the block signature (hash)
func (block *Block) Validate() bool {
	// Declare a big Int version of the hash
	var inthash big.Int

	// Compose the block data
	data := block.Compose(block.Nonce)
	// Generate the hash of the composed data
	hash := sha256.Sum256(data)
	// Set the inthash with the hash slice
	inthash.SetBytes(hash[:])

	// Check if the inthash is lesser than the proof target
	// If the hash of the block data with the given nonce is
	// less than the proof target, the block signature is valid.
	return inthash.Cmp(block.GenerateProofTarget()) == -1
}

// A function that generates and returns the
// Hex/Bytes representation of an int64
func Hexify(number int64) []byte {
	// Construct a new binary buffer
	buff := new(bytes.Buffer)
	// Write the number as a binary into the buffer in Big Endian order
	err := binary.Write(buff, binary.BigEndian, number)
	// Handle any potential error
	Handle(err)

	// Return the bytes from the binary buffer
	return buff.Bytes()
}

// A function to handle errors.
func Handle(err error) {
	// Check if error is non null
	if err != nil {
		// Log the error and throw a panic
		log.Panic(err)
	}
}

// A function to serialize a Block into gob of bytes
func BlockSerialize(block *Block) []byte {
	// Create a bytes buffer
	var gobdata bytes.Buffer
	// Create a new Gob encoder with the bytes buffer
	encoder := gob.NewEncoder(&gobdata)
	// Encode the Block into a gob
	err := encoder.Encode(block)
	// Handle any potential errors
	Handle(err)

	// Return the gob bytes
	return gobdata.Bytes()
}

// A function to deserialize a gob of bytes into a Block
func BlockDeserialize(gobdata []byte) *Block {
	// Declare a Block variable
	var block Block
	// Create a new Gob decoder by reading the gob bytes
	decoder := gob.NewDecoder(bytes.NewReader(gobdata))
	// Decode the gob into a Block
	err := decoder.Decode(&block)
	// Handle any potential errors
	Handle(err)

	// Return the pointer to the Block
	return &block
}
