package bot

import (
	"math/rand"
	"strconv"
)

// chargeUp adds or removes charge points for a user
func (b *Bot) chargeUp(userId string, multiplier int) string {
	chargeValue := (rand.Intn(5) - 1) * multiplier
	b.chargeMap[userId] += chargeValue

	switch {
	case chargeValue < 0:
		return "You lost " + strconv.Itoa(chargeValue*-1) + " charge points :lamo:"
	case chargeValue > 0:
		return "You gained " + strconv.Itoa(chargeValue) + " charge points :hype:"
	default:
		return "You gained " + strconv.Itoa(chargeValue) + " charge points :pepehands:"
	}
}

// getCharge returns the current charge for a user
func (b *Bot) getCharge(userId string) int {
	if val, ok := b.chargeMap[userId]; ok {
		return val
	}
	return 0
}

// getChargeLevelMessage returns a message describing the user's charge level
func (b *Bot) getChargeLevelMessage(userId string) string {
	chargeValue := b.getCharge(userId)
	chargeValueStr := strconv.Itoa(chargeValue)

	switch {
	case chargeValue < 0:
		return "You have " + chargeValueStr + " points charged up. :fuck:"
	case chargeValue == 0:
		return "You have no charge points stored up! :tensepepe:"
	case chargeValue == 69:
		return ":smugpepe:"
	case chargeValue > 0 && chargeValue < 20:
		return "You have " + chargeValueStr + " points charged up. :pogchamp:"
	case chargeValue >= 20 && chargeValue < 100:
		return "You have " + chargeValueStr + " points charged up! :pog:"
	default: // chargeValue >= 100
		return "You have :pogchampignon: points charged up‽"
	}
}
