package main

import (
	"github.com/wan6sta/url-shortener/internal/config"
	"github.com/wan6sta/url-shortener/internal/http-server/handlers/url"
	"github.com/wan6sta/url-shortener/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoadByPath("./config/local.yaml")
	log := setupLogger(cfg.Env)
	st, err := postgres.NewStorage()
	if err != nil {
		log.Error("failed to initialize db")
	}
	u := url.NewUrl(st)

	http.HandleFunc("/", u.UrlHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

/*func httpAddress(cfg *config.Config) string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}*/

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
