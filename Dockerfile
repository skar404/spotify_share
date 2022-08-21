FROM golang:1.17-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN apk add --no-cache --update git tzdata ca-certificates && go mod download

COPY . .

RUN go build -o /app/spotify_share

FROM alpine:3.16.2

RUN apk --no-cache add ca-certificates

# use SSL connect in PORD:
COPY ./mongodb /root/mongodb
COPY --from=builder /app/spotify_share /srv/spotify_share

WORKDIR /srv
CMD ["/srv/spotify_share"]