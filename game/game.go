////////////////////////////////////////////////////////////////////////////////
// Copyright © 2021 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package game

import "sync"

type randomizer interface {
	RandomGeneration(message string, salt []byte) []byte
	Weight(digest []byte) uint
}
type RandomGeneration func(message string) []byte
type Weight func(digest []byte) uint

type Game struct {
	rand     randomizer
	winnings map[string]*Play
}

type Play struct {
	sync.Mutex
	winnings uint
}

func New(current map[string]uint) *Game {
	// TODO: load winnings from file in io, add implementations for RNG &weight, tests for this package
	g := &Game{
		winnings: map[string]*Play{},
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
	return p.play(message, g.rand)
}

func (p *Play) play(message string, rng randomizer) (bool, uint) {
	p.Lock()
	defer p.Unlock()

	if p.winnings == 0 {
		digest := rng.RandomGeneration(message, []byte("salt"))
		weight := rng.Weight(digest)
		p.winnings = weight
		return true, p.winnings
	}
	return false, p.winnings
}
