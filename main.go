package main

import (
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/opendwellers/jujubot/pkg/config"
	"go.uber.org/zap"
)

const (
	SAMPLE_NAME = "Huel420 Bot"

	// BOT_AUTH_TOKEN = "hsb6jqccdfgo7jsa8mcpdccxsy"

	// HOSTNAME         = "chat.polycancer.org"
	// SERVER_URL       = "https://" + HOSTNAME
	// SERVER_WS_URL    = "wss://" + HOSTNAME
	// TEAM_NAME        = "Dwellers"
	// CHANNEL_LOG_NAME = "sandbox"
)

var chargeMap map[string]int = map[string]int{}

var client *model.Client4
var webSocketClient *model.WebSocketClient

var botUser *model.User
var botTeam *model.Team
var debuggingChannel *model.Channel

// Documentation for the Go driver can be found
// at https://godoc.org/github.com/mattermost/platform/model#Client
func main() {
	println(SAMPLE_NAME)
	zap.S().Info("loading configuration file")
	config, error := config.LoadConfig()
	if error != nil {
		zap.S().Error("failed to load configuration file")
		os.Exit(1)
	}

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
	SendMsgToDebuggingChannel("_"+SAMPLE_NAME+" has **started** running_", "")

	// Lets start listening to some channels via the websocket!
	webSocketClient, err := model.NewWebSocketClient4(config.ServerWSURL, client.AuthToken)
	if err != nil {
		println("We failed to connect to the web socket")
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

func SendMsgToDebuggingChannel(msg string, replyToId string) {
	post := &model.Post{}
	post.ChannelId = debuggingChannel.Id
	post.Message = msg

	post.RootId = replyToId

	if _, resp := client.CreatePost(post); resp.Error != nil {
		println("We failed to send a message to the logging channel")
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
	var user, _ = client.GetUser(userId, "")
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

	println("responding to debugging channel msg")

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post != nil {

		// ignore my events
		if post.UserId == botUser.Id {
			return
		}

		var replyToId string
		if post.RootId != "" {
			replyToId = post.RootId
		} else {
			replyToId = ""
		}

		// if you see any word matching 'hello' then respond
		if matched, _ := regexp.MatchString(`(?:salut|allo)`, post.Message); matched {
			var choices = []string{"aaaaaaayyeee", "sup", "yo", "xd"}
			SendMsgToDebuggingChannel(randomChoice(choices), replyToId)
			return
		}

		if matched := regexp.MustCompile(`(xd+)`).FindAllStringSubmatch(post.Message, -1); matched != nil {
			SendMsgToDebuggingChannel("haha "+matched[0][1], replyToId)
			return
		}

		// Named commands
		if matched := regexp.MustCompile("^@huel420-new (.*)$").FindAllStringSubmatch(post.Message, -1); matched != nil {
			var is420 = time.Now().Month() == time.April && time.Now().Day() == 20
			var user, _ = client.GetUser(post.UserId, "")
			var isAdminPoggers = user.Username == "gravufo" || user.Username == "roujo"

			var command = matched[0][1]
			if matched, _ := regexp.MatchString(`^charge up$`, command); matched && (is420 || isAdminPoggers) {
				var chargeValue = rand.Intn(5) - 1
				chargeMap[post.UserId] += chargeValue
				var message string
				switch {
				case chargeValue < 0:
					message = "You lost " + strconv.Itoa(chargeValue*-1) + " charge points :lamo:"
				case chargeValue > 0:
					message = "You gained " + strconv.Itoa(chargeValue) + " charge points :hype:"
				case chargeValue == 0:
					message = "You gained " + strconv.Itoa(chargeValue) + " charge points :pepehands:"
				}
				SendMsgToDebuggingChannel(getUserMention(post.UserId)+": "+message, replyToId)
				return
			}

			if matched, _ := regexp.MatchString("^charge level$", command); matched && (is420 || isAdminPoggers) {
				var chargeValue = chargeMap[post.UserId]
				var chargeValueStr = strconv.Itoa(chargeValue)
				var message string
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
				SendMsgToDebuggingChannel(getUserMention(post.UserId)+": "+message, replyToId)
				return
			}

			// Rolls
			if matched := regexp.MustCompile(`^roll(?: (\d+|:weed:))?$`).FindAllStringSubmatch(command, -1); matched != nil {
				var requestedRoll int
				if matched[0][1] == "" || matched[0][1] == ":weed:" {
					requestedRoll = 420
				} else if matched[0][1] == "dice" {
					requestedRoll = 6
				} else {
					requestedRoll, _ = strconv.Atoi(matched[0][1])
				}
				var message = rollDice(requestedRoll)
				SendMsgToDebuggingChannel(getUserMention(post.UserId)+": "+message, replyToId)
				return
			}

			SendMsgToDebuggingChannel("Kes tu. Veux????", replyToId)
		}
	}
}

func rollDice(dice int) (message string) {
	var roll int
	roll = rand.Intn(dice) + 1
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

			SendMsgToDebuggingChannel("_"+SAMPLE_NAME+" has **stopped** running_", "")
			os.Exit(0)
		}
	}()
}
