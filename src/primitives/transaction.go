/*
This module contains the definition and implementation
of the Transaction struct and its methods
*/
package primitives

// A structure that represents a transaction on the Animus Blockchain
type Transaction struct {
	ID      []byte // Represents the hash of the transaction
	Inputs  []TXI  // Represents the inputs of the transaction
	Outputs []TXO  // Represents the outputs of the transaction
}
