FROM golang:1.15-alpine as builder

WORKDIR /app
COPY . .

RUN go build -o bin/spotify_share

FROM alpine

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/bin/spotify_share /usr/local/bin/

EXPOSE 1323

CMD ["spotify_share"]