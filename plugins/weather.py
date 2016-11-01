import re
from flask import Flask
from flask import Response
from datetime import datetime
from mattermost_bot.bot import listen_to
from mattermost_bot.bot import respond_to
import configparser 
import requests
import time
import json

@respond_to('weather$', re.IGNORECASE)
@respond_to('weather ([A-z_-]+$)', re.IGNORECASE)
def weather(message, location=None):
    text = get_weather(location)
    message.send(text)

@respond_to('weather ([0-9]*$)', re.IGNORECASE)
def weather_id(message, city_id=None):
    text = get_weather_id(city_id)
    message.send(text)

weather.__doc__ = "Get city weather"

app = Flask(__name__)
url = "http://api.openweathermap.org/data/2.5/forecast/daily?q={location}&APPID={key}&cnt=5".format
url_id = "http://api.openweathermap.org/data/2.5/forecast/daily?id={city_id}&APPID={key}&cnt=5".format

def init():
    """Initializes variables"""
    global url
    global url_id
    global config_api_key
    Config = configparser.ConfigParser()
    Config.read("configuration/config")
    Config.sections()

    config_api_key = Config.get('OpenWeather', 'api_key')

init()

def get_weather_id(city_id=None):
    request_url = url_id(city_id=city_id, key=config_api_key)
    r = requests.get(request_url)
    if r.json()['cod'] == "502":
        return "fuck you"
    text = build_response_text(r.json(), r.json()['city']['name'])
    return text

def get_weather(location=None):
    # weird bug where I don't know why len(location) is 1 when no location is passed...
    if (location is not None):
        location.strip(' \t\n\r')
        valid = re.match('^[\w-]+$', location) is not None
        if (valid):
             request_url = url(location=location, key=config_api_key)
        else:
            return 'Bin voyons donc ca existe pas ' + location + ' comme ville'
    else:
        request_url = url(location='Montreal', key=config_api_key)
        location='Montreal'
    request_url = url(location=location, key=config_api_key)
    r = requests.get(request_url)
    text = build_response_text(r.json(), location)
    return text

def get_embedded_icon_url(icon_code, desc):
    """Returns a formatted markdown line to show the weather icon"""
    return '![desc](http://openweathermap.org/img/w/{code}.png "{desc}")'.format(code=icon_code, desc=desc)

def get_day_weather_line(day):
    """Return a markdown formatted line for a weather day"""
    day_weekday = datetime.fromtimestamp(day['dt']).strftime("%A")
    day_month = datetime.fromtimestamp(day['dt']).strftime("%b")
    day_day = datetime.fromtimestamp(day['dt']).strftime("%d")
    day_info_date = """{weekday}, {month}. {day_number}""".ljust(25).format(weekday=day_weekday, month=day_month, day_number=day_day)
    day_desc = day['weather'][0]['description']
    day_desc.ljust(50)
    day_temp_high = int(day['temp']['max'] - 273.15)
    day_temp_low = int(day['temp']['min'] - 273.15)
    day_icon = get_embedded_icon_url(day['weather'][0]['icon'], day_desc)
    return "| {day_info_date_param} | {desc_param}  {day_icon_param} | {day_temp_high_param} °C | {day_temp_low_param} °C  |\n".format(day_info_date_param=day_info_date, desc_param=day_desc, day_icon_param=day_icon, day_temp_high_param=day_temp_high, day_temp_low_param=day_temp_low)

def build_response_text(data, location):
    """Post formatted data to mattermost instance"""
    days = []
    payload_text = """
### Weather in {location} for the next few days

| Day | Description | High | Low |
|:---------------------------|:------------------------------------|:--------|:--------|\n""".format(location=location)
    for day in data['list']:
        payload_text += get_day_weather_line(day)

    return payload_text
