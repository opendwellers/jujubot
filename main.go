package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gojp/kana"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/opendwellers/jujubot/pkg/commands"
	"github.com/opendwellers/jujubot/pkg/config"
	"go.uber.org/zap"
)

const globalRegexOptions = "(?i)"

var (
	chargeMap map[string]int = map[string]int{}

	client          *model.Client4
	webSocketClient *model.WebSocketClient

	botUser          *model.User
	botTeam          *model.Team
	debuggingChannel *model.Channel

	weatherClient commands.Weather

	logger *zap.Logger
)

func initLogger() {
	if os.Getenv("IS_DEBUG") == "true" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	zap.ReplaceGlobals(logger)
}

// Documentation for the Go driver can be found
// at https://godoc.org/github.com/mattermost/platform/model#Client
func main() {
	// Init logger and config
	initLogger()
	defer logger.Sync()
	logger.Sugar().Info("Loading configuration")
	config, error := config.LoadConfig()
	if error != nil {
		logger.Sugar().Fatal("Failed to load configuration")
	}
	logger.Sugar().Info("Configuration loaded")

	SetupGracefulShutdown()

	logger.Sugar().Info("Connecting to Mattermost at " + config.ServerURL)
	client = model.NewAPIv4Client(config.ServerURL)

	// Lets test to see if the mattermost server is up and running
	MakeSureServerIsRunning()

	// lets attempt to login to the Mattermost server as the bot user
	// This will set the token required for all future calls
	// You can get this token with client.AuthToken
	LoginAsTheBotUser(config.AuthToken)

	// Lets find our bot team
	FindBotTeam(config.TeamName)

	// This is an important step.  Lets make sure we use the botTeam
	// for all future web service requests that require a team.
	//client.SetTeamId(botTeam.Id)

	// Lets create a bot channel for logging debug messages into
	CreateBotDebuggingChannelIfNeeded(config.ChannelLogName)

	// Initialize weather client
	var clientErr = error
	weatherClient, clientErr = commands.NewWeatherClient(config.OpenWeatherApiKey)
	if clientErr != nil {
		logger.Sugar().Fatal("Failed to create weather client")
	}

	// Create file at /tmp/ready
	// This is used by the readiness probe to determine if the bot is running
	f, er := os.Create("/tmp/ready")
	if er != nil {
		logger.Sugar().Fatal("Failed to create /tmp/ready")
	}
	f.Close()

	logger.Sugar().Info("Bot is now running and listening to messages.")

	go func() {
		for {
			// Lets start listening to some channels via the websocket!
			var err = error
			webSocketClient, err = model.NewWebSocketClient4(config.ServerWSURL, client.AuthToken)
			if err != nil {
				logger.Sugar().Error("Failed to connect to the web socket", zap.Any("error", err))
				// Sleep for a bit before trying to reconnect
				time.Sleep(5 * time.Second)
				continue
			}
			listen()
		}
	}()

	// You can block forever with
	select {}
}

func listen() {
	webSocketClient.Listen()
	defer webSocketClient.Close()
	for {
		if webSocketClient.ListenError != nil {
			logger.Sugar().Error("Failed to listen to the web socket: ", zap.Any("error", webSocketClient.ListenError))
			logger.Sugar().Info("Reconnecting to the web socket")
			return
		}

		event := <-webSocketClient.EventChannel
		if event == nil {
			continue
		}
		HandleWebSocketResponse(event)
	}
}

func MakeSureServerIsRunning() {
	if props, _, err := client.GetOldClientConfig(context.TODO(), ""); err != nil {
		logger.Sugar().Fatal("There was a problem pinging the Mattermost server.  Are you sure it's running?", zap.Any("error", err))
	} else {
		logger.Sugar().Info("Server detected and is running version " + props["Version"])
	}
}

func LoginAsTheBotUser(token string) {
	client.SetToken(token)
	var user *model.User
	var err error
	if user, _, err = client.GetMe(context.TODO(), ""); err != nil {
		logger.Sugar().Fatal("There was a problem getting the user", zap.Any("error", err))
	}
	botUser = user
	logger.Sugar().Info("Running as " + user.Username)
}

func FindBotTeam(teamName string) {
	if team, _, err := client.GetTeamByName(context.TODO(), teamName, ""); err != nil {
		logger.Sugar().Fatal("Failed to get the initial load or we do not appear to be a member of the team '"+teamName+"'", zap.Any("error", err))
	} else {
		logger.Sugar().Info("Found team " + team.Name)
		botTeam = team
	}
}

