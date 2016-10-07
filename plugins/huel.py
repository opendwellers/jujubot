import re
import time
import random

from random import randint
from mattermost_bot.bot import listen_to
from mattermost_bot.bot import respond_to

@respond_to('salut', re.IGNORECASE)
def hi(message):
    message.reply('aaaaaayyyee')

@respond_to('thanks|merci|ty|thx', re.IGNORECASE)
def hi(message):
    replies = ['de rien la', 'np', 'np ;)'] 
    message.reply(random.choice(replies))

@respond_to('^est-ce qu.* ?$')
def hi(message):
    answers = ['maybe', '??', 'yess', 'no', 'rolf oui', 'omgggg no']
    message.reply(random.choice(answers))
hi.__doc__ = "legit answers try it out!"

@respond_to('I love you', re.IGNORECASE)
def hi(message):
    message.reply('<3')
hi.__doc__="ma vous montrer comment je vous aime"

@listen_to('.* *anime.*', re.IGNORECASE)
def disguted(message):
    message.send('### Disgusting weebs rolf :huel:')

@listen_to('.*vidya|bonshommes.*', re.IGNORECASE)
def quel_age(message):
    message.send('rolf vous avez quel age?')

@listen_to('.*(xd+).*', re.IGNORECASE)
def reply_xd(message, xd):
    message.send('haha '+ xd)

@listen_to('peace|bye|:wave:|alp|see ya|au revoir|ciao|chow|a tantot', re.IGNORECASE)
def reply_xd(message):
    message.send('hey bye la')

@listen_to('bon matin|morning|mornin', re.IGNORECASE)
def reply_xd(message):
    message.send('zzzz kill me now')

@listen_to('velo.*hiver', re.IGNORECASE)
def reply_velo(message):
    message.send('wow cest fukin dangereux faut vraiment etre retarded pour cycler en hiver :huel:') 

@listen_to('^mirin.*?', re.IGNORECASE)
def reply_mirin(message):
    message.reply('fucking mirin')

@listen_to(';\)', re.IGNORECASE)
def reply_wink(message):
    message.send(';)')

@listen_to('^\^$')
def upboat(message):
    message.send('^')

