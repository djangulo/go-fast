# alpine builds fail to install go packages
FROM golang:1.12.5
LABEL maintainer="Denis Angulo <djal@tuta.io>"
WORKDIR $GOPATH/src/github.com/djangulo/go-fast/
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN apk update \
    && apk --no-cache add postgresql-dev ca-certificates postgresql-client

RUN mv ./deployments/production/go/entrypoint /entrypoint
RUN sed 's/\r//g' /entrypoint
RUN chmod +x /entrypoint
