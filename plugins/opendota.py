import re
import requests

from mattermost_bot.bot import listen_to
from mattermost_bot.bot import respond_to

@respond_to('^mmr$', re.IGNORECASE)
@respond_to('^mmr ([0-9]+)$', re.IGNORECASE)
def mmr(message, id=None):
    if id is None:
        url = 'https://api.opendota.com/api/players/{player_id}'.format(player_id=12088460)
        r = requests.get(url)
        mmr = (r.json())['solo_competitive_rank']
        message.reply('lel j''suis rendu ' + str(mmr) + ' ez gaem road to 4k')
    else:
        if (len(id) < 10 and id.isdigit()):
            url = 'https://api.opendota.com/api/players/{player_id}'.format(player_id=id)
            r = requests.get(url)
            array = r.json()
            print(id)
            if array['profile'] is not None or id is 53515020:
                if id is not 53515020:
                    mmr = array['solo_competitive_rank']
                else:
                    mmr = 9000
                personaname = array['profile']['personaname'] 
                if mmr is not None and mmr < 4500:
                    message.reply('lel ' + personaname + ' is only '+ mmr + ' mmr scrub, git gud')
                elif mmr is not None and mmr >= 4500:
                    message.reply('lel ' + personaname + ' is '+ mmr + ' mmr what an amazing player')
                else:
                    message.reply('unranked pleb or hidden mmr')
            else:
                message.reply('rofl ' + id + ' existe meme pas zzz')
        else:
            message.reply('nice fake player id: ' + id)


mmr.__doc__ = "Returns mmr for a given unique id (or mine if none specified) i.e https://www.opendota.com/players/*53515020*"
