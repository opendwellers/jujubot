import requests
import configparser


def init():
    """Initializes variables"""
    global api_key
    Config = configparser.ConfigParser()
    Config.read("configuration/config")
    Config.sections()

    api_key = Config.get('Yandex', 'api_key')

init()


from mattermost_bot.bot import respond_to

yandex_url = 'https://translate.yandex.net/api/v1.5/tr.json/'
yandex_translate_endpoint = 'translate'
yandex_languages_endpoint = 'getLangs'

@respond_to('^translate (\w{2,3}) (\w{2,3}) (.*)')
def translate(message, fromLanguage, toLanguage, sourceText):
    if (areLanguagesValid(fromLanguage,toLanguage)):
        payload =  {'key': api_key, 'lang': fromLanguage + '-' + toLanguage, 'text': sourceText}
        r = requests.get(yandex_url + yandex_translate_endpoint, params=payload)
        if r.status_code == 200:
            response = r.json()
            if 'text' in response:
                message.reply(response['text'][0])
    else:
        message.reply('Invalid translation request :(')


def areLanguagesValid(fromLanguage, toLanguage):
    payload =  {'key': api_key, 'ui': 'en'}
    r = requests.get(yandex_url + yandex_languages_endpoint, params=payload)
    # TODO Handle error (wrong api key etc.)
    if r.status_code == 200:
        response = r.json()
        return True if fromLanguage in response['langs'] and toLanguage in response['langs'] else False
