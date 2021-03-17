////////////////////////////////////////////////////////////////////////////////
// Copyright © 2021 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package game

import "sync"

type RandomGeneration func(message string) []byte
type Weight func(digest []byte) uint

type Game struct {
	rngFunc    RandomGeneration
	weightFunc Weight
	winnings   map[string]*Play
}

type Play struct {
	sync.Mutex
	winnings uint
}

func New() *Game {
	// TODO: load winnings from file in io, add implementations for RNG &weight, tests for this package
	return &Game{
		rngFunc:    nil,
		weightFunc: nil,
		winnings:   map[string]*Play{},
	}
}

func (g *Game) Play(address, message string) (bool, uint) {
	p, ok := g.winnings[address]
	if !ok {
		g.winnings[address] = &Play{
			Mutex:    sync.Mutex{},
			winnings: 0,
		}
		p = g.winnings[address]
	}
	return p.play(message, g.rngFunc, g.weightFunc)
}

func (p *Play) play(message string, rngFunc RandomGeneration, weightFunc Weight) (bool, uint) {
	p.Lock()
	defer p.Unlock()

	if p.winnings == 0 {
		digest := rngFunc(message)
		weight := weightFunc(digest)
		p.winnings = weight
		return true, p.winnings
	}
	return false, p.winnings
}
