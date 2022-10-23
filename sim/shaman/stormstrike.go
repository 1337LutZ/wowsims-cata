package shaman

import (
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/items"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/core/stats"
)

var StormstrikeActionID = core.ActionID{SpellID: 17364}
var TotemOfTheDancingFlame int32 = 45169
var TotemOfDueling int32 = 40322

func (shaman *Shaman) StormstrikeDebuffAura(target *core.Unit) *core.Aura {
	return target.GetOrRegisterAura(core.Aura{
		Label:     "Stormstrike-" + shaman.Label,
		ActionID:  StormstrikeActionID,
		Duration:  time.Second * 12,
		MaxStacks: 4,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shaman.AttackTables[aura.Unit.UnitIndex].NatureDamageTakenMultiplier *= core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfStormstrike), 1.28, 1.2)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shaman.AttackTables[aura.Unit.UnitIndex].NatureDamageTakenMultiplier /= core.TernaryFloat64(shaman.HasMajorGlyph(proto.ShamanMajorGlyph_GlyphOfStormstrike), 1.28, 1.2)

		},
		OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
			if spell.Unit != &shaman.Unit {
				return
			}
			if spell.SpellSchool != core.SpellSchoolNature {
				return
			}
			if !spellEffect.Landed() || spellEffect.Damage == 0 {
				return
			}
			aura.RemoveStack(sim)
		},
	})
}

func (shaman *Shaman) newStormstrikeHitSpell(isMH bool) *core.Spell {
	var flatDamageBonus float64 = 0
	if shaman.Equip[items.ItemSlotRanged].ID == TotemOfTheDancingFlame {
		flatDamageBonus += 155
	}

	var procMask core.ProcMask
	if isMH {
		procMask = core.ProcMaskMeleeMHSpecial
	} else {
		procMask = core.ProcMaskMeleeOHSpecial
	}

	return shaman.RegisterSpell(core.SpellConfig{
		ActionID:    StormstrikeActionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    procMask,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		DamageMultiplier: core.TernaryFloat64(shaman.HasSetBonus(ItemSetWorldbreakerBattlegear, 2), 1.2, 1),
		CritMultiplier:   shaman.DefaultMeleeCritMultiplier(),
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			var baseDamage float64
			if isMH {
				baseDamage = flatDamageBonus +
					spell.Unit.MHWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			} else {
				baseDamage = flatDamageBonus +
					0.5*spell.Unit.OHWeaponDamage(sim, spell.MeleeAttackPower()) +
					spell.BonusWeaponDamage()
			}

			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialCritOnly)
		},
	})
}

func (shaman *Shaman) registerStormstrikeSpell() {
	mhHit := shaman.newStormstrikeHitSpell(true)
	ohHit := shaman.newStormstrikeHitSpell(false)

	baseCost := baseMana * 0.08
	if shaman.Equip[items.ItemSlotRanged].ID == StormfuryTotem {
		baseCost -= 22
	}

	ssDebuffAura := shaman.StormstrikeDebuffAura(shaman.CurrentTarget)

	var skyshatterAura *core.Aura
	if shaman.HasSetBonus(ItemSetSkyshatterHarness, 4) {
		skyshatterAura = shaman.NewTemporaryStatsAura("Skyshatter 4pc AP Bonus", core.ActionID{SpellID: 38432}, stats.Stats{stats.AttackPower: 70}, time.Second*12)
	}
	var totemOfDuelingAura *core.Aura
	if shaman.Equip[items.ItemSlotRanged].ID == TotemOfDueling {
		totemOfDuelingAura = shaman.NewTemporaryStatsAura("Essense of the Storm", core.ActionID{SpellID: 60766},
			stats.Stats{stats.MeleeHaste: 60, stats.SpellHaste: 60}, time.Second*6)
	}

	manaMetrics := shaman.NewManaMetrics(core.ActionID{SpellID: 51522})

	cooldownTime := time.Duration(core.TernaryFloat64(shaman.HasSetBonus(ItemSetGladiatorsEarthshaker, 4), 6, 8))

	shaman.Stormstrike = shaman.RegisterSpell(core.SpellConfig{
		ActionID:     StormstrikeActionID,
		SpellSchool:  core.SpellSchoolPhysical,
		ProcMask:     core.ProcMaskMeleeMHSpecial,
		Flags:        core.SpellFlagMeleeMetrics,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost: baseCost,
				GCD:  core.GCDDefault,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    shaman.NewTimer(),
				Duration: time.Second * cooldownTime,
			},
		},

		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				if shaman.Talents.ImprovedStormstrike > 0 {
					if sim.RandomFloat("Improved Stormstrike") < 0.5*float64(shaman.Talents.ImprovedStormstrike) {
						shaman.AddMana(sim, baseMana*0.2, manaMetrics, true)
					}
				}
				ssDebuffAura.Activate(sim)
				ssDebuffAura.SetStacks(sim, 4)

				if skyshatterAura != nil {
					skyshatterAura.Activate(sim)
				}
				if totemOfDuelingAura != nil {
					totemOfDuelingAura.Activate(sim)
				}

				mhHit.Cast(sim, target)
				ohHit.Cast(sim, target)
				shaman.Stormstrike.SpellMetrics[target.UnitIndex].Casts -= 2
				shaman.Stormstrike.SpellMetrics[target.UnitIndex].Hits--
			}
			spell.DealOutcome(sim, result)
		},
	})
}
