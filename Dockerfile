FROM ubuntu:focal

ENV TZ="Asia/Shanghai"

ARG TARGETOS
ARG TARGETARCH

WORKDIR /nezha-telegram-bot
COPY dist/nezha-telegram-bot-${TARGETOS}-${TARGETARCH} ./app

CMD ["./app"]