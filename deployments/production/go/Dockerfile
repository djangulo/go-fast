# alpine builds fail to install go packages
FROM golang:1.12.5 as builder
LABEL maintainer="Denis Angulo <djal@tuta.io>"
WORKDIR $GOPATH/src/github.com/djangulo/go-fast/
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -p 1 -a -installsuffix cgo -o $GOPATH/bin/webserver ./cmd/webserver/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -p 1 -a -installsuffix cgo -o $GOPATH/bin/entrypoint ./deployments/production/go/entrypoint.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/webserver .
RUN chmod +x ./webserver
COPY --from=builder /go/bin/entrypoint .
RUN chmod +x ./entrypoint
ENTRYPOINT [ "./entrypoint" ]