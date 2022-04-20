package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gojp/kana"
	"github.com/mattermost/mattermost-server/v5/model"
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
	logger.Info("loading configuration file")
	config, error := config.LoadConfig()
	if error != nil {
		logger.Error("failed to load configuration file")
		os.Exit(1)
	}
	logger.Sugar().Infow("configuration file loaded", config)

	SetupGracefulShutdown()

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

	// Initialise weather client
	var clientErr = error
	weatherClient, clientErr = commands.NewWeatherClient(config.OpenWeatherApiKey)
	if clientErr != nil {
		zap.S().Error("failed to create weather client")
		os.Exit(1)
	}

	// Lets start listening to some channels via the websocket!
	webSocketClient, err := model.NewWebSocketClient4(config.ServerWSURL, client.AuthToken)
	if err != nil {
		println("failed to connect to the web socket")
		PrintError(err)
	}

	webSocketClient.Listen()

	go func() {
		for resp := range webSocketClient.EventChannel {
			HandleWebSocketResponse(resp)
		}
	}()

	// You can block forever with
	select {}
}

func MakeSureServerIsRunning() {
	if props, resp := client.GetOldClientConfig(""); resp.Error != nil {
		println("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		println("Server detected and is running version " + props["Version"])
	}
}

func LoginAsTheBotUser(token string) {
	client.SetToken(token)
	var user *model.User
	var resp *model.Response
	if user, resp = client.GetMe(""); resp.Error != nil {
		println("There was a problem getting the user")
		PrintError(resp.Error)
		os.Exit(1)
	}
	botUser = user
	println("Running as " + user.Username)
}

func FindBotTeam(teamName string) {
	if team, resp := client.GetTeamByName(teamName, ""); resp.Error != nil {
		println("We failed to get the initial load")
		println("or we do not appear to be a member of the team '" + teamName + "'")
		PrintError(resp.Error)
		os.Exit(1)
	} else {
		botTeam = team
	}
}

func CreateBotDebuggingChannelIfNeeded(channelName string) {
	if rchannel, resp := client.GetChannelByName(channelName, botTeam.Id, ""); resp.Error != nil {
		println("We failed to get the channels")
		PrintError(resp.Error)
	} else {
		debuggingChannel = rchannel
		return
	}

	// Looks like we need to create the logging channel
	channel := &model.Channel{}
	channel.Name = channelName
	channel.DisplayName = "Debugging For Sample Bot"
	channel.Purpose = "This is used as a test channel for logging bot debug messages"
	channel.Type = model.CHANNEL_OPEN
	channel.TeamId = botTeam.Id
	if rchannel, resp := client.CreateChannel(channel); resp.Error != nil {
		println("We failed to create the channel " + channelName)
		PrintError(resp.Error)
	} else {
		debuggingChannel = rchannel
		println("Looks like this might be the first run so we've created the channel " + channelName)
	}
}

func CreateReply(msg string, replyToId string, replyToUserId string) {
	CreatePost(getUserMention(replyToUserId)+": "+msg, replyToId)
}

func CreatePost(msg string, replyToId string) {
	post := &model.Post{}
	post.ChannelId = debuggingChannel.Id
	post.Message = msg

	post.RootId = replyToId

	if _, resp := client.CreatePost(post); resp.Error != nil {
		println("We failed to send a message to the logging channel")
		PrintError(resp.Error)
	}
}

func CreateReaction(emojiName string, postId string) {
	reaction := &model.Reaction{UserId: botUser.Id, PostId: postId, EmojiName: emojiName}
	if _, resp := client.SaveReaction(reaction); resp.Error != nil {
		println("We failed to add a reaction to the post")
		PrintError(resp.Error)
	}
}

func HandleWebSocketResponse(event *model.WebSocketEvent) {
	HandleMsgFromDebuggingChannel(event)
}

func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

func getUserMention(userId string) string {
	user, _ := client.GetUser(userId, "")
	return "@" + user.Username
}

func HandleMsgFromDebuggingChannel(event *model.WebSocketEvent) {
	// If this isn't the debugging channel then lets ignore it
	if event.Broadcast.ChannelId != debuggingChannel.Id {
		return
	}

	// Lets only reply to messages posted events
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post != nil {

		// ignore my events
		if post.UserId == botUser.Id {
			return
		}

		replyToId := ""
		if post.RootId != "" {
			replyToId = post.RootId
		}

		// if you see any word matching 'hello' then respond
		if matched, _ := regexp.MatchString(globalRegexOptions+`\bsalut|allo\b`, post.Message); matched {
			choices := []string{"aaaaaaayyeee", "sup", "yo"}
			CreateReply(randomChoice(choices), replyToId, post.UserId)
			return
		}

		if matched := regexp.MustCompile(globalRegexOptions+`(xd+)`).FindAllStringSubmatch(post.Message, -1); matched != nil {
			CreatePost("haha "+matched[0][1], replyToId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`\banime|animuh|weeb|weaboo\b`, post.Message); matched {
			CreatePost("### Disgusting weebs rolf :huel:", replyToId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`\bvidya|bonshommes\b`, post.Message); matched {
			CreatePost("rolf vous avez quel age?", replyToId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`\bvelo.*hiver\b`, post.Message); matched {
			CreatePost("wow cest fukin dangereux faut vraiment etre retarded pour cycler en hiver (dans une tempete de verglas) :huel:", replyToId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`\b:disappear:|peace|alp|bye|:wave:|see ya|au revoir|ciao|chow|a tantot\b`, post.Message); matched {
			CreatePost("hey salut la, a prochaine, on se revoit, stait bin lfun", replyToId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`\bbon matin|morning|mornin\b`, post.Message); matched {
			choices := []string{"zzzz kill me now", "omgggggg"}
			CreatePost(randomChoice(choices), replyToId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`\bmirin\b`, post.Message); matched {
			CreateReply("fucking mirin", replyToId, post.UserId)
			return
		}

		if matched, _ := regexp.MatchString(globalRegexOptions+`;-?\)|:wink:`, post.Message); matched {
			CreateReaction("wink", post.Id)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`:-?P|:stuck_out_tongue:`, post.Message); matched {
			CreateReaction("stuck_out_tongue", post.Id)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`:fuck:`, post.Message); matched {
			CreateReaction("fuck", post.Id)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`\^`, post.Message); matched {
			CreatePost("^", replyToId)
			CreateReaction("point_up_2", post.Id)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`\bthis\b`, post.Message); matched {
			CreatePost("this", replyToId)
			CreateReaction("point_up_2", post.Id)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`\breddit\b`, post.Message); matched {
			CreateReply("\\>reddit", replyToId, post.UserId)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`\btumblr\b`, post.Message); matched {
			CreateReply("\\>tumblr", replyToId, post.UserId)
			return
		}
		if matched, _ := regexp.MatchString(globalRegexOptions+`\btgif\b`, post.Message); matched {
			CreateReply("tgiff*", replyToId, post.UserId)
			return
		}

		// Named commands
		if matched := regexp.MustCompile(globalRegexOptions+"^@huel420-new (.*)$").FindAllStringSubmatch(post.Message, -1); matched != nil {
			is420 := time.Now().Month() == time.April && time.Now().Day() == 20
			user, _ := client.GetUser(post.UserId, "")
			isAdminPoggers := user.Username == "gravufo" || user.Username == "roujo"
			command := matched[0][1]

			if matched, _ := regexp.MatchString(globalRegexOptions+`^stfu|fuck you|fuck off|ta yeule|tayeule|shut up|shut the fuck up$`, command); matched {
				choices := []string{"no u?", "no u", ":chuckles:", "rolf"}
				CreateReply(randomChoice(choices), replyToId, post.UserId)
				return
			}
			if matched, _ := regexp.MatchString(globalRegexOptions+`^thanks|merci|ty|thx$`, command); matched {
				choices := []string{"de rien la", "np", "np ;)"}
				CreateReply(randomChoice(choices), replyToId, post.UserId)
				return
			}
			if matched, _ := regexp.MatchString(globalRegexOptions+`^est-ce qu.* ?$`, command); matched {
				choices := []string{"maybe", "??", "yess", "no", "rolf oui", "omgggg no"}
				CreateReply(randomChoice(choices), replyToId, post.UserId)
				return
			}
			if matched, _ := regexp.MatchString(globalRegexOptions+`^I love you$`, command); matched {
				CreateReply("<3", replyToId, post.UserId)
				return
			}

			if matched, _ := regexp.MatchString(globalRegexOptions+`^charge up$`, command); matched && (is420 || isAdminPoggers) {
				chargeValue := rand.Intn(5) - 1
				chargeMap[post.UserId] += chargeValue
				message := ""
				switch {
				case chargeValue < 0:
					message = "You lost " + strconv.Itoa(chargeValue*-1) + " charge points :lamo:"
				case chargeValue > 0:
					message = "You gained " + strconv.Itoa(chargeValue) + " charge points :hype:"
				case chargeValue == 0:
					message = "You gained " + strconv.Itoa(chargeValue) + " charge points :pepehands:"
				}
				CreateReply(message, post.Id, post.UserId)
				return
			}

			if matched, _ := regexp.MatchString(globalRegexOptions+"^charge level$", command); matched && (is420 || isAdminPoggers) {
				chargeValue := chargeMap[post.UserId]
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
				CreateReply(message, post.Id, post.UserId)
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
						CreateReply(message, post.Id, post.UserId)
						return
					}
					from = strings.ToUpper(matched[0][3])
					to = strings.ToUpper(matched[0][4])
				}

				amountStr := strconv.FormatFloat(amount, 'f', 2, 64)
				convertedValue, err := commands.Convert(from, to, amount)
				if err != nil {
					message = "Couldn't convert " + amountStr + " " + from + " to " + to + "."
					CreateReply(message, post.Id, post.UserId)
					return
				}
				message = amountStr + " " + from + " = " + strconv.FormatFloat(convertedValue, 'f', 5, 64) + " " + to

				CreateReply(message, post.Id, post.UserId)
				return
			}

			// Weather
			if matched := regexp.MustCompile(globalRegexOptions+`^weather(?: (.*))?$`).FindAllStringSubmatch(command, -1); matched != nil {
				location := "Montreal"
				if matched[0][1] != "" {
					location = matched[0][1]
				}
				message, err := weatherClient.GetWeather(location)
				if err != nil {
					CreateReply("Couldn't get weather for "+location+".", post.Id, post.UserId)
					return
				}
				CreatePost(message, post.Id)
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
					CreateReply("Couldn't get definition for "+word+".", post.Id, post.UserId)
					return
				}
				message := fmt.Sprintf("%s\n\n_%s_\n\n**by: %s**\n\n`%d`:+1: `%d`:-1:", result.Definition, result.Example, result.Author, result.Upvote, result.Downvote)
				CreatePost(message, post.Id)
				return
			}

			// romaji
			if matched := regexp.MustCompile(globalRegexOptions+`^romaji(?: (.*))?$`).FindAllStringSubmatch(command, -1); matched != nil {
				word := matched[0][1]
				if matched[0][1] == "" {
					CreateReply("Please provide a word to convert.", post.Id, post.UserId)
					return
				}
				message := kana.KanaToRomaji(word)
				CreatePost(message, post.Id)
				return
			}

			// hiragana
			if matched := regexp.MustCompile(globalRegexOptions+`^hiragana(?: (.*))?$`).FindAllStringSubmatch(command, -1); matched != nil {
				word := matched[0][1]
				if matched[0][1] == "" {
					CreateReply("Please provide a word to convert.", post.Id, post.UserId)
					return
				}
				message := kana.RomajiToHiragana(word)
				CreatePost(message, post.Id)
				return
			}

			// katakana
			if matched := regexp.MustCompile(globalRegexOptions+`^katakana(?: (.*))?$`).FindAllStringSubmatch(command, -1); matched != nil {
				word := matched[0][1]
				if matched[0][1] == "" {
					CreateReply("Please provide a word to convert.", post.Id, post.UserId)
					return
				}
				message := kana.RomajiToKatakana(word)
				CreatePost(message, post.Id)
				return
			}

			// wotd japanese
			if matched := regexp.MustCompile(globalRegexOptions+`^wotd japanese.*$`).FindAllStringSubmatch(command, -1); matched != nil {
				message, err := commands.GetWotdJapanese()
				if err != nil {
					CreateReply("Couldn't get WotD Japanese.", post.Id, post.UserId)
					return
				}
				CreatePost(message, post.Id)
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
						CreateReply(fmt.Sprintf("lel nice fake player id: %d.", playerId), post.Id, post.UserId)
						return
					}
				}

				mmr, err := commands.GetDotaMMR(playerId)
				if err != nil {
					CreateReply(fmt.Sprintf("rofl %d existe meme pas zzz", playerId), post.Id, post.UserId)
					return
				}
				if playerId == 12088460 {
					message = fmt.Sprintf("lel j'suis rendu %d ez gaem road to 4k", mmr.SoloCompetitiveRank)
					CreateReply(message, post.Id, post.UserId)
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

				CreateReply(message, post.Id, post.UserId)
				return
			}

			// Rolls
			if matched := regexp.MustCompile(globalRegexOptions+`^roll(?: (\d+|:weed:))?$`).FindAllStringSubmatch(command, -1); matched != nil {
				requestedRoll := 0
				if matched[0][1] == "" || matched[0][1] == ":weed:" {
					requestedRoll = 420
				} else if matched[0][1] == "dice" {
					requestedRoll = 6
				} else {
					requestedRoll, _ = strconv.Atoi(matched[0][1])
				}
				message := rollDice(requestedRoll)
				CreateReply(message, post.Id, post.UserId)
				return
			}

			CreatePost("Kes tu. Veux????", replyToId)
		}
	}
}

func rollDice(dice int) (message string) {
	roll := rand.Intn(dice) + 1
	if dice == 420 {
		if time.Now().Hour()%12 == 4 && time.Now().Minute() == 20 {
			message = strconv.Itoa(roll) + " "

			if roll == 420 {
				message += "BIG WINNER WOW :musk: :weed:"
			} else if roll == 69 {
				message += "_Nice._ :smugpepe:"
			} else {
				message += ":chuckle:"
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

func PrintError(err *model.AppError) {
	println("\tError Details:")
	println("\t\t" + err.Message)
	println("\t\t" + err.Id)
	println("\t\t" + err.DetailedError)
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
