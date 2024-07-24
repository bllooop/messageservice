FROM golang:1.22-bullseye

RUN go version
ENV $GOPATH=/

WORKDIR go/src/app

COPY . .

RUN go mod download
RUN go build -o messageservice ./cmd/main.go

EXPOSE 8000

CMD ["./messageservice"]