FROM golang:1.25-alpine

RUN apk add --no-cache git bash

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

RUN go install github.com/air-verse/air@latest

COPY ./src .

EXPOSE 8080

RUN mkdir -p tmp

CMD ["air", "--build.cmd", "go build -o ./tmp/main ./main.go", "--build.bin", "./tmp/main"]
