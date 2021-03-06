package wallet

import (
	"crypto/ecdsa"

	"github.com/manishmeganathan/weave/utils"
	"github.com/sirupsen/logrus"
)

// A structure that represents a wallet to access the blockchain
type Wallet struct {
	// Represents the private key of the wallet
	PrivateKey ecdsa.PrivateKey

	// Represents the public key of the wallet
	PublicKey utils.PublicKey
}

// A constructor function that generates and returns a Wallet
func NewWallet() *Wallet {
	// Generate private-public key pair
	private, public := utils.KeyGenECDSA()
	// Assign the keys to the wallet fields
	wallet := Wallet{PrivateKey: private, PublicKey: public}

	// Return the wallet
	return &wallet
}

func (w *Wallet) GenerateAddress(prefix byte) *Address {
	// Generate the hash of the public key
	publickeyhash := utils.Hash160(w.PublicKey)
	// Generate the extended hash by appending the public key to the prefix
	extendedhash := append([]byte{prefix}, publickeyhash...)
	// Generate the checksum of the extended hash
	checksum := utils.Hash32(extendedhash)

	// Append the checksum to the end of the extended hash
	finalhash := append(extendedhash, checksum...)

	// Encode the final hash to base58
	address := utils.Base58Encode(finalhash)
	// Generate an Address from the address string
	addr, err := NewAddress(string(address))
	if err != nil {
		// Log a fatal error
		logrus.WithFields(logrus.Fields{"error": err}).Fatalln("failed to generate address for wallet.")
	}

	// Return the address
	return addr
}
