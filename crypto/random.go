////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package crypto

import (
	"crypto/sha256"
)

func (rng *Rng) RandomGeneration(message string, salt []byte) []byte {
	// Generate hash
	h := sha256.New()

	h.Write([]byte(message))
	h.Write(salt)

	// Return a 256 bit digest
	return h.Sum(nil)[:32]
}
