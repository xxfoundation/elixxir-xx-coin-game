package game

import "testing"

func TestNew(t *testing.T) {
	g := New(map[string]uint{"t1": 0, "t2": 0, "t3": 3})
	if g.winnings == nil {
		t.Errorf("Did not initialize winnings map")
	}
}

func TestGame_Play(t *testing.T) {
	g := New(map[string]uint{"t1": 0, "t2": 0, "t3": 3})
	g.Play()
}
