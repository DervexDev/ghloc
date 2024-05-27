package rest

import "net/http"

func Unauthorized(w http.ResponseWriter, r *http.Request)  {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Access-Control-Allow-Origin", "https://github.com")
	w.Write([]byte("Unauthorized"))
}
