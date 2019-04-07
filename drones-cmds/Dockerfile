FROM golang:alpine as builder

RUN apk update && apk add git && apk add ca-certificates
RUN adduser -D -g '' appuser
COPY . $GOPATH/src/github.com/maxsuelmarinho/golang-event-driven-example/drones-cmds/
WORKDIR $GOPATH/src/github.com/maxsuelmarinho/golang-event-driven-example/drones-cmds/
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/drones-cmds

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/drones-cmds /go/bin/drones-cmds

USER appuser
ENTRYPOINT ["/go/bin/drones-cmds"]

