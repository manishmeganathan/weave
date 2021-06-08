/*
This module contains the utility function that are used
for various hashing functionality
*/
package utils

import "crypto/sha256"

// A function that returns a double SHA256
// hash output of a slice of bytes payload
//
// HASH = SHA256(SHA256(PAYLOAD))
func Hash256(payload []byte) []byte {
	var hash [32]byte

	// Generate the hash of the payload
	hash = sha256.Sum256(payload)
	// Generate the hash of the hash
	hash = sha256.Sum256(hash[:])

	// Return the hash as slice
	return hash[:]
}
