import re
import requests

from random import randint
from mattermost_bot.bot import listen_to
from mattermost_bot.bot import respond_to

@respond_to('^roll (\d+)', re.IGNORECASE)
def roll_number(message, number=100):
    random = randint(1,int(number))
    message.reply(random)

roll_number.__doc__ = "roll between 1 and {{number}}"

@respond_to('^roll$', re.IGNORECASE)
def roll(message):
    random = randint(1,100)
    message.reply(random)

roll.__doc__ = "roll between 1 and 100"
