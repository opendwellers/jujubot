import re

from mattermost_bot.bot import listen_to
from mattermost_bot.bot import respond_to

@respond_to('salut', re.IGNORECASE)
def hi(message):
        message.reply('aaaaaayyyee')

@listen_to('.*anime.*', re.IGNORECASE)
def disguted(message):
    message.send('### Disgusting weebs rolf :huel:')

@listen_to('.*vidya|bonshommes.*', re.IGNORECASE)
def quel_age(message):
    message.send('rolf vous avez quel age?')

@listen_to('.*nodame.*', re.IGNORECASE)
def nodame(message):
    message.send('ruffles what a shit show')
