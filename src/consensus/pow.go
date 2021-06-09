package consensus

import (
	"math/big"

	"github.com/manishmeganathan/blockweave/src/primitives"
)

type POW struct {
	Target big.Int
	Nonce  int
}

func NewPOW() *POW {
	return &POW{}
}

func (pow *POW) GenerateTarget() *big.Int {
	return nil
}

func (pow *POW) Mint(block *primitives.Block) error {
	return nil
}

func (pow *POW) Validate(block *primitives.Block) bool {
	return true
}
