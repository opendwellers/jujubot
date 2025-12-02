package bot

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gojp/kana"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/opendwellers/jujubot/pkg/commands"
)

// handleNamedCommands processes commands that start with @botname
func (b *Bot) handleNamedCommands(post *model.Post, replyToId string) bool {
	pattern := globalRegexOptions + "^@" + b.user.Username + " (.*)$"
	matched := regexp.MustCompile(pattern).FindAllStringSubmatch(post.Message, -1)
	if matched == nil {
		return false
	}

	command := matched[0][1]
	is420 := time.Now().Month() == time.April && time.Now().Day() == 20

	// Try each command handler
	handlers := []func(post *model.Post, replyToId, command string, is420 bool) bool{
		b.handleInsultCommands,
		b.handleThanksCommand,
		b.handleQuestionCommand,
		b.handleLoveCommand,
		b.handleChargeCommands,
		b.handleConvertCommand,
		b.handleWeatherCommand,
		b.handleUrbanCommand,
		b.handleJapaneseCommands,
		b.handleDotaCommand,
		b.handleRollCommand,
	}

	for _, handler := range handlers {
		if handler(post, replyToId, command, is420) {
			return true
		}
	}

	// Default response for unrecognized commands
	b.createPost(post.ChannelId, "Kes tu. Veux????", replyToId)
	return true
}

// handleInsultCommands handles insult-type commands
func (b *Bot) handleInsultCommands(post *model.Post, replyToId, command string, _ bool) bool {
	if matched, _ := regexp.MatchString(globalRegexOptions+`^stfu|fuck you|fuck off|ta yeule|tayeule|shut up|shut the fuck up$`, command); matched {
		choices := []string{"no u?", "no u", ":chuckles:", "rolf"}
		b.createReply(post.ChannelId, randomChoice(choices), replyToId, post.UserId)
		return true
	}
	return false
}

// handleThanksCommand handles thank you commands
func (b *Bot) handleThanksCommand(post *model.Post, replyToId, command string, _ bool) bool {
	if matched, _ := regexp.MatchString(globalRegexOptions+`^thanks|merci|ty|thx$`, command); matched {
		choices := []string{"de rien la", "np", "np ;)"}
		b.createReply(post.ChannelId, randomChoice(choices), replyToId, post.UserId)
		return true
	}
	return false
}

// handleQuestionCommand handles question commands
func (b *Bot) handleQuestionCommand(post *model.Post, replyToId, command string, _ bool) bool {
	if matched, _ := regexp.MatchString(globalRegexOptions+`^est-ce qu.*$`, command); matched {
		choices := []string{"maybe", "??", "yess", "no", "rolf oui", "omgggg no"}
		b.createReply(post.ChannelId, randomChoice(choices), replyToId, post.UserId)
		return true
	}
	return false
}

// handleLoveCommand handles love command
func (b *Bot) handleLoveCommand(post *model.Post, replyToId, command string, _ bool) bool {
	if matched, _ := regexp.MatchString(globalRegexOptions+`^I love you$`, command); matched {
		b.createReply(post.ChannelId, "<3", replyToId, post.UserId)
		return true
	}
	return false
}

// handleChargeCommands handles charge-related commands (420 only)
func (b *Bot) handleChargeCommands(post *model.Post, _ string, command string, is420 bool) bool {
	if !is420 {
		return false
	}

	if matched, _ := regexp.MatchString(globalRegexOptions+`^charge up$`, command); matched {
		message := b.chargeUp(post.UserId, 1)
		b.createReply(post.ChannelId, message, post.Id, post.UserId)
		return true
	}

	if matched, _ := regexp.MatchString(globalRegexOptions+`^charge level$`, command); matched {
		message := b.getChargeLevelMessage(post.UserId)
		b.createReply(post.ChannelId, message, post.Id, post.UserId)
		return true
	}

	return false
}

