FROM golang:1.21.3 AS builder

WORKDIR /go/app

COPY . .

RUN go build -o main .

CMD ["./main"]
