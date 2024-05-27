package handler

import (
	"net/http"

	"github.com/DervexDev/ghloc/src/infrastructure/github_files_provider"
	"github.com/DervexDev/ghloc/src/server/github_handler"
	"github.com/DervexDev/ghloc/src/service/github_stat"
)
 
func Handler(w http.ResponseWriter, r *http.Request) {
	github := github_files_provider.New(100)
	service := github_stat.New(github)

	handler := &github_handler.GetStatHandler{Service: service}

	handler.ServeHTTP(w, r)
}
