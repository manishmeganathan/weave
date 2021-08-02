package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/manishmeganathan/weave/utils"
	"github.com/sirupsen/logrus"
)

// A struct that represents a collection of wallets.
// JBOK -> Just a Bunch of Keys.
type JBOK struct {
	// Represents the mapping of wallet address to wallet objects
	Wallets map[string]*Wallet
}

// A function that returns the path to the jbok data file.
// The jbok file is at %HOME%/blockweave/jbok.data
func getjbokfilepath() string {
	// Retrieve the path to the config dir.
	configdir := utils.ConfigDirectory()
	// Return the file location
	return filepath.Join(configdir, "jbok.data")
}

// A constructor function that loads the JBOK data from a file and returns a JBOK object.
// The data is read from the jbok file at %HOME%/blockweave/jbok.data
// If the file does not exist, an empty JBOK is created.
func NewJBOK() *JBOK {
	// Create a new JBOK object
	jbok := JBOK{}
	// Initialize the Wallets field
	jbok.Wallets = make(map[string]*Wallet)
	// Check if the jbok file exists
	if _, err := os.Stat(getjbokfilepath()); errors.Is(err, os.ErrNotExist) {
		// If the file does not exist. Save the empty JBOK object into a file
		jbok.Save()
	}

	// Read the walletstore file into a slice of bytes
	filecontents, err := ioutil.ReadFile(getjbokfilepath())
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

// A method of JBOK that saves the current state of the JBOK to the jbok file
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

	// Write the bytes from the buffer to the jbok file
	err = ioutil.WriteFile(getjbokfilepath(), buff.Bytes(), 0644)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to write jbok data to file.")
	}
}

// A function that purges the JBOK data file by deleting it.
func PurgeJBOK() {
	// Get the path to the jbok file
	file := getjbokfilepath()
	// Remove the JBOK file
	err := os.Remove(file)
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to purge jbok data file.")
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
func (jbok *JBOK) AddWallet(wallet *Wallet) Address {
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
func (jbok *JBOK) CreateWallet() Address {
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

// A method of JBOK that checks if a given address exists in the JBOK
func (jbok *JBOK) CheckWallet(address string) bool {
	// Check if the address exists in the JBOK
	_, ok := jbok.Wallets[address]
	// Return the result of the check
	return ok
}

// A method of JBOK that removes a wallet from the JBOK.
func (jbok *JBOK) RemoveWallet(address string) {
	// Remove the wallet from the JBOK
	delete(jbok.Wallets, address)
	// Save the JBOK to the file
	jbok.Save()
}
