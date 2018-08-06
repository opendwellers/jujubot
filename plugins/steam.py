
import re
import requests

from mmpy_bot.bot import listen_to
from mmpy_bot.bot import respond_to

@respond_to('^steam$', re.IGNORECASE)
def steam_status(message):
        url = 'https://steamgaug.es/api/v2'
        r = requests.get(url).json()
        # print(r.json())
        store = "Up" if r['SteamStore']['online'] is 1 else "*Down*"
        community = "Up" if r['SteamCommunity']['online'] is 1 else "*Down*"
        user = "Up" if r['ISteamUser']['online'] is 1 else "*Down*"
        dota2 = "Up" if r['ISteamGameCoordinator']['570']['online'] is 1 else "*Down*"
        dota2_comments = "Players searching: " + str(r['ISteamGameCoordinator']['570']['stats']['players_searching'])

        reply = """### Steam server status\n

| Service | Status | Comments |
|:---|:---|:---|
|Store|{store}|{store_comments}   |
|Community|{community}|{community_comments}   |
|User api|{user}|{user_comments}   |
|Dota 2 Game Coordinator|{dota2}|{dota2_comments}   |

                    """.format(store=store, store_comments="", community=community, community_comments="", user=user, user_comments="",  dota2=dota2, dota2_comments=dota2_comments)
        message.send(reply)

steam_status.__doc__ = "Get latest steam servers status"
