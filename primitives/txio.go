/*
This module contains the definition and implementation of the
smaller component structures such as TXI, TXO and Address.
*/
package primitives

import (
	"bytes"

	"github.com/manishmeganathan/blockweave/utils"
)

// A structure that represents the inputs in a transaction
// which are really just references to previous outputs
type TXI struct {
	// Represents the transaction ID of which the reference output is a part
	ID Hash

	// Represents the index of reference output in the transaction
	OutIndex int

	// Represents the signature of the transaction
	Signature []byte

	// Represents the public key of the sending address
	PublicKey PublicKey
}

// A structure that represents the outputs in a transaction
type TXO struct {
	// Represents the token value of a given transaction output
	Value int

	// Represents the hash of the public key of the recieving address
	PublicKeyHash Hash
}

// TODO: needs rework when wallet tools are refactored
// A method of TxOutput that locks the output for a given address
func (txo *TXO) Lock(address Address) {
	// Decode the address from base58
	publickeyhash := utils.Base58Decode(address.Bytes)
	// Isolate public key hash from the checksum and version
	publickeyhash = publickeyhash[1 : len(publickeyhash)-4]
	// Assign the output key hash to public hash of the given address
	txo.PublicKeyHash = publickeyhash
}

// A method of TxInput that checks if the input public key is valid for a given public key hash
func (txi *TXI) CheckKey(publickeyhash Hash) bool {
	// Generate the hash of the input public key
	lockhash := utils.Hash160(txi.PublicKey)
	// Check if the locking hash is equal to the given hash
	return bytes.Equal(lockhash, publickeyhash)
}

// A method of TxOutput that checks if the ouput key hash is valid for a given locking hash
func (txo *TXO) CheckLock(lockhash []byte) bool {
	// Check if locking hash is equal to output's key hash
	return bytes.Equal(txo.PublicKeyHash, lockhash)
}
