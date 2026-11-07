package screens

import (
	"fmt"
	"game/internal/domain"
)

// Both character views use the same labels and formatting.
func statRows(st domain.Stats) [][2]string {
	return [][2]string{
		{"Здоровье (HP)", fmt.Sprintf("%d / %d", st.Health, st.MaxHealth)},
		{"Мана (MP)", fmt.Sprintf("%d / %d", st.Mana, st.MaxMana)},
		{"Сила", fmt.Sprint(st.Strength)}, {"Ловкость", fmt.Sprint(st.Dexterity)},
		{"Интеллект", fmt.Sprint(st.Intelligence)}, {"Броня", fmt.Sprint(st.Armor)},
		{"Урон заклинаний", fmt.Sprint(st.SpellDamage)},
		{"Шанс критического удара", fmt.Sprintf("%g%%", st.CritChance)},
		{"Скорость передвижения", fmt.Sprintf("%g%%", st.MoveSpeedPercent)},
	}
}
