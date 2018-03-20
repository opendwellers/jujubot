# jujubot
A bot to replace a long-time lost friend on a Mattermost server

## How to use

1. Install the required dependencies:
`pip install -r requirements.txt`

2. Create configuration files from the templates:
  ```
  cp mattermost_bot_settings.py.template mattermost_bot_settings.py
  cp configuration/config_template configuration/config
  ```
**Note:** Don't forget to update the configuration files to match your environment.

3. Run the bot:
`python run.py`

## Plugins

Plugins are located in the plugins folder (duh!).

## TODO

* Docker image
* ~~Dota mmr command~~
* Reproduce huel cheeky behavior (we're getting there)
