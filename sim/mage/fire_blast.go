package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (mage *Mage) registerFireBlastSpell() {
	mage.FireBlast = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 2136},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          SpellFlagMage | HotStreakSpells | core.SpellFlagAPL,
		ClassSpellMask: MageSpellFireBlast,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.21,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{ // Note: Impact talent triggers CD refresh on spell *land*, not cast
				Timer:    mage.NewTimer(),
				Duration: time.Second * 8,
			},
		},

		BonusCritRating: 0 +
			float64(mage.Talents.ImprovedFireBlast)*4*core.CritRatingPerCritChance,
		DamageMultiplierAdditive: 1 +
			.01*float64(mage.Talents.FirePower),
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.429,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 1.113 * mage.ScalingBaseDamage
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			// impact thing to spread dots goes here most likely
			// not working, at least on dummies. will need to test if duration refresh
		},
	})
}
