# Usamos la imagen oficial estándar (basada en Debian)
# Esto evita los problemas de red/DNS de Alpine en tu entorno
FROM golang:1.23

WORKDIR /app

# Copiar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código (incluyendo las vistas generadas por templ)
COPY . .

# Compilar el binario
# No necesitamos instalar curl extra, la imagen base ya tiene lo necesario
RUN go build -o peliculas-api .

# Variables de entorno y puerto
ENV PORT=8080
EXPOSE 8080

# Ejecutar
CMD ["./peliculas-api"]