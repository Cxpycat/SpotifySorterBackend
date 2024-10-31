FROM golang:1.23.2-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go install github.com/air-verse/air@latest

EXPOSE 8080

CMD ["air"]
