package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	animus "github.com/manishmeganathan/animus-blockchain/animus"
)

// A function to generate a random 20 character case-sensitive alpha-numeric string.
func generaterandomtext() string {
	var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, 20)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}

	return string(b)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	chain := animus.NewBlockChain()

	for i := 0; i < 5; i++ {
		chain.AddBlock(generaterandomtext())
	}

	fmt.Println()
	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(block.Validate()))
	}
}
