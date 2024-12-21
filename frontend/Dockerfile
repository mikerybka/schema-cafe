FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add go
RUN go install github.com/mikerybka/schema-cafe@latest

ENTRYPOINT ["/root/go/bin/schema-cafe"]
