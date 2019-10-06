package router

import "github.com/gorilla/mux"

func addAppRoutes(r *mux.Router) {
	r.HandleFunc("/", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/health-check", HealthCheckHandler).Methods("GET")
}
