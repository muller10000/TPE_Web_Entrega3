package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/muller10000/TPE_Web_Entrega3/repository"
	"github.com/muller10000/TPE_Web_Entrega3/views"
)

func NewHandlerPeliculas(queries *repository.Queries) *HandlerPeliculas {
	return &HandlerPeliculas{queries: queries}
}

type HandlerPeliculas struct {
	queries *repository.Queries
}

// GET / Carga Inicial (Renderiza toda la página)
func (h *HandlerPeliculas) HandleIndex(w http.ResponseWriter, r *http.Request) {
	movies, err := h.queries.ListMovies(r.Context())
	if err != nil {
		http.Error(w, "Error DB", 500)
		return
	}
	// Renderizamos Layout + Form + Lista
	views.IndexPage(movies).Render(r.Context(), w)
}

// POST /peliculas Crear (HTMX)
func (h *HandlerPeliculas) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	// Mapeo de datos
	year, _ := strconv.Atoi(r.FormValue("year"))
	params := repository.CreateMovieParams{
		Title:    r.FormValue("title"),
		Director: sql.NullString{String: r.FormValue("director"), Valid: true},
		Year:     sql.NullInt32{Int32: int32(year), Valid: true},
		Genre:    sql.NullString{String: r.FormValue("genre"), Valid: true},
		Rating:   sql.NullString{String: r.FormValue("rating"), Valid: true},
	}

	// Guardar en BD
	if _, err := h.queries.CreateMovie(r.Context(), params); err != nil {
		http.Error(w, "Error al crear", 500)
		return
	}

	// LÓGICA HTMX
	// Recuperamos la lista actualizada de la BD
	movies, err := h.queries.ListMovies(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener lista", 500)
		return
	}

	// Renderizamos SOLO el componente de la lista (Fragmento)
	// HTMX tomará este HTML y reemplazará el #movie-list-container viejo
	views.MovieList(movies).Render(r.Context(), w)
}

// DELETE /peliculas/{id} Borrar (HTMX)
func (h *HandlerPeliculas) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la URL
	idStr := strings.TrimPrefix(r.URL.Path, "/peliculas/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", 400)
		return
	}

	// Borrar de la BD
	if err := h.queries.DeleteMovie(r.Context(), int32(id)); err != nil {
		http.Error(w, "Error al borrar", 500)
		return
	}

	// lÓGICA HTMX
	// Recuperamos y devolvemos la lista actualizada
	movies, err := h.queries.ListMovies(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener lista", 500)
		return
	}

	views.MovieList(movies).Render(r.Context(), w)
}
