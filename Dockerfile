FROM golang:alpine3.15 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /vault

FROM alpine:3.15.0

WORKDIR /

COPY --from=build /vault /vault

USER guest

CMD ["/vault"]
