FROM golang:1.24-alpine AS build

RUN apk add --no-cache make git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM alpine
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=build /app/bin/jujubot /app
ENV CONFIG_PATH="/config"
ENTRYPOINT [ "/app/jujubot" ]
