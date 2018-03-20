FROM python:3.6.4

RUN git clone https://github.com/opendwellers/jujubot && \
    cd jujubot && \
    pip install --no-cache-dir -r requirements.txt && \
    pip install --no-cache-dir git+https://github.com/attzonko/mattermost_bot.git

CMD [ "python", "/jujubot/run.py" ]

