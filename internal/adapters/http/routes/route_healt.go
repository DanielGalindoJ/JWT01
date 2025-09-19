package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// setupHealthRoutes configura las rutas de health check
func setupHealthRoutes(r *mux.Router) {
	// @Summary Verifica estado del servicio
	// @Description Retorna estado, versión y timestamp del servicio.
	// @Tags Health
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /health [get]
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// @Summary Página principal
	// @Description Información general del Auth Service.
	// @Tags Health
	// @Produce plain
	// @Success 200 {string} string "Información del servicio"
	// @Router / [get]
	r.HandleFunc("/", homeHandler).Methods("GET")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"service":   "auth-service",
		"version":   "1.0.0",
		"timestamp": "2024-01-01T00:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	info := `Auth Service API`
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(info))
}
