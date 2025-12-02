package bot

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/opendwellers/jujubot/pkg/commands"
	"github.com/opendwellers/jujubot/pkg/config"
	"go.uber.org/zap"
)

// Bot represents the Mattermost bot instance
type Bot struct {
	config          config.Config
	client          *model.Client4
	webSocketClient *model.WebSocketClient
	user            *model.User
	team            *model.Team
	debugChannel    *model.Channel
	weatherClient   commands.Weather
	chargeMap       map[string]int
}

// New creates a new Bot instance
func New(cfg config.Config) (*Bot, error) {
	b := &Bot{
		config:    cfg,
		client:    model.NewAPIv4Client(cfg.ServerURL),
		chargeMap: make(map[string]int),
	}

	// Initialize weather client
	weatherClient, err := commands.NewWeatherClient(cfg.OpenWeatherApiKey)
	if err != nil {
		return nil, err
	}
	b.weatherClient = weatherClient

	return b, nil
}

// Start initializes the bot and starts listening for events
func (b *Bot) Start() error {
	zap.S().Info("Connecting to Mattermost at " + b.config.ServerURL)

	// Login
	if err := b.login(); err != nil {
		return err
	}

	// Find team
	if err := b.findTeam(); err != nil {
		return err
	}

	// Setup debugging channel
	b.setupDebuggingChannel()

	// Create readiness file
	if err := b.createReadinessFile(); err != nil {
		return err
	}

	// Setup graceful shutdown
	b.setupGracefulShutdown()

	zap.S().Info("Bot is now running and listening to messages.")

	// Start WebSocket listener in a goroutine
	go b.startWebSocketListener()

	return nil
}

// login authenticates the bot with the Mattermost server
func (b *Bot) login() error {
	b.client.SetToken(b.config.AuthToken)

	user, _, err := b.client.GetMe(context.TODO(), "")
	if err != nil {
		zap.S().Error("Failed to get bot user", zap.Error(err))
		return err
	}

	b.user = user
	zap.S().Info("Running as " + user.Username)
	return nil
}

// findTeam locates the bot's team
func (b *Bot) findTeam() error {
	team, _, err := b.client.GetTeamByName(context.TODO(), b.config.TeamName, "")
	if err != nil {
		zap.S().Error("Failed to find team '"+b.config.TeamName+"'", zap.Error(err))
		return err
	}

	b.team = team
	zap.S().Info("Found team " + team.Name)
	return nil
}

// setupDebuggingChannel creates or finds the debugging channel
func (b *Bot) setupDebuggingChannel() {
	rchannel, _, err := b.client.GetChannelByName(context.TODO(), b.config.ChannelLogName, b.team.Id, "")
	if err == nil {
		b.debugChannel = rchannel
		return
	}

	zap.S().Info("Creating debugging channel " + b.config.ChannelLogName)

	channel := &model.Channel{
		Name:        b.config.ChannelLogName,
		DisplayName: "Debugging For Sample Bot",
		Purpose:     "This is used as a test channel for logging bot debug messages",
		Type:        model.ChannelTypeOpen,
		TeamId:      b.team.Id,
	}

	rchannel, _, err = b.client.CreateChannel(context.TODO(), channel)
	if err != nil {
		zap.S().Error("Failed to create channel "+b.config.ChannelLogName, zap.Error(err))
		return
	}

	b.debugChannel = rchannel
	zap.S().Info("Created debugging channel " + b.config.ChannelLogName)
}

// createReadinessFile creates the readiness probe file
func (b *Bot) createReadinessFile() error {
	f, err := os.Create("/tmp/ready")
	if err != nil {
		zap.S().Error("Failed to create /tmp/ready", zap.Error(err))
		return err
	}
	f.Close()
	return nil
}

// startWebSocketListener manages the WebSocket connection
func (b *Bot) startWebSocketListener() {
	for {
		wsClient, err := model.NewWebSocketClient4(b.config.ServerWSURL, b.client.AuthToken)
		if err != nil {
			zap.S().Error("Failed to connect to WebSocket", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}

		b.webSocketClient = wsClient
		b.listen()
	}
}

// listen handles incoming WebSocket events
func (b *Bot) listen() {
	b.webSocketClient.Listen()
	defer b.webSocketClient.Close()

	for {
		if b.webSocketClient.ListenError != nil {
			zap.S().Error("WebSocket listen error", zap.Error(b.webSocketClient.ListenError))
			zap.S().Info("Reconnecting to WebSocket")
			return
		}

		event := <-b.webSocketClient.EventChannel
		if event == nil {
			continue
		}

		b.handleEvent(event)
	}
}

// setupGracefulShutdown handles OS interrupt signals
func (b *Bot) setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if b.webSocketClient != nil {
				b.webSocketClient.Close()
			}
			os.Exit(0)
		}
	}()
}
