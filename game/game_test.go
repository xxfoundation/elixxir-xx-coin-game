package game

import "testing"

func TestNew(t *testing.T) {
	g := New(map[string]uint{"t1": 0, "t2": 0, "t3": 3}, []byte("salt"), nil)
	if g.winnings == nil {
		t.Errorf("Did not initialize winnings map")
	}
}

func TestGame_Play(t *testing.T) {
	g := New(map[string]uint{"t1": 0, "t2": 0, "t3": 3}, []byte("salt"), nil)
	ok, _ := g.Play("addr", "i'm a message")
	if ok {
		t.Error("Should not have been able to play with unknown address")
	}
	ok, _ = g.Play("t1", "i'm not a message")
	if !ok {
		t.Error("Failed to play game with address t1")
	}
}
