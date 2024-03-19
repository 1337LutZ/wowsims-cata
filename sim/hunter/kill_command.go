package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
)

func (hunter *Hunter) registerKillCommandSpell() {
	if hunter.Pet == nil {
		return
	}

	actionID := core.ActionID{SpellID: 34026}

	hunter.KillCommand = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMelee,
		Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagAPL,

		FocusCost: core.FocusCostOptions{
			Cost: 40 - core.TernaryFloat64(hunter.HasPrimeGlyph(proto.HunterPrimeGlyph_GlyphOfKillCommand), 3, 0), // Todo: Check if changed by other stuff
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second,
			},
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Second * 6,
			},
		},
		DamageMultiplierAdditive: 1 + (float64(hunter.Talents.ImprovedKillCommand) * 0.05),
		CritMultiplier:           hunter.CritMultiplier(false, false, false),
		ThreatMultiplier:         1,
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if hunter.Talents.KillingStreak > 0 {
				if !hunter.KillingStreakCounterAura.IsActive() {
					hunter.KillingStreakCounterAura.Activate(sim)
				}
			}
			baseDamage := 0.516*spell.RangedAttackPower(target) + 918 //Todo: fix calc
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeRangedHitAndCrit)
		},
	})
}
