#!/bin/bash

# Script de Automatización para TPE Web (Versión Estable)
# Autor: Matías Muller

set -e

echo "============================================"
echo "🚀 INICIANDO CONSTRUCCIÓN Y EJECUCIÓN"
echo "============================================"

# 1. Generación de Código (Requisito de la nueva entrega)
# Esto debe ocurrir ANTES de que Docker intente compilar
echo ""
echo "🔨 1. Generando código Go (Templ y SQLC)..."

# Generar SQLC si está instalado
if command -v sqlc &> /dev/null; then
    echo "   -> Ejecutando sqlc generate..."
    sqlc generate
fi

# Generar Templ (CRÍTICO)
if command -v templ &> /dev/null; then
    echo "   -> Ejecutando templ generate..."
    templ generate
else
    echo "❌ ERROR: 'templ' no encontrado."
    echo "   Es necesario para compilar las vistas."
    echo "   Instálalo con: go install github.com/a-h/templ/cmd/templ@latest"
    exit 1
fi

# 2. Limpieza del entorno previo
echo ""
echo "🧹 2. Limpiando entorno Docker anterior..."
docker compose down -v

# 3. Construcción de la imagen
echo ""
echo "🐳 3. Construyendo imagen Docker..."
# Usamos --no-cache para asegurar que tome los archivos _templ.go recién generados
docker compose build --no-cache

# 4. Levantamiento de servicios
echo ""
echo "▶️  4. Levantando servicios en segundo plano..."
docker compose up -d

# 5. Espera de arranque
echo ""
echo "⏳ 5. Esperando servicios (5s)..."
sleep 5

# 6. Verificación de Salud (Health Check)
# Comprobamos que la página de inicio (SSR) responda correctamente
echo ""
echo "🔍 6. Verificando estado..."
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