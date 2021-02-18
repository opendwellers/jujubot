FROM python:3.9-alpine

COPY ./ /jujubot/
WORKDIR /jujubot
RUN apk add --no-cache cargo rust git build-base libffi libffi-dev openssl openssl-dev && \
    pip3 install --no-cache-dir -r requirements.txt

CMD [ "python3", "./run.py" ]
