package handlers

import (
	"time"
)

// Estructura de Sesión
type Session struct {
	Username string
	Expiry   time.Time
}

// Mapa global para almacenar las sesiones activas en memoria
// La clave es el "session_token" (UUID)
var sessions = map[string]Session{}

// Método para verificar si la sesión ha expirado
func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
