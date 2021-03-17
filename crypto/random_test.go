////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package crypto

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/crypto/csprng"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	Init()
	os.Exit(m.Run())
}

func TestRandomGeneration(t *testing.T) {
	rng := NewRng()

	salt := make([]byte, 32)
	_, err := csprng.NewSystemRNG().Read(salt)
	if err != nil {
		jww.FATAL.Panicf(err.Error())
	}

	jww.INFO.Printf("Pre-committed output with message \"test\": %v",
		rng.Weight(rng.RandomGeneration("test", salt)))

	digest := rng.RandomGeneration("test", salt)

	if len(digest) != 32 {
		t.Errorf("RandomGeneration did not output a digest against the spec."+
			"\n\tExpected length: %v"+
			"\n\tReceived Lenth: %v", 32, len(digest))
	}
	winnings := rng.Weight(digest)

	t.Logf("resultLookup: %v", resultLookup)
	t.Logf("Salt: %v", salt)
	if winnings < 32 || winnings > 1024 {
		t.Errorf("Winnings out of bound of 32 to 1024."+
			"Winning value: %v", winnings)
	}
}
