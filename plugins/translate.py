from microsofttranslator import Translator
from microsofttranslator import TranslateApiException
from microsofttranslator import ArgumentOutOfRangeException
from mattermost_bot.bot import respond_to
import configparser 
import re

@respond_to('translate ([a-z]{2,5}) ([a-z]{2,5}) (.*)', re.IGNORECASE)
def translate(message, lang1='fr', lang2='ht', string='criminel'):
    try:
        translation = translator.translate(text=string, from_lang=lang1, to_lang=lang2)
    except (TranslateApiException, ArgumentOutOfRangeException):
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


