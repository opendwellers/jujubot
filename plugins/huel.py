import re
import time
import random

from random import randint
from mmpy_bot.bot import listen_to
from mmpy_bot.bot import respond_to

@respond_to('salut|allo', re.IGNORECASE)
def hi(message):
    hi_messages = ['aaaaaaayyeee', 'sup', 'yo']
    message.reply(random.choice(hi_messages))

@respond_to('stfu|fuck you|fuck off|ta yeule|tayeule|shut up|shut the fuck up', re.IGNORECASE)
def tayeule_reply(message):
    tayeule_messages = ['no u?', 'NO U', ':chuckles:', 'rolf']
    message.reply(random.choice(tayeule_messages))

@respond_to('^thanks|^merci|^ty|^thx', re.IGNORECASE)
def hi(message):
    replies = ['de rien la', 'np', 'np ;)'] 
    message.reply(random.choice(replies))

@respond_to('^est-ce qu.* ?$')
def hi(message):
    answers = ['maybe', '??', 'yess', 'no', 'rolf oui', 'omgggg no']
    message.reply(random.choice(answers))

@respond_to('I love you', re.IGNORECASE)
def hi(message):
    message.reply('<3')

@listen_to('anime|animuh|weeb|weaboo', re.IGNORECASE)
def disguted(message):
    message.send('### Disgusting weebs rolf :huel:')

@listen_to('.*vidya|bonshommes.*', re.IGNORECASE)
def quel_age(message):
    message.send('rolf vous avez quel age?')

@listen_to('(?<!bagel)(xd+)', re.IGNORECASE)
def reply_xd(message, xd):
    message.send('haha '+ xd)

@respond_to(':disappear:|peace|alp|bye|:wave:|see ya|au revoir|ciao|chow|a tantot', re.IGNORECASE)
@listen_to(':disappear:|peace|alp|bye|:wave:|see ya|au revoir|ciao|chow|a tantot', re.IGNORECASE)
def reply_bye(message):
    message.send('hey salut la, a prochaine, on se revoit, stait bin lfun')

@respond_to('bon matin|morning|mornin', re.IGNORECASE)
@listen_to('bon matin|morning|mornin', re.IGNORECASE)
def reply_morning(message):
    morning_messages = ['zzzz kill me now', 'omgggggg']
    message.send(random.choice(morning_messages))

@listen_to('velo.*hiver', re.IGNORECASE)
def reply_velo(message):
    message.send('wow cest fukin dangereux faut vraiment etre retarded pour cycler en hiver (dans une tempete de verglas) :huel:')

@respond_to('^mirin.*?', re.IGNORECASE)
@listen_to('^mirin.*?', re.IGNORECASE)
def reply_mirin(message):
    message.reply('fucking mirin')

@listen_to('(?<!\w)(;-?\)|:wink:)(?!\w)', re.IGNORECASE)
def reply_wink(message):
    message.send(';)')

@listen_to('(?<!\w)(:-?P|:stuck_out_tongue:)(?!\w)', re.IGNORECASE)
def reply_wink(message):
    message.send(':P')

@listen_to('^\^$')
def upboat(message):
    message.send('^')

@listen_to('^this$')
def upboat(message):
    message.send('this')

@listen_to('reddit', re.IGNORECASE)
def reddit(message):
    message.reply('\>reddit')

@listen_to('tumblr')
def reddit(message):
    message.reply('\>tumblr')

@listen_to('tgif(?!f)', re.IGNORECASE)
def reddit(message):
    message.reply('tgiff*')
