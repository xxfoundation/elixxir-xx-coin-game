////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package game

import (
	"github.com/pkg/errors"
	"gitlab.com/elixxir/xx-coin-game/crypto"
	"sync"
)

type Game struct {
	crypto   crypto.Crypto
	salt     []byte
	winnings map[string]*Play
}

type Play struct {
	sync.Mutex
	winnings uint
}

func New(current map[string]uint64, salt []byte, crypto crypto.Crypto) *Game {
	// TODO: load winnings from file in io, add implementations for RNG &weight, tests for this package
	g := &Game{
		winnings: map[string]*Play{},
		salt:     salt,
		crypto:   crypto,
	}
	for k, v := range current {
		g.winnings[k] = &Play{
			Mutex:    sync.Mutex{},
			winnings: uint(v),
		}
	}
	return g
}

func (g *Game) Play(address, message string) (bool, uint, error) {
	p, ok := g.winnings[address]
	if !ok {
		return false, 0, errors.Errorf("Could not find eth address %s, does it have xx coins?", address)
	}
	new, value := p.play(message, g.crypto, g.salt)
	return new, value, nil
}

func (p *Play) play(message string, crypto crypto.Crypto, salt []byte) (bool, uint) {
	p.Lock()
	defer p.Unlock()

	if p.winnings == 0 {
		digest := crypto.RandomGeneration(message, salt)
		weight := crypto.Weight(digest)
		p.winnings = weight
		return true, p.winnings
	}
	return false, p.winnings
}
