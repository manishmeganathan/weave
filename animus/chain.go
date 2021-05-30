package animus

type BlockChain struct {
	Blocks []*Block
}

// A constructor function that generates a new BlockChain
// and initializes the chain with a Genesis Block.
func NewBlockChain() *BlockChain {
	// Generate a Genesis Block for the chain
	genesisblock := NewBlock("Genesis", []byte{})
	// Construct a new BlockChain and add the Genesis Block to it
	blockchain := BlockChain{Blocks: []*Block{genesisblock}}
	// Return the blockchain
	return &blockchain
}

// A method of BlockChain that adds a new block to the chain
func (chain *BlockChain) AddBlock(blockdata string) {
	// Retrieve the previous block on the chain
	prevblock := chain.Blocks[len(chain.Blocks)-1]
	// Generate a new block from the block data and the has of the previous block
	newblock := NewBlock(blockdata, prevblock.Hash)
	// Add the new blocj to the chain
	chain.Blocks = append(chain.Blocks, newblock)
}
