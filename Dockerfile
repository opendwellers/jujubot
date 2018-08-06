FROM python:3.6.5

COPY ./ /jujubot/
WORKDIR /jujubot
RUN pip install --no-cache-dir -r requirements.txt

CMD [ "python", "./run.py" ]
