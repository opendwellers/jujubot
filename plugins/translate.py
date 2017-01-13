from microsofttranslator import Translator
from microsofttranslator import TranslateApiException
from mattermost_bot.bot import respond_to
import configparser 
import re

@respond_to('translate ([a-z]{5}) ([a-z]{5}) (.*)', re.IGNORECASE)
def translate(message, lang1='fr', lang2='ht', string='criminel'):
    try:
        # print(string)
        # print(lang1)
        # print(lang2)
        translation = translator.translate(text=string, from_lang=lang1, to_lang=lang2)
    except TranslateApiException:
        message.reply('microsoft api sucks azure zzz')
        return
    message.reply(translation)


translate.__doc__ = "translate!"

def init():
    global translator
    Config = configparser.ConfigParser()
    Config.read("configuration/config")
    Config.sections()
    clientId = Config.get('Microsoft', 'clientId')
    clientSecret = Config.get('Microsoft', 'clientSecret')
    # print(clientId)
    # print(clientSecret)

    translator = Translator(clientId, clientSecret)

init()


