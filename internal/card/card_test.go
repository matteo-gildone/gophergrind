package card_test

import (
	"fmt"
	"gophergrind/internal/card"
	"testing"
)

// fakeRNG is a scripted RNG: it replays values in order and records the n
// passed to each Intn call so tests can assert consumption order. It
// ignores n when choosing the return value, and panics if a card reads
// past the script so over-consumption fails loudly.
type fakeRNG struct {
	values []int
	i      int
}

func (f *fakeRNG) Intn(n int) int {
	if f.i >= len(f.values) {
		panic(fmt.Sprintf("fakeRNG exhausted after %d calls (last n=%d)", f.i, n))
	}
	v := f.values[f.i]
	f.i++
	return v
}

func TestOffer(t *testing.T) {
	t.Parallel()
	rng := &fakeRNG{values: []int{1, 1, 1}}
	offer := card.Offer(rng)

	if len(offer) != 3 {
		t.Fatalf("Offer returned %d cards, want 3", len(offer))
	}

	if offer[0].ID != "damage" {
		t.Errorf("want: damage, got: %s ", offer[0].ID)
	}

	if offer[1].ID != "hp" {
		t.Errorf("want: hp, got: %s ", offer[1].ID)
	}

	if offer[2].ID != "defence" {
		t.Errorf("want: defence, got: %s ", offer[2].ID)
	}
}
