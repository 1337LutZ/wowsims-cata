package survival

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (hunter *SurvivalHunter) registerBlackArrowSpell(timer *core.Timer) {
	if !hunter.Talents.BlackArrow {
		return
	}

	actionID := core.ActionID{SpellID: 3674}

	hunter.Hunter.BlackArrow = hunter.Hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,
		ProcMask:    core.ProcMaskRangedSpecial,

		Flags: core.SpellFlagMeleeMetrics | core.SpellFlagAPL,
		FocusCost: core.FocusCostOptions{
			Cost: 35,
		},
		MissileSpeed: 40,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			IgnoreHaste: true, // Hunter GCD is locked at 1.5s
			CD: core.Cooldown{
				Timer:    timer,
				Duration: time.Second*30 - time.Second*2*time.Duration(hunter.Talents.Resourcefulness),
			},
		},

		DamageMultiplier: 1 + .10*float64(hunter.Talents.TrapMastery),
		ThreatMultiplier: 1,
		CritMultiplier:   hunter.SpellCritMultiplier(1, float64(hunter.Talents.Toxicology)*0.5),

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Black Arrow Dot",
			},
			NumberOfTicks: 10,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				// https://web.archive.org/web/20120207222124/http://elitistjerks.com/f74/t110306-hunter_faq_cataclysm_edition_read_before_asking_questions/
				//  66.5% RAP + 2849 (total damage) - changed 6/28 in 4.2 (based off spell crit multiplier, modified by toxicology)
				// https://wago.tools/db2/SpellEffect?build=4.4.0.53750&filter[SpellID]=exact%3A3674&page=1
				baseDamage := 2852.0
				rap := dot.Spell.RangedAttackPower(target)
				percentageOfRAP := 0.665

				// SnapshotBaseDamage calculation for the DoT, divided by 10 to spread across all ticks
				dot.SnapshotBaseDamage = (baseDamage + (percentageOfRAP * rap)) / 10

				attackTable := dot.Spell.Unit.AttackTables[target.UnitIndex]
				dot.SnapshotCritChance = dot.Spell.PhysicalCritChance(attackTable)
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex], true)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeRangedHitAndCritSnapshot)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeRangedHit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			})
		},
	})
}
