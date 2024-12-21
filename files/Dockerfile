FROM alpine:latest

RUN apk update
RUN apk add go

COPY main.go main.go
RUN go build -o /bin/app main.go
RUN mkdir data

ENTRYPOINT ["/bin/app"]
