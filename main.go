package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/muller10000/TPE_Web_Entrega3/handlers"
	"github.com/muller10000/TPE_Web_Entrega3/repository"
)

func connectDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return sql.Open("postgres", connStr)
}

func main() {
	// Servir estáticos (CSS)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	db, err := connectDB()
	if err != nil {
		panic(err)
	}
	queries := repository.New(db)

	h := handlers.NewHandlerPeliculas(queries)

	// Definición de Rutas

	// 1. Carga Inicial (Página completa)
	http.HandleFunc("/", h.HandleIndex)

	// 2. Creación (AJAX vía HTMX)
	http.HandleFunc("POST /peliculas", h.HandleCreate)

	// 3. Eliminación (AJAX vía HTMX)
	// Manejo manual del método DELETE para compatibilidad
	http.HandleFunc("/peliculas/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			h.HandleDelete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Servidor HTMX corriendo en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error al iniciar servidor:", err)
	}
}
