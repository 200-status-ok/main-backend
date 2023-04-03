FROM golang:alpine3.17 AS build

WORKDIR /app

ENV APP_ENV=development
RUN apk add --no-cache bash && apk add make

RUN go install github.com/cosmtrek/air@latest \
   && go install github.com/swaggo/swag/cmd/swag@latest

COPY ../../go.mod .
COPY ../../go.sum .

RUN go mod tidy && go mod download
COPY ../../ .

EXPOSE 8080

ENTRYPOINT air
