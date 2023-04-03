FROM golang:alpine3.17 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod tidy && go mod download
COPY . .

RUN go build -o main . && go build -o migrate ./cmd/migrate/migration.go
#RUN go run cmd/migrate/migration.go

EXPOSE 8080

ENTRYPOINT ./main

#FROM alpine:3.14
#
#WORKDIR /app
#
#COPY --from=build /app/main .
#
#EXPOSE 8080
#
#ENTRYPOINT ./main
#RUN go run cmd/migrate/migration.go