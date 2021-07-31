package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/manishmeganathan/blockweave/primitives"
	"github.com/manishmeganathan/blockweave/utils"
)

// A structure that represents a store of wallets
type WalletStore struct {
	Wallets map[string]*Wallet // Represents the map of wallet address to wallet pointers
}

// A constructor function to generate and return a new
// WalletStore by reading from the walletstore data file.
func NewWalletStore() *WalletStore {
	// Create a new WalletStore object
	walletstore := WalletStore{}
	// Initialize the Wallets field
	walletstore.Wallets = make(map[string]*Wallet)

	// Read the walletstore file
	err := walletstore.Load()
	// Handle any potential errors
	utils.HandleErrorLog(err, "wallet store creation failed")

	// Return the wallet store
	return &walletstore
}

// A method of WalletStore that adds a new wallet to the
// store and returns the address of the new Wallet.
func (ws *WalletStore) AddWallet() primitives.Address {
	// Construct a new Wallet
	wallet := NewWallet()
	// Generate the address of the wallet
	address := wallet.GenerateAddress(byte(0x00))

	// Assign the wallet to the store with its address as the key
	ws.Wallets[address.String] = wallet
	// Return the address of the wallet
	return *address
}

// A method of WalletStore that fetches a wallet from
// the wallet store given a wallet address.
func (ws *WalletStore) FetchWallet(address string) *Wallet {
	return ws.Wallets[address]
}

// A method of WalletStore that returns the
// address of wallets in the WalletStore.
func (ws *WalletStore) GetAddresses() []string {
	// Declare a slice of strings
	var addresses []string

	// Iterate over the wallets in the store
	for address := range ws.Wallets {
		// Add the address to the slice
		addresses = append(addresses, address)
	}

	// Return the slice of addresses
	return addresses
}

// A method of WalletStore that saves the current state
// of the WalletStore to the walletstore file.
func (ws *WalletStore) Save() {
	// Declare a bytes buffer
	var buff bytes.Buffer
	// Register the gob library to use the elliptic sepc256r1 curve
	gob.Register(elliptic.P256())

	// Create a new gob encoder with the bytes buffer
	encoder := gob.NewEncoder(&buff)
	// Encode the wallet store to the buffer
	err := encoder.Encode(ws)
	// Handle any potential errors
	utils.HandleErrorLog(err, "wallet store encode failed")

	// Write the bytes from the buffer to the walletstore file
	err = ioutil.WriteFile(utils.WalletDB, buff.Bytes(), 0644)
	// Handle any potential errors
	utils.HandleErrorLog(err, "wallet store save failed")
}

// A method of WalletStore that loads the walletstore file
// into the current state of the WalletStore
func (ws *WalletStore) Load() error {
	// Check if the walletstore file exists
	if _, err := os.Stat(utils.WalletDB); errors.Is(err, os.ErrNotExist) {
		// If the file does not exist. Save the current state of the walletstore into a file
		ws.Save()
	}

	// Declare a WalletStore object
	var walletstore WalletStore

	// Read the walletstore file into a slice of bytes
	filecontents, err := ioutil.ReadFile(utils.WalletDB)
	if err != nil {
		return fmt.Errorf("walletstore file could not be read")
	}

	// Register the gob library to use the elliptic sepc256r1 curve
	gob.Register(elliptic.P256())

	// Create a new gob decoder for the contents of the walletstore file
	decoder := gob.NewDecoder(bytes.NewReader(filecontents))
	// Decode the gob data into the WalletStore object
	err = decoder.Decode(&walletstore)
	if err != nil {
		return fmt.Errorf("walletstore file could not be decoded")
	}

	// Copy the Wallets field from the new WalletStore object
	ws.Wallets = walletstore.Wallets
	// Return with no errors
	return nil
}
