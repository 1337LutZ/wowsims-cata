package unholy

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/death_knight"
)

func RegisterUnholyDeathKnight() {
	core.RegisterAgentFactory(
		proto.Player_UnholyDeathKnight{},
		proto.Spec_SpecUnholyDeathKnight,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewUnholyDeathKnight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_UnholyDeathKnight)
			if !ok {
				panic("Invalid spec value for Unholy Death Knight!")
			}
			player.Spec = playerSpec
		},
	)
}

type UnholyDeathKnight struct {
	*death_knight.DeathKnight

	lastScourgeStrikeDamage float64
}

func NewUnholyDeathKnight(character *core.Character, player *proto.Player) *UnholyDeathKnight {
	unholyOptions := player.GetUnholyDeathKnight().Options

	uhdk := &UnholyDeathKnight{
		DeathKnight: death_knight.NewDeathKnight(character, death_knight.DeathKnightInputs{
			Spec: proto.Spec_SpecUnholyDeathKnight,

			StartingRunicPower: unholyOptions.ClassOptions.StartingRunicPower,
			PetUptime:          unholyOptions.ClassOptions.PetUptime,
			IsDps:              true,

			UseAMS:            unholyOptions.UseAms,
			AvgAMSSuccessRate: unholyOptions.AvgAmsSuccessRate,
			AvgAMSHit:         unholyOptions.AvgAmsHit,
		}, player.TalentsString, 56835),
	}

	uhdk.Inputs.UnholyFrenzyTarget = unholyOptions.UnholyFrenzyTarget

	uhdk.EnableAutoAttacks(uhdk, core.AutoAttackOptions{
		MainHand:       uhdk.WeaponFromMainHand(uhdk.DefaultMeleeCritMultiplier()),
		OffHand:        uhdk.WeaponFromOffHand(uhdk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	return uhdk
}

func (uhdk UnholyDeathKnight) getMasteryShadowBonus() float64 {
	return 0.2 + 0.025*uhdk.GetMasteryPoints()
}

func (uhdk *UnholyDeathKnight) GetDeathKnight() *death_knight.DeathKnight {
	return uhdk.DeathKnight
}

func (uhdk *UnholyDeathKnight) Initialize() {
	uhdk.DeathKnight.Initialize()

	uhdk.registerScourgeStrikeSpell()
}

func (uhdk *UnholyDeathKnight) ApplyTalents() {
	uhdk.DeathKnight.ApplyTalents()

	masteryMod := uhdk.AddDynamicMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		FloatValue: uhdk.getMasteryShadowBonus(),
		School:     core.SpellSchoolShadow,
	})

	uhdk.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		masteryMod.UpdateFloatValue(uhdk.getMasteryShadowBonus())
	})

	core.MakePermanent(uhdk.GetOrRegisterAura(core.Aura{
		Label:    "Dreadblade",
		ActionID: core.ActionID{SpellID: 77515},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Deactivate()
		},
	}))

	// Unholy Might
	uhdk.MultiplyStat(stats.Strength, 1.25)
	core.MakePermanent(uhdk.GetOrRegisterAura(core.Aura{
		Label:    "Unholy Might",
		ActionID: core.ActionID{SpellID: 91107},
	}))

	// Master of Ghouls
}

func (uhdk *UnholyDeathKnight) Reset(sim *core.Simulation) {
	uhdk.DeathKnight.Reset(sim)
}
