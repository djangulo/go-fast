# alpine builds fail to install go packages
FROM golang:1.12.5 as builder
LABEL maintainer="Denis Angulo <djal@tuta.io>"
WORKDIR $GOPATH/src/github.com/djangulo/go-fast/
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -p 1 -a -installsuffix cgo -o $GOPATH/bin/webserver ./cmd/webserver/main.go

# wait for postgres connection
RUN git clone https://github.com/eficode/wait-for.git /wait-for



FROM alpine:latest

RUN addgroup -S golang \
    && adduser -S -G golang golang

RUN apk update \
    && apk --no-cache add ca-certificates

COPY --from=builder /go/bin/webserver /webserver
RUN chmod +x ./webserver
RUN chown golang /webserver

COPY --from=builder /wait-for/wait-for /wait-for
RUN sed -i 's/\r//g' /wait-for
RUN chmod +x /wait-for
RUN chown golang /wait-for
