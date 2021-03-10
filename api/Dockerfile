FROM golang:1.15-alpine AS GO_BUILD

COPY . /app
WORKDIR /app
RUN go build -o /go/bin/app/server

FROM alpine:latest
WORKDIR /app
COPY --from=GO_BUILD /go/bin/app .
EXPOSE 8080
ENTRYPOINT ./server
