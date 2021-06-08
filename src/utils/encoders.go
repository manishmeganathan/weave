/*
This module contains the utility function that are used to
encode binary objects to different formats and to convert
structs into their binary gob formats
*/
package utils

import (
	"github.com/mr-tron/base58"
	log "github.com/sirupsen/logrus"
)

// A function to encode a bytes payload into a Base58 bytes payload
func Base58Encode(payload []byte) []byte {
	// Encode the payload to base58
	encode := base58.Encode(payload)
	// Cast the encoded string to a slice of bytes and return it
	return []byte(encode)
}

// A function to decode a Base58 bytes payload into a regular bytes payload
func Base58Decode(encodeddata []byte) []byte {
	// Cast the base 58 encoded data into a string and decode it
	decode, err := base58.Decode(string(encodeddata[:]))
	// Handle any potential errors
	log.Fatal("base58 decode failed!", err)
	// Return the decoded bytes
	return decode
}
