FROM golang:alpine3.17 AS build

WORKDIR /app
RUN go install github.com/cosmtrek/air@latest \
   && go install github.com/swaggo/swag/cmd/swag@latest

COPY ../../go.mod .
COPY ../../go.sum .

RUN go mod tidy && go mod download
COPY ../../ .

EXPOSE 8080

ENTRYPOINT air