// handleConvertCommand handles currency conversion
func (b *Bot) handleConvertCommand(post *model.Post, _ string, command string, _ bool) bool {
	matched := regexp.MustCompile(globalRegexOptions+`^convert( (\d+)? ?(\w{3}) (?:to )?(\w{3}))?$`).FindAllStringSubmatch(command, -1)
	if matched == nil {
		return false
	}

	// Default to 1 CAD to USD
	from := "CAD"
	to := "USD"
	amount := 1.0

	// If a value was provided
	if matched[0][1] != "" {
		var err error
		amount, err = strconv.ParseFloat(matched[0][2], 64)
		if err != nil {
			b.createReply(post.ChannelId, "Couldn't convert "+matched[0][2]+" to an integer.", post.Id, post.UserId)
			return true
		}
		from = strings.ToUpper(matched[0][3])
		to = strings.ToUpper(matched[0][4])
	}

	amountStr := strconv.FormatFloat(amount, 'f', 2, 64)
	convertedValue, err := commands.Convert(from, to, amount)
	if err != nil {
		b.createReply(post.ChannelId, "Couldn't convert "+amountStr+" "+from+" to "+to+".", post.Id, post.UserId)
		return true
	}

	message := amountStr + " " + from + " = " + strconv.FormatFloat(convertedValue, 'f', 5, 64) + " " + to
	b.createReply(post.ChannelId, message, post.Id, post.UserId)
	return true
}

// handleWeatherCommand handles weather queries
func (b *Bot) handleWeatherCommand(post *model.Post, _ string, command string, _ bool) bool {
	matched := regexp.MustCompile(globalRegexOptions+`^weather ?((now) (.*)|(.*))$`).FindAllStringSubmatch(command, -1)
	if matched == nil {
		return false
	}

	location := ""
	subcommand := strings.ToLower(matched[0][1])

	// Default to Montreal if no location is provided
	if subcommand == "" || (subcommand == "now" && matched[0][2] == "") {
		location = "Montreal"
	} else if subcommand != "" && subcommand != "now" {
		location = subcommand
	} else {
		location = matched[0][2]
	}

	var message string
	var err error

	if subcommand == "now" {
		message, err = b.weatherClient.GetCurrentWeather(location)
	} else {
		message, err = b.weatherClient.GetWeather(location)
	}

	if err != nil {
		b.createReply(post.ChannelId, "Couldn't get weather for "+location+".", post.Id, post.UserId)
		return true
	}

	b.createPost(post.ChannelId, message, post.Id)
	return true
}

// handleUrbanCommand handles Urban Dictionary lookups
func (b *Bot) handleUrbanCommand(post *model.Post, _ string, command string, _ bool) bool {
	matched := regexp.MustCompile(globalRegexOptions+`^urban(?: (.*))?$`).FindAllStringSubmatch(command, -1)
	if matched == nil {
		return false
	}

	word := "huel"
	if matched[0][1] != "" {
		word = matched[0][1]
	}

	result, err := commands.GetUrbanDictionaryDefinition(word)
	if err != nil {
		b.createReply(post.ChannelId, "Couldn't get definition for "+word+".", post.Id, post.UserId)
		return true
	}

	message := fmt.Sprintf("%s\n\n_%s_\n\n**by: %s**\n\n`%d`:+1: `%d`:-1:",
		result.Definition, result.Example, result.Author, result.Upvote, result.Downvote)
	b.createPost(post.ChannelId, message, post.Id)
	return true
}

// handleJapaneseCommands handles Japanese language commands
func (b *Bot) handleJapaneseCommands(post *model.Post, _ string, command string, _ bool) bool {
	// Romaji conversion
	if matched := regexp.MustCompile(globalRegexOptions+`^romaji(?: (.*))?$`).FindAllStringSubmatch(command, -1); matched != nil {
		if matched[0][1] == "" {
			b.createReply(post.ChannelId, "Please provide a word to convert.", post.Id, post.UserId)
			return true
		}
		b.createPost(post.ChannelId, kana.KanaToRomaji(matched[0][1]), post.Id)
		return true
	}

	// Hiragana conversion
	if matched := regexp.MustCompile(globalRegexOptions+`^hiragana(?: (.*))?$`).FindAllStringSubmatch(command, -1); matched != nil {
		if matched[0][1] == "" {
			b.createReply(post.ChannelId, "Please provide a word to convert.", post.Id, post.UserId)
			return true
		}
		b.createPost(post.ChannelId, kana.RomajiToHiragana(matched[0][1]), post.Id)
		return true
	}

	// Katakana conversion
	if matched := regexp.MustCompile(globalRegexOptions+`^katakana(?: (.*))?$`).FindAllStringSubmatch(command, -1); matched != nil {
		if matched[0][1] == "" {
			b.createReply(post.ChannelId, "Please provide a word to convert.", post.Id, post.UserId)
			return true
		}
		b.createPost(post.ChannelId, kana.RomajiToKatakana(matched[0][1]), post.Id)
		return true
	}

	// Word of the day (Japanese)
	if matched := regexp.MustCompile(globalRegexOptions+`^wotd japanese.*$`).FindAllStringSubmatch(command, -1); matched != nil {
		message, err := commands.GetWotdJapanese()
		if err != nil {
			b.createReply(post.ChannelId, "Couldn't get WotD Japanese.", post.Id, post.UserId)
			return true
		}
		b.createPost(post.ChannelId, message, post.Id)
		return true
	}

	return false
}

