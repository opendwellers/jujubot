package bot

import (
	"context"

	"github.com/mattermost/mattermost/server/public/model"
	"go.uber.org/zap"
)

// createPost creates a new post in the specified channel
func (b *Bot) createPost(channelId, message, replyToId string) {
	post := &model.Post{
		ChannelId: channelId,
		Message:   message,
		RootId:    replyToId,
	}

	if _, _, err := b.client.CreatePost(context.TODO(), post); err != nil {
		zap.S().Error("Failed to send message", zap.Error(err))
	}
}

// createReply creates a reply mentioning a specific user
func (b *Bot) createReply(channelId, message, replyToId, replyToUserId string) {
	mention := b.getUserMention(replyToUserId)
	b.createPost(channelId, mention+": "+message, replyToId)
}

// createReaction adds a reaction to a post
func (b *Bot) createReaction(emojiName, postId string) {
	reaction := &model.Reaction{
		UserId:    b.user.Id,
		PostId:    postId,
		EmojiName: emojiName,
	}

	if _, _, err := b.client.SaveReaction(context.TODO(), reaction); err != nil {
		zap.S().Error("Failed to add reaction", zap.Error(err))
	}
}

// getUserMention returns the @mention string for a user
func (b *Bot) getUserMention(userId string) string {
	user, _, err := b.client.GetUser(context.TODO(), userId, "")
	if err != nil {
		zap.S().Error("Failed to get user", zap.Error(err))
		return "@unknown"
	}
	return "@" + user.Username
}
