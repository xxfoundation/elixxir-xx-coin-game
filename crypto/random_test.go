////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package crypto

import (
	"testing"
)

func TestRandomGeneration(t *testing.T) {
	digest := RandomGeneration("test")
	t.Logf("Pre-committed digest: %v", digest)

	winnings := Weight(digest)
	t.Logf("Winning: %v", winnings)
}
