import re

from mattermost_bot.bot import listen_to
from mattermost_bot.bot import respond_to

@respond_to('salut', re.IGNORECASE)
def hi(message):
        message.reply('aaaaaayyyee')

@respond_to('thanks|merci|ty|thx', re.IGNORECASE)
def hi(message):
        message.reply('de rien la')

@listen_to('.* *anime.*', re.IGNORECASE)
def disguted(message):
    message.send('### Disgusting weebs rolf :huel:')

@listen_to('.*vidya|bonshommes.*', re.IGNORECASE)
def quel_age(message):
    message.send('rolf vous avez quel age?')

@listen_to('.*(xd+).*', re.IGNORECASE)
def reply_xd(message, xd):
    message.send('haha '+ xd)

@listen_to('peace|bye|:wave:|alp|see ya|au revoir|ciao|chow', re.IGNORECASE)
def reply_xd(message):
    message.send('hey bye la')

@listen_to('bon matin|morning|mornin', re.IGNORECASE)
def reply_xd(message):
    message.send('zzzz kill me now')
