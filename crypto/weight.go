////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package crypto

import "encoding/binary"

const (
	TotalRange              = 1000000
	LowestWinningThreshold  = 515625
	SecondWinningThreshold  = 250000
	ThirdWinningThreshold   = 125000
	FourthWinningThreshold  = 625000
	FifthWinningThreshold   = 312500
	HighestWinningThreshold = 156250
)

const (
	LowestWinnings  = 32
	SecondWinnings  = 64
	ThirdWinnings   = 128
	FourthWinnings  = 256
	FifthWinnings   = 512
	HighestWinnings = 1024
)

// Weights the random value to determine the winnings
// Digest is converted to a uint and modded to fit within a TotalRange
// This value is named x, and we determine which range it is in:
// 		[LowestWinningThreshold, x, TotalRange] -> LowestWinnings
// 		[SecondWinningThreshold, x, LowestWinningThreshold] -> SecondWinnings
// 		[ThirdWinningThreshold, x, SecondWinningThreshold] -> ThirdWinnings
// 		[FourthWinningThreshold, x, ThirdWinningThreshold] -> FourthWinnings
// 		[FifthWinningThreshold, x, FourthWinningThreshold] -> FifthWinnings
// 		[HighestWinningThreshold, x, FifthWinningThreshold] -> HighestWinnings
func Weight(digest []byte) uint {
	data := binary.BigEndian.Uint64(digest)

	// Set the random value to the range of winnings
	winnings := data % TotalRange

	if between(LowestWinningThreshold, TotalRange, winnings) {
		winnings = LowestWinnings
	} else if between(SecondWinningThreshold, LowestWinningThreshold, winnings) {
		winnings = SecondWinnings
	} else if between(ThirdWinningThreshold, SecondWinningThreshold, winnings) {
		winnings = ThirdWinnings
	} else if between(FourthWinningThreshold, ThirdWinningThreshold, winnings) {
		winnings = FourthWinnings
	} else if between(FifthWinningThreshold, FourthWinningThreshold, winnings){
		winnings = FifthWinnings
	} else if between(HighestWinningThreshold, FifthWinningThreshold, winnings) {
		winnings = HighestWinnings
	}

	return uint(winnings)
}

func between(min, max, value uint64) bool {
	return value > min && value <= max
}