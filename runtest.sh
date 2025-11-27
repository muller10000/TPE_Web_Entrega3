#!/bin/bash

# Script de Automatización Full Docker (Entrega 6 - Autogeneración)
# Autor: Matías Muller

set -e

echo "============================================"
echo "🚀 INICIANDO CONSTRUCCIÓN Y EJECUCIÓN (Autogeneración de código)"
echo "============================================"

# 1. Generación de Código SQLC (Opcional, si tienes sqlc local)
echo ""
echo "🔨 1. Ejecutando generación de código SQLC (Si la herramienta está disponible localmente)..."
if command -v sqlc &> /dev/null; then
    sqlc generate
fi
# NOTA: La generación de TEMPL ahora ocurre EXCLUSIVAMENTE dentro del Dockerfile.

# 2. Limpieza
echo ""
echo "🧹 2. Limpiando entorno Docker anterior..."
docker compose down -v

# 3. Construcción de la Imagen (Aquí Docker ejecuta sqlc generate y templ generate)
echo ""
echo "🐳 3. Construyendo imagen Docker (Generando vistas dentro del contenedor)..."
docker compose build --no-cache

# 4. Levantamiento
echo ""
echo "▶️  4. Levantando servicios en segundo plano..."
docker compose up -d

# 5. Espera
echo ""
echo "⏳ 5. Esperando servicios (5s)..."
sleep 5

# 6. Verificación (Health Check simple en la ruta de Login)
echo ""
echo "🔍 6. Verificando estado en /login..."
HTTP_STATUS=$(curl -o /dev/null -s -w "%{http_code}\n" http://localhost:8080/login)

if [ "$HTTP_STATUS" == "200" ]; then
    echo "✅ Servidor respondiendo correctamente (HTTP 200 OK)."
else
    echo "⚠️  El servidor respondió con estado: $HTTP_STATUS. Revise logs con 'docker compose logs'."
fi

# ---------------------------------------------------
# 🧪 PRUEBAS DE INTEGRACIÓN DEL FLUJO DE AUTENTICACIÓN
# ---------------------------------------------------

TEST_USERNAME="testuser_temp"
TEST_PASSWORD="password123"

echo ""
echo "==================================================="
echo "🧪 PRUEBA DE AUTENTICACIÓN (Simulación de Usuario)"
echo "==================================================="

# 1. REGISTRO (Es necesario para que exista el usuario en la BD)
echo "-> 1. Intentando registrar usuario..."
REGISTER_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/register \
  -d "username=${TEST_USERNAME}&password=${TEST_PASSWORD}")

if [ "$REGISTER_STATUS" == "303" ]; then
    echo "   ✅ Registro exitoso (Status 303 Redirect)."
elif [ "$REGISTER_STATUS" == "200" ]; then
    echo "   ⚠️ Usuario ya existía o fallo de validación (Status 200 OK). Continuando..."
else
    echo "   ❌ ERROR CRÍTICO: Registro fallido. Status: $REGISTER_STATUS"
    exit 1
fi

# 2. INICIO DE SESIÓN Y CAPTURA DE COOKIE
echo "-> 2. Iniciando sesión y capturando la cookie..."
# -c guarda la cookie en el archivo sesion_data
# -D guarda los headers en headers_output
curl -s -c sesion_data -D headers_output -X POST http://localhost:8080/signin \
  -d "username=${TEST_USERNAME}&password=${TEST_PASSWORD}" > /dev/null

LOGIN_STATUS=$(grep 'HTTP/' headers_output | tail -1 | awk '{print $2}')
SESSION_TOKEN=$(grep 'session_token' sesion_data | awk '{print $NF}')

# Limpieza de archivos temporales (no queremos que se queden en el disco)
rm headers_output
rm sesion_data

if [ "$LOGIN_STATUS" == "303" ] && [[ "$SESSION_TOKEN" != "" ]]; then
    echo "   ✅ Login exitoso. Token de sesión capturado."
else
    echo "   ❌ LOGIN FALLIDO. Status del Login: $LOGIN_STATUS"
    exit 1
fi

# 3. ACCESO A RUTA PROTEGIDA (Home)
echo "-> 3. Accediendo a ruta protegida (Home) con la cookie..."
# -b envía la cookie capturada
PROTECTED_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X GET http://localhost:8080/ \
  -b "session_token=${SESSION_TOKEN}")

if [ "$PROTECTED_STATUS" == "200" ]; then
    echo "   ✅ Acceso protegido exitoso (HTTP 200 OK)."
else
    echo "   ❌ Fallo de acceso. Status: $PROTECTED_STATUS (Debería ser 200)."
    exit 1
fi

# 4. Cierre de Sesión
echo "-> 4. Cerrando sesión..."
LOGOUT_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X GET http://localhost:8080/logout)

if [ "$LOGOUT_STATUS" == "303" ]; then
    echo "   ✅ Logout exitoso (Redirección)."
else
    echo "   ❌ Logout fallido. Status: $LOGOUT_STATUS"
fi

echo ""
echo "============================================"
echo "🎉 LISTO PARA USAR"
echo "============================================"
echo "👉 http://localhost:8080/login"
echo "🛑 Para detener: docker compose down"
echo ""