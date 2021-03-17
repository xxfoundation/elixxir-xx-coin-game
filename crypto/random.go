////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package crypto

import (
	"crypto/sha256"
)

func RandomGeneration(message string, salt []byte) []byte {
	// Generate hash
	h := sha256.New()

	h.Write([]byte(message))
	h.Write(salt)

	// Return a 256 bit digest
	return h.Sum(nil)[:32]
}
