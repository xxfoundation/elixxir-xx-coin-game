////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package crypto

import (
	"math/big"
)

// Weights the random value to determine the winnings
func (rng *Rng) Weight(digest []byte) uint {
	data := big.NewInt(1)
	data.SetBytes(digest)
	mod := big.NewInt(1000)
	data.Mod(data, mod)

	return resultLookup[data.Uint64()]
}
