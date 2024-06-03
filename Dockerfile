FROM golang:1.20-alpine as builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/auction cmd/auction/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/auction /app/auction
COPY --from=builder /app/cmd/auction/.env .
COPY mongo.sh /root/mongo.sh

RUN apk add --no-cache libc6-compat

EXPOSE 8080

ENTRYPOINT ["/app/auction"]

