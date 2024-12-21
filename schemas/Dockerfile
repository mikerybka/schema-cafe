FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add go
RUN go install github.com/mikerybka/schemas@latest

ENTRYPOINT ["/root/go/bin/schemas"]
