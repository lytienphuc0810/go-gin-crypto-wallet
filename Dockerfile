# syntax=docker/dockerfile:1
FROM golang:1.18.3
WORKDIR /app
COPY . .
COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN go run . migrate
CMD go run .
EXPOSE 8080