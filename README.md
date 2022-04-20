# jujubot
A bot to replace a long-time lost friend on a Mattermost server

## How to use

1. Make sure you have go 1.18 or higher _or_ docker.


1. Create configuration files from the templates:
  ```
  mkdir -p data
  cp config.yaml.template data/config.yaml
  ```
**Note:** Don't forget to update the configuration file to match your environment.

3. Run the bot:
`make`

## Structure

The entrypoint is in main.go.
The pkg directory contains the config definition and generic functions which are useful for commands processing.

## TODO

* ~~Docker image~~
* Helm chart
* ~~Dota mmr command~~
* Reproduce huel cheeky behavior (we're getting there)
