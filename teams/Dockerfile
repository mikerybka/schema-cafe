FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add go
RUN go install github.com/mikerybka/teams@latest

ENTRYPOINT ["/root/go/bin/teams"]
