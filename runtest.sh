#!/bin/bash

# Script de Automatización para TPE Web

set -e

echo "============================================"
echo "🚀 INICIANDO CONSTRUCCIÓN Y EJECUCIÓN (Dockerizado)"
echo "============================================"

# 1. Generación de Código SQLC (Opcional, si tienes sqlc local)
# Si no tienes sqlc, asumimos que el código repo/ ya fue commiteado o generado.
if command -v sqlc &> /dev/null; then
    echo "🔨 Generando código DB (SQLC)..."
    sqlc generate
else
    echo "⚠️  SQLC no encontrado localmente. Se usará el código existente en repository/."
fi


# 2. Limpieza del entorno previo
echo ""
echo "🧹 1. Limpiando entorno Docker anterior..."
docker compose down -v

# 3. Construcción de la imagen
echo ""
echo "🐳 2. Construyendo imagen Docker (Generando vistas dentro del contenedor)..."
# Usamos --no-cache para forzar la regeneración de templ dentro del build
docker compose build --no-cache

# 4. Levantamiento de servicios
echo ""
echo "▶️  3. Levantando servicios en segundo plano..."
docker compose up -d

# 5. Espera de arranque
echo ""
echo "⏳ 4. Esperando servicios (5s)..."
sleep 5

# 6. Verificación de Salud (Health Check)
echo ""
echo "🔍 5. Verificando estado..."
HTTP_STATUS=$(curl -o /dev/null -s -w "%{http_code}\n" http://localhost:8080)

if [ "$HTTP_STATUS" == "200" ]; then
    echo "✅ Servidor respondiendo correctamente (HTTP 200 OK)."
else
    echo "⚠️  El servidor respondió con estado: $HTTP_STATUS."
    echo "    Revisa los logs con 'docker compose logs' para ver detalles."
fi

echo ""
echo "============================================"
echo "🎉 LISTO PARA USAR"
echo "============================================"
echo "👉 http://localhost:8080"
echo "🛑 Para detener: docker compose down"
echo ""