// handleDotaCommand handles Dota MMR lookups
func (b *Bot) handleDotaCommand(post *model.Post, _ string, command string, _ bool) bool {
	matched := regexp.MustCompile(globalRegexOptions+`^mmr(?: (\d+))?$`).FindAllStringSubmatch(command, -1)
	if matched == nil {
		return false
	}

	playerId := 12088460
	if matched[0][1] != "" {
		var err error
		playerId, err = strconv.Atoi(matched[0][1])
		if err != nil || len(strconv.Itoa(playerId)) > 10 {
			b.createReply(post.ChannelId, fmt.Sprintf("lel nice fake player id: %d.", playerId), post.Id, post.UserId)
			return true
		}
	}

	mmr, err := commands.GetDotaMMR(playerId)
	if err != nil {
		b.createReply(post.ChannelId, fmt.Sprintf("rofl %d existe meme pas zzz", playerId), post.Id, post.UserId)
		return true
	}

	var message string
	if playerId == 12088460 {
		message = fmt.Sprintf("lel j'suis rendu %d ez gaem road to 4k", mmr.SoloCompetitiveRank)
		b.createReply(post.ChannelId, message, post.Id, post.UserId)
		return true
	} else if playerId == 53515020 {
		mmr.SoloCompetitiveRank = 9000
	}

	switch {
	case mmr.SoloCompetitiveRank <= 0:
		message = "unranked pleb or hidden mmr"
	case mmr.SoloCompetitiveRank < 4500:
		message = fmt.Sprintf("lel %s is only %d mmr scrub, git gud", mmr.Profile.Personaname, mmr.SoloCompetitiveRank)
	default:
		message = fmt.Sprintf("lel %s is %d mmr what an amazing player", mmr.Profile.Personaname, mmr.SoloCompetitiveRank)
	}

	b.createReply(post.ChannelId, message, post.Id, post.UserId)
	return true
}

// handleRollCommand handles dice rolling
func (b *Bot) handleRollCommand(post *model.Post, _ string, command string, _ bool) bool {
	matched := regexp.MustCompile(globalRegexOptions+`^roll(?: (\d+|:weed:))?\s*$`).FindAllStringSubmatch(command, -1)
	if matched == nil {
		return false
	}

	requestedRoll := 0
	if matched[0][1] == "" || matched[0][1] == ":weed:" {
		requestedRoll = 420
	} else if matched[0][1] == "dice" {
		requestedRoll = 6
	} else {
		requestedRoll, _ = strconv.Atoi(matched[0][1])
	}

	message := b.rollDice(requestedRoll, post.UserId)
	b.createReply(post.ChannelId, message, post.Id, post.UserId)
	return true
}

// rollDice performs a dice roll with special 420 logic
func (b *Bot) rollDice(dice int, userId string) string {
	roll := rand.Intn(dice) + 1

	if dice == 420 {
		now := time.Now()
		if now.Hour()%12 == 4 && now.Minute() == 20 {
			chargeBonus := b.getCharge(userId)
			if now.Month() == time.April && now.Day() == 20 && chargeBonus != 0 {
				actualRoll := roll
				roll = actualRoll + chargeBonus
				message := strconv.Itoa(actualRoll) + " + " + strconv.Itoa(chargeBonus) + " charge bonus = " + strconv.Itoa(roll) + " "

				if roll == 420 {
					return message + "BIG WINNER WOW :musk: :weed:"
				} else if roll == 69 {
					return message + "_Nice._ :smugpepe:"
				}
				return message + ":chuckles:"
			}

			if roll == 420 {
				return strconv.Itoa(roll) + " BIG WINNER WOW :musk: :weed:"
			} else if roll == 69 {
				return strconv.Itoa(roll) + " _Nice._ :smugpepe:"
			}
			return strconv.Itoa(roll) + " :chuckles:"
		}
		return "Spa leur smh"
	}

	if dice == 1 {
		return ":99:"
	}

	return strconv.Itoa(roll)
}
