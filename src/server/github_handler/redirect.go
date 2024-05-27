package github_handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/DervexDev/ghloc/src/server/rest"
	"github.com/DervexDev/ghloc/src/util"
	"github.com/go-chi/chi/v5"
)

type RedirectHandler struct {
}

func (h *RedirectHandler) RegisterOn(router chi.Router) {
	router.Get("/{user}/{repo}", h.ServeHTTP)
}

func getDefaultBranch(user, repo, token string) (_ string, err error) {
	defer util.WrapErr("get default branch", &err)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%v/%v", user, repo), nil)
	if err != nil {
		return "", err
	}
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return "", rest.NotFound
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %s (body: %s)", resp.Status, string(body))
	}

	repoInfo := struct {
		DefaultBranch string `json:"default_branch"`
	}{}
	if err = json.Unmarshal(body, &repoInfo); err != nil {
		return "", err
	}
	if repoInfo.DefaultBranch == "" {
		return "", fmt.Errorf("empty branch (body: %s)", string(body))
	}

	return repoInfo.DefaultBranch, nil
}

func (h RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "user")
	repo := chi.URLParam(r, "repo")
	token := r.Header.Get("Authorization")

	branch, err := getDefaultBranch(user, repo, token)
	if err != nil {
		rest.WriteResponse(w, r, err, true)
		return
	}

	url := *r.URL
	url.Path += "/" + branch
	http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
}
