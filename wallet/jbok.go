package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"os"

	"github.com/manishmeganathan/blockweave/primitives"
	"github.com/manishmeganathan/blockweave/utils"
	"github.com/sirupsen/logrus"
)

// A struct that represents a collection of wallets.
// JBOK -> Just a Bunch of Keys.
type JBOK struct {
	// Represents the mapping of wallet address to wallet objects
	Wallets map[string]*Wallet
}

// A constructor function that loads the JBOK data from a file and returns a JBOK object.
// The data is read from a file specified in the 'jbokfile' key in the config file.
// If the file does not exist, an empty JBOK is created.
func NewJBOK() *JBOK {
	// Get the Config data
	config := utils.ReadConfigFile()

	// Create a new JBOK object
	jbok := JBOK{}
	// Initialize the Wallets field
	jbok.Wallets = make(map[string]*Wallet)
	// Check if the jbok file exists
	if _, err := os.Stat(config.JBOKFile); errors.Is(err, os.ErrNotExist) {
		// If the file does not exist. Save the empty JBOK object into a file
		jbok.Save()
	}

	// Read the walletstore file into a slice of bytes
	filecontents, err := ioutil.ReadFile(config.JBOKFile)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to read jbok data from file.")
	}

	// Register the gob library to use the elliptic sepc256r1 curve
	gob.Register(elliptic.P256())

	// Create a new gob decoder for the contents of the walletstore file
	decoder := gob.NewDecoder(bytes.NewReader(filecontents))
	// Decode the gob data into the WalletStore object
	err = decoder.Decode(&jbok)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to decode jbok data.")
	}

	// Return the JBOK object
	return &jbok
}

// A method of JBOK that saves the current state of the
// JBOK to the jbok file specified in the config file.
func (jbok *JBOK) Save() {
	// Declare a bytes buffer
	var buff bytes.Buffer
	// Register the gob library to use the elliptic
	// sepc256r1 curve while encoding the JBOK data
	gob.Register(elliptic.P256())

	// Create a new gob encoder with the bytes buffer
	encoder := gob.NewEncoder(&buff)
	// Encode the jbok to the buffer
	err := encoder.Encode(jbok)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to encode jbok data.")
	}

	// Get the Config data
	config := utils.ReadConfigFile()
	// Write the bytes from the buffer to the jbok file
	err = ioutil.WriteFile(config.JBOKFile, buff.Bytes(), 0644)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to write jbok data to file.")
	}
}

// A method of JBOK that retrieves the addresses of all wallets in the JBOK.
func (jbok *JBOK) GetAddresses() []string {
	// Declare a slice of strings
	var addrs []string

	// Iterate over the wallets in the JBOK
	for address := range jbok.Wallets {
		// Add the address to the slice
		addrs = append(addrs, address)
	}

	// Return the slice of addresses
	return addrs
}

// A method of JBOK that adds a given wallet to the JBOK and returns its address.
func (jbok *JBOK) AddWallet(wallet *Wallet) primitives.Address {
	// Generate the address of the wallet
	address := wallet.GenerateAddress(byte(0x00))
	// Assign the wallet to the JBOK with its address as the key
	jbok.Wallets[address.String] = wallet
	// Save the JBOK to the file
	jbok.Save()

	// Return the address of the wallet
	return *address
}

// A method of JBOK that creates and adds a new wallet to the JBOK.
// The address of the newly created wallet is returned.
func (jbok *JBOK) CreateWallet() primitives.Address {
	// Construct a new Wallet
	wallet := NewWallet()
	// Add the wallet to the JBOK
	address := jbok.AddWallet(wallet)

	// Return the address of the wallet
	return address
}

// A method of JBOK that retrieves a wallet from the JBOK for a given address string.
func (jbok *JBOK) FetchWallet(address string) *Wallet {
	return jbok.Wallets[address]
}
