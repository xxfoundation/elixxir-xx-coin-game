////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package crypto

import "crypto"

func RandomGeneration(message string) []byte {
	// Generate hash
	sha := crypto.SHA256
	h := sha.New()

	h.Write([]byte(message))
	h.Write(Salt)

	// Return a 256 bit digest
	return h.Sum(nil)[:32]
}
