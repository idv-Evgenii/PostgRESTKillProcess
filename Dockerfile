FROM golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /postgrest-bin

FROM alpine:3.16.2

WORKDIR /
COPY --from=build /postgrest-bin /postgrest-bin
ENTRYPOINT [ "./postgrest-bin" ]
