package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GulanRiasSpeedGuardian ...
func GulanRiasSpeedGuardian(c *match.Card) {

	c.Name = "Gulan Rias, Speed Guardian"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(
		fx.Creature,
		fx.CantBeBlockedIf(func(blocker *match.Card) bool {
			return blocker.Civ == civ.Darkness
		}),
		fx.CantBeAttackedIf(func(attacker *match.Card) bool {
			return attacker.Civ == civ.Darkness
		}),
	)
}

// MistRiasSonicGuardian ...
func MistRiasSonicGuardian(c *match.Card) {

	c.Name = "Mist Rias, Sonic Guardian"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			if card.Zone != match.BATTLEZONE {

				exit()
				return
			}

			if event, ok := ctx2.Event.(*match.CardMoved); ok && event.To == match.BATTLEZONE && event.CardID != card.ID {

				card.Player.DrawCards(1)

			}

		})

	}))

}
