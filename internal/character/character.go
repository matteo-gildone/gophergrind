// Package character models a gopher's identity: the stats set by the
// creation quiz, current level, the permanently locked path, and whether
// the character is hardcore (permadeath) or softcore. Card offers and the
// XP curve live in other packages; character only holds and guards state.
package character

import "fmt"

// Stat identifies which of the three combat stats a quiz answer feeds.
type Stat int

const (
	Accuracy Stat = iota
	Damage
	Defence
)

// Path is the combat career a character locks into from level 3 onward.
// A character starts as PathNone and, once a path is locked, can never
// switch to another.
type Path int

const (
	PathNone    Path = iota
	PathFighter      // Defence
	PathRanger       // Accuracy
	PathMage         // Damage
)

// Mode is the death rule chosen at creation: Hardcore is permadeath,
// Softcore allows recovery. It is fixed for the character's lifetime.
type Mode int

const (
	Softcore Mode = iota
	Hardcore
)

// StartingHP is the fixed HP every character begins with, untouched by the quiz.
const StartingHP = 50
const statIncrease = 10
const pathUnlockLevel = 3

// Stats holds a character's three combat stats as flat point values.
type Stats struct{ Accuracy, Damage, Defence int }

// Character is a gopher's full state: name, HP, combat stats, level,
// locked path, and death mode.
type Character struct {
	Name  string
	HP    int
	MaxHP int
	Stats Stats
	Level int
	Path  Path
	Mode  Mode
}

// New builds a level-1 character from exactly three quiz answers, scoring
// statIncrease onto the chosen stat for each. It returns an error, and
// no usable character, if any answer is not a valid Stat.
func New(name string, mode Mode, answers [3]Stat) (Character, error) {
	c := Character{
		Name:  name,
		Mode:  mode,
		Level: 1,
		HP:    StartingHP,
		MaxHP: StartingHP,
		Path:  PathNone,
	}

	for _, answer := range answers {
		switch answer {
		case Accuracy:
			c.Stats.Accuracy += statIncrease
		case Damage:
			c.Stats.Damage += statIncrease
		case Defence:
			c.Stats.Defence += statIncrease
		default:
			return Character{}, fmt.Errorf("invalid quiz answer: %d", answer)
		}
	}

	return c, nil
}

// LevelUp advances the character by one level. It only changes the level
// counter; deciding when to level and offering cards is the caller's job.
func (c *Character) LevelUp() {
	c.Level++
}

// LockPath permanently sets the character's path. It errors if the
// character is below the path-unlock level, if p is PathNone, or if a
// different path is already locked. Re-locking the current path is a
// no-op and returns nil.
func (c *Character) LockPath(p Path) error {
	switch {
	case c.Level < pathUnlockLevel:
		return fmt.Errorf("locking path too early")
	case c.Path == p && p != PathNone:
		return nil
	case p == PathNone:
		return fmt.Errorf("you must pick a path")
	case c.Path != PathNone:
		return fmt.Errorf("you can't respec path")
	}

	c.Path = p
	return nil
}
