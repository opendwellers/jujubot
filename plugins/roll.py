import re
import requests
import datetime

from random import randint
from mmpy_bot.bot import listen_to
from mmpy_bot.bot import respond_to

@respond_to('^roll ?(\d+|:weed:)?$', re.IGNORECASE)
def roll_number(message, number):
   
    if number == ":weed:":
        number = 420
    
    number = 420 if number is None else int(number)

    if number == 420 :

        now = datetime.datetime.now()

        if (now.hour == 4 or now.hour == 16) and now.minute == 20 :

            random = randint(1, number)

            if random == 420:
                message.reply("{} BIG WINNER WOW :musk: :weed:".format(random))
            elif random == 69:
                message.reply("{} _Nice._ :smugpepe:".format(random))
            else:
                message.reply("{} :chuckles:".format(random))
        else:
            message.reply('Spa leur smh')
    elif number > 1:
        random = randint(1, number)
        message.reply(random)
    else:
        message.reply(number + ' is not a valid number for roll command')

roll_number.__doc__ = "Roll between 1 and {{number}}, between 1 and 100 if nothing specified"
