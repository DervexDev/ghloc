package handler

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DervexDev/ghloc/src/infrastructure/github_files_provider"
	"github.com/DervexDev/ghloc/src/server/github_handler"
	"github.com/DervexDev/ghloc/src/server/rest"
	"github.com/DervexDev/ghloc/src/service/github_stat"
	"github.com/caarlos0/env/v9"
	"github.com/rs/zerolog"
)

type Config struct {
	MaxRepoSizeMB   int     `env:"MAX_REPO_SIZE_MB" envDefault:"100"`
	MaxAge		    int     `env:"MAX_AGE" envDefault:"300"`
}
 
func Handler(w http.ResponseWriter, r *http.Request) {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().Round(time.Microsecond).UTC()
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	log.SetFlags(0)
	log.SetOutput(logger)

	config := &Config{}
	if err := env.Parse(config); err != nil {
		logger.Fatal().Err(err).Msg("Error parsing config")
	}

	if r.Header.Get("Origin") != "https://github.com" {
		rest.Unauthorized(w, r)
		return
	}

	github := github_files_provider.New(config.MaxRepoSizeMB)
	service := github_stat.New(github)

	handler := &github_handler.GetStatHandler{Service: service, MaxAge: config.MaxAge}

	handler.ServeHTTP(w, r)
}
