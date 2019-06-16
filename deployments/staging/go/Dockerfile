# alpine builds fail to install go packages
FROM golang:1.12.5
LABEL maintainer="Denis Angulo <djal@tuta.io>"
WORKDIR $GOPATH/src/github.com/djangulo/go-fast/
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

