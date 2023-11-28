FROM golang:alpine3.17 AS build

WORKDIR /app

ENV APP_ENV2=production
RUN apk add --no-cache bash && apk add make

COPY go.mod .
COPY go.sum .

RUN go mod tidy && go mod download
COPY . .

RUN go build -o main . && go build -o Migrate ./Cmd/Migrate/migration.go \
    && go build -o elastic ./Cmd/Migrate/esSetup.go && go build -o bulk ./Cmd/Migrate/bulkInsert.go

EXPOSE 8080

ENTRYPOINT ["./main"]
