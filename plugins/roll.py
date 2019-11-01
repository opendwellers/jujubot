import re
import requests
import datetime

from random import randint
from mmpy_bot.bot import listen_to
from mmpy_bot.bot import respond_to

@respond_to('^roll ?(\d+)?$', re.IGNORECASE)
def roll_number(message, number=100):
    number = 100 if number is None else number
    
    if number is 420 :
        
        now = datetime.datetime.now()
        
        if (now.hour is 4 or now.hour is 16) and now.minute is 20 :
            
            random = randint(1,int(number))
            
            if random is 420:
                message.reply('###BIG WINNER WOW ' + random)
            else:
                message.reply(random + ' :chuckle:')
        else:
            message.reply('Spa leur smh')
        
    if int(number) > 1:
        random = randint(1,int(number))
        message.reply(random)
    else:
        message.reply(number + ' is not a valid number for roll command')

roll_number.__doc__ = "Roll between 1 and {{number}}, between 1 and 100 if nothing specified"
