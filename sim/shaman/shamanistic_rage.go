package shaman

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (shaman *Shaman) registerShamanisticRageCD() {
	if !shaman.Talents.ShamanisticRage {
		return
	}

	actionID := core.ActionID{SpellID: 30823}
	srAura := shaman.RegisterAura(core.Aura{
		Label:    "Shamanistic Rage",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier *= 0.7
			shaman.PseudoStats.CostMultiplier -= 1.0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageTakenMultiplier /= 0.7
			shaman.PseudoStats.CostMultiplier += 1.0
		},
	})

	spell := shaman.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Minute * 1,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			srAura.Activate(sim)
		},
	})

	shaman.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeMana,
		//TODO: This might change since it works a bit different now
		ShouldActivate: func(sim *core.Simulation, character *core.Character) bool {
			return character.CurrentManaPercent() <= 0.2
		},
	})
}
