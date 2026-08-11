// Package combat_test is a black-box test of the fight resolver, so it
// doubles as a design check on the exported API.
//
// Contract (see spec.md, Combat):
//   - RNG is consumed as: one Intn(2) for swing order (0 = player first,
//     1 = opponent first), then per swing Intn(20) for the swinger's roll
//     followed by Intn(20) for the target's roll.
//   - A hit lands when swingerRoll+Accuracy > targetRoll+Defence (strict;
//     an exact tie is a miss, favouring the defender).
//   - Damage is flat on a hit; Defence never mitigates a landed hit.
//     Round.DefenderHP is the target's HP after the swing, floored at 0.
//   - A round is two swings; the fight ends on a KO or after 10 rounds
//     (20 swings). At the cap the higher remaining HP% wins, and an exact
//     HP% tie goes to the opponent (the player, as fixed attacker, loses).
package combat_test

import (
	"fmt"
	"slices"
	"testing"

	"gophergrind/internal/combat"
)

func TestFight_HitResolution(t *testing.T) {
	t.Parallel()

	// Neither side deals damage, so every fight runs the full 20-swing
	// cap and Rounds[0] is always the player's opening swing (order 0).
	cases := []struct {
		name    string
		acc     int // player Accuracy
		def     int // opponent Defence
		ar, dr  int // rolls for the opening swing
		wantHit bool
	}{
		{name: "equal totals is a miss (tie favours defender)", acc: 0, def: 0, ar: 10, dr: 10, wantHit: false},
		{name: "attacker one higher hits", acc: 0, def: 0, ar: 10, dr: 9, wantHit: true},
		{name: "accuracy tips a tie into a hit", acc: 1, def: 0, ar: 10, dr: 10, wantHit: true},
		{name: "defence tips a hit into a miss", acc: 0, def: 1, ar: 10, dr: 10, wantHit: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			player := combat.Combatant{Name: "player", HP: 100, MaxHP: 100, Accuracy: tc.acc}
			opp := combat.Combatant{Name: "opp", HP: 100, MaxHP: 100, Defence: tc.def}
			script := capScript(0, 20, func(i int) (int, int) {
				if i == 0 {
					return tc.ar, tc.dr
				}
				return 0, 0
			})

			log := combat.Fight(player, opp, &fakeRNG{values: script})

			if got := log.Rounds[0].Hit; got != tc.wantHit {
				t.Errorf("Rounds[0].Hit = %v, want %v", got, tc.wantHit)
			}
		})
	}
}

func TestFight_KillEndsFightAndFloorsHP(t *testing.T) {
	t.Parallel()

	player := combat.Combatant{Name: "player", HP: 100, MaxHP: 100, Accuracy: 10, Damage: 100}
	opp := combat.Combatant{Name: "opp", HP: 10, MaxHP: 10}

	// order 0 (player first); opening swing 5 vs 5 -> 5+10 > 5 hits.
	log := combat.Fight(player, opp, &fakeRNG{values: []int{0, 5, 5}})

	if got, want := len(log.Rounds), 1; got != want {
		t.Fatalf("len(Rounds) = %d, want %d (fight must stop on the KO)", got, want)
	}
	want := combat.Round{
		Attacker:   "player",
		Defender:   "opp",
		Hit:        true,
		Damage:     100, // raw stat dealt, even though it overkills 10 HP
		DefenderHP: 0,   // floored, not -90
	}
	if got := log.Rounds[0]; got != want {
		t.Errorf("Rounds[0] = %+v, want %+v", got, want)
	}
	if got, want := log.Winner, "player"; got != want {
		t.Errorf("Winner = %q, want %q", got, want)
	}
}

func TestFight_AlternatesAndOpponentCanWin(t *testing.T) {
	t.Parallel()

	player := combat.Combatant{Name: "player", HP: 10, MaxHP: 10, Accuracy: 10, Damage: 1}
	opp := combat.Combatant{Name: "opp", HP: 10, MaxHP: 10, Accuracy: 10, Damage: 100}

	// order 0: player chips for 1, then opponent lands the kill.
	log := combat.Fight(player, opp, &fakeRNG{values: []int{0, 10, 0, 10, 0}})

	if got, want := len(log.Rounds), 2; got != want {
		t.Fatalf("len(Rounds) = %d, want %d", got, want)
	}
	if got, want := log.Rounds[0].Attacker, "player"; got != want {
		t.Errorf("Rounds[0].Attacker = %q, want %q", got, want)
	}
	if got, want := log.Rounds[1].Attacker, "opp"; got != want {
		t.Errorf("Rounds[1].Attacker = %q, want %q (swings alternate)", got, want)
	}
	if got, want := log.Rounds[0].DefenderHP, 9; got != want {
		t.Errorf("Rounds[0].DefenderHP = %d, want %d", got, want)
	}
	if got, want := log.Rounds[1].DefenderHP, 0; got != want {
		t.Errorf("Rounds[1].DefenderHP = %d, want %d", got, want)
	}
	if got, want := log.Winner, "opp"; got != want {
		t.Errorf("Winner = %q, want %q", got, want)
	}
}

