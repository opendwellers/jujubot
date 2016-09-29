import re

from mattermost_bot.bot import listen_to
from mattermost_bot.bot import respond_to


@respond_to('hi', re.IGNORECASE)
def hi(message):
        message.reply('I can understand hi or HI!')


@respond_to('I love you')
def love(message):
        message.reply('I love you too!')


@listen_to('Can someone help me?')
def help_me(message):
        # Message is replied to the sender (prefixed with @user)
        message.reply('Yes, I can!')

                # Message is sent on the channel
                    # message.send('I can help everybody!')
