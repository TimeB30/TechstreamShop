FROM golang:1.24-alpine AS builder

RUN apk add --no-cache build-base librdkafka-dev

WORKDIR /techstreamshop

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /techstreamshop/telegrambot ./cmd/main.go


FROM alpine:latest

RUN apk add --no-cache librdkafka ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /techstreamshop/telegrambot .

COPY ./config/local.yaml .

EXPOSE 8080

CMD ["./telegrambot"]

