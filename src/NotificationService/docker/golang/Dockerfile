FROM golang:alpine3.17 AS build

WORKDIR /app2

ENV APP_ENV3=development
RUN apk add --no-cache bash && apk add make

COPY ../../go.mod .
COPY ../../go.sum .

RUN go mod tidy && go mod download
COPY ../../. .
RUN go build -o main .

RUN chmod +x main

ENTRYPOINT ["./main"]
