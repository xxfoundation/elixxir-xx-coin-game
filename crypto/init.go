////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package crypto

import (
	"math"
)

var resultLookup []uint

func init() {
	resultLookup = make([]uint, 1000)

	base := 5
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
}
