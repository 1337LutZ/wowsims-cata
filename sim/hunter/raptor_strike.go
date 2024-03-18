package hunter

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

// Todo: Introduce weaving again? With the glyph, might be interesting
func (hunter *Hunter) registerRaptorStrikeSpell() {
	hunter.RaptorStrike = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 48996},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHAuto | core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		FocusCost: core.FocusCostOptions{
			Cost: 0,
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
		DamageMultiplier: 1,
		CritMultiplier:   hunter.CritMultiplier(false, false, false),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := 335 +
				spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
				spell.BonusWeaponDamage()

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

// Returns true if the regular melee swing should be used, false otherwise.
func (hunter *Hunter) TryRaptorStrike(_ *core.Simulation, mhSwingSpell *core.Spell) *core.Spell {
	//if hunter.Rotation.Weave == proto.Hunter_Rotation_WeaveAutosOnly || !hunter.RaptorStrike.IsReady(sim) || hunter.CurrentMana() < hunter.RaptorStrike.DefaultCast.Cost {
	//	return nil
	//}
	return mhSwingSpell
}