func CreateBotDebuggingChannelIfNeeded(channelName string) {
	if rchannel, _, err := client.GetChannelByName(context.TODO(), channelName, botTeam.Id, ""); err != nil {
		logger.Sugar().Error("Failed to get the channels", zap.Any("error", err))
	} else {
		debuggingChannel = rchannel
		return
	}

	// Looks like we need to create the logging channel
	channel := &model.Channel{}
	channel.Name = channelName
	channel.DisplayName = "Debugging For Sample Bot"
	channel.Purpose = "This is used as a test channel for logging bot debug messages"
	channel.Type = model.ChannelTypeOpen
	channel.TeamId = botTeam.Id
	if rchannel, _, err := client.CreateChannel(context.TODO(), channel); err != nil {
		logger.Sugar().Error("Failed to create the channel "+channelName, zap.Any("error", err))
	} else {
		debuggingChannel = rchannel
		logger.Sugar().Info("Looks like this might be the first run so we've created the channel " + channelName)
	}
}

func CreateReply(channelId string, msg string, replyToId string, replyToUserId string) {
	CreatePost(channelId, getUserMention(replyToUserId)+": "+msg, replyToId)
}

func CreatePost(channelId string, msg string, replyToId string) {
	post := &model.Post{}
	post.ChannelId = channelId
	post.Message = msg

	post.RootId = replyToId

	if _, _, err := client.CreatePost(context.TODO(), post); err != nil {
		logger.Sugar().Error("Failed to send a message to the channel", zap.Any("error", err))
	}
}

func CreateReaction(emojiName string, postId string) {
	reaction := &model.Reaction{UserId: botUser.Id, PostId: postId, EmojiName: emojiName}
	if _, _, err := client.SaveReaction(context.TODO(), reaction); err != nil {
		logger.Sugar().Error("Failed to add a reaction to the post", zap.Any("error", err))
	}
}

func HandleWebSocketResponse(event *model.WebSocketEvent) {
	HandleMessage(event)
}

func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

func getUserMention(userId string) string {
	var user *model.User
	var err error
	if user, _, err = client.GetUser(context.TODO(), userId, ""); err != nil {
		logger.Sugar().Error("Failed to get user", zap.Any("error", err))
	}

	return "@" + user.Username
}

func chargeUp(userId string, multiplier int) string {
	chargeValue := (rand.Intn(5) - 1) * multiplier
	chargeMap[userId] += chargeValue
	message := ""
	switch {
	case chargeValue < 0:
		message = "You lost " + strconv.Itoa(chargeValue*-1) + " charge points :lamo:"
	case chargeValue > 0:
		message = "You gained " + strconv.Itoa(chargeValue) + " charge points :hype:"
	case chargeValue == 0:
		message = "You gained " + strconv.Itoa(chargeValue) + " charge points :pepehands:"
	}
	return message
}

func getCharge(userId string) int {
	if val, ok := chargeMap[userId]; ok {
		return val
	} else {
		return 0
	}
}

