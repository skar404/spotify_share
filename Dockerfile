FROM golang:1.17 as builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o bin/spotify_share

FROM alpine:3.12

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# use SSL connect in PORD:
COPY ./mongodb ./mongodb
COPY --from=builder /app/bin/spotify_share /usr/local/bin/

EXPOSE 1323

CMD ["spotify_share"]