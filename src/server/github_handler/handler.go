package github_handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/DervexDev/ghloc/src/server/rest"
	"github.com/DervexDev/ghloc/src/service/github_stat"
	"github.com/DervexDev/ghloc/src/service/loc_count"
)

type Service interface {
	GetStat(ctx context.Context, user, repo, branch, token string, filter, matcher *string, noLOCProvider bool, tempStorage github_stat.TempStorage) (*loc_count.StatTree, error)
}

type GetStatHandler struct {
	Service    Service
	MaxAge     int
}

func (h GetStatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")

	if len(path) < 4 {
		rest.WriteResponse(w, r, rest.BadRequest{Msg: "Invalid path"}, true)
		return
	}

	user := path[1]
	repo := path[2]
	branch := path[3]

	token := r.Header.Get("Authorization")

	r.ParseForm()

	noLOCProvider := false
	tempStorage := github_stat.TempStorageFile

	var filter *string
	if filters := r.Form["filter"]; len(filters) >= 1 {
		filter = &filters[0]
	}

	var matcher *string
	if matchers := r.Form["match"]; len(matchers) >= 1 {
		matcher = &matchers[0]
	}

	stat, err := h.Service.GetStat(r.Context(), user, repo, branch, token, filter, matcher, noLOCProvider, tempStorage)
	if err != nil {
		rest.WriteResponse(w, r, err, true)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Ghloc-Authorization, Authorization")
	w.Header().Add("Cache-Control", "public, max-age=" + strconv.Itoa(h.MaxAge))
	rest.WriteResponse(w, r, (*rest.SortedStat)(stat), r.FormValue("pretty") != "false")
}
