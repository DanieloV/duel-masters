package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// DomeShell ...
func DomeShell(c *match.Card) {

	c.Name = "Dome Shell"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}

// StormShell ...
func StormShell(c *match.Card) {

	c.Name = "Storm Shell"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		opponent := ctx.Match.Opponent(card.Player)

		battlezone, err := opponent.Container(match.BATTLEZONE)

		if err != nil {
			return
		}

		if len(battlezone) < 1 {
			return
		}

		ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")

		ctx.Match.NewAction(opponent, battlezone, 1, 1, "Storm Shell: Select 1 card from your battlezone that will be sent to your manazone", false)

		defer func() {
			ctx.Match.EndWait(card.Player)
			ctx.Match.CloseAction(opponent)
		}()

		for {

			action := <-opponent.Action

			if len(action.Cards) != 1 || !match.AssertCardsIn(battlezone, action.Cards...) {
				ctx.Match.ActionWarning(opponent, "Your selection of cards does not fulfill the requirements")
				continue
			}

			movedCard, err := opponent.MoveCard(action.Cards[0], match.BATTLEZONE, match.MANAZONE, card.ID)

			if err != nil {
				break
			}

			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved from %s's battlezone to their manazone", movedCard.Name, opponent.Username()))

			break

		}

	}))

}

// TowerShell ...
func TowerShell(c *match.Card) {

	c.Name = "Tower Shell"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.BlockerSelectionStep, fx.CantBeBlockedByPowerUpTo4000))

}
