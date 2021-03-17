////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package crypto

type Crypto interface {
	RandomGeneration(message string, salt []byte) []byte
	Weight(digest []byte) uint
}

type Rng struct {}

func NewRng() *Rng  {
	return &Rng{}
}