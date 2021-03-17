package crypto

import (
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/crypto/csprng"
	"math"
)

var resultLookup []uint
var salt []byte

func InitCrypto() {
	resultLookup = make([]uint, 1000)

	base := 5

	entree := 0
	for i := 0; i < 6; i++ {
		coinValue := math.Pow(2, float64(i+base))
		lastEntree := entree + int(1/(math.Pow(2, float64(1+i))))
		for ; entree <= lastEntree; entree++ {
			resultLookup[entree] = uint(coinValue)
		}
	}
	for ; entree < 1000; entree++ {
		resultLookup[entree] = 32
	}

	salt = make([]byte, 32)
	_, err := csprng.NewSystemRNG().Read(salt)
	if err != nil {
		jww.FATAL.Panicf(err.Error())
	}

	jww.INFO.Printf("Pre-committed output with message \"test\": %v", Weight(RandomGeneration("test", salt)))

}
