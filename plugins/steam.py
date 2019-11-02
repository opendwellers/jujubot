
import re
import requests
import cloudscraper

from mmpy_bot.bot import listen_to
from mmpy_bot.bot import respond_to

@respond_to('^steam$', re.IGNORECASE)
def steam_status(message):
        scraper = cloudscraper.create_scraper()
        r = scraper.get('https://crowbar.steamstat.us/Barney')
        r = r.json()
        store = "Up" if r['services']['store']['status'] == 'good' else "*Down*"
        community = "Up" if r['services']['community']['status'] == 'good' else "*Down*"
        user = "Up" if r['services']['webapi']['status'] == 'good' else "*Down*"
        dota2 = "Up" if r['services']['dota2']['status'] == 'good' else "*Down*"
        dota2_comments = "Load: " + str(r['services']['dota2']['title'])

        reply = """### Steam server status\n
| Service | Status | Comments |
|:---|:---|:---|
|Store|{store}|{store_comments}   |
|Community|{community}|{community_comments}   |
|User api|{user}|{user_comments}   |
|Dota 2 Game Coordinator|{dota2}|{dota2_comments}   |""".format(store=store, store_comments="", community=community, community_comments="", user=user, user_comments="",  dota2=dota2, dota2_comments=dota2_comments)
        message.send(reply)

steam_status.__doc__ = "Get latest steam servers status"
