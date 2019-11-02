FROM python:3.7-alpine

COPY ./ /jujubot/
WORKDIR /jujubot
RUN apk add --no-cache git build-base libffi libffi-dev openssl openssl-dev && \
    pip install --no-cache-dir -r requirements.txt

CMD [ "python", "./run.py" ]
