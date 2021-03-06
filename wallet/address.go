package wallet

import (
	"bytes"
	"fmt"

	"github.com/manishmeganathan/weave/utils"
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
	Checksum utils.Hash

	// Represents the hash of the public key of the address
	PublicKeyHash utils.Hash
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
	// Calculate the breakpoint between the checksum and extended hash
	hashlen := len(fullhash) - 4

	// Isolate the checksum hash from the full hash (after the breakpoint)
	addr.Checksum = fullhash[hashlen:]
	// Isolate the extended hash from the full hash (before the breakpoint)
	extendedhash := fullhash[:hashlen]

	// Assign the prefix from the extended hash (first byte)
	addr.Prefix = extendedhash[0]
	// Assign the public key hash from the extended hash (remove prefix)
	addr.PublicKeyHash = extendedhash[1:]

	//Check if the address is valid
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
	newsum := utils.Hash32(extendedhash)

	// Check if the new checksum is equal to the isolated checksum
	return bytes.Equal(addr.Checksum, newsum)
}
