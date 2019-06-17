# alpine builds fail to install go packages
FROM golang:1.12.5 as builder
LABEL maintainer="Denis Angulo <djal@tuta.io>"
WORKDIR $GOPATH/src/github.com/djangulo/go-fast/
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -x -p 1 -a -installsuffix cgo -o $GOPATH/bin/webserver


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /go/bin/webserver .
RUN chmod +x /webserver
COPY --from=builder /go/src/github.com/djangulo/go-fast/deployments/production/go/start .
RUN sed -i 's/\r//' /start
RUN chmod +x /start

