package character_test

import (
	"gophergrind/internal/card"
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

func TestLevelUp(t *testing.T) {
	t.Parallel()

	answers := [3]character.Stat{character.Accuracy, character.Damage, character.Defence}
	c, err := character.New("gopher", character.Softcore, answers)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	err = c.LevelUp()
	if err != nil {
		t.Fatalf("LevelUp returned error: %v", err)
	}
	if c.Level != 2 {
		t.Errorf("Level after one LevelUp = %d, want 2", c.Level)
	}
}

func TestLevelUp_RequiresPathAtLevel3(t *testing.T) {
	t.Parallel()

	c := atLevel(t, 3)
	want := c

	err := c.LevelUp()
	if err == nil {
		t.Fatal("expected error got nil")
	}
	if c != want {
		t.Errorf("character = %+v, want %+v", c, want)
	}
}

func TestApplyEffect(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		effect card.Effect
		want   character.Character
	}{
		{
			name: "apply HP card", effect: card.Effect{Target: card.HP, Amount: 10}, want: character.Character{
				Name:  "gopher",
				Mode:  character.Softcore,
				Level: 1,
				HP:    60,
				MaxHP: 60,
				Stats: character.Stats{
					Accuracy: 10,
					Damage:   10,
					Defence:  10,
				},
				Path: character.PathNone,
			},
		},
		{
			name: "apply Accuracy card", effect: card.Effect{Target: card.Accuracy, Amount: 2}, want: character.Character{
				Name:  "gopher",
				Mode:  character.Softcore,
				Level: 1,
				HP:    50,
				MaxHP: 50,
				Stats: character.Stats{
					Accuracy: 12,
					Damage:   10,
					Defence:  10,
				},
				Path: character.PathNone,
			},
		},
		{
			name: "apply Defence card", effect: card.Effect{Target: card.Defence, Amount: 2}, want: character.Character{
				Name:  "gopher",
				Mode:  character.Softcore,
				Level: 1,
				HP:    50,
				MaxHP: 50,
				Stats: character.Stats{
					Accuracy: 10,
					Damage:   10,
					Defence:  12,
				},
				Path: character.PathNone,
			},
		},
		{
			name: "apply Damage card", effect: card.Effect{Target: card.Damage, Amount: 2}, want: character.Character{
				Name:  "gopher",
				Mode:  character.Softcore,
				Level: 1,
				HP:    50,
				MaxHP: 50,
				Stats: character.Stats{
					Accuracy: 10,
					Damage:   12,
					Defence:  10,
				},
				Path: character.PathNone,
			},
		},
	}

	answers := [3]character.Stat{character.Accuracy, character.Damage, character.Defence}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c, err := character.New("gopher", character.Softcore, answers)
			if err != nil {
				t.Fatalf("New returned error: %v", err)
			}

			err = c.ApplyEffect(tc.effect)
			if err != nil {
				t.Fatalf("ApplyEffect(Effect) returned error: %v", err)
			}

			if c != tc.want {
				t.Errorf("character = %+v, want %+v", c, tc.want)
			}
		})
	}
}

func TestApplyEffect_InvalidTarget(t *testing.T) {
	t.Parallel()
	answers := [3]character.Stat{character.Accuracy, character.Damage, character.Defence}
	c, err := character.New("gopher", character.Softcore, answers)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	want := c
	err = c.ApplyEffect(card.Effect{Target: 99, Amount: 2})

	if err == nil {
		t.Error("ApplyEffect(invalid target) returned nil error, want error")
	}

	if c != want {
		t.Errorf("character = %+v, want %+v", c, want)
	}
}

func TestLevellingFlow(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		wantPath character.Path
		want     character.Character
	}{
		{
			name:     "pick ranger path",
			wantPath: character.PathRanger,
			want: character.Character{
				Name:  "gopher",
				Mode:  character.Softcore,
				Level: 4,
				HP:    50,
				MaxHP: 50,
				Stats: character.Stats{
					Accuracy: 20,
					Damage:   10,
					Defence:  10,
				},
				Path: character.PathRanger,
			},
		},
		{
			name:     "pick fighter path",
			wantPath: character.PathFighter,
			want: character.Character{
				Name:  "gopher",
				Mode:  character.Softcore,
				Level: 4,
				HP:    50,
				MaxHP: 50,
				Stats: character.Stats{
					Accuracy: 10,
					Damage:   10,
					Defence:  20,
				},
				Path: character.PathFighter,
			},
		},
		{
			name:     "pick mage path",
			wantPath: character.PathMage,
			want: character.Character{
				Name:  "gopher",
				Mode:  character.Softcore,
				Level: 4,
				HP:    50,
				MaxHP: 50,
				Stats: character.Stats{
					Accuracy: 10,
					Damage:   20,
					Defence:  10,
				},
				Path: character.PathMage,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c := atLevel(t, 3)
			if err := c.ChoosePath(tc.wantPath); err != nil {
				t.Fatalf("locking starting path %v: %v", tc.wantPath, err)
			}

			err := c.LevelUp()
			if err != nil {
				t.Fatalf("LevelUp returned error: %v", err)
			}

			if c != tc.want {
				t.Errorf("character = %+v, want %+v", c, tc.want)
			}
		})
	}
}

func TestChoosePath(t *testing.T) {
	t.Parallel()

	t.Run("rejects locking below level 3", func(t *testing.T) {
		t.Parallel()
		c := atLevel(t, 2)
		want := c

		if err := c.ChoosePath(character.PathMage); err == nil {
			t.Fatal("expected error got nil")
		}

		if c != want {
			t.Errorf("character = %+v, want %+v", c, want)
		}
	})

	t.Run("locking no path is not possible", func(t *testing.T) {
		t.Parallel()
		c := atLevel(t, 3)
		want := c

		if err := c.ChoosePath(character.PathNone); err == nil {
			t.Error("locking no path should return error, got nil")
		}

		if c != want {
			t.Errorf("character = %+v, want %+v", c, want)
		}
	})

	t.Run("rejects crossing to a different path", func(t *testing.T) {
		t.Parallel()
		c := atLevel(t, 3)
		if err := c.ChoosePath(character.PathMage); err != nil {
			t.Fatalf("initial Mage choice returned error: %v", err)
		}

		want := c

		if err := c.ChoosePath(character.PathRanger); err == nil {
			t.Error("choosing Ranger after Mage returned nil error, want error")
		}

		if c != want {
			t.Errorf("character = %+v, want %+v", c, want)
		}
	})

	t.Run("rejects locking to invalid path", func(t *testing.T) {
		t.Parallel()
		c := atLevel(t, 3)
		want := c

		if err := c.ChoosePath(character.Path(90)); err == nil {
			t.Fatal("expected error got nil")
		}

		if c != want {
			t.Errorf("character = %+v, want %+v", c, want)
		}
	})

}

// atLevel builds a character advanced to the given level via LevelUp.
func atLevel(t *testing.T, level int) character.Character {
	t.Helper()
	answers := [3]character.Stat{character.Accuracy, character.Damage, character.Defence}
	c, err := character.New("gopher", character.Softcore, answers)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	for c.Level < level {
		err = c.LevelUp()
		if err != nil {
			t.Fatalf("LevelUp returned error: %v", err)
		}
	}
	return c
}
