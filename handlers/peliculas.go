package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings" // Necesario para el delete

	"github.com/muller10000/TPE_Web_Entrega3/repository"
	"github.com/muller10000/TPE_Web_Entrega3/views" // Importamos tus vistas generadas
)

// Helper para manejar nulls en el POST
func valueOrEmpty(s string) string {
	return s
}

// Constructor del Handler
func NewHandlerPeliculas(queries *repository.Queries) *HandlerPeliculas {
	return &HandlerPeliculas{queries: queries}
}

type HandlerPeliculas struct {
	queries *repository.Queries
}

// 1. GET / -> Renderizar la página completa
func (h *HandlerPeliculas) HandleIndex(w http.ResponseWriter, r *http.Request) {
	// Obtener datos de la DB
	movies, err := h.queries.ListMovies(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener películas", http.StatusInternalServerError)
		return
	}

	// Renderizar componente TEMPL
	w.Header().Set("Content-Type", "text/html")
	views.IndexPage(movies).Render(r.Context(), w)
}

// 2. POST /peliculas -> Procesar Formulario y Redirigir
func (h *HandlerPeliculas) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// Parsear formulario
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error en el formulario", http.StatusBadRequest)
		return
	}

	// Extraer valores (Nota: r.FormValue devuelve string vacío si no existe)
	title := r.FormValue("title")
	director := r.FormValue("director")
	genre := r.FormValue("genre")
	rating := r.FormValue("rating")
	yearStr := r.FormValue("year")

	year, _ := strconv.Atoi(yearStr) // Manejar error en prod

	// Crear Params para sqlc
	params := repository.CreateMovieParams{
		Title:    title,
		Director: sql.NullString{String: director, Valid: director != ""},
		Year:     sql.NullInt32{Int32: int32(year), Valid: year != 0},
		Genre:    sql.NullString{String: genre, Valid: genre != ""},
		Rating:   sql.NullString{String: rating, Valid: rating != ""},
	}

	// Guardar en DB
	_, err := h.queries.CreateMovie(r.Context(), params)
	if err != nil {
		http.Error(w, "Error al guardar", http.StatusInternalServerError)
		return
	}

	// PATRÓN PRG: Redirigir al Index
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// 3. POST /peliculas/delete/{id} -> Eliminar y Redirigir
func (h *HandlerPeliculas) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// Extraer ID de la URL
	idStr := strings.TrimPrefix(r.URL.Path, "/peliculas/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Eliminar
	err = h.queries.DeleteMovie(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar", http.StatusInternalServerError)
		return
	}

	// Redirigir
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
