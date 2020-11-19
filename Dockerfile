FROM golang:1.15 as builder

WORKDIR /app
COPY . .

RUN go build -o bin/spotify_share

FROM alpine as app

WORKDIR /app
COPY --from=builder /app/bin/spotify_share /usr/local/bin/

EXPOSE 1323

CMD ["spotify_share"]