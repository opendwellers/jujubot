import re
import time
import random

from random import randint
from mmpy_bot.bot import listen_to
from mmpy_bot.bot import respond_to

@respond_to('(?<!\w)(?:salut|allo)(?!\w)', re.IGNORECASE)
def hi(message):
    hi_messages = ['aaaaaaayyeee', 'sup', 'yo']
    message.reply(random.choice(hi_messages))

@respond_to('(?<!\w)(?:stfu|fuck you|fuck off|ta yeule|tayeule|shut up|shut the fuck up)(?!\w)', re.IGNORECASE)
def tayeule_reply(message):
    tayeule_messages = ['no u?', 'NO U', ':chuckles:', 'rolf']
    message.reply(random.choice(tayeule_messages))

@respond_to('^(?:thanks|merci|ty|thx)(?!\w)', re.IGNORECASE)
def thanks(message):
    replies = ['de rien la', 'np', 'np ;)'] 
    message.reply(random.choice(replies))

@respond_to('^est-ce qu.* ?$')
def question(message):
    answers = ['maybe', '??', 'yess', 'no', 'rolf oui', 'omgggg no']
    message.reply(random.choice(answers))

@respond_to('(?<!\w)I love you(?!\w)', re.IGNORECASE)
def love(message):
    message.reply('<3')

@listen_to('(?<!\w)(?:anime|animuh|weeb|weaboo)(?!\w)', re.IGNORECASE)
def disguted(message):
    message.send('### Disgusting weebs rolf :huel:')

@listen_to('(?<!\w)(?:vidya|bonshommes)(?!\w)', re.IGNORECASE)
def quel_age(message):
    message.send('rolf vous avez quel age?')

@listen_to('(?<!\w)(xd+)(?!\w)', re.IGNORECASE)
def reply_xd(message, xd):
    message.send('haha '+ xd)

@respond_to('(?<!\w)(?::disappear:|peace|alp|bye|:wave:|see ya|au revoir|ciao|chow|a tantot)(?!\w)', re.IGNORECASE)
@listen_to('(?<!\w)(?::disappear:|peace|alp|bye|:wave:|see ya|au revoir|ciao|chow|a tantot)(?!\w)', re.IGNORECASE)
def reply_bye(message):
    message.send('hey salut la, a prochaine, on se revoit, stait bin lfun')

@respond_to('(?<!\w)(?:bon matin|morning|mornin)(?!\w)', re.IGNORECASE)
@listen_to('(?<!\w)(?:bon matin|morning|mornin)(?!\w)', re.IGNORECASE)
def reply_morning(message):
    morning_messages = ['zzzz kill me now', 'omgggggg']
    message.send(random.choice(morning_messages))

@listen_to('(?<!\w)(?:velo.*hiver)(?!\w)', re.IGNORECASE)
def reply_velo(message):
    message.send('wow cest fukin dangereux faut vraiment etre retarded pour cycler en hiver (dans une tempete de verglas) :huel:')

@respond_to('(?<!\w)mirin(?!\w)', re.IGNORECASE)
@listen_to('(?<!\w)mirin(?!\w)', re.IGNORECASE)
def reply_mirin(message):
    message.reply('fucking mirin')

@listen_to('(?<!\w)(?:;-?\)|:wink:)(?!\w)', re.IGNORECASE)
def reply_wink(message):
    message.send(';)')

@listen_to('(?<!\w)(?::-?P|:stuck_out_tongue:)(?!\w)', re.IGNORECASE)
def reply_tongue(message):
    message.send(':P')

@listen_to('(?<!\w)(?::fuck:)(?!\w)', re.IGNORECASE)
def react_fuck(message):
    message.react('fuck')

@listen_to('^\^$')
def upboat(message):
    message.send('^')

@listen_to('^this$')
def upvote(message):
    message.send('this')

@listen_to('reddit', re.IGNORECASE)
def reddit(message):
    message.reply('\>reddit')

@listen_to('tumblr')
def tumblr(message):
    message.reply('\>tumblr')

@listen_to('tgif(?!f)', re.IGNORECASE)
def friday(message):
    message.reply('tgiff*')
