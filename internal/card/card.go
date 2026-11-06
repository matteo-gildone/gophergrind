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

var catalogue = []Card{
	{ID: "accuracy", Effect: Effect{Target: Accuracy, Amount: 2}},
	{ID: "damage", Effect: Effect{Target: Damage, Amount: 2}},
	{ID: "defence", Effect: Effect{Target: Defence, Amount: 2}},
	{ID: "hp", Effect: Effect{Target: HP, Amount: 10}},
}
