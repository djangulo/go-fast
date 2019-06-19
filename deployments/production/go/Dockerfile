# alpine builds fail to install go packages
FROM golang:1.12.5 as builder
LABEL maintainer="Denis Angulo <djal@tuta.io>"
WORKDIR $GOPATH/src/github.com/djangulo/go-fast/
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -p 1 -a -installsuffix cgo -o $GOPATH/bin/webserver ./cmd/webserver/main.go


FROM alpine:latest

RUN addgroup -S golang \
    && adduser -S -G golang golang

RUN apk update \
    && apk --no-cache add postgresql-dev ca-certificates postgresql-client
COPY --from=builder /go/bin/webserver /webserver
RUN chmod +x ./webserver
RUN chown golang /webserver

COPY --from=builder /go/src/github.com/djangulo/go-fast/deployments/production/go/entrypoint /
RUN sed -i 's/\r//g' /entrypoint
RUN chmod +x /entrypoint
RUN chown golang /entrypoint
ENTRYPOINT [ "/entrypoint" ]