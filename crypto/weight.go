////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
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
