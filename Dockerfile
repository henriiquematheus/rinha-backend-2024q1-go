# Estágio de compilação
FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# Estágio final
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .

# Garante que o binário seja executável
RUN chmod +x /app/server

# O comando `CMD` define o comando padrão a ser executado quando o contêiner é iniciado
CMD ["/app/server"]