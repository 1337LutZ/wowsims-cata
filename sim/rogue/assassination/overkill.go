package assassination

import (
	"time"

	"github.com/wowsims/cata/sim/core"
)

func (sinRogue *AssassinationRogue) registerOverkill() {
	if !sinRogue.Talents.Overkill {
		return
	}

	sinRogue.OverkillAura = sinRogue.RegisterAura(core.Aura{
		Label:    "Overkill",
		ActionID: core.ActionID{SpellID: 58427},
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			sinRogue.ApplyEnergyTickMultiplier(0.3)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			sinRogue.ApplyEnergyTickMultiplier(-0.3)
		},
	})
}
