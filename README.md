# Trabajo Práctico Especial - Programación Web

# Autor: Matías Muller

# Proyecto: Películas 3ra Entrega

# Dominio de la aplicación

El dominio elegido es Películas.
Cada película cuenta con los siguientes atributos:

id → identificador único

title → título de la película

director → director de la película

year → año de estreno

genre → género de la película

rating → calificación de la película

# Requisitos previos

-Linux 

-Go 1.22 o superior

-Docker instalado (Para levantar contenedores)

# Ejecución del proyecto

1) Clonar este repositorio.

2) Copiar el contenido del archivo ".env.example" en un nuevo archivo ".env" reemplazando con las credenciales reales.
En mi caso:
DB_NAME=peliculas_tp3
DB_USER=peliculas_user
DB_PASSWORD=peliculas_pass

# Dar permisos de ejecución al script
chmod +x runtest.sh

3) Ejecucion de script en consola linux

./runtest.sh

- Construye la app (build del binario con Docker).
- Levanta los contenedores (DB + API).
- Ejecuta los tests CRUD automáticamente.

# 💻 Acceso al Frontend (TP4)

La aplicación (API + Frontend) se sirve desde el mismo contenedor Go.

Una vez que el script runtest.sh termine (o si levantas los servicios manualmente con docker compose up), la aplicación quedará corriendo en segundo plano.

Para acceder a la aplicación web, abre tu navegador y visita:

http://localhost:8080

Podrás ver el formulario, agregar películas a la base de datos, ver la lista y eliminarlas, todo interactuando con la API de Go.

Para detener la aplicación:

docker compose down