func HandleMessage(event *model.WebSocketEvent) {
	logger.Sugar().Debug("Got event: ", event)
	// If this isn't the debugging channel then lets ignore it
	// if event.Broadcast.ChannelId != debuggingChannel.Id {
	// 	return
	// }

	// Lets only reply to messages posted events
	if event.EventType() != model.WebsocketEventPosted {
		return
	}
	var post *model.Post
	json.NewDecoder(strings.NewReader(event.GetData()["post"].(string))).Decode(&post)

	if post != nil {
		// ignore my events
		if post.UserId == botUser.Id {
			return
		}

		user, _, err := client.GetUser(context.TODO(), post.UserId, "")
		if err != nil {
			logger.Sugar().Info("Failed to get user ", post.UserId, zap.Any("error", err))
		}
		logger.Sugar().Info("Processing message from user ", user.Username, ": ", post.Message)

		replyToId := ""
		if post.RootId != "" {
			replyToId = post.RootId
		}

		// if you see any word matching 'hello' then respond
		if matched, _ := regexp.MatchString(globalRegexOptions+`\bsalut|allo\b`, post.Message); matched {
			choices := []string{"aaaaaaayyeee", "sup", "yo"}
			CreateReply(post.ChannelId, randomChoice(choices), replyToId, post.UserId)
			return
		}

		if matched := regexp.MustCompile(globalRegexOptions+`(xd+)`).FindAllStringSubmatch(post.Message, -1); matched != nil {
			CreatePost(post.ChannelId, "haha "+matched[0][1], replyToId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`\banime|animuh|weeb|weaboo\b`, post.Message); matched {
			CreatePost(post.ChannelId, "### Disgusting weebs rolf :huel:", replyToId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`\bvidya|bonshommes\b`, post.Message); matched {
			CreatePost(post.ChannelId, "rolf vous avez quel age?", replyToId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`\bvelo.*hiver\b`, post.Message); matched {
			CreatePost(post.ChannelId, "wow cest fukin dangereux faut vraiment etre retarded pour cycler en hiver (dans une tempete de verglas) :huel:", replyToId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`\b:disappear:|peace|alp|bye|:wave:|see ya|au revoir|ciao|chow|a tantot\b`, post.Message); matched {
			CreatePost(post.ChannelId, "hey salut la, a prochaine, on se revoit, stait bin lfun", replyToId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`\bbon matin|morning|mornin\b`, post.Message); matched {
			choices := []string{"zzzz kill me now", "omgggggg"}
			CreatePost(post.ChannelId, randomChoice(choices), replyToId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`\bmirin\b`, post.Message); matched {
			CreateReply(post.ChannelId, "fucking mirin", replyToId, post.UserId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`;-?\)(\s|$)|:wink:`, post.Message); matched {
			CreateReaction("wink", post.Id)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`:-?P(\s|$)|:stuck_out_tongue:`, post.Message); matched {
			CreateReaction("stuck_out_tongue", post.Id)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`:fuck:`, post.Message); matched {
			CreateReaction("fuck", post.Id)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`(\s|^)\^(\s|$)`, post.Message); matched {
			CreatePost(post.ChannelId, "^", replyToId)
			CreateReaction("point_up_2", post.Id)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`^this$`, post.Message); matched {
			CreatePost(post.ChannelId, "this", replyToId)
			CreateReaction("point_up_2", post.Id)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`\breddit\b`, post.Message); matched {
			CreateReply(post.ChannelId, "\\>reddit", replyToId, post.UserId)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`\btumblr\b`, post.Message); matched {
			CreateReply(post.ChannelId, "\\>tumblr", replyToId, post.UserId)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`\btgif\b`, post.Message); matched {
			CreateReply(post.ChannelId, "tgiff*", replyToId, post.UserId)
			return
		}
		// charging up
		if matched := regexp.MustCompile(globalRegexOptions+`(a{5,}h{2,}!*)|:charging_up:`).FindAllStringSubmatch(post.Message, -1); matched != nil {
			match := matched[0][0]
			length := len(match)
			if match == ":charging_up:" {
				length = rand.Intn(50) + 1
			}
			// x1 at 8 characters
			// +1 multiplier every time you add 15 characters
			multiplier := (length-8)/15 + 1
			message := chargeUp(post.UserId, multiplier)
			CreateReply(post.ChannelId, message, replyToId, post.UserId)
			return
		}

		// Named commands
		if matched := regexp.MustCompile(globalRegexOptions+"^@"+botUser.Username+" (.*)$").FindAllStringSubmatch(post.Message, -1); matched != nil {
			is420 := time.Now().Month() == time.April && time.Now().Day() == 20
			command := matched[0][1]

			if matched, _ := regexp.MatchString(globalRegexOptions+`^stfu|fuck you|fuck off|ta yeule|tayeule|shut up|shut the fuck up$`, command); matched {
				choices := []string{"no u?", "no u", ":chuckles:", "rolf"}
				CreateReply(post.ChannelId, randomChoice(choices), replyToId, post.UserId)
				return
			}
			if matched, _ := regexp.MatchString(globalRegexOptions+`^thanks|merci|ty|thx$`, command); matched {
				choices := []string{"de rien la", "np", "np ;)"}
				CreateReply(post.ChannelId, randomChoice(choices), replyToId, post.UserId)
				return
			}
			if matched, _ := regexp.MatchString(globalRegexOptions+`^est-ce qu.*$`, command); matched {
				choices := []string{"maybe", "??", "yess", "no", "rolf oui", "omgggg no"}
				CreateReply(post.ChannelId, randomChoice(choices), replyToId, post.UserId)
				return
			}
			if matched, _ := regexp.MatchString(globalRegexOptions+`^I love you$`, command); matched {
				CreateReply(post.ChannelId, "<3", replyToId, post.UserId)
				return
			}

			if matched, _ := regexp.MatchString(globalRegexOptions+`^charge up$`, command); matched && is420 {
				message := chargeUp(post.UserId, 1)
				CreateReply(post.ChannelId, message, post.Id, post.UserId)
				return
			}

			if matched, _ := regexp.MatchString(globalRegexOptions+"^charge level$", command); matched && is420 {
				chargeValue := getCharge(post.UserId)
				chargeValueStr := strconv.Itoa(chargeValue)
				message := ""
				switch {
				case chargeValue < 0:
					message = "You have " + chargeValueStr + " points charged up. :fuck:"
				case chargeValue == 0:
					message = "You have no charge points stored up! :tensepepe:"
				case chargeValue > 0 && chargeValue < 20:
					message = "You have " + chargeValueStr + " points charged up. :pogchamp:"
				case chargeValue == 69:
					message = ":smugpepe:"
				case chargeValue >= 20 && chargeValue < 100:
					message = "You have " + chargeValueStr + " points charged up! :pog:"
				case chargeValue >= 100:
					message = "You have :pogchampignon: points charged up‽"
				}
				CreateReply(post.ChannelId, message, post.Id, post.UserId)
				return
			}

			// Currency conversion
			if matched := regexp.MustCompile(globalRegexOptions+`^convert( (\d+)? ?(\w{3}) (?:to )?(\w{3}))?$`).FindAllStringSubmatch(command, -1); matched != nil {
				// Default to 1 CAD to USD
				from := "CAD"
				to := "USD"
				amount := 1.0
				message := ""
				var err error

				// If a value was provided
				if matched[0][1] != "" {
					amount, err = strconv.ParseFloat(matched[0][2], 64)
					if err != nil {
						message = "Couldn't convert " + matched[0][2] + " to an integer."
						CreateReply(post.ChannelId, message, post.Id, post.UserId)
						return
					}
					from = strings.ToUpper(matched[0][3])
					to = strings.ToUpper(matched[0][4])
				}

				amountStr := strconv.FormatFloat(amount, 'f', 2, 64)
				convertedValue, err := commands.Convert(from, to, amount)
				if err != nil {
					message = "Couldn't convert " + amountStr + " " + from + " to " + to + "."
					CreateReply(post.ChannelId, message, post.Id, post.UserId)
					return
				}
				message = amountStr + " " + from + " = " + strconv.FormatFloat(convertedValue, 'f', 5, 64) + " " + to

				CreateReply(post.ChannelId, message, post.Id, post.UserId)
				return
			}

			// Weather
			if matched := regexp.MustCompile(globalRegexOptions+`^weather ?((now) (.*)|(.*))$`).FindAllStringSubmatch(command, -1); matched != nil {
				message := ""
				location := ""
				var err error
				subcommand := strings.ToLower(matched[0][1])

				// Default to Montreal if no location is provided
				if subcommand == "" || (subcommand == "now" && matched[0][2] == "") {
					location = "Montreal"
				} else if subcommand != "" && subcommand != "now" {
					// If a location was provided
					location = subcommand
				} else {
					// If a location was provided for the now subcommand
					location = matched[0][2]
				}

				if subcommand == "now" {
					message, err = weatherClient.GetCurrentWeather(location)
				} else {
					message, err = weatherClient.GetWeather(location)
				}

				if err != nil {
					CreateReply(post.ChannelId, "Couldn't get weather for "+location+".", post.Id, post.UserId)
					return
				}
				CreatePost(post.ChannelId, message, post.Id)
				return
			}

			// Urban Dictionary
			if matched := regexp.MustCompile(globalRegexOptions+`^urban(?: (.*))?$`).FindAllStringSubmatch(command, -1); matched != nil {
				word := "huel"
				if matched[0][1] != "" {
					word = matched[0][1]
				}
				result, err := commands.GetUrbanDictionaryDefinition(word)
				if err != nil {
					CreateReply(post.ChannelId, "Couldn't get definition for "+word+".", post.Id, post.UserId)
					return
				}
				message := fmt.Sprintf("%s\n\n_%s_\n\n**by: %s**\n\n`%d`:+1: `%d`:-1:", result.Definition, result.Example, result.Author, result.Upvote, result.Downvote)
				CreatePost(post.ChannelId, message, post.Id)
				return
			}

			// romaji
			if matched := regexp.MustCompile(globalRegexOptions+`^romaji(?: (.*))?$`).FindAllStringSubmatch(command, -1); matched != nil {
				word := matched[0][1]
				if matched[0][1] == "" {
					CreateReply(post.ChannelId, "Please provide a word to convert.", post.Id, post.UserId)
					return
				}
				message := kana.KanaToRomaji(word)
				CreatePost(post.ChannelId, message, post.Id)
				return
			}

			// hiragana
			if matched := regexp.MustCompile(globalRegexOptions+`^hiragana(?: (.*))?$`).FindAllStringSubmatch(command, -1); matched != nil {
				word := matched[0][1]
				if matched[0][1] == "" {
					CreateReply(post.ChannelId, "Please provide a word to convert.", post.Id, post.UserId)
					return
				}
				message := kana.RomajiToHiragana(word)
				CreatePost(post.ChannelId, message, post.Id)
				return
			}

			// katakana
			if matched := regexp.MustCompile(globalRegexOptions+`^katakana(?: (.*))?$`).FindAllStringSubmatch(command, -1); matched != nil {
				word := matched[0][1]
				if matched[0][1] == "" {
					CreateReply(post.ChannelId, "Please provide a word to convert.", post.Id, post.UserId)
					return
				}
				message := kana.RomajiToKatakana(word)
				CreatePost(post.ChannelId, message, post.Id)
				return
			}

			// wotd japanese
			if matched := regexp.MustCompile(globalRegexOptions+`^wotd japanese.*$`).FindAllStringSubmatch(command, -1); matched != nil {
				message, err := commands.GetWotdJapanese()
				if err != nil {
					CreateReply(post.ChannelId, "Couldn't get WotD Japanese.", post.Id, post.UserId)
					return
				}
				CreatePost(post.ChannelId, message, post.Id)
				return
			}

			// Dota MMR
			if matched := regexp.MustCompile(globalRegexOptions+`^mmr(?: (\d+))?$`).FindAllStringSubmatch(command, -1); matched != nil {
				playerId := 12088460
				message := ""
				if matched[0][1] != "" {
					var err error
					playerId, err = strconv.Atoi(matched[0][1])
					if err != nil || len(strconv.Itoa(playerId)) > 10 {
						CreateReply(post.ChannelId, fmt.Sprintf("lel nice fake player id: %d.", playerId), post.Id, post.UserId)
						return
					}
				}

				mmr, err := commands.GetDotaMMR(playerId)
				if err != nil {
					CreateReply(post.ChannelId, fmt.Sprintf("rofl %d existe meme pas zzz", playerId), post.Id, post.UserId)
					return
				}
				if playerId == 12088460 {
					message = fmt.Sprintf("lel j'suis rendu %d ez gaem road to 4k", mmr.SoloCompetitiveRank)
					CreateReply(post.ChannelId, message, post.Id, post.UserId)
					return
				} else if playerId == 53515020 {
					mmr.SoloCompetitiveRank = 9000
				}

				switch {
				case mmr.SoloCompetitiveRank <= 0:
					message = "unranked pleb or hidden mmr"
				case mmr.SoloCompetitiveRank < 4500:
					message = fmt.Sprintf("lel %s is only %d mmr scrub, git gud", mmr.Profile.Personaname, mmr.SoloCompetitiveRank)
				case mmr.SoloCompetitiveRank >= 4500:
					message = fmt.Sprintf("lel %s is %d mmr what an amazing player", mmr.Profile.Personaname, mmr.SoloCompetitiveRank)
				}

				CreateReply(post.ChannelId, message, post.Id, post.UserId)
				return
			}

			// Rolls
			if matched := regexp.MustCompile(globalRegexOptions+`^roll(?: (\d+|:weed:))?\s*$`).FindAllStringSubmatch(command, -1); matched != nil {
				requestedRoll := 0
				if matched[0][1] == "" || matched[0][1] == ":weed:" {
					requestedRoll = 420
				} else if matched[0][1] == "dice" {
					requestedRoll = 6
				} else {
					requestedRoll, _ = strconv.Atoi(matched[0][1])
				}
				message := rollDice(requestedRoll, post.UserId)
				CreateReply(post.ChannelId, message, post.Id, post.UserId)
				return
			}

			CreatePost(post.ChannelId, "Kes tu. Veux????", replyToId)
		}
	}
}

func rollDice(dice int, userId string) (message string) {
	roll := rand.Intn(dice) + 1
	if dice == 420 {
		now := time.Now()
		if now.Hour()%12 == 4 && now.Minute() == 20 {
			chargeBonus := getCharge(userId)
			if now.Month() == time.April && now.Day() == 20 && chargeBonus != 0 {
				actualRoll := roll
				roll := actualRoll + chargeBonus
				message = strconv.Itoa(actualRoll) + " + " + strconv.Itoa(chargeBonus) + " charge bonus = " + strconv.Itoa(roll) + " "
			} else {
				message = strconv.Itoa(roll) + " "
			}

			if roll == 420 {
				message += "BIG WINNER WOW :musk: :weed:"
			} else if roll == 69 {
				message += "_Nice._ :smugpepe:"
			} else {
				message += ":chuckles:"
			}
		} else {
			message = "Spa leur smh"
		}
	} else if dice == 1 {
		message = ":99:"
	} else {
		message = strconv.Itoa(roll)
	}

	return message
}

func SetupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if webSocketClient != nil {
				webSocketClient.Close()
			}

			os.Exit(0)
		}
	}()
}
