import re
import requests

from mattermost_bot.bot import listen_to
from mattermost_bot.bot import respond_to

@respond_to('^mmr$', re.IGNORECASE)
@respond_to('^mmr (.*)', re.IGNORECASE)
def hi(message, id=None):
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
            mmr = array['solo_competitive_rank']
            personaname = array['profile']['personaname'] 
            if mmr is not None:
                message.reply('lel ' + personaname + ' is only '+ mmr + ' mmr scrub, git gud')
            else:
                message.reply('unranked pleb or hidden mmr')
        else:
            message.reply('nice fake player id: ' + id)

