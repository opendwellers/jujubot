import re
import requests

from mattermost_bot.bot import listen_to
from mattermost_bot.bot import respond_to

@respond_to('currency (\w{3}) (\w{3})$', re.IGNORECASE)
def rates(message, base, otherCurrency):
        url = 'http://api.fixer.io/latest?base={base}'.format(base=base)
        r = requests.get(url)
        if r.status_code == 422:
            message.reply('Invalid base :huel:')
        elif r.status_code == 200:
            jsonResponse = r.json()
            if otherCurrency.upper() in jsonResponse['rates']:
                # print(jsonResponse['rates'][otherCurrency.upper()])
                rate = jsonResponse['rates'][otherCurrency.upper()] 
                message.reply('Conversion rate for {base} to {currency} is: {rate}'.format(base=base.upper(), currency=otherCurrency.upper(), rate=rate))
            else:
                print('NOT')

rates.__doc__ = 'Displays conversion rate between two currencies'
