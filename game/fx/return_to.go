package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// ReturnToHand returns the card to the players hand instead of the graveyard
func ReturnToHand(card *match.Card, ctx *match.Context) {

	// When destroyed
	if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

		if event.Card == card {

			ctx.InterruptFlow()

			card.Player.MoveCard(card.ID, match.BATTLEZONE, match.HAND, card.ID)
			ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by %s and returned to the hand", event.Card.Name, event.Source.Name))

		}

	}

}

// ReturnToMana returns the card to the players manazone instead of the graveyard
func ReturnToMana(card *match.Card, ctx *match.Context) {

	// When destroyed
	if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

		if event.Card == card {

			ctx.InterruptFlow()

			card.Player.MoveCard(card.ID, match.BATTLEZONE, match.MANAZONE, card.ID)
			card.Tapped = false
			ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by %s and moved to the mana zone", event.Card.Name, event.Source.Name))

		}

	}

}

// ReturnToShield returns the card to the players shield zone instead of the graveyard
func ReturnToShield(card *match.Card, ctx *match.Context) {

	// When destroyed
	if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

		if event.Card == card {

			ctx.InterruptFlow()

			card.Player.MoveCard(card.ID, match.BATTLEZONE, match.SHIELDZONE, card.ID)
			card.Tapped = false
			ctx.Match.Chat("Server", fmt.Sprintf("%s was destroyed by %s and moved to the shield zone", event.Card.Name, event.Source.Name))

		}

	}

}
