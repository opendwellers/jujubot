import urbandictionary as ud
import re
from mattermost_bot.bot import listen_to
from mattermost_bot.bot import respond_to

@respond_to('^urban ([a-zA-Z\-\_ 0-9\&]*)$', re.IGNORECASE)
def urban_def(message, query):
    defs = ud.define(query)
    if len(defs) > 0:
        message.send(defs[0].definition)
    else:
        message.send('No definition found :\'(')
