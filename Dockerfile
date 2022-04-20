FROM golang:alpine as build

RUN apk add --no-cache make git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM alpine
WORKDIR /app
COPY --from=build /app/bin/jujubot /app
ENV CONFIG_PATH /config
ENTRYPOINT [ "/app/jujubot" ]
