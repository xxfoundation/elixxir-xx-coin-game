////////////////////////////////////////////////////////////////////////////////
// Copyright © 2021 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package game

import (
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

func New(current map[string]uint, salt []byte, crypto crypto.Crypto) *Game {
	// TODO: load winnings from file in io, add implementations for RNG &weight, tests for this package
	g := &Game{
		winnings: map[string]*Play{},
		salt:     salt,
		crypto:   crypto,
	}
	for k, v := range current {
		g.winnings[k] = &Play{
			Mutex:    sync.Mutex{},
			winnings: v,
		}
	}
	return g
}

func (g *Game) Play(address, message string) (bool, uint) {
	p, ok := g.winnings[address]
	if !ok {
		return false, 0
	}
	return p.play(message, g.crypto)
}

func (p *Play) play(message string, crypto crypto.Crypto) (bool, uint) {
	p.Lock()
	defer p.Unlock()

	if p.winnings == 0 {
		digest := crypto.RandomGeneration(message)
		weight := crypto.Weight(digest)
		p.winnings = weight
		return true, p.winnings
	}
	return false, p.winnings
}
