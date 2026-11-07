package domain

import "testing"

func TestAdjustStats(t *testing.T) {
	c := NewArcanist()
	original := c.Stats
	for stat := Health; stat < StatCount; stat++ {
		c.Stats.Adjust(stat, 1)
	}
	want := Stats{Health: 101, MaxHealth: 101, Mana: 81, MaxMana: 81, Strength: 6, Dexterity: 6, Intelligence: 11, Armor: 1, SpellDamage: 13, CritChance: 6, MoveSpeedPercent: 101}
	if c.Stats != want {
		t.Fatalf("increment: got %+v, want %+v", c.Stats, want)
	}
	for stat := Health; stat < StatCount; stat++ {
		c.Stats.Adjust(stat, -1)
	}
	if c.Stats != original {
		t.Fatalf("decrement: got %+v", c.Stats)
	}
}

func TestStatLimits(t *testing.T) {
	s := Stats{Health: 1, MaxHealth: 1}
	for stat := Health; stat < StatCount; stat++ {
		s.Adjust(stat, -1)
	}
	if s != (Stats{Health: 1, MaxHealth: 1}) {
		t.Fatalf("lower bounds: %+v", s)
	}
	s = Stats{Health: 9999, MaxHealth: 9999, Mana: 9999, MaxMana: 9999, Strength: 9999, Dexterity: 9999, Intelligence: 9999, Armor: 9999, SpellDamage: 9999, CritChance: 100, MoveSpeedPercent: 999}
	want := s
	for stat := Health; stat < StatCount; stat++ {
		s.Adjust(stat, 1)
	}
	if s != want {
		t.Fatalf("upper bounds: %+v", s)
	}
}

func TestResourceAdjustmentPreservesMissingAmount(t *testing.T) {
	s := Stats{Health: 50, MaxHealth: 100, Mana: 0, MaxMana: 1}
	s.Adjust(Health, 1)
	s.Adjust(Mana, -1)
	if s.Health != 51 || s.MaxHealth != 101 || s.Mana != 0 || s.MaxMana != 0 {
		t.Fatalf("resources: %+v", s)
	}
}
