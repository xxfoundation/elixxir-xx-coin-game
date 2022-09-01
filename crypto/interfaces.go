////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package crypto

type Crypto interface {
	RandomGeneration(message string, salt []byte) []byte
	Weight(digest []byte) uint
}

type Rng struct{}

func NewRng() *Rng {
	return &Rng{}
}
