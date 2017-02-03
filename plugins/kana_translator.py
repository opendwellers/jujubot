import regex as re
import romkan
import requests
import time

from bs4 import BeautifulSoup
from mattermost_bot.bot import respond_to


@respond_to('hiragana ([\u30A0-\u30FF\u3040-\u309Fー\u4E00-\u9FFF]+)', re.U)
def translateJapanese(message, kanas=None):
    message.reply(romkan.to_roma(kanas))

translateJapanese.__doc__ = "Translate japanese hiragana symbols"


@respond_to('wotd japanese', re.IGNORECASE)
def wotdJapanese(message):
    url ='http://www.coscom.co.jp/learnjapanese101/words2000today/words2000today.html'
    r = requests.get(url)
    r.encoding = 'utf-8'
    if r.status_code == 200:
        soup = BeautifulSoup(r.text, 'html.parser')
        romajiJap = soup.find('p', {'class':'wdjpr'}).getText()
        kanaJap = soup.findAll('p', {'class':'wdjpk'})[0].getText()
        romajiEng = soup.find('p', {'class':'wdeng'}).getText()
        expJap = soup.find('td', {'class':'exjpr'}).getText()
        expKana = soup.find('td', {'class':'exjpk'}).getText()
        expEng = soup.find('td', {'class':'exeng'}).getText()
        kanjiJap = soup.findAll('p', {'class':'wdjpk'})[1].getText()



        currentDate = time.strftime("%d/%m/%Y")
        payload = """
#### Japanese word of the day for %s

# **%s**
*%s* (%s) - %s
%s
%s
%s
        """ % (currentDate, kanjiJap, romajiJap, kanaJap, romajiEng, expJap, expKana, expEng)

        message.send(payload)



    



