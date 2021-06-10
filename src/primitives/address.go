package primitives

import (
	"bytes"
	"fmt"

	"github.com/manishmeganathan/blockweave/src/utils"
)

// A struct that represents the Address of a User/Wallet
type Address struct {
	// Represents the bytes value of the address
	Bytes []byte

	// Represents the string value of the address
	String string

	// Represents the version/weave prefix of the address
	Prefix byte

	// Represents the adddress checksum
	Checksum Hash

	// Represents the hash of the public key of the address
	PublicKeyHash Hash
}

// A constructor function that generates and returns
// a new Address object from a given address string.
func NewAddress(address string) (*Address, error) {
	// Create an empty address
	addr := Address{}

	// Assign the bytes and string value of the address
	addr.Bytes = []byte(address)
	addr.String = address

	// Decode the address from base58 to get the full hash
	fullhash := utils.Base58Decode(addr.Bytes)

	// Isolate the checksum of the address from the full hash (last 4 bytes)
	addr.Checksum = fullhash[len(fullhash)-4:]
	// Isolate the extended hash from the full hash (remove checksum)
	extendedhash := fullhash[:4]

	// Assign the prefix from the extended hash (first byte)
	addr.Prefix = extendedhash[0]
	// Assign the public key hash from the extended hash (remove prefix)
	addr.PublicKeyHash = extendedhash[1:]

	// Check if the address is valid
	if !addr.IsValid() {
		// Return an empty address and an error
		return &Address{}, fmt.Errorf("invalid address")
	}

	// Return the address
	return &addr, nil
}

// A method of Address that checks if the address
// checksum is valid and returns a boolean
func (addr *Address) IsValid() bool {
	// Create the extended hash by concatenatin the public key hash to the prefix
	extendedhash := append([]byte{addr.Prefix}, addr.PublicKeyHash...)
	// Generate a checksum32 for the extended hash
	newsum := utils.CheckSum32(extendedhash)

	// Check if the new checksum is equal to the isolated checksum
	if !bytes.Equal(addr.Checksum, newsum) {
		// Return false if unequal
		return false
	}
	return true
}