func TestFight_SwingOrderChosenByRNG(t *testing.T) {
	t.Parallel()

	// Symmetric one-shotters: whoever swings first wins in one swing, so
	// the outcome reveals which combatant Intn(2) picked.
	cases := []struct {
		name         string
		order        int
		wantAttacker string
		wantWinner   string
	}{
		{name: "0 picks the player", order: 0, wantAttacker: "player", wantWinner: "player"},
		{name: "1 picks the opponent", order: 1, wantAttacker: "opp", wantWinner: "opp"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			player := combat.Combatant{Name: "player", HP: 10, MaxHP: 10, Accuracy: 10, Damage: 100}
			opp := combat.Combatant{Name: "opp", HP: 10, MaxHP: 10, Accuracy: 10, Damage: 100}

			log := combat.Fight(player, opp, &fakeRNG{values: []int{tc.order, 5, 0}})

			if got := log.Rounds[0].Attacker; got != tc.wantAttacker {
				t.Errorf("Rounds[0].Attacker = %q, want %q", got, tc.wantAttacker)
			}
			if got := log.Winner; got != tc.wantWinner {
				t.Errorf("Winner = %q, want %q", got, tc.wantWinner)
			}
		})
	}
}

func TestFight_CapDecidedByHPPercent(t *testing.T) {
	t.Parallel()

	// evenChip lands a 1-damage hit on even swings, oddChip on odd swings;
	// the other swings whiff. Player always swings first (order 0).
	evenChip := func(i int) (int, int) {
		if i%2 == 0 {
			return 10, 0
		}
		return 0, 0
	}
	oddChip := func(i int) (int, int) {
		if i%2 == 1 {
			return 10, 0
		}
		return 0, 0
	}

	cases := []struct {
		name       string
		player     combat.Combatant
		opp        combat.Combatant
		roll       func(i int) (int, int)
		wantWinner string
	}{
		{
			name:       "player ends with higher HP% and wins",
			player:     combat.Combatant{Name: "player", HP: 100, MaxHP: 100, Accuracy: 10, Damage: 1},
			opp:        combat.Combatant{Name: "opp", HP: 100, MaxHP: 100},
			roll:       evenChip,
			wantWinner: "player",
		},
		{
			name:       "exact HP% tie goes to the opponent, player loses draws",
			player:     combat.Combatant{Name: "player", HP: 100, MaxHP: 100},
			opp:        combat.Combatant{Name: "opp", HP: 100, MaxHP: 100},
			roll:       func(int) (int, int) { return 0, 0 },
			wantWinner: "opp",
		},
		{
			name:       "opponent ends with higher HP% and wins",
			player:     combat.Combatant{Name: "player", HP: 100, MaxHP: 100},
			opp:        combat.Combatant{Name: "opp", HP: 100, MaxHP: 100, Accuracy: 10, Damage: 1},
			roll:       oddChip,
			wantWinner: "opp",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			log := combat.Fight(tc.player, tc.opp, &fakeRNG{values: capScript(0, 20, tc.roll)})

			if got, want := len(log.Rounds), 20; got != want {
				t.Fatalf("len(Rounds) = %d, want %d (cap not honoured)", got, want)
			}
			if got := log.Winner; got != tc.wantWinner {
				t.Errorf("Winner = %q, want %q", got, tc.wantWinner)
			}
		})
	}
}

func TestFight_RNGConsumptionOrder(t *testing.T) {
	t.Parallel()

	player := combat.Combatant{Name: "player", HP: 10, MaxHP: 10, Accuracy: 10, Damage: 1}
	opp := combat.Combatant{Name: "opp", HP: 10, MaxHP: 10, Accuracy: 10, Damage: 100}

	// Two swings: Intn(2) for order, then Intn(20) twice per swing,
	// swinger's roll first.
	rng := &fakeRNG{values: []int{0, 10, 0, 10, 0}}
	combat.Fight(player, opp, rng)

	want := []int{2, 20, 20, 20, 20}
	if !slices.Equal(rng.calls, want) {
		t.Errorf("Intn call arguments = %v, want %v", rng.calls, want)
	}
}

// fakeRNG is a scripted RNG: it replays values in order and records the n
// passed to each Intn call so tests can assert consumption order. It
// ignores n when choosing the return value, and panics if a fight reads
// past the script so over-consumption fails loudly.
type fakeRNG struct {
	values []int
	calls  []int
	i      int
}

func (f *fakeRNG) Intn(n int) int {
	f.calls = append(f.calls, n)
	if f.i >= len(f.values) {
		panic(fmt.Sprintf("fakeRNG exhausted after %d calls (last n=%d)", f.i, n))
	}
	v := f.values[f.i]
	f.i++
	return v
}

// capScript builds an RNG script for a fight that runs the full cap
// without a KO: one Intn(2) for swing order, then swings pairs of rolls
// from roll(i) for swing i.
func capScript(order, swings int, roll func(i int) (int, int)) []int {
	s := []int{order}
	for i := 0; i < swings; i++ {
		a, d := roll(i)
		s = append(s, a, d)
	}
	return s
}
