package domain

type Character struct {
	Name, Class string
	Level       int
	Stats       Stats
	Inventory   [12]*Item
}

// A nil slot is empty. Items can be populated by future gameplay systems.
type Item struct {
	Name     string
	Quantity int
}

type Stats struct {
	Health, MaxHealth                 int
	Mana, MaxMana                     int
	Strength, Dexterity, Intelligence int
	Armor, SpellDamage                int
	CritChance, MoveSpeedPercent      float64
}

type Stat int

const (
	Health Stat = iota
	Mana
	Strength
	Dexterity
	Intelligence
	Armor
	SpellDamage
	CritChance
	MoveSpeed
	StatCount
)

// Adjust changes a base stat by one step. Resource adjustments preserve the
// missing amount, while keeping current values within the new maximum.
func (s *Stats) Adjust(stat Stat, delta int) {
	if delta != -1 && delta != 1 {
		return
	}
	change := func(value *int, minimum int) {
		*value = max(minimum, min(9999, *value+delta))
	}
	switch stat {
	case Health:
		old := s.MaxHealth
		change(&s.MaxHealth, 1)
		s.Health = max(0, min(s.MaxHealth, s.Health+s.MaxHealth-old))
	case Mana:
		old := s.MaxMana
		change(&s.MaxMana, 0)
		s.Mana = max(0, min(s.MaxMana, s.Mana+s.MaxMana-old))
	case Strength:
		change(&s.Strength, 0)
	case Dexterity:
		change(&s.Dexterity, 0)
	case Intelligence:
		change(&s.Intelligence, 0)
	case Armor:
		change(&s.Armor, 0)
	case SpellDamage:
		change(&s.SpellDamage, 0)
	case CritChance:
		s.CritChance = max(0, min(100, s.CritChance+float64(delta)))
	case MoveSpeed:
		s.MoveSpeedPercent = max(0, min(999, s.MoveSpeedPercent+float64(delta)))
	}
}

func NewArcanist() Character {
	return Character{Name: "Арканист", Class: "Маг", Level: 1, Stats: Stats{
		Health: 100, MaxHealth: 100, Mana: 80, MaxMana: 80,
		Strength: 5, Dexterity: 5, Intelligence: 10,
		Armor: 0, SpellDamage: 12, CritChance: 5, MoveSpeedPercent: 100,
	}}
}
