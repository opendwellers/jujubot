import re
import requests

from mattermost_bot.bot import listen_to
from mattermost_bot.bot import respond_to

@respond_to('^rapture index$', re.IGNORECASE)
def rapture_index(message):
        url = 'https://rapture-index-cors-api.appspot.com/'
        r = requests.get(url)
        index = (r.json())['raptureIndexValue']
        message.reply('[current rapture index](http://www.raptureready.com/rap2.html): ' + index)

rapture_index.__doc__ = "Get latest rapture index"
