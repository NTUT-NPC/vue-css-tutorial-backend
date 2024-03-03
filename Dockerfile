FROM golang:1.22-alpine as builder
WORKDIR /usr/src/app
COPY . .
RUN go build -o /usr/local/bin/app

FROM alpine
WORKDIR /srv
COPY --from=builder /usr/local/bin/app app
ENTRYPOINT ["/srv/app"]
