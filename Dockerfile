FROM golang:1.15 as build

WORKDIR /app
COPY . .

RUN go build
FROM alpine as app

WORKDIR /app
COPY --from=build /app/spotify_share .

RUN chmod +x spotify_share
