# docker build -t goPost .
# docker run -it goPost

FROM golang:1.15.10-alpine3.13 AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

RUN mkdir /pro
ADD ./usePost05.go /pro/
WORKDIR /pro
RUN go get -d -v ./...
RUN go build -o server usePost05.go

FROM alpine:latest

RUN mkdir /pro
COPY --from=builder /pro/server /pro/server
WORKDIR /pro
CMD ["/pro/server"]
