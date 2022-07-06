FROM golang:alpine AS binarybuilder
RUN apk --no-cache --no-progress add \
    gcc git musl-dev
WORKDIR /nezha-telegram-bot
COPY . .
RUN go build -o app -ldflags="-s -w"

FROM alpine:latest

WORKDIR /nezha-telegram-bot
COPY --from=binarybuilder /nezha-telegram-bot/app ./app

CMD ["./app"]