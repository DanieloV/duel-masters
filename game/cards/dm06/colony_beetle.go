package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func CarrierShell(c *match.Card) {

	c.Name = "Carrier Shell"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker3000)
}

func SlumberShell(c *match.Card) {

	c.Name = "Slumber Shell"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature)
}

func FactoryShellQ(c *match.Card) {

	c.Name = "Factory Shell Q"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle, family.Survivor}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Survivor, func(card *match.Card, ctx *match.Context) {

		if !ctx.Match.IsPlayerTurn(card.Player) || card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CardMoved); ok && event.To == match.BATTLEZONE {

			creature, err := card.Player.GetCard(event.CardID, match.BATTLEZONE)

			if err != nil {
				return
			}

			if !creature.HasFamily(family.Survivor) {
				return
			}

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.DECK,
				"Select 1 survivor from your deck that will be shown to your opponent and sent to your hand",
				1,
				1,
				true,
				func(x *match.Card) bool { return x.HasCondition(cnd.Survivor) },
			).Map(func(x *match.Card) {
				card.Player.MoveCard(x.ID, match.DECK, match.HAND, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's deck to their hand", x.Name, card.Player.Username()))
				card.Player.ShuffleDeck()
			})

		}

	})

}
