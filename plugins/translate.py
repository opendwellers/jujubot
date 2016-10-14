from microsofttranslator import Translator
from mattermost_bot.bot import respond_to
import configparser 
import re

@respond_to('translate ([a-z]{2}) ([a-z]{2}) (.*)', re.IGNORECASE)
def translate(message, lang1='fr', lang2='ht', string='criminel'):
    translation = translator.translate(text=string, from_lang=lang1, to_lang=lang2)
    message.reply(translation)


translate.__doc__ = "translate!"

def init():
    global translator
    Config = configparser.ConfigParser()
    Config.read("configuration/config")
    Config.sections()
    clientId = Config.get('Microsoft', 'clientId')
    clientSecret = Config.get('Microsoft', 'clientSecret')
    print(clientId)
    print(clientSecret)

    translator = Translator(clientId, clientSecret)

init()


