# alpine builds fail to install go packages
FROM golang:1.12.5 as builder
LABEL maintainer="Denis Angulo <djal@tuta.io>"
WORKDIR $GOPATH/src/github.com/djangulo/go-fast/
COPY . .
RUN go get -d -v ./...
COPY ./deployments/go/start /start
RUN sed -i 's/\r//' /start
RUN chmod +x /start
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $GOPATH/bin/go-fast


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/go-fast .
EXPOSE 9000

