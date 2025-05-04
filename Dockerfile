FROM golang:1.24-alpine

WORKDIR /app


COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o 1337b04rd ./cmd/main.go

EXPOSE 8080

CMD ["./1337b04rd"]
