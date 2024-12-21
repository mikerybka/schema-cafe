FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add go
RUN go install github.com/mikerybka/files@latest

ENTRYPOINT ["/root/go/bin/files"]
