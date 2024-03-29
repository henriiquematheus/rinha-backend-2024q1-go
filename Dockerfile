# Estágio de compilação
FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .


RUN chmod +x /app/server

CMD ["/app/server"]