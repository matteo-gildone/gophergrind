package combat

// RNG is the randomness seam for fights. *rand.Rand satisfies it.
type RNG interface {
	Intn(n int) int
}

// Combatant is a fighter's stat block. It is combat-local to keep this
// package independent of character and avoid an import cycle.
type Combatant struct {
	Name     string
	HP       int
	MaxHP    int
	Accuracy int
	Damage   int
	Defence  int
}

// Round records a single swing.
type Round struct {
	Attacker   string // name of the swinging combatant
	Defender   string // name of the defender
	Hit        bool
	Damage     int // dealt this swing (0 on miss)
	DefenderHP int // defender HP remaining after the swing (floored at 0)
}

// Log is the full record of a fight.
type Log struct {
	Winner string
	Rounds []Round // one entry per swing; capped at 20 swings
}

const maxRound = 10

// Fight resolves a full fight. player is the fixed attacker: on a
// round-cap HP% tie, player loses. Swing order is chosen randomly via
// rng.
func Fight(player, opponent Combatant, rng RNG) Log {
	log := Log{}
	a, d := &player, &opponent
	if rng.Intn(2) == 1 {
		a, d = d, a
	}
	for range maxRound * 2 {
		hit, damage := attack(*a, *d, rng)
		if damage > 0 {
			d.HP -= damage
			if d.HP < 0 {
				d.HP = 0
			}
		}

		log.Rounds = append(log.Rounds, Round{
			Attacker:   a.Name,
			Defender:   d.Name,
			Hit:        hit,
			Damage:     damage,
			DefenderHP: d.HP,
		})

		if d.HP == 0 {
			log.Winner = a.Name
			return log
		}

		a, d = d, a
	}
	log.Winner = higherHP(player, opponent)
	return log
}

func attack(attacker, defender Combatant, rng RNG) (hit bool, damage int) {
	if rng.Intn(20)+attacker.Accuracy > rng.Intn(20)+defender.Defence {
		return true, attacker.Damage
	}
	return false, 0
}

func higherHP(player, opponent Combatant) string {
	if player.HP*opponent.MaxHP > opponent.HP*player.MaxHP {
		return player.Name
	}
	return opponent.Name // tie falls here: player loses draws
}
