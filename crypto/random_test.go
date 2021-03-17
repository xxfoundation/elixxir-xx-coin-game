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
	Init()
	os.Exit(m.Run())
}

func TestRandomGeneration(t *testing.T) {
	rng := NewRng()
	digest := rng.RandomGeneration("test", Salt)

	if len(digest) != 32 {
		t.Errorf("RandomGeneration did not output a digest against the spec."+
			"\n\tExpected length: %v"+
			"\n\tReceived Lenth: %v", 32, len(digest))
	}
	winnings := rng.Weight(digest)

	t.Logf("resultLookup: %v", resultLookup)
	t.Logf("Salt: %v", Salt)
	if winnings < 32 || winnings > 1024 {
		t.Errorf("Winnings out of bound of 32 to 1024."+
			"Winning value: %v", winnings)
	}
}
