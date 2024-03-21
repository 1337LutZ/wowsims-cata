package subtlety

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/rogue"
)

func (subRogue *SubtletyRogue) registerSanguinaryVein() {
	if subRogue.Talents.SanguinaryVein == 0 {
		return
	}

	svBonus := 1 + 0.08*float64(subRogue.Talents.SanguinaryVein)
	isApplied := false

	svDebuffArray := subRogue.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			Label:    "Sanguinary Vein Debuff",
			Duration: core.NeverExpires,
			// Action ID Suppressed to not fill debuff log
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				if !isApplied {
					isApplied = true
					subRogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier *= svBonus
				}
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				if !aura.Unit.HasAuraWithTag(rogue.RogueBleedTag) {
					subRogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= svBonus
					isApplied = false
				}
			},
			OnReset: func(aura *core.Aura, sim *core.Simulation) {
				if isApplied {
					isApplied = false
					subRogue.AttackTables[aura.Unit.UnitIndex].DamageTakenMultiplier /= svBonus
				}
			},
		})
	})

	subRogue.Env.RegisterPreFinalizeEffect(func() {
		if subRogue.Rupture != nil {
			subRogue.Rupture.RelatedAuras = append(subRogue.Rupture.RelatedAuras, svDebuffArray)
		}
		if subRogue.Hemorrhage != nil && subRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfHemorrhage) {
			subRogue.Hemorrhage.RelatedAuras = append(subRogue.Hemorrhage.RelatedAuras, svDebuffArray)
		}
	})

	subRogue.RegisterAura(core.Aura{
		Label:    "Sanguinary Vein Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() {
				return
			}

			if spell == subRogue.Rupture || (spell == subRogue.Hemorrhage && subRogue.HasPrimeGlyph(proto.RoguePrimeGlyph_GlyphOfHemorrhage)) {
				aura := svDebuffArray.Get(result.Target)
				dot := spell.Dot(result.Target)
				aura.Duration = dot.TickLength * time.Duration(dot.NumberOfTicks)
				aura.Activate(sim)
			}
		},
	})
}
