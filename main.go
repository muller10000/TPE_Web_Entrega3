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
	// 1. Servir archivos estáticos (CSS, Imágenes)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 2. Conexión DB
	db, err := connectDB()
	if err != nil {
		panic(err)
	}
	queries := repository.New(db)

	// 3. Inicializar Handlers
	h := handlers.NewHandlerPeliculas(queries)

	// 4. Definir Rutas (SSR)
	http.HandleFunc("/", h.HandleIndex)                   // GET: Ver lista y form
	http.HandleFunc("/peliculas", h.HandleCreate)         // POST: Crear
	http.HandleFunc("/peliculas/delete/", h.HandleDelete) // POST: Eliminar

	fmt.Println("Servidor SSR corriendo en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error al iniciar servidor:", err)
	}
}
