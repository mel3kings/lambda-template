package router

import (
	"encoding/json"
	"net/http"
	"strings"
)

type LocalHandler func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func (fn LocalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response, _ := fn(w, r)
	writeResponse(w, r, response)
}

func writeResponse(w http.ResponseWriter, r *http.Request, result interface{}) {
	if strings.EqualFold(r.Method, "POST") {
		write(w, http.StatusCreated, result)
	} else {
		write(w, http.StatusOK, result)
	}
}

func write(w http.ResponseWriter, statusCode int, result interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(result)
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Lambda server running"))
}
