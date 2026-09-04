FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server
RUN go build -o seed-admin ./cmd/seed-admin
RUN cp .env.example .env

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server .
COPY --from=builder /app/seed-admin .

COPY migrations ./migrations

EXPOSE 8080

CMD ["./server"]