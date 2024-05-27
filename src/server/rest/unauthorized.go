package rest

import (
	"net/http"
)

func  Unauthorized(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusUnauthorized)
}
