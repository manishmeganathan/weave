/*
This module contains the definition and implementation of
the Transaction input and output struct and their methods
*/
package primitives

import (
	"bytes"

	"github.com/manishmeganathan/blockweave/src/utils"
)

// A structure that represents the inputs in a transaction
// which are really just references to previous outputs
type TXI struct {
	ID        []byte // Represents the reference transaction of which the output is a part
	OutIndex  int    // Represents the index of output in the reference transaction
	Signature []byte // Represents the signature of the transaction
	PublicKey []byte // Represents the public key of the sending address
}

// A structure that represents the outputs in a transaction
type TXO struct {
	Value         int    // Represents the token value of a given transaction output
	PublicKeyHash []byte // Represents the hash of the public key of the recieving address
}

// A type alias for a slice of transaction inputs
type TXIList []TXI

// A type alias for a slice of transaction outputs
type TXOList []TXO

// A constructor function that generates and returns a new
// transaction output given a token value and address
func NewTxOutput(value int, address string) *TXO {
	txo := TXO{Value: value, PublicKeyHash: nil}
	txo.Lock([]byte(address))

	return &txo
}

// A method of TxOutput that locks the output for a given address
func (txo *TXO) Lock(address []byte) {
	// Decode the address from base58
	publickeyhash := utils.Base58Decode(address)
	// Isolate public key hash from the checksum and version
	publickeyhash = publickeyhash[1 : len(publickeyhash)-4]
	// Assign the output key hash to public hash of the given address
	txo.PublicKeyHash = publickeyhash
}

// A method of TxOutput that checks if the ouput key hash is valid for a given locking hash
func (txo *TXO) CheckLock(lockhash []byte) bool {
	// Check if locking hash is equal to output's key hash
	return bytes.Equal(txo.PublicKeyHash, lockhash)
}
