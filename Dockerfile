FROM golang:alpine as build

COPY ./ /app/
WORKDIR /app
RUN apk add --no-cache make git
RUN make build

FROM alpine
WORKDIR /app
COPY --from=build /app/bin/jujubot /app
ENV CONFIG_PATH /config
ENTRYPOINT [ "/app/jujubot" ]
