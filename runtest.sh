#!/bin/bash

# Este script construye, levanta y prueba la aplicación completa.

# Salir inmediatamente si un comando falla
set -e

echo "1. Limpiando entorno anterior (contenedores y volúmenes)..."
docker compose down -v

echo "2. Construyendo la imagen de la API Go (sin caché)..."
docker compose build --no-cache

echo "▶3. Levantando servicios (API + DB) en SEGUNDO PLANO..."
docker compose up -d

echo "4. Esperando 5 segundos para que la API y DB inicien correctamente..."
sleep 5

echo "5. Ejecutando tests CRUD (API) con cURL..."
echo "------------------------------------------------"

echo " Creando 3 películas (POST)..."
curl -s -X POST http://localhost:8080/peliculas \
  -H "Content-Type: application/json" \
  -d '{"title":"Matrix","director":"Wachowski","year":1999,"genre":"Sci-Fi", "rating":"9.0"}'
echo ""
curl -s -X POST http://localhost:8080/peliculas \
  -H "Content-Type: application/json" \
  -d '{"title":"Inception","director":"Christopher Nolan","year":2010,"genre":"Sci-Fi", "rating":"8.8"}'
echo ""
curl -s -X POST http://localhost:8080/peliculas \
  -H "Content-Type: application/json" \
  -d '{"title":"The Godfather","director":"Francis Ford Coppola","year":1972,"genre":"Crime", "rating":"9.2"}'
echo ""

echo " Listando todas las películas (GET)..."
curl -s http://localhost:8080/peliculas
echo ""
echo ""

echo " Actualizando película (ID=2) (PUT)..."
curl -s -X PUT http://localhost:8080/peliculas/2 \
  -H "Content-Type: application/json" \
  -d '{"title":"Inception Updated","director":"Christopher Nolan","year":2010,"genre":"Thriller", "rating":"9.1"}'
echo ""

echo " Obteniendo película por ID (ID=2) (GET)..."
curl -s http://localhost:8080/peliculas/2
echo ""
echo ""

echo " Eliminando película (ID=1) (DELETE)..."
curl -s -X DELETE http://localhost:8080/peliculas/1
echo ""

echo " Listando películas después de la eliminación (GET)..."
curl -s http://localhost:8080/peliculas
echo ""
echo ""

echo "------------------------------------------------"
echo "✅ ¡Pruebas de API completadas!"
echo ""
echo "ℹ️  La aplicación sigue corriendo en segundo plano."
echo "    (Instrucciones de acceso para TP4)"
echo ""
echo "    Puedes acceder al frontend en:"
echo "    ➡️  http://localhost:8080"
echo ""
echo "    Para detener la aplicación, ejecuta: docker compose down"