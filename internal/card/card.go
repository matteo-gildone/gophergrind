package card

type Target int

const (
	Accuracy Target = iota
	Damage
	Defence
	HP
)

type Effect struct {
	Target Target
	Amount int
}

type Card struct {
	ID     string
	Effect Effect
}

// RNG is the randomness seam for cards. *rand.Rand satisfies it.
type RNG interface {
	Intn(n int) int
}

const maxCards = 3

var catalogue = []Card{
	{ID: "accuracy", Effect: Effect{Target: Accuracy, Amount: 2}},
	{ID: "damage", Effect: Effect{Target: Damage, Amount: 2}},
	{ID: "defence", Effect: Effect{Target: Defence, Amount: 2}},
	{ID: "hp", Effect: Effect{Target: HP, Amount: 10}},
}

func Offer(rng RNG) []Card {
	pool := make([]Card, len(catalogue))
	copy(pool, catalogue)

	offers := make([]Card, 0, maxCards)

	for range maxCards {
		index := rng.Intn(len(pool))
		offers = append(offers, pool[index])

		lastIndex := len(pool) - 1
		pool[index], pool[lastIndex] = pool[lastIndex], pool[index]
		pool = pool[:lastIndex]
	}
	return offers
}
