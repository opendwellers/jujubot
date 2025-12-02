package bot

import (
	"math/rand"
	"regexp"

	"github.com/mattermost/mattermost/server/public/model"
)

const globalRegexOptions = "(?i)"

// patternReaction defines a pattern and its response
type patternReaction struct {
	pattern     string
	handler     func(b *Bot, post *model.Post, replyToId string, matched [][]string) bool
	useSubmatch bool
}

// patternReactions contains all pattern-based reactions
var patternReactions = []patternReaction{
	// Greetings
	{
		pattern: `\bsalut|allo\b`,
		handler: func(b *Bot, post *model.Post, replyToId string, _ [][]string) bool {
			choices := []string{"aaaaaaayyeee", "sup", "yo"}
			b.createReply(post.ChannelId, randomChoice(choices), replyToId, post.UserId)
			return true
		},
	},
	// XD reaction
	{
		pattern:     `(xd+)`,
		useSubmatch: true,
		handler: func(b *Bot, post *model.Post, replyToId string, matched [][]string) bool {
			b.createPost(post.ChannelId, "haha "+matched[0][1], replyToId)
			return true
		},
	},
	// Anime/weeb reaction
	{
		pattern: `\banime|animuh|weeb|weaboo\b`,
		handler: func(b *Bot, post *model.Post, replyToId string, _ [][]string) bool {
			b.createPost(post.ChannelId, "### Disgusting weebs rolf :huel:", replyToId)
			return true
		},
	},
	// Vidya reaction
	{
		pattern: `\bvidya|bonshommes\b`,
		handler: func(b *Bot, post *model.Post, replyToId string, _ [][]string) bool {
			b.createPost(post.ChannelId, "rolf vous avez quel age?", replyToId)
			return true
		},
	},
	// Winter cycling
	{
		pattern: `\bvelo.*hiver\b`,
		handler: func(b *Bot, post *model.Post, replyToId string, _ [][]string) bool {
			b.createPost(post.ChannelId, "wow cest fukin dangereux faut vraiment etre retarded pour cycler en hiver (dans une tempete de verglas) :huel:", replyToId)
			return true
		},
	},
	// Goodbye
	{
		pattern: `\b:disappear:|peace|alp|bye|:wave:|see ya|au revoir|ciao|chow|a tantot\b`,
		handler: func(b *Bot, post *model.Post, replyToId string, _ [][]string) bool {
			b.createPost(post.ChannelId, "hey salut la, a prochaine, on se revoit, stait bin lfun", replyToId)
			return true
		},
	},
	// Good morning
	{
		pattern: `\bbon matin|morning|mornin\b`,
		handler: func(b *Bot, post *model.Post, replyToId string, _ [][]string) bool {
			choices := []string{"zzzz kill me now", "omgggggg"}
			b.createPost(post.ChannelId, randomChoice(choices), replyToId)
			return true
		},
	},
	// Mirin
	{
		pattern: `\bmirin\b`,
		handler: func(b *Bot, post *model.Post, replyToId string, _ [][]string) bool {
			b.createReply(post.ChannelId, "fucking mirin", replyToId, post.UserId)
			return true
		},
	},
	// Wink emoji
	{
		pattern: `;-?\)(\s|$)|:wink:`,
		handler: func(b *Bot, post *model.Post, _ string, _ [][]string) bool {
			b.createReaction("wink", post.Id)
			return true
		},
	},
	// Tongue emoji
	{
		pattern: `:-?P(\s|$)|:stuck_out_tongue:`,
		handler: func(b *Bot, post *model.Post, _ string, _ [][]string) bool {
			b.createReaction("stuck_out_tongue", post.Id)
			return true
		},
	},
	// Fuck emoji
	{
		pattern: `:fuck:`,
		handler: func(b *Bot, post *model.Post, _ string, _ [][]string) bool {
			b.createReaction("fuck", post.Id)
			return true
		},
	},
	// Caret (^)
	{
		pattern: `(\s|^)\^(\s|$)`,
		handler: func(b *Bot, post *model.Post, replyToId string, _ [][]string) bool {
			b.createPost(post.ChannelId, "^", replyToId)
			b.createReaction("point_up_2", post.Id)
			return true
		},
	},
	// "this"
	{
		pattern: `^this$`,
		handler: func(b *Bot, post *model.Post, replyToId string, _ [][]string) bool {
			b.createPost(post.ChannelId, "this", replyToId)
			b.createReaction("point_up_2", post.Id)
			return true
		},
	},
	// Reddit
	{
		pattern: `\breddit\b`,
		handler: func(b *Bot, post *model.Post, replyToId string, _ [][]string) bool {
			b.createReply(post.ChannelId, "\\>reddit", replyToId, post.UserId)
			return true
		},
	},
	// Tumblr
	{
		pattern: `\btumblr\b`,
		handler: func(b *Bot, post *model.Post, replyToId string, _ [][]string) bool {
			b.createReply(post.ChannelId, "\\>tumblr", replyToId, post.UserId)
			return true
		},
	},
	// TGIF
	{
		pattern: `\btgif\b`,
		handler: func(b *Bot, post *model.Post, replyToId string, _ [][]string) bool {
			b.createReply(post.ChannelId, "tgiff*", replyToId, post.UserId)
			return true
		},
	},
	// Charging up
	{
		pattern:     `(a{5,}h{2,}!*)|:charging_up:`,
		useSubmatch: true,
		handler: func(b *Bot, post *model.Post, replyToId string, matched [][]string) bool {
			match := matched[0][0]
			length := len(match)
			if match == ":charging_up:" {
				length = rand.Intn(50) + 1
			}
			// x1 at 8 characters
			// +1 multiplier every time you add 15 characters
			multiplier := (length-8)/15 + 1
			message := b.chargeUp(post.UserId, multiplier)
			b.createReply(post.ChannelId, message, replyToId, post.UserId)
			return true
		},
	},
}

// handlePatternReactions checks all pattern reactions against the message
func (b *Bot) handlePatternReactions(post *model.Post, replyToId string) bool {
	for _, reaction := range patternReactions {
		pattern := globalRegexOptions + reaction.pattern
		re := regexp.MustCompile(pattern)

		if reaction.useSubmatch {
			if matched := re.FindAllStringSubmatch(post.Message, -1); matched != nil {
				if reaction.handler(b, post, replyToId, matched) {
					return true
				}
			}
		} else {
			if matched, _ := regexp.MatchString(pattern, post.Message); matched {
				if reaction.handler(b, post, replyToId, nil) {
					return true
				}
			}
		}
	}
	return false
}
