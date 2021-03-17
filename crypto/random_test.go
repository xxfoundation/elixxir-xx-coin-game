////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package crypto

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InitCrypto()
	os.Exit(m.Run())
}

func TestRandomGeneration(t *testing.T) {
	digest := RandomGeneration("test", salt)
	t.Logf("Pre-committed digest: %v", digest)

	winnings := Weight(digest)
	t.Logf("Winning: %v", winnings)

	if winnings < 32 || winnings > 1024 {
		t.Errorf("Winnings out of bound of 32 to 1024." +
			"Winning value: %v", winnings)
	}
}
