package crypto

import "testing"

func TestRandomGeneration(t *testing.T) {
	digest := RandomGeneration("test")
	t.Logf("Pre-committed digest: %v", digest)
}
