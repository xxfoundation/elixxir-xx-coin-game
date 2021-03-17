package crypto

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/crypto/csprng"
	"math"
)

var resultLookup []uint
var Salt []byte

func Init() {
	resultLookup = make([]uint, 1000)

	base := 5
	rng := NewRng()
	entree := 0
	for i := 0; i < 6; i++ {
		coinValue := math.Pow(2, float64(i+base))
		lastEntree := entree + (1000 / (int(math.Pow(2, float64(i+1)))))
		for ; entree <= lastEntree; entree++ {
			resultLookup[entree] = uint(coinValue)
		}
	}
	for ; entree < 1000; entree++ {
		resultLookup[entree] = 32
	}

	Salt = make([]byte, 32)
	_, err := csprng.NewSystemRNG().Read(Salt)
	if err != nil {
		jww.FATAL.Panicf(err.Error())
	}

	jww.INFO.Printf("Pre-committed output with message \"test\": %v",
		rng.Weight(rng.RandomGeneration("test", Salt)))

}
