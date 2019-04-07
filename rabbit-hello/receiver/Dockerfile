FROM golang:alpine as builder

RUN apk update && apk add git && apk add ca-certificates
RUN adduser -D -g '' appuser
COPY . $GOPATH/src/github.com/maxsuelmarinho/golang-event-driven-example/rabbit-hello/receiver/
WORKDIR $GOPATH/src/github.com/maxsuelmarinho/golang-event-driven-example/rabbit-hello/receiver/
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/rabbit-hello-receiver

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/rabbit-hello-receiver /go/bin/rabbit-hello-receiver

USER appuser
ENTRYPOINT ["/go/bin/rabbit-hello-receiver"]

