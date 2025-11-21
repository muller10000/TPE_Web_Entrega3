Trabajo Práctico Especial - Programación Web (Entrega 5)

📋 Descripción del Proyecto

Este repositorio contiene la 5ta entrega del Trabajo Práctico Especial, centrada en la refactorización arquitectónica hacia Server-Side Rendering (SSR).

El objetivo principal de esta etapa fue eliminar la dependencia de JavaScript en el cliente (SPA/AJAX) y migrar toda la lógica de presentación al servidor utilizando Go y la librería de templating Templ.

-------------------------------------------------------

🚀 Evolución de la Arquitectura

En la entrega anterior (TP4), el renderizado se realizaba en el cliente mediante JavaScript manipulando el DOM. La interacción dependía de llamadas asíncronas (AJAX/fetch), los datos viajaban en formato JSON y el estado de la aplicación era efímero en el cliente. Existía una alta dependencia de archivos JavaScript complejos.

En la entrega actual (TP5), el renderizado ocurre completamente en el servidor utilizando Go y Templ. La interacción se basa en estándares web clásicos como formularios HTML y redirecciones (patrón PRG). Los datos viajan directamente como HTML listo para mostrar. El estado es persistente, residiendo en la URL y la base de datos. Se ha logrado una dependencia nula de JavaScript en el cliente (0% JS).

-------------------------------------------------------- 

⚙️ Instrucciones de Ejecución (Todo en Uno)

Para facilitar la corrección y el despliegue, se ha automatizado todo el ciclo de vida del proyecto (generación de código, construcción de imagen y levantamiento de servicios) en un único script.

Requisitos Previos:

- Docker y Docker Compose instalados.

- (Opcional) go, templ y sqlc si se desea ejecutar localmente sin Docker.

-----------------------------------------------------------

▶️ Paso a Paso

1) Clonar el repositorio y ubicarse en la rama correspondiente.

2) Crear archivo de entorno: Copiar el contenido de .env.example en un nuevo archivo llamado .env.

3) Ejecutar el script maestro:

Abre una terminal en la raíz del proyecto y ejecuta los siguientes comandos:

chmod +x runtest.sh
./runtest.sh

¿Qué hace este script?

- Generación de Código: Ejecuta sqlc generate y templ generate para asegurar que los modelos de base de datos y las vistas HTML estén actualizados y compilados a Go antes de construir la aplicación.

- Limpieza: Ejecuta docker compose down -v para garantizar un entorno de pruebas limpio, eliminando contenedores y volúmenes de ejecuciones anteriores (la base de datos inicia vacía).

- Construcción: Crea la imagen de Docker optimizada utilizando el Dockerfile del proyecto.

- Despliegue: Levanta los servicios de base de datos y la aplicación en el puerto 8080 en segundo plano.

- Verificación: Realiza un health-check automático mediante curl para confirmar que el servidor SSR está respondiendo con un código HTTP 200 OK.

------------------------------------------------------------

🌐 Acceso a la Aplicación

Una vez que el script finalice y muestre el mensaje de éxito, la aplicación estará disponible en tu navegador web.

Dirección de acceso: http://localhost:8080

Desde allí podrá realizar las siguientes acciones:

Listar: Ver la tabla de películas renderizada directamente desde el servidor.

Crear: Usar el formulario para agregar nuevas películas. Al enviar, el servidor procesará los datos y redirigirá a la lista actualizada (patrón Post-Redirect-Get).

Eliminar: Borrar registros mediante los botones de eliminar, que funcionan como formularios POST embebidos.
