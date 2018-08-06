import urbandictionary as ud
import re
from mmpy_bot.bot import listen_to
from mmpy_bot.bot import respond_to

@respond_to('^urban ([a-zA-Z\-\_ 0-9\&]*)$', re.IGNORECASE)
def urban_def(message, query):
    defs = ud.define(query)
    if len(defs) > 0:
        definition = defs[0]
        if definition.example:
            payload = """{0}\r
*\r
{1}\r
*""".format(definition.definition, definition.example)
        else:
            payload = "{0}".format(definition.definition)

        message.send(payload)
    else:
        message.send('No definition found :\'(')
