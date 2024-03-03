# Use a imagem base desejada para Go (por exemplo, golang)
FROM golang:1.22

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie apenas os arquivos necessários
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copie o restante dos arquivos
COPY . .


RUN ls -l /app  # Adiciona comando para listar arquivos no diretório
RUN go build -o server . && chmod +x server
RUN ls -l /app  # Adiciona comando para listar arquivos novamente

# Exponha a porta da aplicação
EXPOSE 8080

# Comando para iniciar a aplicação
CMD ["./server"]
