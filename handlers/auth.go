package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/muller10000/TPE_Web_Entrega3/views"
)

// GET /login -> Muestra el formulario de login
func (h *HandlerPeliculas) HandleLoginShow(w http.ResponseWriter, r *http.Request) {
	// Si ya está autenticado, redirigir al home
	if IsAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	views.LoginPage("").Render(r.Context(), w)
}

// POST /signin -> Procesa el login
func (h *HandlerPeliculas) HandleSignin(w http.ResponseWriter, r *http.Request) {
	// 1. Leer credenciales
	username := r.FormValue("username")
	password := r.FormValue("password")

	// 2. Validar usuario contra la BD
	user, err := h.queries.GetUser(r.Context(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			views.LoginPage("Usuario no encontrado").Render(r.Context(), w)
			return
		}
		http.Error(w, "Error interno", http.StatusInternalServerError)
		return
	}

	// 2.1 Validar contraseña (texto plano para este ejemplo simple)
	if user.Password != password {
		views.LoginPage("Contraseña incorrecta").Render(r.Context(), w)
		return
	}

	// 3. Crear token de sesión único (UUID)
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second) // 2 minutos de duración

	// 4. Almacenar sesión en el mapa de memoria
	sessions[sessionToken] = Session{
		Username: username,
		Expiry:   expiresAt,
	}

	// 5. Enviar cookie al cliente
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true, // Importante para seguridad
		Path:     "/",  // Disponible en toda la aplicación
	})

	// Redirigir al home (Welcome)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// POST /logout -> Cierra la sesión
func (h *HandlerPeliculas) HandleLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// 1. Borrar la sesión del mapa
	delete(sessions, sessionToken)

	// 2. Invalidar la cookie en el cliente (fecha en el pasado)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Helper: Verifica si la petición tiene una sesión válida
func IsAuthenticated(r *http.Request) bool {
	// Obtener cookie
	c, err := r.Cookie("session_token")
	if err != nil {
		return false
	}
	sessionToken := c.Value

	// Validar sesión en el mapa
	userSession, exists := sessions[sessionToken]
	if !exists {
		return false
	}

	// Validar expiración
	if userSession.IsExpired() {
		delete(sessions, sessionToken)
		return false
	}

	return true
}

// Middleware: Protege las rutas requiriendo autenticación
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
