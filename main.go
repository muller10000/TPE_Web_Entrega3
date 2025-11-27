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
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	db, err := connectDB()
	if err != nil {
		panic(err)
	}
	queries := repository.New(db)

	h := handlers.NewHandlerPeliculas(queries)

	// --- RUTAS PÚBLICAS ---
	http.HandleFunc("/login", h.HandleLoginShow)
	http.HandleFunc("/signin", h.HandleSignin)
	http.HandleFunc("/logout", h.HandleLogout)

	// --- RUTAS PROTEGIDAS (Middleware) ---
	// La raíz ahora está protegida
	http.HandleFunc("/", handlers.AuthMiddleware(h.HandleIndex))

	// Operaciones de Películas protegidas
	http.HandleFunc("POST /peliculas", handlers.AuthMiddleware(h.HandleCreate))

	// Manejo manual de DELETE protegido
	http.HandleFunc("/peliculas/", func(w http.ResponseWriter, r *http.Request) {
		// Validar sesión antes de nada
		if !handlers.IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if r.Method == http.MethodDelete {
			h.HandleDelete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("🚀 Servidor con Auth corriendo en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error al iniciar servidor:", err)
	}
}
