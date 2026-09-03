FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o simulator .


FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/simulator .

CMD ["./simulator"]