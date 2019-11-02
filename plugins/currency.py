import re
import requests

from mmpy_bot.bot import listen_to
from mmpy_bot.bot import respond_to

@respond_to('^convert$', re.IGNORECASE)
@respond_to('^convert (\d+)? ?(\w{3}) (?:to )?(\w{3})$', re.IGNORECASE)
def rates(message, amount=None, baseCurrency=None, targetCurrency=None):
        amount = '1' if amount is None else amount
        baseCurrency = 'CAD' if baseCurrency is None else baseCurrency
        targetCurrency = 'USD' if targetCurrency is None else targetCurrency
        url = 'https://frankfurter.app/latest?from={baseCurrency}&to={targetCurrency}&amount={amount}'.format(baseCurrency=baseCurrency, targetCurrency=targetCurrency, amount=amount)
        r = requests.get(url)
        if r.status_code == 200:
            jsonResponse = r.json()
            if targetCurrency.upper() in jsonResponse['rates']:
                message.reply('{amount} {baseCurrency} = {convertedAmount} {targetCurrency}'.format(amount=amount, baseCurrency=baseCurrency.upper(), convertedAmount=jsonResponse['rates'][targetCurrency.upper()], targetCurrency=targetCurrency.upper()))
            else:
                print('Error converting currency.')

rates.__doc__ = 'Displays conversion rate between two currencies'
