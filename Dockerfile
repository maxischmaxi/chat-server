FROM golang:1.22.3 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN apt-get update && apt-get install -y \
  openssl

RUN openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/C=DE/ST=Bavaria/L=Munich/O=Jeschek/OU=IT/CN=jeschek.dev"

COPY . .
COPY key.pem key.pem
COPY cert.pem cert.pem

RUN go build -o chat-server .

EXPOSE 8080

CMD ["./chat-server"]
