package character_test

import (
	"testing"

	"gophergrind/internal/character"
)

func TestNew_StatsFromAnswers(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		answers [3]character.Stat
		want    character.Stats
	}{
		{
			name:    "all accuracy",
			answers: [3]character.Stat{character.Accuracy, character.Accuracy, character.Accuracy},
			want:    character.Stats{Accuracy: 30},
		},
		{
			name:    "one each is the neutral 10/10/10",
			answers: [3]character.Stat{character.Accuracy, character.Damage, character.Defence},
			want:    character.Stats{Accuracy: 10, Damage: 10, Defence: 10},
		},
		{
			name:    "two accuracy one defence",
			answers: [3]character.Stat{character.Accuracy, character.Accuracy, character.Defence},
			want:    character.Stats{Accuracy: 20, Defence: 10},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, err := character.New("gopher", character.Softcore, tc.answers)
			if err != nil {
				t.Fatalf("New returned error: %v", err)
			}
			if c.Stats != tc.want {
				t.Errorf("Stats = %+v, want %+v", c.Stats, tc.want)
			}
		})
	}
}

func TestNew_Invariants(t *testing.T) {
	t.Parallel()

	stats := []character.Stat{character.Accuracy, character.Damage, character.Defence}

	// Every one of the 27 possible answer combinations must hold the
	// invariants: 30 points total, none negative, fixed HP, level 1, no path.
	for _, a := range stats {
		for _, b := range stats {
			for _, d := range stats {
				answers := [3]character.Stat{a, b, d}

				c, err := character.New("gopher", character.Softcore, answers)
				if err != nil {
					t.Fatalf("New(%v) returned error: %v", answers, err)
				}

				total := c.Stats.Accuracy + c.Stats.Damage + c.Stats.Defence
				if total != 30 {
					t.Errorf("New(%v): stat total = %d, want 30", answers, total)
				}
				if c.Stats.Accuracy < 0 || c.Stats.Damage < 0 || c.Stats.Defence < 0 {
					t.Errorf("New(%v): negative stat in %+v", answers, c.Stats)
				}
				if c.HP != character.StartingHP || c.MaxHP != character.StartingHP {
					t.Errorf("New(%v): HP=%d MaxHP=%d, want %d for both", answers, c.HP, c.MaxHP, character.StartingHP)
				}
				if c.Level != 1 {
					t.Errorf("New(%v): Level = %d, want 1", answers, c.Level)
				}
				if c.Path != character.PathNone {
					t.Errorf("New(%v): Path = %v, want PathNone", answers, c.Path)
				}
			}
		}
	}
}

func TestNew_StoresMode(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		mode character.Mode
	}{
		{name: "softcore", mode: character.Softcore},
		{name: "hardcore", mode: character.Hardcore},
	}

	answers := [3]character.Stat{character.Accuracy, character.Damage, character.Defence}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, err := character.New("gopher", tc.mode, answers)
			if err != nil {
				t.Fatalf("New returned error: %v", err)
			}
			if c.Mode != tc.mode {
				t.Errorf("Mode = %v, want %v", c.Mode, tc.mode)
			}
		})
	}
}
func TestNew_RejectsInvalidAnswer(t *testing.T) {
	t.Parallel()

	// A Stat outside the defined set must be rejected, not silently scored.
	answers := [3]character.Stat{character.Accuracy, character.Stat(99), character.Defence}

	if _, err := character.New("gopher", character.Softcore, answers); err == nil {
		t.Fatal("New with an out-of-range Stat returned nil error, want error")
	}
}

func TestLockPath(t *testing.T) {
	t.Parallel()

	// atLevel builds a character advanced to the given level via LevelUp.
	atLevel := func(t *testing.T, level int, start character.Path) character.Character {
		t.Helper()
		answers := [3]character.Stat{character.Accuracy, character.Damage, character.Defence}
		c, err := character.New("gopher", character.Softcore, answers)
		if err != nil {
			t.Fatalf("New returned error: %v", err)
		}
		for c.Level < level {
			c.LevelUp()
		}
		if start != character.PathNone {
			if err := c.LockPath(start); err != nil {
				t.Fatalf("locking starting path %v: %v", start, err)
			}
		}
		return c
	}

	t.Run("rejects locking below level 3", func(t *testing.T) {
		t.Parallel()
		c := atLevel(t, 2, character.PathNone)

		if err := c.LockPath(character.PathMage); err == nil {
			t.Error("LockPath at level 2 returned nil error, want error")
		}
		if c.Path != character.PathNone {
			t.Errorf("Path = %v after rejected lock, want PathNone", c.Path)
		}
	})

	t.Run("locks a path from level 3", func(t *testing.T) {
		t.Parallel()
		c := atLevel(t, 3, character.PathNone)

		if err := c.LockPath(character.PathMage); err != nil {
			t.Fatalf("LockPath returned error: %v", err)
		}
		if c.Path != character.PathMage {
			t.Errorf("Path = %v, want PathMage", c.Path)
		}
	})

	t.Run("re-locking the same path is idempotent", func(t *testing.T) {
		t.Parallel()
		c := atLevel(t, 3, character.PathMage)

		if err := c.LockPath(character.PathMage); err != nil {
			t.Errorf("re-locking same path returned error: %v", err)
		}
		if c.Path != character.PathMage {
			t.Errorf("Path = %v, want PathMage", c.Path)
		}
	})

	t.Run("rejects crossing to a different path", func(t *testing.T) {
		t.Parallel()
		c := atLevel(t, 3, character.PathMage)

		if err := c.LockPath(character.PathRanger); err == nil {
			t.Error("crossing path returned nil error, want error")
		}
		if c.Path != character.PathMage {
			t.Errorf("Path = %v after rejected cross, want PathMage (unchanged)", c.Path)
		}
	})

	t.Run("rejects locking PathNone", func(t *testing.T) {
		t.Parallel()
		c := atLevel(t, 3, character.PathNone)

		if err := c.LockPath(character.PathNone); err == nil {
			t.Error("LockPath(PathNone) returned nil error, want error")
		}
	})
}
func TestLevelUp(t *testing.T) {
	t.Parallel()

	answers := [3]character.Stat{character.Accuracy, character.Damage, character.Defence}
	c, err := character.New("gopher", character.Softcore, answers)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	c.LevelUp()
	if c.Level != 2 {
		t.Errorf("Level after one LevelUp = %d, want 2", c.Level)
	}
}
