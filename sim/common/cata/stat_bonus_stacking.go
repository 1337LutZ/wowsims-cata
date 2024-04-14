package cata

import (
	"time"

	"github.com/wowsims/cata/sim/common/shared"
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/stats"
)

func init() {
	shared.NewStackingStatBonusCD(shared.StackingStatBonusCD{
		Name:                  "World-Queller Focus",
		ID:                    63842,
		AuraID:                90900,
		Bonus:                 stats.Stats{stats.SpellPower: 313},
		MaxStacks:             5,
		ProcMask:              core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Duration:              time.Second * 20,
		CD:                    time.Minute * 2,
		Callback:              core.CallbackOnCastComplete,
		Harmful:               false,
		TrinketLimitsDuration: true,
		ProcChance:            1,
		IsDefensive:           false,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Gale of Shadows",
		ID:         56138,
		AuraID:     90943,
		Bonus:      stats.Stats{stats.SpellPower: 15},
		MaxStacks:  20,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt,
		Harmful:    false,
		ProcChance: 1,
		Icd:        time.Millisecond * 500,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Gale of Shadows (Heroic)",
		ID:         56462,
		AuraID:     90986,
		Bonus:      stats.Stats{stats.SpellPower: 17},
		MaxStacks:  20,
		ProcMask:   core.ProcMaskSpellDamage | core.ProcMaskSpellHealing,
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnPeriodicDamageDealt | core.CallbackOnPeriodicHealDealt,
		Harmful:    false,
		ProcChance: 1,
		Icd:        time.Millisecond * 500,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Tia's Grace",
		ID:         55874,
		AuraID:     92086,
		Bonus:      stats.Stats{stats.Agility: 30},
		MaxStacks:  10,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		Harmful:    false,
		ProcChance: 1,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Tia's Grace (Heroic)",
		ID:         56394,
		AuraID:     92090,
		Bonus:      stats.Stats{stats.Agility: 34},
		MaxStacks:  10,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		Harmful:    false,
		ProcChance: 1,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Tia's Grace (Heroic)",
		ID:         62050,
		AuraID:     92090,
		Bonus:      stats.Stats{stats.Spirit: 80},
		MaxStacks:  5,
		ProcMask:   core.ProcMaskSpellHealing,
		Duration:   time.Second * 20,
		Callback:   core.CallbackOnHealDealt | core.CallbackOnPeriodicHealDealt,
		Harmful:    false,
		ProcChance: 1,
		Icd:        time.Second * 2,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "Fluid Death",
		ID:         58181,
		AuraID:     92105,
		Bonus:      stats.Stats{stats.Agility: 28},
		MaxStacks:  10,
		ProcMask:   core.ProcMaskMeleeOrRanged,
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		Harmful:    false,
		ProcChance: 1,
	})

	shared.NewStackingStatBonusEffect(shared.StackingStatBonusEffect{
		Name:       "License to Slay",
		ID:         58180,
		AuraID:     91810,
		Bonus:      stats.Stats{stats.Agility: 38},
		MaxStacks:  10,
		ProcMask:   core.ProcMaskMelee,
		Duration:   time.Second * 15,
		Callback:   core.CallbackOnSpellHitDealt,
		Harmful:    false,
		ProcChance: 1,
	})
}
