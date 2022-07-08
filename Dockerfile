FROM ubuntu:focal

ENV TZ="Asia/Shanghai"

ARG TARGETOS
ARG TARGETARCH

RUN export DEBIAN_FRONTEND="noninteractive" && \
    apt update && apt install -y ca-certificates tzdata && \
    update-ca-certificates && \
    ln -fs /usr/share/zoneinfo/$TZ /etc/localtime && \
    dpkg-reconfigure tzdata

WORKDIR /nezha-telegram-bot
COPY dist/nezha-telegram-bot-${TARGETOS}-${TARGETARCH} ./app

CMD ["./app"]