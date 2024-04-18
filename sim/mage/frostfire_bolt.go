package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (mage *Mage) registerFrostfireBoltSpell() {

	hasGlyph := mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfFrostfire)

	mage.FrostfireBolt = mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 44614},
		SpellSchool:    core.SpellSchoolFire,
		ProcMask:       core.ProcMaskSpellDamage,
		Flags:          SpellFlagMage | ArcaneMissileSpells | HotStreakSpells | core.SpellFlagAPL,
		ClassSpellMask: MageSpellFrostfireBolt,
		MissileSpeed:   28,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.09,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 2500,
			},
		},

		DamageMultiplier: 1,
		CritMultiplier:   mage.DefaultSpellCritMultiplier(),
		BonusCoefficient: 0.977,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label:     "FrostfireBolt",
				MaxStacks: 3,
				Duration:  time.Second * 12,
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, _ bool) {
				dot.SnapshotBaseDamage = 0.00712 * mage.ScalingBaseDamage
				dot.SnapshotBaseDamage *= float64(dot.Aura.GetStacks())
			},

			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.Spell.OutcomeAlwaysHit)
			},
			BonusCoefficient: 0.00733,
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 0.949 * mage.ScalingBaseDamage
			// Not sure if double dipping exists in Cata. Removed for now.
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				dot := spell.Dot(target)
				if result.Landed() {
					if dot.IsActive() {
						dot.Refresh(sim)
						dot.AddStack(sim)
						dot.TakeSnapshot(sim, true)
					} else if hasGlyph {
						dot.Apply(sim)
						dot.SetStacks(sim, 1)
						dot.TakeSnapshot(sim, true)
					}
					spell.DealDamage(sim, result)
				}
			})
		},
	})
}
