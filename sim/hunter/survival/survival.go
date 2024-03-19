package survival

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/hunter"
)

func RegisterSurvivalHunter() {
	core.RegisterAgentFactory(
		proto.Player_SurvivalHunter{},
		proto.Spec_SpecSurvivalHunter,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewSurvivalHunter(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_SurvivalHunter)
			if !ok {
				panic("Invalid spec value for Survival Hunter!")
			}
			player.Spec = playerSpec
		},
	)
}

func (hunter *SurvivalHunter) Initialize() {
	// Initialize global Hunter spells
	hunter.Hunter.Initialize()

	hunter.registerExplosiveShotSpell()
	hunter.registerBlackArrowSpell(hunter.FireTrapTimer)
}
func (hunter *SurvivalHunter) getMasteryBonus(mastery float64) float64 {
	return 1.08 + (mastery * 0.01)
}

func NewSurvivalHunter(character *core.Character, options *proto.Player) *SurvivalHunter {
	survivalOptions := options.GetSurvivalHunter().Options

	svHunter := &SurvivalHunter{
		Hunter: hunter.NewHunter(character, options, survivalOptions.ClassOptions),
	}
	schoolsAffectedBySurvivalMastery := []stats.SchoolIndex{
		stats.SchoolIndexNature,
		stats.SchoolIndexFire,
		stats.SchoolIndexArcane,
		stats.SchoolIndexFrost,
	}
	svHunter.SurvivalOptions = survivalOptions
	// Todo: Is there a better way to do this?
	baseMastery := svHunter.GetStat(stats.Mastery)
	for _, school := range schoolsAffectedBySurvivalMastery {
		svHunter.PseudoStats.SchoolDamageDealtMultiplier[school] *= svHunter.getMasteryBonus(baseMastery)
	}

	svHunter.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery float64, newMastery float64) {
		for _, school := range schoolsAffectedBySurvivalMastery {
			svHunter.PseudoStats.SchoolDamageDealtMultiplier[school] /= svHunter.getMasteryBonus(oldMastery)
			svHunter.PseudoStats.SchoolDamageDealtMultiplier[school] *= svHunter.getMasteryBonus(newMastery)
		}
	})
	return svHunter
}

type SurvivalHunter struct {
	*hunter.Hunter
}

func (svHunter *SurvivalHunter) GetHunter() *hunter.Hunter {
	return svHunter.Hunter
}

func (svHunter *SurvivalHunter) Reset(sim *core.Simulation) {
	svHunter.Hunter.Reset(sim)
}
