# Usamos la imagen oficial estándar de Go
FROM golang:1.23

WORKDIR /app

# 1. Instalar la herramienta templ dentro del contenedor
RUN go install github.com/a-h/templ/cmd/templ@latest

# 2. Copiar dependencias
COPY go.mod go.sum ./
RUN go mod download

# 3. Copiar todo el código fuente (archivos .go, .templ, sql, etc.)
COPY . .

# 4. Generar el código templ DENTRO del build
# Esto crea los archivos _templ.go antes de compilar
RUN templ generate

# 5. Compilar el binario
RUN go build -o peliculas-api .

# Variables de entorno y puerto
ENV PORT=8080
EXPOSE 8080

# Ejecutar
CMD ["./peliculas-api"]