FROM golang:1.23-alpine3.19 AS builder

COPY . /github.com/mirshodNasilloyev/tg-bot-youtube-go/
WORKDIR /github.com/mirshodNasilloyev/tg-bot-youtube-go/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/mirshodNasilloyev/tg-bot-youtube-go/bin/bot/ .
COPY --from=0 /github.com/mirshodNasilloyev/tg-bot-youtube-go/configs/ configs/

EXPOSE 80

CMD ["./bot"]


