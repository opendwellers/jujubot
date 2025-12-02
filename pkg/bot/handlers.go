package bot

import (
	"context"
	"encoding/json"
	"math/rand"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"go.uber.org/zap"
)

// handleEvent routes WebSocket events to appropriate handlers
func (b *Bot) handleEvent(event *model.WebSocketEvent) {
	zap.S().Debug("Got event: ", event)

	// Only handle posted messages
	if event.EventType() != model.WebsocketEventPosted {
		return
	}

	var post *model.Post
	json.NewDecoder(strings.NewReader(event.GetData()["post"].(string))).Decode(&post)

	if post == nil {
		return
	}

	// Ignore own messages
	if post.UserId == b.user.Id {
		return
	}

	b.handleMessage(post)
}

// handleMessage processes an incoming message
func (b *Bot) handleMessage(post *model.Post) {
	user, _, err := b.client.GetUser(context.TODO(), post.UserId, "")
	if err != nil {
		zap.S().Info("Failed to get user ", post.UserId, zap.Error(err))
	}
	zap.S().Info("Processing message from user ", user.Username, ": ", post.Message)

	replyToId := post.RootId

	// Check for named commands first (messages starting with @botname)
	if b.handleNamedCommands(post, replyToId) {
		return
	}

	// Then check for pattern-based reactions
	b.handlePatternReactions(post, replyToId)
}

// randomChoice returns a random element from a slice
func randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}
