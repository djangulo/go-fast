# alpine builds fail to install go packages
FROM golang:1.12.5
LABEL maintainer="Denis Angulo <djal@tuta.io>"
WORKDIR $GOPATH/src/github.com/djangulo/go-fast/
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -p 1 -a -installsuffix cgo -o /webserver ./cmd/webserver/main.go
RUN git clone https://github.com/eficode/wait-for.git ./wait
RUN mv ./wait/wait-for /wait-for

RUN sed -i 's/\r//g' /wait-for
RUN chmod +x /wait-for
