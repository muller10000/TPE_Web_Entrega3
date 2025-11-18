/*
  PASO 1: CONECTAR APP.JS
  Este archivo se conecta al HTML porque en `index.html` tengo
  <script src="/static/app.js"></script> al final del <body>.
  
  Definimos la URL de nuestra API. Usamos una ruta relativa porque
  el backend Go sirve este JS y la API desde el mismo puerto (localhost:8080).
*/
const API_URL = "/peliculas";

// Elementos del DOM que usaremos varias veces
const listaPeliculas = document.getElementById("lista-peliculas");
const statsTotal = document.getElementById("stats-total");
const formPelicula = document.getElementById("form-pelicula");

/*
  PASO 2: FETCH INICIAL (GET)
  Usamos el evento 'DOMContentLoaded' para asegurarnos de que el HTML
  esté completamente cargado antes de intentar manipularlo
  En ese momento, llamamos a la función para cargar las películas
*/

document.addEventListener("DOMContentLoaded", cargarPeliculas);

// Carga la lista de películas desde la API (GET) y las renderiza en el DOM.

async function cargarPeliculas() {
    try {
        const res = await fetch(API_URL);
        if (!res.ok) throw new Error("Error al cargar la lista");

        // Gracias al 'MovieResponse' en Go, el JSON viene limpio.
        const peliculas = await res.json(); 

        // Actualizar estadísticas
        statsTotal.innerText = peliculas.length;

        // Limpiar la lista (quitar el "Cargando...")
        listaPeliculas.innerHTML = "";

        if (peliculas.length === 0) {
            listaPeliculas.innerHTML = "<p>No hay películas en el catálogo. ¡Agrega una!</p>";
            return;
        }

        /*
          PASO 3: RENDERIZAR DINÁMICAMENTE (GET)
          Recorremos el array de películas y creamos el HTML
          para cada una.
        */
        peliculas.forEach(peli => {
            const div = document.createElement("div");
            div.classList.add("film"); // Para que 'estilos.css' le aplique el formato

            // Formateamos la fecha que viene de Go (ej: "2025-11-17T15:00:00Z")
            const fechaCreacion = new Date(peli.created_at).toLocaleDateString("es-AR");

            div.innerHTML = `
                <!-- <h3>${peli.title} (${peli.year})</h3> -->
                <h3>${peli.title}</h3>
                <p><strong>Director:</strong> ${peli.director}</p>
                <p><strong>Año:</strong> ${peli.year}</p>
                <p><strong>Género:</strong> ${peli.genre}</p>
                <p><strong>Calificación:</strong> <span class="rating-badge">${peli.rating} / 10</span></p>
                <p class="created-at"><small>Agregada el: ${fechaCreacion}</small></p>
                
                <!-- 
                  PASO 4 (Preparación): Botón de borrado
                  Añadimos la clase 'delete-btn' y 'data-id'
                  para la delegación de eventos.
                -->
                <button class="delete-btn" data-id="${peli.id}">Eliminar</button>
            `;
            listaPeliculas.appendChild(div);
        });

    } catch (error) {
        console.error("Error en cargarPeliculas:", error);
        listaPeliculas.innerHTML = "<p style='color: red;'>Error al cargar películas.</p>";
    }
}

/*
  PASO 3: CAPTURAR FORMULARIO (POST)
  Escuchamos el evento 'submit' del formulario.
*/
formPelicula.addEventListener("submit", async (e) => {
    // Evitamos que el formulario recargue la página
    e.preventDefault(); 

    // Tomamos los datos de los inputs
    const title = document.getElementById("title").value.trim();
    const director = document.getElementById("director").value.trim();
    const year = parseInt(document.getElementById("year").value);
    const genre = document.getElementById("genre").value.trim();
    const ratingInput = document.getElementById("rating").value; // Lo tomamos como string

    if (!title || !director || !year || !genre || !ratingInput) {
        alert("Por favor, completa todos los campos.");
        return;
    }

    try {
        const response = await fetch(API_URL, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                title,
                director,
                year,
                genre,
                // Tu backend (CreateMovieRequest) espera el rating como *string
                rating: ratingInput 
            })
        });

        if (response.status === 201) { // 201 Creado
            e.target.reset();      // Limpiamos el formulario
            cargarPeliculas();     // Recargamos la lista
        } else {
            alert("Error al crear la película. Código: " + response.status);
        }
    } catch (error) {
        console.error("Error en el POST:", error);
    }
});

/*
  PASO 4: IMPLEMENTAR BORRADO (DELETE)
  Usamos "Delegación de Eventos". Escuchamos clics en TODO el documento,
  pero solo reaccionamos si el clic fue en un botón 'delete-btn'.
*/
document.addEventListener("click", async (e) => {
    if (e.target.classList.contains("delete-btn")) {
        // Obtenemos el ID guardado en el 'data-id'
        const id = e.target.dataset.id; 

        if (confirm("¿Estás seguro de que quieres eliminar esta película?")) {
            try {
                const response = await fetch(`${API_URL}/${id}`, { 
                    method: "DELETE" 
                });

                if (response.status === 204) { // 204 No Content (Éxito)
                    cargarPeliculas(); // Refrescar la vista
                } else {
                    alert("Error al eliminar la película.");
                }
            } catch (error) {
                console.error("Error en el DELETE:", error);
            }
        }
    }
});
