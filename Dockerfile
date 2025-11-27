# Usamos la imagen oficial estándar de Go 1.23
FROM golang:1.23

WORKDIR /app

# 1. Instalar herramientas necesarias (templ y sqlc)
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# 2. Copiar dependencias
COPY go.mod go.sum ./
RUN go mod download

# 3. Copiar todo el código fuente
COPY . .

# 4. Generar código (Orden importante: SQLC primero, luego Templ)
RUN sqlc generate
RUN templ generate

# 5. Compilar el binario
RUN go build -o peliculas-api .

# Variables de entorno y puerto
ENV PORT=8080
EXPOSE 8080

# Ejecutar
CMD ["./peliculas-api"]