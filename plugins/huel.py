import re
import time

from random import randint
from mattermost_bot.bot import listen_to
from mattermost_bot.bot import respond_to

@respond_to('salut', re.IGNORECASE)
def hi(message):
    message.reply('aaaaaayyyee')

@respond_to('thanks|merci|ty|thx', re.IGNORECASE)
def hi(message):
    message.reply('de rien la')

@respond_to('^est-ce qu.* ?$')
def hi(message):
    choice = randint(0,5)
    if choice is 0:
        message.reply('maybe') 
    elif choice is 1:
        message.reply('??')
    elif choice is 2:
        message.reply('yess')
    elif choice is 3:
        message.reply('no')
    elif choice is 4:
        message.reply('rolf oui')
    elif choice is 5:
        message.reply('omggggg no')
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


@respond_to('^tayeule$')
def sleep_reply(message):
    message.reply('Ok brb 5 mins')
    time.sleep(int(300))
    message.reply('hey guys :)')
sleep_reply.__doc__ = "je vais shut up pour 5 mins"
