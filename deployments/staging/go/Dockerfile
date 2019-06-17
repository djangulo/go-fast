# alpine builds fail to install go packages
FROM golang:1.12.5 as builder
LABEL maintainer="Denis Angulo <djal@tuta.io>"
WORKDIR $GOPATH/src/github.com/djangulo/go-fast/
COPY . .
RUN go get -d -v ./...
# CGO_ENABLED=1 due to the sqlite implementation
# remove after migrating to postgres
RUN CGO_ENABLED=1 GOOS=linux go build -p 1 -a -installsuffix cgo -o $GOPATH/bin/webserver ./cmd/webserver/main.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /go/bin/webserver /usr/local/bin/webserver
RUN chmod +x /webserver
COPY --from=builder /go/src/github.com/djangulo/go-fast/deployments/staging/go/start .
RUN sed -i 's/\r//' /start
RUN chmod +x /start

