package destruction

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/warlock"
)

func RegisterDestructionWarlock() {
	core.RegisterAgentFactory(
		proto.Player_DestructionWarlock{},
		proto.Spec_SpecDestructionWarlock,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewDestructionWarlock(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_DestructionWarlock)
			if !ok {
				panic("Invalid spec value for Destruction Warlock!")
			}
			player.Spec = playerSpec
		},
	)
}

func NewDestructionWarlock(character *core.Character, options *proto.Player) *DestructionWarlock {
	destroOptions := options.GetDestructionWarlock().Options
	destruction := &DestructionWarlock{
		Warlock: warlock.NewWarlock(character, options, destroOptions.ClassOptions),
	}

	return destruction
}

type DestructionWarlock struct {
	*warlock.Warlock
}

func (destruction DestructionWarlock) getMasteryBonus() float64 {
	return 0.11 + 0.0135*destruction.GetMasteryPoints()
}

func (destruction *DestructionWarlock) GetWarlock() *warlock.Warlock {
	return destruction.Warlock
}

func (destruction *DestructionWarlock) Initialize() {
	destruction.Warlock.Initialize()

	destruction.registerConflagrateSpell()
}

func (destruction *DestructionWarlock) ApplyTalents() {
	destruction.Warlock.ApplyTalents()

	// Mastery: Fiery Apocalypse
	masteryMod := destruction.AddDynamicMod(core.SpellModConfig{
		Kind:      core.SpellMod_DamageDone_Pct,
		ClassMask: warlock.WarlockFireDamage,
	})

	destruction.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		masteryMod.UpdateFloatValue(destruction.getMasteryBonus())
	})

	core.MakePermanent(destruction.GetOrRegisterAura(core.Aura{
		Label:    "Mastery: Fiery Apocalypse",
		ActionID: core.ActionID{SpellID: 77220},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.UpdateFloatValue(destruction.getMasteryBonus())
			masteryMod.Activate()
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			masteryMod.Deactivate()
		},
	}))

	// Cataclysm
	destruction.AddStaticMod(core.SpellModConfig{
		Kind:       core.SpellMod_DamageDone_Pct,
		ClassMask:  warlock.WarlockFireDamage,
		FloatValue: 0.25,
	})
}

func (destruction *DestructionWarlock) Reset(sim *core.Simulation) {
	destruction.Warlock.Reset(sim)
}
