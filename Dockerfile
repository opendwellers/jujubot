FROM python:3.9

COPY ./ /jujubot/
WORKDIR /jujubot
RUN pip3 install --no-cache-dir -r requirements.txt

CMD [ "python3", "./run.py" ]
