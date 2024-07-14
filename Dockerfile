# Fase 1: Costruzione del binario
FROM golang:1.22 AS builder

# Imposta la directory di lavoro
WORKDIR /app

# Copia i file di moduli e scarica le dipendenze
COPY go.mod go.sum ./
RUN go mod download

# Copia il codice sorgente
COPY . .

# Assicura che tutte le dipendenze siano aggiornate
RUN go get -d github.com/openzipkin/zipkin-go/propagation/b3@v0.4.3
RUN go get -d google.golang.org/grpc/metadata@latest

# Esegui go mod tidy per garantire che tutte le dipendenze siano corrette
RUN go mod tidy

# Compila l'applicazione
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp cmd/myapp/main.go

# Fase 2: Creazione dell'immagine leggera per l'esecuzione
FROM alpine:latest

# Installa le dipendenze necessarie
RUN apk --no-cache add ca-certificates

# Installa dockerize
RUN apk add --no-cache wget \
    && wget -O dockerize-linux-amd64-v0.6.1.tar.gz https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-v0.6.1.tar.gz \
    && rm dockerize-linux-amd64-v0.6.1.tar.gz \
    && chmod +x /usr/local/bin/dockerize

# Imposta la directory di lavoro
WORKDIR /root/

# Copia l'applicazione compilata dalla fase di compilazione
COPY --from=builder /app/myapp .

# Imposta le variabili d'ambiente
ENV MONGO_URI=mongodb://mongodb:27017
ENV MONGO_DATABASE=myapp
ENV SERVICE_NAME=myapp_service

# Espone la porta dell'applicazione
EXPOSE 8080

# Comando di avvio dell'applicazione con dockerize
CMD ["dockerize", "-wait", "tcp://mongodb:27017", "-timeout", "30s", "./myapp"]
