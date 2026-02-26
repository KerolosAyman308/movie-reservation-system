FROM golang:1.26.0-alpine3.23 as BASE

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN go build -o ./app ./cmd

FROM alpine3.23 as PROD
WORKDIR /prod
COPY --from=BASE /build/app ./app

EXPOSE 8000
CMD ["./app"]