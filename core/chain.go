package core

import (
	"github.com/manishmeganathan/blockweave/merkle"
	"github.com/manishmeganathan/blockweave/persistence"
	"github.com/manishmeganathan/blockweave/utils"
	"github.com/manishmeganathan/blockweave/wallet"
	"github.com/sirupsen/logrus"
)

// A structure that represents the blockchain
type BlockChain struct {
	// Represents the database bucket for the chain state
	State *persistence.DatabaseBucket

	// Represents the database bucket for the chain blocks
	Blocks *persistence.DatabaseBucket

	// Represents the hash of the latest block
	ChainHead utils.Hash

	// Represents the number of block on the chain (last block height+1)
	ChainHeight int
}

// A constructor function that creates a new BlockChain object.
// Checks if the chain database is already configured and initializes
// the object based on that, otherwise configures a new chain database.
func NewBlockChain() *BlockChain {
	// Create a null blockchain
	blockchain := BlockChain{}

	// Check if a blockchain db already exists
	if persistence.CheckDatabase() {
		// Setup existing blockchain db
		blockchain.setup_oldchain()
	} else {
		// Setup new blockchain db
		blockchain.setup_newchain()
	}

	// Return the blockchain
	return &blockchain
}

// A method of BlockChain that configures an existing chain database.
func (chain *BlockChain) setup_oldchain() {
	chain.OpenBuckets()

	// Get the chain head from the state bucket
	chainhead, err := chain.State.GetKey(utils.ChainHeadKey)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to get chain head from state.")
	}

	// Get the chain height from the state bucket
	chainheight, err := chain.State.GetKey(utils.ChainHeightKey)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to get chain height from state.")
	}

	// Assign the current chain head
	chain.ChainHead = chainhead
	// Assign the current chain height
	chain.ChainHeight = utils.HexDecode(chainheight)
}

// A method of BlockChain that configures a new chain database.
func (chain *BlockChain) setup_newchain() {
	// Open the database clients for all buckets
	chain.OpenBuckets()

	// Get the Config data
	config := utils.ReadConfigFile()
	// Get the wallet address of the miner coinbase from the config
	address, err := wallet.NewAddress(config.JBOK.Default)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to get address for coinbase.")
	}

	// Generate a coinbase transaction for the genesis block
	coinbase := NewCoinbaseTransaction(*address)

	// Create a merkle builder
	merkletree := merkle.NewMerkleTree()
	// Build the merkle tree for the coinbase transaction
	merkletree.BuildFull([]utils.GobEncodable{coinbase})

	// Generate a Genesis Block for the chain with a coinbase transaction
	genesisblock := NewBlock(merkletree, []byte{}, 0, *address)
	// Log the minting of the genesis block
	logrus.WithFields(logrus.Fields{"address": address.String, "reward": coinbase.Outputs[0].Value}).Info("genesis block has been minted!")

	// Set the genesis block to the blocks bukcet
	if err = chain.Blocks.SetKey(genesisblock.BlockHash, genesisblock.Serialize()); err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to add genesis block to blocks.")
	}

	// Set the genesis block hash as the chain head in the state bucket
	if err = chain.State.SetKey(utils.ChainHeadKey, genesisblock.BlockHash); err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to add chain head to state.")
	}

	// Set the chain height as 1 in the state bucket
	if err = chain.State.SetKey(utils.ChainHeightKey, utils.HexEncode(1)); err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to add chain height to state.")
	}

	// Assign the current chain head
	chain.ChainHead = genesisblock.BlockHash
	// Assign the current chain height
	chain.ChainHeight = 1
}

// A method of BlockChain that adds a new Block to the chain and returns it
func (chain *BlockChain) AddBlock(blocktxns []*Transaction, addr wallet.Address) *Block {
	// Create a merkle builder
	merkletree := merkle.NewMerkleTree()
	// Start the merkle builder
	go merkletree.Build()
	// Iterate over the block transactions
	for _, txn := range blocktxns {
		// Send each transaction to the merkle build queue
		merkletree.BuildQueue <- txn
	}
	// Close the build queue
	close(merkletree.BuildQueue)

	// Generate a new Block
	block := NewBlock(merkletree, chain.ChainHead, chain.ChainHeight, addr)

	// Assign the hash of the block as the chain head
	chain.ChainHead = block.BlockHash
	// Increment the chain height
	chain.ChainHeight++

	// Set the block to the blocks bucket
	if err := chain.Blocks.SetKey(block.BlockHash, block.Serialize()); err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to add block head to blocks.")
	}

	// Set the block hash as the chain head in the state bucket
	if err := chain.State.SetKey(utils.ChainHeadKey, chain.ChainHead); err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to update chain head state.")
	}

	// Set the chain height as the current chain height in the state bucket
	if err := chain.State.SetKey(utils.ChainHeightKey, utils.HexEncode(chain.ChainHeight)); err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to update chain height state.")
	}

	// Return the block
	return block
}

// A method of BlockChain that opens the client for all database buckets.
// The method also sets up the exit handler to automatically close the clients.
func (chain *BlockChain) OpenBuckets() {
	// Set up the database client for the state bucket
	chain.State = persistence.NewDatabaseBucket(persistence.STATE)
	// Set up the database client for the blocks bucket
	chain.Blocks = persistence.NewDatabaseBucket(persistence.BLOCKS)

	logrus.RegisterExitHandler(func() {
		// Close the database client
		chain.CloseBuckets()
	})
}

// A method of BlockChain that closes the client for all database buckets.
func (chain *BlockChain) CloseBuckets() {
	// Close the database client for the state bucket
	chain.State.Close()
	// Close the database client for the blocks bucket
	chain.Blocks.Close()
}
