package assassination

import (
	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
	"github.com/wowsims/cata/sim/rogue"
)

const masteryDamagePerPercent = .035
const masteryBaseEffect = .28

func RegisterAssassinationRogue() {
	core.RegisterAgentFactory(
		proto.Player_AssassinationRogue{},
		proto.Spec_SpecAssassinationRogue,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewAssassinationRogue(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_AssassinationRogue)
			if !ok {
				panic("Invalid spec value for Assassination Rogue!")
			}
			player.Spec = playerSpec
		},
	)
}

func (sinRogue *AssassinationRogue) Initialize() {
	sinRogue.Rogue.Initialize()

	sinRogue.registerMutilateSpell()
	sinRogue.registerOverkill()
	sinRogue.registerColdBloodCD()
	sinRogue.applySealFate()
	sinRogue.registerVenomousWounds()
	sinRogue.registerVendetta()

	// Apply Mastery
	masteryPercent := sinRogue.GetStat(stats.Mastery) / core.MasteryRatingPerMasteryPercent
	masteryEffect := masteryBaseEffect + masteryPercent*masteryDamagePerPercent
	for _, spell := range sinRogue.InstantPoison {
		spell.DamageMultiplier += masteryEffect
	}
	for _, spell := range sinRogue.WoundPoison {
		spell.DamageMultiplier += masteryEffect
	}
	sinRogue.DeadlyPoison.DamageMultiplier += masteryEffect
	sinRogue.Envenom.DamageMultiplier += masteryEffect
	sinRogue.VenomousWounds.DamageMultiplier += masteryEffect

	sinRogue.AddOnMasteryStatChanged(func(sim *core.Simulation, oldMastery, newMastery float64) {
		masteryPercentOld := oldMastery / core.MasteryRatingPerMasteryPercent
		masteryPercentNew := newMastery / core.MasteryRatingPerMasteryPercent
		masteryEffectChange := masteryBaseEffect + (masteryPercentNew-masteryPercentOld)*masteryDamagePerPercent
		for _, spell := range sinRogue.InstantPoison {
			spell.DamageMultiplier += masteryEffectChange
		}
		for _, spell := range sinRogue.WoundPoison {
			spell.DamageMultiplier += masteryEffectChange
		}
		sinRogue.DeadlyPoison.DamageMultiplier += masteryEffectChange
		sinRogue.Envenom.DamageMultiplier += masteryEffectChange
		sinRogue.VenomousWounds.DamageMultiplier += masteryEffectChange
	})
}

func NewAssassinationRogue(character *core.Character, options *proto.Player) *AssassinationRogue {
	sinOptions := options.GetAssassinationRogue().Options

	sinRogue := &AssassinationRogue{
		Rogue: rogue.NewRogue(character, sinOptions.ClassOptions, options.TalentsString),
	}
	sinRogue.AssassinationOptions = sinOptions

	return sinRogue
}

type AssassinationRogue struct {
	*rogue.Rogue
}

func (sinRogue *AssassinationRogue) GetRogue() *rogue.Rogue {
	return sinRogue.Rogue
}

func (sinRogue *AssassinationRogue) Reset(sim *core.Simulation) {
	sinRogue.Rogue.Reset(sim)
}
