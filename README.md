Trabajo Práctico Especial - Programación Web (Entrega 6)

📋 Descripción del Proyecto

Este repositorio contiene la 6ta entrega del Trabajo Práctico Especial, enfocada en la evolución de la aplicación hacia una arquitectura de Interfaz Reactiva con HTMX.

El objetivo principal de esta etapa fue transformar la aplicación SSR (Server-Side Rendering) tradicional de la entrega anterior en una experiencia similar a una SPA (Single Page Application), eliminando las recargas completas de página al crear o eliminar entidades, pero manteniendo la simplicidad del backend en Go.

🚀 Evolución de la Arquitectura

En la entrega anterior (TP5), la aplicación dependía de recargas completas de página para cada interacción (patrón Post-Redirect-Get). Cada vez que se creaba o eliminaba una película, el navegador debía volver a cargar todos los recursos (CSS, Scripts, Layout).

En la entrega actual (TP6), se ha integrado la librería HTMX. Ahora, las interacciones ocurren mediante AJAX transparente. El servidor ya no responde con redirecciones, sino con fragmentos de HTML específicos (componentes) que actualizan solo las partes necesarias del DOM (la lista de películas), logrando una experiencia de usuario fluida e instantánea sin parpadeos.

⚙️ Instrucciones de Ejecución (Todo en Uno)

Para facilitar la corrección, se ha automatizado todo el ciclo de vida del proyecto en un único script.

Requisitos Previos

Docker y Docker Compose instalados.

(Opcional) go, templ y sqlc si se desea ejecutar localmente sin contenedores.

▶️ Paso a Paso para Clonar y Ejecutar

Clonar el repositorio:
Asegúrese de ubicarse en la rama entrega6 tras clonar.

git clone <URL_DEL_REPOSITORIO>
cd <NOMBRE_DEL_PROYECTO>
git checkout entrega6

Crear archivo de entorno:
Copie el contenido de .env.example en un nuevo archivo llamado .env en la raíz del proyecto.

Ejecutar el script maestro:
Desde la terminal en la raíz del proyecto, ejecute:

chmod +x runtest.sh
./runtest.sh

¿Qué realiza este script?

Generación de Código: Ejecuta sqlc generate y templ generate para asegurar que los binarios coincidan con las últimas definiciones de vistas y base de datos.

Limpieza Profunda: Ejecuta docker compose down -v para eliminar contenedores previos y volúmenes, garantizando que la base de datos inicie desde cero con el esquema limpio.

Construcción: Crea la imagen de Docker optimizada.

Despliegue: Levanta los servicios (API + DB) en segundo plano.

Validación: Realiza un health-check y una prueba de creación automática para verificar que el sistema responde correctamente.

🌐 Acceso a la Aplicación

Una vez que el script finalice exitosamente, la aplicación estará disponible en:

👉 http://localhost:8080

Prueba de Interactividad (Validación de HTMX)

Para comprobar que la implementación de HTMX es correcta:

Abra las herramientas de desarrollador del navegador (F12) y vaya a la pestaña "Network" (Red).

Complete el formulario y haga clic en "Agregar Película".

Verá que se realiza una petición POST, pero la página no se recarga (el icono de carga del navegador no gira).

La respuesta de esa petición será únicamente el fragmento HTML de la lista de películas, no la página completa.

Lo mismo ocurrirá al hacer clic en "Eliminar".