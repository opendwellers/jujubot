import regex as re
import romkan
from jNlp.jTokenize import jTokenize
from mattermost_bot.bot import listen_to


@listen_to('testtest ([\u30A0-\u30FF\u3040-\u309Fー\u4E00-\u9FFF]+)', re.U)
def translateJapanese(message, kanas=None):
    print(kanas)
    print(romkan.to_roma(kanas))
    message.reply(romkan.to_roma(kanas))


translateJapanese.__doc__ = "Translate japanese kanas on the fly"



